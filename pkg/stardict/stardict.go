package stardict

import (
	"fmt"
	"html"
	"os"
	"path"

	"github.com/ilius/ayandict/pkg/common"
	stardict "github.com/ilius/go-stardict"
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
}

func LookupHTML(query string, title bool) []*common.QueryResult {
	results := []*common.QueryResult{}
	for _, dic := range dicList {
		definitions := []string{}
		for _, res := range dic.SearchAuto(query) {
			defi := ""
			if title {
				defi = fmt.Sprintf(
					"<b>%s</b>\n",
					html.EscapeString(res.Keyword),
				)
			}
			for _, item := range res.Items {
				if item.Type == 'h' {
					defi += string(item.Data) + "<br/>\n"
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
