package application

import (
	"log/slog"

	"github.com/ilius/ayandict/v2/pkg/activity"
	"github.com/ilius/qt/gui"
	"github.com/ilius/qt/widgets"
)

type HistoryView struct {
	storage *activity.ActivityStorage
	maxSize int

	*widgets.QListWidget

	doQuery func(string)
}

func NewHistoryView(
	storage *activity.ActivityStorage,
	maxSize int,
) *HistoryView {
	widget := widgets.NewQListWidget(nil)
	return &HistoryView{
		storage:     storage,
		maxSize:     maxSize,
		QListWidget: widget,
	}
}

func (h *HistoryView) Load() error {
	history, err := h.storage.LoadHistory()
	if err != nil {
		return err
	}
	h.AddHistoryList(history)
	return nil
}

func (h *HistoryView) Save() {
	err := h.storage.SaveHistory()
	if err != nil {
		slog.Error("error saving history: " + err.Error())
	}
}

func (h *HistoryView) Add(query string) {
	slog.Debug("HistoryView: Add", "query", query)
	if !h.storage.AddHistory(query) {
		return
	}

	h.InsertItem2(0, query)
	h.TrimHistory(h.maxSize)
	if conf.HistoryAutoSave {
		h.Save()
	}
}

func (h *HistoryView) AddHistoryList(list []string) {
	for _, query := range list {
		h.InsertItem2(0, query)
	}
}

func (h *HistoryView) TrimHistory(maxSize int) {
	count := h.Count()
	if count <= maxSize {
		return
	}
	for i := maxSize; i < count; i++ {
		h.TakeItem(maxSize)
	}
}

func (h *HistoryView) ClearHistory() {
	h.storage.ClearHistory()
	h.Clear()
	h.Save()
}

func (h *HistoryView) SetupCustomHandlers() {
	doQuery := h.doQuery
	if doQuery == nil {
		panic("doQuery is not set")
	}

	// view.SelectedItems() panics
	// and even after fixing panic, doesn't return anything
	// you have to use view.CurrentIndex()
	h.ConnectMousePressEvent(func(event *gui.QMouseEvent) {
		h.MousePressEventDefault(event)
		index := h.CurrentIndex()
		if index == nil {
			return
		}
		if index.Row() < 0 {
			return
		}
		h.Activated(index)
	})

	h.ConnectItemActivated(func(item *widgets.QListWidgetItem) {
		doQuery(item.Text())
	})

	// we are doing query on MousePressEvent (before release, with any button)
	// so we don't need ConnectItemClicked
	// unless we decide to have right-click do something else
	// view.ConnectItemClicked(func(item *widgets.QListWidgetItem) {
	// 	doQuery(item.Text())
	// })
}
