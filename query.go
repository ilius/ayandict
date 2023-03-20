package main

import (
	"fmt"
	"strings"

	"github.com/ilius/ayandict/pkg/stardict"
	"github.com/therecipe/qt/widgets"
)

func onQuery(query string, webview *widgets.QTextBrowser) {
	fmt.Printf("Query: %s\n", query)
	results := stardict.LookupHTML(query)
	parts := []string{}
	for _, res := range results {
		parts = append(parts, res.Definitions...)
	}
	webview.Clear()
	webview.SetHtml(strings.Join(parts, "\n"))
}
