package stardict

import (
	"bytes"
	"fmt"
	std_html "html"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/ilius/ayandict/pkg/common"
	"github.com/ilius/ayandict/pkg/config"
	"github.com/ilius/ayandict/pkg/html"
)

var (
	srcRE       = regexp.MustCompile(` src="[^<>"]*?"`)
	hrefSoundRE = regexp.MustCompile(` href="sound://[^<>"]*?"`)
	audioRE     = regexp.MustCompile(`<audio[ >].*?</audio>`)
	linkRE      = regexp.MustCompile(`<link [^<>]+>`)

	hrefBwordSpaceRE = regexp.MustCompile(` href="bword://[^<>"]*?( |%20)[^<>" ]*?"`)
)

var dicList []*Dictionary

type QueryResultImp struct {
	*SearchResult
	dic    *Dictionary
	conf   *config.Config
	hDefis []string
}

func (r *QueryResultImp) DictName() string {
	return r.dic.DictName()
}

func (r *QueryResultImp) Score() uint8 {
	return r.score
}

func (r *QueryResultImp) Terms() []string {
	return r.terms
}

func (r *QueryResultImp) DefinitionsHTML() []string {
	if r.hDefis != nil {
		return r.hDefis
	}
	definitions := []string{}
	resURL := r.dic.ResourceURL()
	for _, item := range r.items() {
		if item.Type == 'h' {
			itemDefi := string(item.Data)
			itemDefi = fixDefiHTML(itemDefi, resURL, r.conf, r.dic)
			definitions = append(definitions, itemDefi+"<br/>\n")
			continue
		}
		definitions = append(definitions, fmt.Sprintf(
			"<pre>%s</pre>\n<br/>\n",
			std_html.EscapeString(string(item.Data)),
		))
	}
	r.hDefis = definitions
	return definitions
}

func (r *QueryResultImp) ResourceDir() string {
	return r.dic.resDir
}

type DicListSorter struct {
	Order map[string]int
	List  []*Dictionary
}

func (s DicListSorter) Len() int {
	return len(s.List)
}

func (s DicListSorter) Swap(i, j int) {
	s.List[i], s.List[j] = s.List[j], s.List[i]
}

func absInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func (s DicListSorter) Less(i, j int) bool {
	return absInt(s.Order[s.List[i].DictName()]) < absInt(s.Order[s.List[j].DictName()])
}

func Init(directoryList []string, order map[string]int) {
	t := time.Now()
	var err error
	dicList, err = Open(directoryList, order)
	if err != nil {
		panic(err)
	}
	fmt.Println("Loading dictionaries took", time.Now().Sub(t))
	if order != nil {
		Reorder(order)
	}
}

func GetInfoList() []Info {
	infos := make([]Info, len(dicList))
	for i, dic := range dicList {
		infos[i] = *dic.info
	}
	return infos
}

func Reorder(order map[string]int) {
	sort.Sort(DicListSorter{
		List:  dicList,
		Order: order,
	})
}

func ApplyDictsOrder(order map[string]int) {
	Reorder(order)
	for _, dic := range dicList {
		disabled := dic.disabled
		dic.disabled = order[dic.DictName()] < 0
		if disabled && !dic.disabled {
			dic.load()
		}
	}
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
			fmt.Printf("linkSub: did not find quote end in q_href=%#v\n", q_href)
			return match
		}
		href := q_href[1 : j+1]
		// fmt.Printf("linkSub: href=%#v\n", href)
		if strings.Contains(href, "://") {
			// TODO: download?
			return match
		}
		data, err := ioutil.ReadFile(filepath.Join(resDir, href))
		if err != nil {
			if !os.IsNotExist(err) {
				fmt.Println(err)
			}
			return match
		}
		return fmt.Sprintf("<style>\n%s\n</style>", string(data))
	}

	defi = linkRE.ReplaceAllStringFunc(defi, subFunc)
	// fmt.Println(defi)
	return defi
}

func fixDefiHTML(
	defi string,
	resURL string,
	conf *config.Config,
	dic *Dictionary,
) string {
	if resURL != "" {
		defi = fixFileSrc(defi, resURL)
		defi = fixHrefSound(defi, resURL)
	}
	if conf.Audio {
		defi = fixAudioTag(defi, resURL)
	}
	if conf.EmbedExternalStylesheet {
		defi = embedExternalStyle(defi, dic.resDir)
	}
	defi = hrefBwordSpaceRE.ReplaceAllStringFunc(defi, hrefBwordSpaceSub)
	return defi
}

func LookupHTML(
	query string,
	conf *config.Config,
	dictsOrder map[string]int,
) []common.QueryResult {
	results := []common.QueryResult{}
	maxResultsPerDict := conf.MaxResultsPerDict
	for _, dic := range dicList {
		if dic.disabled {
			continue
		}
		for _, res := range dic.Search(query, maxResultsPerDict) {
			results = append(results, &QueryResultImp{
				SearchResult: res,
				dic:          dic,
				conf:         conf,
			})
		}
	}
	sort.Slice(results, func(i, j int) bool {
		res1 := results[i]
		res2 := results[j]
		score1 := res1.Score()
		score2 := res2.Score()
		if score1 != score2 {
			return score1 > score2
		}
		return dictsOrder[res1.DictName()] < dictsOrder[res2.DictName()]
	})
	cutoff := conf.MaxResultsTotal
	if cutoff > 0 && len(results) > cutoff {
		results = results[:cutoff]
	}
	return results
}
