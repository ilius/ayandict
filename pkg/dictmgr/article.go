package dictmgr

import (
	"bytes"
	"fmt"
	std_html "html"
	"log/slog"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/ilius/ayandict/v2/pkg/config"
	"github.com/ilius/ayandict/v2/pkg/html"
	"github.com/ilius/ayandict/v2/pkg/qtcommon/qerr"
	common "github.com/ilius/go-dict-commons"
)

var (
	srcRE = regexp.MustCompile(` src="[^<>"]*?"`)

	emptySoundRE = regexp.MustCompile(`<a [^<>]*href="sound://[^<>"]*?"></a>`)
	hrefSoundRE  = regexp.MustCompile(` href="sound://[^<>"]*?"`)
	audioRE      = regexp.MustCompile(`<audio[ >].*?</audio>`)

	linkRE = regexp.MustCompile(`<link [^<>]+>`)

	colorRE      = regexp.MustCompile(` color=["']#?[a-zA-Z0-9]+["']`)
	styleColorRE = regexp.MustCompile(`color:#?[a-zA-Z0-9]+`)

	hrefBwordRE = regexp.MustCompile(` href="bword://[^<>"]*?"`)
)

const (
	webPlayImage  = "/web/audio-play.png"
	playImageName = "audio-play.png"
)

var playImagePath = filepath.Join(config.GetCacheDir(), playImageName)

type DictProcessor struct {
	common.Dictionary
	conf  *config.Config
	flags uint32

	playImageMutex sync.Mutex
}

func (p *DictProcessor) dictResURL(relPath string) string {
	if p.flags&common.ResultFlag_Web > 0 {
		values := url.Values{}
		values.Add("dictName", p.DictName())
		values.Add("path", relPath)
		return DictResPathBase + "?" + values.Encode()
	}
	return p.ResourceURL() + "/" + relPath
}

func (p *DictProcessor) fixResURL(quoted string) (bool, string) {
	urlStr, err := strconv.Unquote(quoted)
	if err != nil {
		slog.Error("error", "err", err)
		return false, ""
	}
	if urlStr == webPlayImage {
		return false, ""
	}
	_url, err := url.Parse(urlStr)
	if err != nil {
		slog.Error("error", "err", err)
		return false, ""
	}
	if _url.Scheme != "" || _url.Host != "" {
		return false, ""
	}
	return true, p.dictResURL(_url.Path)
}

func (p *DictProcessor) fixSoundURL(quoted string) (bool, string) {
	urlStr, err := strconv.Unquote(quoted)
	if err != nil {
		slog.Error("error", "err", err)
		return false, ""
	}
	return true, p.dictResURL(urlStr[len("sound://"):])
}

func (*DictProcessor) fixEmptySoundLink(defi string, playImg string) string {
	subFunc := func(match string) string {
		return match[:len(match)-4] + playImg + "</a>"
	}
	return emptySoundRE.ReplaceAllStringFunc(defi, subFunc)
}

func (p *DictProcessor) fixHrefSound(defi string) string {
	subFunc := func(match string) string {
		// slog.Info("hrefSoundSub: match:", match)
		ok, _url := p.fixSoundURL(match[6:])
		if !ok {
			return match
		}
		newStr := " href=" + strconv.Quote(_url)
		// slog.Info("hrefSoundSub:", newStr)
		return newStr
	}
	return hrefSoundRE.ReplaceAllStringFunc(defi, subFunc)
}

