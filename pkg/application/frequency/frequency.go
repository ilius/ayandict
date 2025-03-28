package frequency

import (
	"log/slog"
	"strconv"

	"github.com/ilius/ayandict/v2/pkg/activity"
	"github.com/ilius/qt/core"
	"github.com/ilius/qt/widgets"
)

func NewFrequencyView(
	storage *activity.ActivityStorage,
	maxSize int,
) *FrequencyTable {
	widget := widgets.NewQTableWidget(nil)
	widget.SetColumnCount(2)

	widget.ConnectItemClicked(func(item *widgets.QTableWidgetItem) {
		widget.ItemActivated(item)
	})

	return &FrequencyTable{
		QTableWidget: widget,
		storage:      storage,
		maxSize:      maxSize,
		KeyMap:       map[string]int{},
		Counts:       map[string]int{},
	}
}

type FrequencyTable struct {
	*widgets.QTableWidget

	storage *activity.ActivityStorage

	maxSize int

	Counts map[string]int
	Keys   []string
	KeyMap map[string]int
}

func (view *FrequencyTable) Clear() {
	view.storage.ClearFrequency()
	view.Counts = map[string]int{}
	view.Keys = []string{}
	view.KeyMap = map[string]int{}
	view.QTableWidget.SetRowCount(0)
}

func (view *FrequencyTable) newItem(text string) *widgets.QTableWidgetItem {
	item := widgets.NewQTableWidgetItem2(text, 0)
	item.SetFlags(core.Qt__ItemIsSelectable | core.Qt__ItemIsEnabled)
	return item
}

func (view *FrequencyTable) addNew(key string, count int) {
	index := len(view.Keys)
	view.KeyMap[key] = index
	view.Keys = append(view.Keys, key)
	view.Counts[key] = count

	for index > 0 && count >= view.Counts[view.Keys[index-1]] {
		index = view.moveUp(key)
	}

	view.InsertRow(index)
	view.setItemForKey(index, key)
	view.Trim()
}

func (view *FrequencyTable) moveUp(key string) int {
	index := view.KeyMap[key]
	prevKey := view.Keys[index-1]

	view.Keys[index-1] = key
	view.Keys[index] = prevKey
	view.KeyMap[key] = index - 1
	view.KeyMap[prevKey] = index

	return index - 1
}

func (view *FrequencyTable) setItemForKey(index int, key string) {
	view.SetItem(index, 0, view.newItem(key))
	view.SetItem(index, 1, view.newItem(
		strconv.FormatInt(int64(view.Counts[key]), 10),
	))
}

func (view *FrequencyTable) Add(key string, plus int) {
	view.storage.AddFrequency(key, plus)
	index, ok := view.KeyMap[key]
	if !ok {
		view.addNew(key, plus)
		return
	}
	count := view.Counts[key] + plus
	view.Counts[key] = count

	view.RemoveRow(index)

	for index > 0 && count >= view.Counts[view.Keys[index-1]] {
		index = view.moveUp(key)
	}

	view.InsertRow(index)
	view.setItemForKey(index, key)
}

func (view *FrequencyTable) Load() error {
	countList, err := view.storage.LoadFrequency()
	if err != nil {
		return err
	}
	for _, item := range countList {
		view.addNew(item.Word, item.Count)
	}
	return nil
}

func (view *FrequencyTable) Trim() {
	if len(view.Counts) <= view.maxSize {
		return
	}
	maxSize := view.maxSize
	// to avoid trimming on every new item
	if maxSize > 20 {
		maxSize -= 10
	} else {
		maxSize -= maxSize / 3
	}
	newKeys := view.Keys[:maxSize]
	newKeyMap := map[string]int{}
	newCounts := map[string]int{}
	for _, key := range newKeys {
		newKeyMap[key] = view.KeyMap[key]
		newCounts[key] = view.Counts[key]
	}
	view.Keys = newKeys
	view.KeyMap = newKeyMap
	view.Counts = newCounts
	view.SetRowCount(maxSize)
}

func (view *FrequencyTable) Save() error {
	return view.storage.SaveFrequency()
}

func (view *FrequencyTable) SaveNoError() {
	err := view.Save()
	if err != nil {
		slog.Error("Error saving frequency: " + err.Error())
	}
}
