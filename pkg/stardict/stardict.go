package stardict

import (
	"fmt"
	"html"
	"os"
	"path"

	"github.com/ilius/ayandict/pkg/common"
	stardict "github.com/ilius/go-stardict"
	"github.com/ozeidan/fuzzy-patricia/patricia"
)

var dicList []*stardict.Dictionary

func Init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	dicDir := path.Join(homeDir, ".stardict", "dic")
	dicList, err = stardict.Open(dicDir)
	if err != nil {
		panic(err)
	}
	BuildFuzzyTrie(dicList)
}

func LookupHTML(query string, title bool) []*common.QueryResult {
	results := []*common.QueryResult{}
	caseInsensitive := true

	checked := map[string]bool{}

	visitor := func(prefix patricia.Prefix, dicIn patricia.Item) error {
		keyword := string(prefix)
		if checked[keyword] {
			return nil
		}
		checked[keyword] = true
		dic := dicIn.(*stardict.Dictionary)
		defi := ""
		if title {
			defi = fmt.Sprintf(
				"<b>%s</b>\n",
				html.EscapeString(keyword),
			)
		}
		for _, translation := range dic.Translate(keyword) {
			for _, item := range translation.Parts {
				if item.Type == 'h' {
					defi += string(item.Data) + "<br/>\n"
					continue
				}
				defi += fmt.Sprintf(
					"<pre>%s</pre>\n<br/>\n",
					html.EscapeString(string(item.Data)),
				)
			}
		}
		results = append(results, &common.QueryResult{
			DictName:    dic.GetBookName(),
			Definitions: []string{defi},
		})
		return nil
	}

	// trie.VisitFuzzy(
	// 	patricia.Prefix(query),
	// 	caseInsensitive,
	// 	func(prefix patricia.Prefix, dicIn patricia.Item, skipped int) error {
	// 		return visitor(prefix, dicIn)
	// 	},
	// )
	trie.VisitSubtree(
		patricia.Prefix(query),
		visitor,
	)
	trie.VisitSubstring(
		patricia.Prefix(query),
		caseInsensitive,
		visitor,
	)

	return results
}
