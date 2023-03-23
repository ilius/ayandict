package stardict

import (
	"bytes"
	"fmt"
	std_html "html"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/ilius/ayandict/pkg/common"
	"github.com/ilius/ayandict/pkg/config"
	"golang.org/x/net/html"
)

var dicList []*Dictionary

var (
	srcRE       = regexp.MustCompile(` src="[^<>"]*?"`)
	hrefSoundRE = regexp.MustCompile(` href="sound://[^<>"]*?"`)
	audioRE     = regexp.MustCompile(`<audio[ >].*?</audio>`)
	sourceRE    = regexp.MustCompile(`<source [^<>]*?>`)
)

func Init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	dicDir := path.Join(homeDir, ".stardict", "dic")
	t := time.Now()
	dicList, err = Open(dicDir)
	if err != nil {
		panic(err)
	}
	fmt.Println("Loading dictionaries took", time.Now().Sub(t))
}

func fixResURL(quoted string, resURL string) (bool, string) {
	urlStr, err := strconv.Unquote(quoted)
	if err != nil {
		fmt.Println(err)
		return false, ""
	}
	_url, err := url.Parse(urlStr)
	if err != nil {
		fmt.Println(err)
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
		fmt.Println(err)
		return false, ""
	}
	return true, resURL + "/" + urlStr[len("sound://"):]
}

func fixHrefSound(defi string, resURL string) string {
	subFunc := func(match string) string {
		// fmt.Println("hrefSoundSub: match:", match)
		ok, _url := fixSoundURL(match[6:], resURL)
		if !ok {
			return match
		}
		newStr := " href=" + strconv.Quote(_url)
		// fmt.Println("hrefSoundSub:", newStr)
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

func fixAudioTag(defi string, resURL string) string {
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
				">"+filepath.Base(src)+"</a>")
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
		// fmt.Println("srcSub:", newStr)
		return newStr
	}
	return srcRE.ReplaceAllStringFunc(defi, srcSub)
}

func fixDefiHTML(defi string, resURL string, conf *config.Config) string {
	if resURL != "" {
		defi = fixFileSrc(defi, resURL)
		defi = fixHrefSound(defi, resURL)
	}
	if conf.Audio {
		defi = fixAudioTag(defi, resURL)
	}
	return defi
}

func LookupHTML(query string, title bool, conf *config.Config) []*common.QueryResult {
	results := []*common.QueryResult{}
	for _, dic := range dicList {
		definitions := []string{}
		for _, res := range dic.SearchAuto(query) {
			resURL := dic.ResourceURL()
			defi := ""
			if title {
				defi = fmt.Sprintf(
					"<b>%s</b>\n",
					std_html.EscapeString(res.Keyword),
				)
			}
			for _, item := range res.Items {
				if item.Type == 'h' {
					itemDefi := string(item.Data)
					itemDefi = fixDefiHTML(itemDefi, resURL, conf)
					defi += itemDefi + "<br/>\n"
					continue
				}
				defi += fmt.Sprintf(
					"<pre>%s</pre>\n<br/>\n",
					std_html.EscapeString(string(item.Data)),
				)
			}
			definitions = append(definitions, defi)
		}
		fmt.Printf("%d results from %s\n", len(definitions), dic.GetBookName())
		if len(definitions) == 0 {
			continue
		}
		results = append(results, &common.QueryResult{
			DictName:    dic.GetBookName(),
			Definitions: definitions,
		})
	}
	return results
}
