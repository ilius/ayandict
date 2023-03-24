package main

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
	fmt.Printf("Query: %s\n", query)
	t := time.Now()
	results := stardict.LookupHTML(query, conf)
	fmt.Println("LookupHTML took", time.Now().Sub(t))
	if len(results) == 0 {
		if !isAuto {
			setHtml(fmt.Sprintf("No results for %#v", query))
		}
		return
	}
	addHistory(query)
	parts := []string{}
	for _, res := range results {
		header := conf.HeaderTag
		if header == "" {
			header = "b"
		}
		// TODO: configure style of res.Term and res.DictName
		// with <span style=...>
		parts = append(parts, fmt.Sprintf(
			"<%s>%s (from %s)</%s>\n",
			header,
			html.EscapeString(res.Term),
			html.EscapeString(res.DictName),
			header,
		))
		parts = append(parts, res.Definitions...)
	}
	htmlStr := strings.Join(parts, "\n<br/>\n")
	setHtml(htmlStr)
	// fmt.Println(htmlStr)
}
