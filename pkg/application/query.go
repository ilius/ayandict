package application

import (
	"bytes"
	"fmt"
	"html"
	"strings"
	"time"

	"github.com/ilius/ayandict/pkg/common"
	"github.com/ilius/ayandict/pkg/stardict"
	"github.com/therecipe/qt/widgets"
)

type QueryWidgets struct {
	Webview     *widgets.QTextBrowser
	ResultList  *ResultListWidget
	HeaderLabel *widgets.QLabel
}

func NewResultListWidget(
	webview *widgets.QTextBrowser,
	headerLabel *widgets.QLabel,
) *ResultListWidget {
	widget := widgets.NewQListWidget(nil)
	resultList := &ResultListWidget{
		QListWidget: widget,
		HeaderLabel: headerLabel,
		Webview:     webview,
	}
	widget.ConnectCurrentRowChanged(func(row int) {
		if row < 0 {
			return
		}
		resultList.OnActivate(row)
	})
	widget.ConnectItemActivated(func(item *widgets.QListWidgetItem) {
		row := widget.Row(item)
		if row < 0 {
			return
		}
		resultList.OnActivate(row)
	})
	return resultList
}

type ResultListWidget struct {
	*widgets.QListWidget
	HeaderLabel *widgets.QLabel
	Webview     *widgets.QTextBrowser
	results     []common.QueryResult
}

func (w *ResultListWidget) SetResults(results []common.QueryResult) {
	w.QListWidget.Clear()
	w.results = results
	for _, res := range results {
		terms := res.Terms()
		var text string
		switch len(terms) {
		case 0:
			text = ""
			fmt.Printf("empty terms, res=%#v\n", res)
		case 1:
			text = terms[0]
		case 2:
			text = strings.Join(terms, ", ")
		default:
			text += fmt.Sprintf("%s (+%d)", terms[0], len(terms)-1)
		}
		ds := dictSettingsMap[res.DictName()]
		if ds != nil && ds.Symbol != "" {
			text = fmt.Sprintf("%s %s", text, ds.Symbol)
		}
		w.AddItem(text)
	}
	if len(results) > 0 {
		w.SetCurrentRow(0)
	}
}

type HeaderTemplateInput struct {
	Terms    []string
	Term     string
	DictName string
	Score    uint8
}

func (w *ResultListWidget) OnActivate(row int) {
	if row >= len(w.results) {
		fmt.Printf("ResultListWidget: OnActivate: row index %v out of range\n", row)
		return
	}
	res := w.results[row]
	terms := res.Terms()
	term := html.EscapeString(strings.Join(terms, " | "))
	headerBuf := bytes.NewBuffer(nil)
	err := headerTpl.Execute(headerBuf, HeaderTemplateInput{
		Terms:    terms,
		Term:     term,
		DictName: res.DictName(),
		Score:    res.Score() / 2,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	w.HeaderLabel.SetText(headerBuf.String())
	text := strings.Join(
		res.DefinitionsHTML(),
		"\n<br/>\n",
	)
	if definitionStyleString != "" {
		text = definitionStyleString + text
	}
	w.Webview.SetHtml(text)
	resDir := res.ResourceDir()
	if resDir == "" {
		w.Webview.SetSearchPaths([]string{})
	} else {
		w.Webview.SetSearchPaths([]string{resDir})
	}
}

func (w *ResultListWidget) Clear() {
	w.QListWidget.Clear()
	w.results = nil
}

func onQuery(
	query string,
	queryWidgets *QueryWidgets,
	isAuto bool,
) {
	if query == "" {
		if !isAuto {
			queryWidgets.Webview.SetHtml("")
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
	queryWidgets.ResultList.SetResults(results)
	if len(results) == 0 {
		if !isAuto {
			queryWidgets.Webview.SetHtml(fmt.Sprintf("No results for %#v", query))
			queryWidgets.HeaderLabel.SetText("")
			addHistoryAndFrequency(query)
		}
		return
	}
	addHistoryAndFrequency(query)
	// fmt.Println(htmlStr)
}
