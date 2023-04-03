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

func (view *HistoryView) AddHistory(query string) {
	if len(history) > 0 && query == history[len(history)-1] {
		return
	}
	addHistoryLow(query)
	view.InsertItem2(0, query)
	view.TrimHistory(historyMaxSize)
	if conf.HistoryAutoSave {
		SaveHistory()
	}
}

func (view *HistoryView) AddHistoryList(list []string) {
	for _, query := range list {
		view.InsertItem2(0, query)
	}
}

func (view *HistoryView) TrimHistory(maxSize int) {
	count := view.Count()
	if count <= maxSize {
		return
	}
	for i := maxSize; i < count; i++ {
		view.TakeItem(maxSize)
	}
}

func (view *HistoryView) ClearHistory() {
	historyMutex.Lock()
	history = []string{}
	historyMutex.Unlock()

	SaveHistory()
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

	// we are doing query on MousePressEvent (before release, with any button)
	// so we don't need ConnectItemClicked
	// unless we decide to have right-click do something else
	// view.ConnectItemClicked(func(item *widgets.QListWidgetItem) {
	// 	doQuery(item.Text())
	// })
}
