package main

import (
	"fmt"
	"strings"

	"github.com/ilius/ayandict/pkg/stardict"
	"github.com/therecipe/qt/widgets"
)

func onQuery(query string, textview *widgets.QPlainTextEdit) {
	fmt.Printf("Query: %s\n", query)
	results := stardict.LookupPlaintext(query)
	textview.Clear()
	parts := []string{}
	for _, res := range results {
		parts = append(parts, res.Definitions...)
	}
	textview.Document().SetPlainText(strings.Join(parts, "\n"))
}
