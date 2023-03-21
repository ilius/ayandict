package main

import (
	"fmt"
	"html"
	"strings"

	"github.com/ilius/ayandict/pkg/stardict"
)

func onQuery(
	query string,
	setHtml func(string),
	isAuto bool,
) {
	fmt.Printf("Query: %s\n", query)
	results := stardict.LookupHTML(query, false)
	if len(results) == 0 {
		if !isAuto {
			setHtml(fmt.Sprintf("No results for %#v", query))
		}
		return
	}
	addHistory(query)
	parts := []string{}
	for _, res := range results {
		parts = append(parts, fmt.Sprintf(
			"<h4>Dictionary: %s</h4>\n",
			html.EscapeString(res.DictName),
		))
		parts = append(parts, res.Definitions...)
	}
	htmlStr := strings.Join(parts, "\n<br/>\n")
	setHtml(htmlStr)
	// fmt.Println(htmlStr)
}
