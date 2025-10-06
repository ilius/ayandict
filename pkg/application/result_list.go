package application

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/ilius/ayandict/v3/pkg/dictmgr"
	common "github.com/ilius/go-dict-commons"
	qt "github.com/mappu/miqt/qt6"
)

func NewResultListWidget(
	articleView *ArticleView,
	headerLabel *HeaderLabel,
	onResultDisplay func(terms []string),
) *ResultListWidget {
	widget := qt.NewQListWidget(nil)
	resultList := &ResultListWidget{
		QListWidget:     widget,
		HeaderLabel:     headerLabel,
		ArticleView:     articleView,
		onResultDisplay: onResultDisplay,
	}
	widget.OnCurrentRowChanged(func(row int) {
		if row < 0 {
			return
		}
		resultList.OnActivate(row)
	})
	widget.OnItemActivated(func(item *qt.QListWidgetItem) {
		row := widget.Row(item)
		if row < 0 {
			return
		}
		resultList.OnActivate(row)
	})
	return resultList
}

type ResultListWidget struct {
	*qt.QListWidget

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
