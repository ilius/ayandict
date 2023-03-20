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
func LookupHTML(query string) []*common.QueryResult {
	results := []*common.QueryResult{}
	for _, dic := range dicList {
		definitions := []string{}
		for _, res := range dic.SearchAuto(query) {
			defi := fmt.Sprintf(
				"<b>%s</b>",
				html.EscapeString(res.Keyword),
			)
			for _, item := range res.Items {
				defi += string(item.Data) + "<br/>"
			}
			definitions = append(definitions, defi)
		}
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
