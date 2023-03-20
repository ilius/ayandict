package main

import (
	"fmt"
	"html"
	"strings"

	"github.com/ilius/ayandict/pkg/stardict"
	"github.com/therecipe/qt/widgets"
)

func onQuery(query string, webview *widgets.QTextBrowser) {
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
	// webview.Clear()
	htmlStr := strings.Join(parts, "\n<br/>\n")
	fmt.Println(htmlStr)
	webview.SetHtml(htmlStr)
}
