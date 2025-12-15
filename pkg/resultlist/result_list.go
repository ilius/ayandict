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

func NewResultListWidget(
	articleView ArticleViewIface,
	headerLabel HeaderLabelIface,
	app ApplicationIface,
) *ResultListWidget {
	widget := qt.NewQListWidget(nil)
	resultList := &ResultListWidget{
		QListWidget: widget,
		HeaderLabel: headerLabel,
		ArticleView: articleView,
		app:         app,
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

	HeaderLabel HeaderLabelIface
	ArticleView ArticleViewIface

	app ApplicationIface
}

func (w *ResultListWidget) SetResults(results []common.SearchResultIface) {
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
	w.app.OnResultDisplay(res.Terms())
	w.Active = res
}

func (w *ResultListWidget) Clear() {
	w.QListWidget.Clear()
	w.results = nil
}

func (w *ResultListWidget) GoNext() {
	row := w.QListWidget.CurrentRow() + 1
	if row >= len(w.results) {
		return
	}
	w.SetCurrentRow(row)
}

func (w *ResultListWidget) GoPrevious() {
	row := w.QListWidget.CurrentRow() - 1
	if row < 0 {
		return
	}
	w.SetCurrentRow(row)
}
