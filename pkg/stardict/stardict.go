package stardict

import (
	"fmt"
	"html"
	"net/url"
	"os"
	"path"
	"regexp"
	"runtime"
	"strconv"
	"strings"
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

func fixResURL(quoted string, resDir string) *url.URL {
	// resDir must be Unix-style
	urlStr, err := strconv.Unquote(quoted)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	_url, err := url.Parse(urlStr)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	if _url.Scheme != "" || _url.Host != "" {
		return nil
	}
	_url.Scheme = "file"
	_url.Path = resDir + "/" + _url.Path
	return _url
}

func fixSoundURL(quoted string, resDir string) *url.URL {
	// resDir must be Unix-style
	urlStr, err := strconv.Unquote(quoted)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	_url, err := url.Parse(urlStr)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	switch _url.Scheme {
	case "", "sound":
	default:
		return nil
	}
	_url.Scheme = "file"
	host := _url.Host
	_url.Host = ""
	_url.Path = resDir + "/" + host + "/" + _url.Path
	return _url
}

func fixDefiHTML(defi string, resDir string) string {
	// resDir must be Unix-style
	srcSub := func(match string) string {
		_url := fixResURL(match[5:], resDir)
		if _url == nil {
			return match
		}
		newStr := " src=" + strconv.Quote(_url.String())
		fmt.Println("srcSub:", newStr)
		return newStr
	}
	hrefSoundSub := func(match string) string {
		fmt.Println("hrefSoundSub: match:", match)
		_url := fixSoundURL(match[6:], resDir)
		if _url == nil {
			return match
		}
		newStr := " href=" + strconv.Quote(strings.TrimRight(_url.String(), "/"))
		fmt.Println("hrefSoundSub:", newStr)
		return newStr
	}

	defi = srcRE.ReplaceAllStringFunc(defi, srcSub)
	defi = hrefSoundRE.ReplaceAllStringFunc(defi, hrefSoundSub)
	return defi
}

func resourceDirUnix(dic *stardict.Dictionary) string {
	resDir := dic.ResourceDir()
	if runtime.GOOS != "windows" {
		return resDir
	}
	return "/" + strings.Replace(resDir, `\`, `/`, -1)
}

func LookupHTML(query string, title bool) []*common.QueryResult {
	results := []*common.QueryResult{}
	for _, dic := range dicList {
		definitions := []string{}
		for _, res := range dic.SearchAuto(query) {
			resDir := resourceDirUnix(dic)
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
					if resDir != "" {
						itemDefi = fixDefiHTML(itemDefi, resDir)
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
