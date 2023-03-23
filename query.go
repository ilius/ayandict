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
		header := conf.DictHeaderTag
		if header != "" {
			parts = append(parts, fmt.Sprintf(
				"<%s>Dictionary: %s</%s>\n",
				header,
				html.EscapeString(res.DictName),
				header,
			))
		}
		parts = append(parts, res.Definitions...)
	}
	htmlStr := strings.Join(parts, "\n<br/>\n")
	setHtml(htmlStr)
	// fmt.Println(htmlStr)
}
