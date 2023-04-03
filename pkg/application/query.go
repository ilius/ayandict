package application

import (
	"bytes"
	"fmt"
	"html"
	"log"
	"strings"
	"time"

	"github.com/ilius/ayandict/pkg/common"
	"github.com/ilius/ayandict/pkg/stardict"
	"github.com/therecipe/qt/widgets"
)

type QueryWidgets struct {
	ArticleView *ArticleView
	ResultList  *ResultListWidget
	HeaderLabel *widgets.QLabel
	HistoryView *HistoryView
}

func (w *QueryWidgets) AddHistoryAndFrequency(query string) {
	if !conf.HistoryDisable {
		w.HistoryView.AddHistory(query)
	}
	if !conf.MostFrequentDisable {
		frequencyTable.Add(query, 1)
		if conf.MostFrequentAutoSave {
			SaveFrequency()
		}
	}
}

func NewResultListWidget(
	articleView *ArticleView,
	headerLabel *widgets.QLabel,
) *ResultListWidget {
	widget := widgets.NewQListWidget(nil)
	resultList := &ResultListWidget{
		QListWidget: widget,
		HeaderLabel: headerLabel,
		ArticleView: articleView,
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
	ArticleView *ArticleView
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
			log.Printf("empty terms, res=%#v\n", res)
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
		log.Printf("ResultListWidget: OnActivate: row index %v out of range\n", row)
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
		log.Println(err)
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
	w.ArticleView.SetHtml(text)
	resDir := res.ResourceDir()
	if resDir == "" {
		w.ArticleView.SetSearchPaths([]string{})
	} else {
		w.ArticleView.SetSearchPaths([]string{resDir})
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
			queryWidgets.ArticleView.SetHtml("")
			queryWidgets.HeaderLabel.SetText("")
		}
		return
	}
	log.Printf("Query: %s\n", query)
	t := time.Now()
	results := stardict.LookupHTML(
		query,
		conf,
		dictsOrder,
	)
	log.Println("LookupHTML took", time.Now().Sub(t))
	queryWidgets.ResultList.SetResults(results)
	if len(results) == 0 {
		if !isAuto {
			queryWidgets.ArticleView.SetHtml(fmt.Sprintf("No results for %#v", query))
			queryWidgets.HeaderLabel.SetText("")
			queryWidgets.AddHistoryAndFrequency(query)
		}
		return
	}
	queryWidgets.AddHistoryAndFrequency(query)
	// log.Println(htmlStr)
}
