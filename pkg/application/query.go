package application

import (
	"fmt"
	"html"
	"strings"
	"time"

	"github.com/ilius/ayandict/pkg/stardict"
)

func onQuery(
	query string,
	setHtml func(string),
	isAuto bool,
) {
	if query == "" {
		if !isAuto {
			setHtml("")
		}
		return
	}
	fmt.Printf("Query: %s\n", query)
	t := time.Now()
	results := stardict.LookupHTML(
		query,
		conf,
		dictsOrder,
	)
	fmt.Println("LookupHTML took", time.Now().Sub(t))
	if len(results) == 0 {
		if !isAuto {
			setHtml(fmt.Sprintf("No results for %#v", query))
			addHistoryAndFrequency(query)
		}
		return
	}
	addHistoryAndFrequency(query)
	parts := []string{}
	for _, res := range results {
		header := conf.HeaderTag
		if header == "" {
			header = "b"
		}
		term := html.EscapeString(strings.Join(res.Terms(), " | "))
		if conf.ShowScore {
			term += fmt.Sprintf(" [%%%d]", res.Score()/2)
		}
		// TODO: configure style of res.Term and res.DictName
		// with <span style=...>
		parts = append(parts, fmt.Sprintf(
			"<%s>%s (from %s)</%s>\n",
			header,
			term,
			html.EscapeString(res.DictName()),
			header,
		))
		parts = append(parts, res.DefinitionsHTML()...)
	}
	htmlStr := strings.Join(parts, "\n<br/>\n")
	setHtml(htmlStr)
	// fmt.Println(htmlStr)
}
