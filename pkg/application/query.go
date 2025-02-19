package application

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/ilius/ayandict/v2/pkg/dictmgr"
	common "github.com/ilius/go-dict-commons"
	"github.com/ilius/qt/widgets"
)

const resultFlags = uint32(
	common.ResultFlag_FixAudio |
		common.ResultFlag_FixFileSrc |
		common.ResultFlag_FixWordLink |
		common.ResultFlag_ColorMapping)

type QueryArgs struct {
	ArticleView *ArticleView
	ResultList  *ResultListWidget
	HeaderLabel *HeaderLabel
	HistoryView *HistoryView
	PostQuery   func(string)
	Entry       *widgets.QLineEdit
	ModeCombo   *widgets.QComboBox
}

func (w *QueryArgs) AddHistoryAndFrequency(query string) {
	if !conf.HistoryDisable {
		w.HistoryView.Add(query)
	}
	if !conf.MostFrequentDisable {
		frequencyTable.Add(query, 1)
		if conf.MostFrequentAutoSave {
			frequencyTable.SaveNoError()
		}
	}
}

func (w *QueryArgs) SetNoResult(query string) {
	w.ArticleView.SetHtml(fmt.Sprintf("No results for %#v", query))
	w.HeaderLabel.SetText("")
	w.AddHistoryAndFrequency(query)
}

func NewResultListWidget(
	articleView *ArticleView,
	headerLabel *HeaderLabel,
	onResultDisplay func(terms []string),
) *ResultListWidget {
	widget := widgets.NewQListWidget(nil)
	resultList := &ResultListWidget{
		QListWidget:     widget,
		HeaderLabel:     headerLabel,
		ArticleView:     articleView,
		onResultDisplay: onResultDisplay,
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

	results []common.SearchResultIface

	Active common.SearchResultIface

	HeaderLabel *HeaderLabel
	ArticleView *ArticleView

	onResultDisplay func(terms []string)
}

func (w *ResultListWidget) SetResults(results []common.SearchResultIface) {
	w.QListWidget.Clear()
	w.results = results
	for _, res := range results {
		terms := res.Terms()
		var text string
		switch len(terms) {
		case 0:
			text = ""
			slog.Error("empty terms", "res", res)
		case 1:
			text = terms[0]
		case 2:
			text = strings.Join(terms, ", ")
		default:
			text += fmt.Sprintf("%s (+%d)", terms[0], len(terms)-1)
		}
		symbol := dictmgr.DictSymbol(res.DictName())
		if symbol != "" {
			text = fmt.Sprintf("%s %s", text, symbol)
		}
		w.AddItem(text)
	}
	if len(results) > 0 {
		w.SetCurrentRow(0)
	}
}

func (w *ResultListWidget) OnActivate(row int) {
	if row >= len(w.results) {
		slog.Error("ResultListWidget: OnActivate: row index out of range", "row", row)
		return
	}
	res := w.results[row]
	w.HeaderLabel.SetResult(res)
	w.ArticleView.SetResult(res)
	resDir := res.ResourceDir()
	if resDir == "" {
		w.ArticleView.SetSearchPaths([]string{})
	} else {
		w.ArticleView.SetSearchPaths([]string{resDir})
	}
	w.onResultDisplay(res.Terms())
	w.Active = res
}

func (w *ResultListWidget) Clear() {
	w.QListWidget.Clear()
	w.results = nil
}

func onQuery(
	query string,
	queryArgs *QueryArgs,
	isAuto bool,
) {
	if query == "" {
		if !isAuto {
			queryArgs.ArticleView.SetHtml("")
			queryArgs.HeaderLabel.SetText("")
		}
		return
	}
	t := time.Now()
	mode := dictmgr.QueryModeFuzzy
	switch queryArgs.ModeCombo.CurrentIndex() {
	case 1:
		mode = dictmgr.QueryModeStartWith
	case 2:
		mode = dictmgr.QueryModeRegex
	case 3:
		mode = dictmgr.QueryModeGlob
	}
	results := dictmgr.LookupHTML(query, conf, mode, resultFlags, 0)
	slog.Debug("LookupHTML running time", "dt", time.Since(t), "query", query)
	queryArgs.ResultList.SetResults(results)
	if len(results) == 0 {
		if !isAuto {
			queryArgs.SetNoResult(query)
		}
	}
	if isAuto {
		if len(results) > 0 {
			if results[0].Score() == 200 {
				queryArgs.AddHistoryAndFrequency(query)
			}
		}
	} else {
		queryArgs.AddHistoryAndFrequency(query)
	}
	queryArgs.PostQuery(query)
}
