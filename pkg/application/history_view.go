package application

import (
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

type HistoryView struct {
	*widgets.QListWidget

	doQuery func(string)
}

func NewHistoryView() *HistoryView {
	widget := widgets.NewQListWidget(nil)
	return &HistoryView{
		QListWidget: widget,
	}
}

func (view *HistoryView) SetupCustomHandlers() {
	doQuery := view.doQuery
	if doQuery == nil {
		panic("doQuery is not set")
	}

	// view.SelectedItems() panics
	// and even after fixing panic, doesn't return anything
	// you have to use view.CurrentIndex()
	view.ConnectMousePressEvent(func(event *gui.QMouseEvent) {
		view.MousePressEventDefault(event)
		index := view.CurrentIndex()
		if index == nil {
			return
		}
		view.Activated(index)
	})

	view.ConnectItemActivated(func(item *widgets.QListWidgetItem) {
		doQuery(item.Text())
	})

}
