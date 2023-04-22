package dictmgr

import (
	"bytes"
	"fmt"
	std_html "html"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/ilius/ayandict/pkg/config"
	"github.com/ilius/ayandict/pkg/html"
	common "github.com/ilius/go-dict-commons"
	"github.com/ilius/qt/core"
)

var (
	srcRE = regexp.MustCompile(` src="[^<>"]*?"`)

	emptySoundRE = regexp.MustCompile(`<a [^<>]*href="sound://[^<>"]*?"></a>`)
	hrefSoundRE  = regexp.MustCompile(` href="sound://[^<>"]*?"`)
	audioRE      = regexp.MustCompile(`<audio[ >].*?</audio>`)

	linkRE = regexp.MustCompile(`<link [^<>]+>`)

	colorRE      = regexp.MustCompile(` color=["']#?[a-zA-Z0-9]+["']`)
	styleColorRE = regexp.MustCompile(`color:#?[a-zA-Z0-9]+`)

	hrefBwordSpaceRE = regexp.MustCompile(` href="bword://[^<>"]*?( |%20)[^<>" ]*?"`)
)

func fixResURL(quoted string, resURL string) (bool, string) {
	urlStr, err := strconv.Unquote(quoted)
	if err != nil {
		log.Println(err)
		return false, ""
	}
	_url, err := url.Parse(urlStr)
	if err != nil {
		log.Println(err)
		return false, ""
	}
	if _url.Scheme != "" || _url.Host != "" {
		return false, ""
	}
	return true, resURL + "/" + _url.Path
}

func fixSoundURL(quoted string, resURL string) (bool, string) {
	urlStr, err := strconv.Unquote(quoted)
	if err != nil {
		log.Println(err)
		return false, ""
	}
	return true, resURL + "/" + urlStr[len("sound://"):]
}

func fixEmptySoundLink(defi string, playImg string) string {
	subFunc := func(match string) string {
		return match[:len(match)-4] + playImg + "</a>"
	}
	return emptySoundRE.ReplaceAllStringFunc(defi, subFunc)
}

func fixHrefSound(defi string, resURL string) string {
	subFunc := func(match string) string {
		// log.Println("hrefSoundSub: match:", match)
		ok, _url := fixSoundURL(match[6:], resURL)
		if !ok {
			return match
		}
		newStr := " href=" + strconv.Quote(_url)
		// log.Println("hrefSoundSub:", newStr)
		return newStr
	}
	return hrefSoundRE.ReplaceAllStringFunc(defi, subFunc)
}

func findParsedTags(node *html.Node, tagName string) []*html.Node {
	result := []*html.Node{}

	var recurse func(argNode *html.Node)

	recurse = func(argNode *html.Node) {
		if argNode.Data == tagName {
			result = append(result, argNode)
			return
		}
		child := argNode.FirstChild
		for child != nil {
			recurse(child)
			child = child.NextSibling
		}
	}

	recurse(node)
	return result
}

func getAttr(node *html.Node, attrName string) string {
	for _, attr := range node.Attr {
		if attr.Key == attrName {
			return attr.Val
		}
	}
	return ""
}

func fixAudioTag(
	defi string,
	resURL string,
	playImage string,
) string {
	// fix <audio ...><source src="..."></audio>
	// value for src= is already fixed
	// just need to replace `<source ...>` with `<a ...>Audio</a>`
	// and remove <audio ...> and </audio>
	// but I decided to use an html parsing library
	// I only found golang.org/x/net/html
	// which always starts from <html><body>...
	// so I had to resursivly find all <source> tags
	// and the extract src attributes
	// there might be multiple <source> tags, with mp3, ogg etc
	// but QMediaPlayer does not play ogg for me
	subFunc := func(match string) string {
		root, err := html.Parse(bytes.NewBufferString(match))
		if err != nil {
			return match
		}
		parts := []string{}
		for _, sourceNode := range findParsedTags(root, "source") {
			src := getAttr(sourceNode, "src")
			if src == "" {
				continue
			}
			parts = append(parts, "<a href="+
				strconv.Quote(std_html.EscapeString(src))+
				">"+playImage+"</a>")
		}
		return strings.Join(parts, ", ")
	}
	defi = audioRE.ReplaceAllStringFunc(defi, subFunc)
	return defi
}