func (*DictProcessor) findParsedTags(node *html.Node, tagName string) []*html.Node {
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

func (*DictProcessor) getAttr(node *html.Node, attrName string) string {
	for _, attr := range node.Attr {
		if attr.Key == attrName {
			return attr.Val
		}
	}
	return ""
}

func (p *DictProcessor) fixAudioTag(
	defi string,
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
		for _, sourceNode := range p.findParsedTags(root, "source") {
			src := p.getAttr(sourceNode, "src")
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

func (p *DictProcessor) fixFileSrc(defi string) string {
	srcSub := func(match string) string {
		ok, _url := p.fixResURL(match[5:])
		if !ok {
			return match
		}
		newStr := " src=" + strconv.Quote(_url)
		// slog.Info("srcSub:", newStr)
		return newStr
	}
	return srcRE.ReplaceAllStringFunc(defi, srcSub)
}

// hrefBwordSub fixes several problems, working around qt bugs
// on handling links
// problem 1: href value has space
// for example: <a href="bword://abscisic acid">
// clicking on these link do not work
// ConnectAnchorClicked will get an empty url
// link.ToString(core.QUrl__None) == ""
// unless I remove `bword://` prefix
// also tried replacing space with %20
// problem 2: href value has quoted unicode characters, using &#...;
// like "fl&#x205;k" for "flȅk", when you click on link, qt redirects to
// a non-sense term, and does not even emit AnchorClicked signal
func (p *DictProcessor) hrefBwordSub(match string) string {
	ref := match[len(` href="bword://`):]
	ref = html.UnescapeString(ref)
	return ` href="` + ref
}

func (p *DictProcessor) embedExternalStyle(defi string) string {
	const pre = len(` href=`)
	resDir := p.ResourceDir()

	subFunc := func(match string) string {
		i := strings.Index(match, ` href=`)
		if i < 0 {
			return match
		}
		q_href := match[i+pre:]
		j := strings.Index(q_href[1:], q_href[:1])
		if j < 0 {
			slog.Error("linkSub: did not find quote end in q_href", "q_href", q_href)
			return match
		}
		href := q_href[1 : j+1]
		if strings.Contains(href, "://") {
			// TODO: download?
			return match
		}
		data, err := os.ReadFile(filepath.Join(resDir, href))
		if err != nil {
			if !os.IsNotExist(err) {
				slog.Error("error", "err", err)
			}
			return match
		}
		return fmt.Sprintf("<style>\n%s\n</style>", string(data))
	}

	defi = linkRE.ReplaceAllStringFunc(defi, subFunc)
	return defi
}

func (p *DictProcessor) applyColorMapping(defi string) string {
	colorMapping := p.conf.ColorMapping
	if len(colorMapping) == 0 {
		return defi
	}
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

func (p *DictProcessor) createPlayImage() bool {
	_, statErr := os.Stat(playImagePath)
	if statErr == nil {
		return true
	}

	p.playImageMutex.Lock()
	defer p.playImageMutex.Unlock()

	data, err := res.ReadFile("res/" + playImageName)
	if err != nil {
		qerr.Error(err)
		return false
	}
	err = os.WriteFile(playImagePath, data, 0o644)
	if err != nil {
		qerr.Error(err)
		return false
	}
	return true
}

func (p *DictProcessor) getPlayImage() string {
	if p.flags&common.ResultFlag_Web > 0 {
		return fmt.Sprintf(
			`<img src="%s" />`, webPlayImage,
		)
	}
	if !p.createPlayImage() {
		return ""
	}

	_url := url.URL{}
	_url.Scheme = "file"
	_url.Path = playImagePath
	_urlStr := _url.String()
	return fmt.Sprintf(
		`<img src=%s />`,
		strconv.Quote(_urlStr),
	)
}

func (p *DictProcessor) FixDefiHTML(defi string) string {
	conf := p.conf
	flags := p.flags
	var playImage string
	hasResource := p.ResourceDir() != ""
	_fixAudio := conf.Audio && flags&common.ResultFlag_FixAudio > 0
	if _fixAudio {
		playImage = p.getPlayImage()
		defi = p.fixEmptySoundLink(defi, playImage)
		if hasResource {
			defi = p.fixHrefSound(defi)
		}
	}
	if hasResource && flags&common.ResultFlag_FixFileSrc > 0 {
		defi = p.fixFileSrc(defi)
	}
	if _fixAudio {
		defi = p.fixAudioTag(defi, playImage)
	}
	if conf.EmbedExternalStylesheet {
		defi = p.embedExternalStyle(defi)
	}
	if flags&common.ResultFlag_FixWordLink > 0 {
		defi = hrefBwordRE.ReplaceAllStringFunc(defi, p.hrefBwordSub)
	}
	if flags&common.ResultFlag_ColorMapping > 0 {
		defi = p.applyColorMapping(defi)
	}
	return defi
}
