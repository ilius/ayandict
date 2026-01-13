package resulttree

import (
	"fmt"
	"log/slog"
	"slices"

	common "codeberg.org/ilius/go-dict-commons"
	"github.com/ilius/ayandict/v3/pkg/dictmgr"
	qt "github.com/mappu/miqt/qt6"
)

type ApplicationIface interface {
	OnResultDisplay(terms []string)
	HasFavorite(term string) bool
	ResultsLabel() *qt.QLabel
}

type HeaderLabelIface interface {
	SetResult(common.SearchResultIface)
}

type ArticleViewIface interface {
	SetResult(common.SearchResultIface)
	SetSearchPaths([]string)
}

func NewResultTree(
	articleView ArticleViewIface,
	headerLabel HeaderLabelIface,
	app ApplicationIface,
) *ResultTree {
	widget := qt.NewQTreeWidget(nil)
	resultTree := &ResultTree{
		QTreeWidget: widget,
		HeaderLabel: headerLabel,
		ArticleView: articleView,
		app:         app,
	}
	widget.SetHeaderHidden(true)
	widget.OnCurrentItemChanged(func(current, previous *qt.QTreeWidgetItem) {
		resultTree.ItemActivated(current, 0)
	})
	widget.OnItemActivated(func(item *qt.QTreeWidgetItem, column int) {
		if item == nil {
			return
		}
		termIndex := 0
		for item.Parent() != nil {
			termIndex = item.Parent().IndexOfChild(item) + 1
			item = item.Parent()
		}
		resultIndex := widget.IndexOfTopLevelItem(item)
		if resultIndex < 0 {
			return
		}
		resultTree.onActivate(resultIndex, termIndex)
	})
	return resultTree
}

type ResultTree struct {
	*qt.QTreeWidget

	results []common.SearchResultIface

	currentResult common.SearchResultIface

	HeaderLabel HeaderLabelIface
	ArticleView ArticleViewIface

	app ApplicationIface
}

func (w *ResultTree) QWidget() *qt.QWidget {
	return w.QTreeWidget.QWidget
}

func (w *ResultTree) CurrentResult() common.SearchResultIface {
	return w.currentResult
}

func (w *ResultTree) SetResults(results []common.SearchResultIface) {
	w.QTreeWidget.Clear()
	w.results = results
	if len(results) == 0 {
		return
	}
	for _, res := range results {
		if res == nil {
			slog.Warn("ResultTreeWidget: SetResults: res == nil")
			continue
		}
		terms := res.Terms()
		if len(terms) == 0 {
			slog.Error("empty terms", "res", res)
			continue
		}
		text := terms[0]
		if len(terms) > 1 {
			text = fmt.Sprintf("%s (+%d)", text, len(terms)-1)
		}
		symbol := dictmgr.DictSymbol(res.DictName())
		if symbol != "" {
			text = fmt.Sprintf("%s %s", text, symbol)
		}
		if w.app.HasFavorite(terms[0]) {
			text += " ★"
		}
		toplevelItem := qt.NewQTreeWidgetItem2([]string{text})
		w.AddTopLevelItem(toplevelItem)
		for _, term := range terms[1:] {
			if w.app.HasFavorite(term) {
				term += " ★"
			}
			toplevelItem.AddChild(qt.NewQTreeWidgetItem2([]string{term}))
		}
	}
	w.SetCurrentItem(w.TopLevelItem(0))
}

func (w *ResultTree) Reload() {
	count := w.TopLevelItemCount()
	expanded := map[int]struct{}{}
	for index := range count {
		if w.TopLevelItem(index).IsExpanded() {
			expanded[index] = struct{}{}
		}
	}
	current := w.indexListFromItem(w.CurrentItem())
	w.SetResults(w.results)
	w.SetCurrentItem(w.indexListToItem(current))
	for index := range expanded {
		w.TopLevelItem(index).SetExpanded(true)
	}
}

func (w *ResultTree) SetCurrentResult(resultIndex int) {
	w.SetCurrentItem(w.TopLevelItem(resultIndex))
}

func (w *ResultTree) onActivate(resultIndex int, termIndex int) {
	if resultIndex >= len(w.results) {
		slog.Error("ResultTreeWidget: OnActivate: row index out of range", "row", resultIndex)
		return
	}
	res := w.results[resultIndex]
	if res == w.currentResult {
		slog.Debug("ResultTree.onActivate: ignoring same result")
		return
	}
	w.HeaderLabel.SetResult(res)
	w.ArticleView.SetResult(res)
	resDir := res.ResourceDir()
	if resDir == "" {
		w.ArticleView.SetSearchPaths([]string{})
	} else {
		w.ArticleView.SetSearchPaths([]string{resDir})
	}
	w.app.OnResultDisplay(res.Terms())
	w.currentResult = res
	label := w.app.ResultsLabel()
	if label != nil {
		if termIndex > 0 {
			label.SetText(fmt.Sprintf(
				"Results: %2d / %d (term %d)",
				resultIndex+1, len(w.results), termIndex+1,
			))
		} else {
			label.SetText(fmt.Sprintf(
				"Results: %2d / %d",
				resultIndex+1, len(w.results),
			))
		}
	}
}

func (w *ResultTree) Clear() {
	w.QTreeWidget.Clear()
	w.results = nil
}

func (w *ResultTree) indexListFromItem(item *qt.QTreeWidgetItem) []int {
	if item == nil {
		return nil
	}
	indexList := []int{}
	for item.Parent() != nil {
		parent := item.Parent()
		indexList = append(indexList, parent.IndexOfChild(item))
		item = parent
	}
	indexList = append(indexList, w.IndexOfTopLevelItem(item))
	slices.Reverse(indexList)
	return indexList
}

func (w *ResultTree) indexListToItem(indexList []int) *qt.QTreeWidgetItem {
	if indexList == nil {
		return nil
	}
	item := w.TopLevelItem(indexList[0])
	indexList = indexList[1:]
	for index := range indexList {
		item = item.Child(index)
	}
	return item
}

func (w *ResultTree) GoNext() {
	top := w.CurrentItem()
	if top == nil {
		return
	}
	for top.Parent() != nil {
		top = top.Parent()
	}
	topRow := w.IndexOfTopLevelItem(top)
	if topRow < 0 {
		return
	}
	nextTopRow := topRow + 1
	if nextTopRow >= w.TopLevelItemCount() {
		return
	}
	item := w.TopLevelItem(nextTopRow)
	w.SetCurrentItem(item)
	w.ScrollToItem(item)
}

func (w *ResultTree) GoPrevious() {
	top := w.CurrentItem()
	if top == nil {
		return
	}
	for top.Parent() != nil {
		top = top.Parent()
	}
	topRow := w.IndexOfTopLevelItem(top)
	if topRow < 1 {
		return
	}
	item := w.TopLevelItem(topRow - 1)
	w.SetCurrentItem(item)
	w.ScrollToItem(item)
}
