package stardict

import (
	"fmt"
	"html"
	"net/url"
	"os"
	"path"
	"regexp"
	"strconv"
	"time"

	"github.com/ilius/ayandict/pkg/common"
	stardict "github.com/ilius/go-stardict"
)

var dicList []*stardict.Dictionary

var (
	srcRE       = regexp.MustCompile(` src="[^<>"]*?"`)
	hrefSoundRE = regexp.MustCompile(` href="sound://[^<>"]*?"`)
)

func Init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	dicDir := path.Join(homeDir, ".stardict", "dic")
	t := time.Now()
	dicList, err = stardict.Open(dicDir)
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
	_url, err := url.Parse(urlStr)
	if err != nil {
		fmt.Println(err)
		return false, ""
	}
	switch _url.Scheme {
	case "", "sound":
	default:
		return false, ""
	}
	return true, resURL + "/" + _url.Host + "/" + _url.Path
}

func fixDefiHTML(defi string, resURL string) string {
	srcSub := func(match string) string {
		ok, _url := fixResURL(match[5:], resURL)
		if !ok {
			return match
		}
		newStr := " src=" + strconv.Quote(_url)
		fmt.Println("srcSub:", newStr)
		return newStr
	}
	hrefSoundSub := func(match string) string {
		fmt.Println("hrefSoundSub: match:", match)
		ok, _url := fixSoundURL(match[6:], resURL)
		if !ok {
			return match
		}
		newStr := " href=" + strconv.Quote(_url)
		fmt.Println("hrefSoundSub:", newStr)
		return newStr
	}

	defi = srcRE.ReplaceAllStringFunc(defi, srcSub)
	defi = hrefSoundRE.ReplaceAllStringFunc(defi, hrefSoundSub)
	return defi
}

func LookupHTML(query string, title bool) []*common.QueryResult {
	results := []*common.QueryResult{}
	for _, dic := range dicList {
		definitions := []string{}
		for _, res := range dic.SearchAuto(query) {
			resURL := dic.ResourceURL()
			defi := ""
			if title {
				defi = fmt.Sprintf(
					"<b>%s</b>\n",
					html.EscapeString(res.Keyword),
				)
			}
			for _, item := range res.Items {
				if item.Type == 'h' {
					itemDefi := string(item.Data)
					if resURL != "" {
						itemDefi = fixDefiHTML(itemDefi, resURL)
					}
					defi += itemDefi + "<br/>\n"
					continue
				}
				defi += fmt.Sprintf(
					"<pre>%s</pre>\n<br/>\n",
					html.EscapeString(string(item.Data)),
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
