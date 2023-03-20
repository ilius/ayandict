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
) {
	fmt.Printf("Query: %s\n", query)
	results := stardict.LookupHTML(query)
	parts := []string{}
	for _, res := range results {
		parts = append(parts, fmt.Sprintf(
			"<h3>Dictionary: %s</h3>\n",
			html.EscapeString(res.DictName),
		))
		parts = append(parts, res.Definitions...)
	}
	htmlStr := strings.Join(parts, "\n<br/>\n")
	setHtml(htmlStr)
	fmt.Println(htmlStr)
}
