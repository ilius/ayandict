package resultlist

import (
	"fmt"
	"log/slog"
	"strings"

	common "codeberg.org/ilius/go-dict-commons"
	"github.com/ilius/ayandict/v3/pkg/dictmgr"
	qt "github.com/mappu/miqt/qt6"
)

type ApplicationIface interface {
	OnResultDisplay(terms []string)
	HasFavorite(term string) bool
}

type HeaderLabelIface interface {
	SetResult(common.SearchResultIface)
}

type ArticleViewIface interface {
	SetResult(common.SearchResultIface)
	SetSearchPaths([]string)
}

func NewResultList(
	articleView ArticleViewIface,
	headerLabel HeaderLabelIface,
	app ApplicationIface,
) *ResultList {
	widget := qt.NewQListWidget(nil)
	resultList := &ResultList{
		QListWidget: widget,
		HeaderLabel: headerLabel,
		ArticleView: articleView,
		app:         app,
	}
	widget.OnCurrentRowChanged(func(row int) {
		if row < 0 {
			return
		}
		resultList.onActivate(row)
	})
	widget.OnItemActivated(func(item *qt.QListWidgetItem) {
		row := widget.Row(item)
		if row < 0 {
			return
		}
		resultList.onActivate(row)
	})
	return resultList
}

type ResultList struct {
	*qt.QListWidget

	results []common.SearchResultIface

	active common.SearchResultIface

	HeaderLabel HeaderLabelIface
	ArticleView ArticleViewIface

	app ApplicationIface
}

func (w *ResultList) QWidget() *qt.QWidget {
	return w.QListWidget.QWidget
}

func (w *ResultList) Active() common.SearchResultIface {
	return w.active
}

func (w *ResultList) SetResults(results []common.SearchResultIface) {
	w.QListWidget.Clear()
	w.results = results
	if len(results) == 0 {
		return
	}
	for _, res := range results {
		if res == nil {
			slog.Warn("ResultListWidget: SetResults: res == nil")
			continue
		}
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
		if w.app.HasFavorite(terms[0]) {
			text += " ★"
		}
		w.AddItem(text)
	}
	w.SetCurrentRow(0)
}

func (w *ResultList) Reload() {
	row := w.CurrentRow()
	w.SetResults(w.results)
	w.SetCurrentRow(row)
}

func (w *ResultList) SetCurrentResult(resultIndex int) {
	w.SetCurrentRow(resultIndex)
}

func (w *ResultList) onActivate(row int) {
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
	w.app.OnResultDisplay(res.Terms())
	w.active = res
}

func (w *ResultList) Clear() {
	w.QListWidget.Clear()
	w.results = nil
}

func (w *ResultList) GoNext() {
	row := w.QListWidget.CurrentRow() + 1
	if row >= len(w.results) {
		return
	}
	w.SetCurrentRow(row)
}

func (w *ResultList) GoPrevious() {
	row := w.QListWidget.CurrentRow() - 1
	if row < 0 {
		return
	}
	w.SetCurrentRow(row)
}
