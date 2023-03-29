package application

import (
	"fmt"
	"html"
	"strings"
	"time"

	"github.com/ilius/ayandict/pkg/common"
	"github.com/ilius/ayandict/pkg/stardict"
	"github.com/therecipe/qt/widgets"
)

type QueryWidgets struct {
	Webview    *widgets.QTextBrowser
	ResultList *ResultListWidget
}

func NewResultListWidget(
	webview *widgets.QTextBrowser,
	titleLabel *widgets.QLabel,
) *ResultListWidget {
	widget := widgets.NewQListWidget(nil)
	resultList := &ResultListWidget{
		QListWidget: widget,
		TitleLabel:  titleLabel,
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
	TitleLabel *widgets.QLabel
	Webview    *widgets.QTextBrowser
	results    []common.QueryResult
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
		w.AddItem(text)
	}
	if len(results) > 0 {
		w.SetCurrentRow(0)
	}
}

func (w *ResultListWidget) OnActivate(row int) {
	if row >= len(w.results) {
		fmt.Printf("ResultListWidget: OnActivate: row index %v out of range\n", row)
		return
	}
	// row := item.
	res := w.results[row]
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
	w.TitleLabel.SetText(fmt.Sprintf(
		"<%s>%s (from %s)</%s>\n",
		header,
		term,
		html.EscapeString(res.DictName()),
		header,
	))
	w.Webview.SetHtml(strings.Join(
		res.DefinitionsHTML(),
		"\n<br/>\n",
	))
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
			addHistoryAndFrequency(query)
		}
		return
	}
	addHistoryAndFrequency(query)
	// fmt.Println(htmlStr)
}