func fixFileSrc(defi string, resURL string) string {
	srcSub := func(match string) string {
		ok, _url := fixResURL(match[5:], resURL)
		if !ok {
			return match
		}
		newStr := " src=" + strconv.Quote(_url)
		// log.Println("srcSub:", newStr)
		return newStr
	}
	return srcRE.ReplaceAllStringFunc(defi, srcSub)
}

// work around qt bug on internal entry links with space
// for example: <a href="bword://abscisic acid">
// clicking on these link do not work
// ConnectAnchorClicked will get an empty url
// link.ToString(core.QUrl__None) == ""
// unless I remove `bword://` prefix
// also tried replacing space with %20
func hrefBwordSpaceSub(match string) string {
	return ` href="` + match[len(` href="bword://`):]
}

func embedExternalStyle(defi string, resDir string) string {
	const pre = len(` href=`)

	subFunc := func(match string) string {
		i := strings.Index(match, ` href=`)
		if i < 0 {
			return match
		}
		q_href := match[i+pre:]
		j := strings.Index(q_href[1:], q_href[:1])
		if j < 0 {
			log.Printf("linkSub: did not find quote end in q_href=%#v\n", q_href)
			return match
		}
		href := q_href[1 : j+1]
		// log.Printf("linkSub: href=%#v\n", href)
		if strings.Contains(href, "://") {
			// TODO: download?
			return match
		}
		data, err := ioutil.ReadFile(filepath.Join(resDir, href))
		if err != nil {
			if !os.IsNotExist(err) {
				log.Println(err)
			}
			return match
		}
		return fmt.Sprintf("<style>\n%s\n</style>", string(data))
	}

	defi = linkRE.ReplaceAllStringFunc(defi, subFunc)
	// log.Println(defi)
	return defi
}

func applyColorMapping(defi string, colorMapping map[string]string) string {
	colorSub := func(match string) string {
		key := match[len(` color="`) : len(match)-1]
		if key == "" {
			return match
		}
		if key[0] == '#' {
			key = key[1:]
		}
		color, ok := colorMapping[key]
		if !ok || color == "" {
			return match
		}
		return match[:len(` color="`)] + color + match[len(match)-1:]
	}
	styleColorSub := func(match string) string {
		key := match[len("color:"):]
		if key == "" {
			return match
		}
		if key[0] == '#' {
			key = key[1:]
		}
		color, ok := colorMapping[key]
		if !ok || color == "" {
			return match
		}
		return "color:" + color
	}

	defi = colorRE.ReplaceAllStringFunc(defi, colorSub)
	defi = styleColorRE.ReplaceAllStringFunc(defi, styleColorSub)
	return defi
}

func getPlayImage() string {
	imgPath, err := loadPNGFile("audio-play.png")
	if err != nil {
		log.Println(err)
	}
	qUrl := core.NewQUrl()
	qUrl.SetScheme("file")
	qUrl.SetPath(imgPath, core.QUrl__TolerantMode)
	return fmt.Sprintf(
		`<img src=%s />`,
		strconv.Quote(qUrl.ToString(core.QUrl__None)),
	)
}

func fixDefiHTML(
	defi string,
	resURL string,
	conf *config.Config,
	dic common.Dictionary,
) string {
	var playImage string
	if conf.Audio {
		playImage = getPlayImage()
		defi = fixEmptySoundLink(defi, playImage)
		if resURL != "" {
			defi = fixHrefSound(defi, resURL)
		}
	}
	if resURL != "" {
		defi = fixFileSrc(defi, resURL)
	}
	if conf.Audio {
		defi = fixAudioTag(defi, resURL, playImage)
	}
	if conf.EmbedExternalStylesheet {
		defi = embedExternalStyle(defi, dic.ResourceDir())
	}
	defi = hrefBwordSpaceRE.ReplaceAllStringFunc(defi, hrefBwordSpaceSub)
	if len(conf.ColorMapping) > 0 {
		defi = applyColorMapping(defi, conf.ColorMapping)
	}
	return defi
}
