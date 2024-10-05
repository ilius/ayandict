package frequency

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
	"sync"

	"github.com/ilius/ayandict/v2/pkg/qtcommon/qerr"
	"github.com/ilius/qt/core"
	"github.com/ilius/qt/widgets"
)

func NewFrequencyView(
	filePath string,
	maxSize int,
) *FrequencyTable {
	widget := widgets.NewQTableWidget(nil)
	widget.SetColumnCount(2)

	widget.ConnectItemClicked(func(item *widgets.QTableWidgetItem) {
		widget.ItemActivated(item)
	})

	return &FrequencyTable{
		QTableWidget: widget,
		filePath:     filePath,
		maxSize:      maxSize,
		KeyMap:       map[string]int{},
		Counts:       map[string]int{},
	}
}

type FrequencyTable struct {
	*widgets.QTableWidget

	filePath string

	maxSize int

	Counts map[string]int
	Keys   []string
	KeyMap map[string]int

	saveMutex sync.Mutex
}

func (view *FrequencyTable) Clear() {
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

func (view *FrequencyTable) LoadFromFile(pathStr string) error {
	jsonBytes, err := os.ReadFile(pathStr)
	if err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("error loading frequency: %w", err)
		}
		return nil
	}
	countMap := map[string]int{}
	err = json.Unmarshal(jsonBytes, &countMap)
	if err != nil {
		return fmt.Errorf("bad frequency file %#v: %w", pathStr, err)
	}
	countList := [][2]any{}
	for key, count := range countMap {
		countList = append(countList, [2]any{key, count})
	}
	sort.Slice(countList, func(i, j int) bool {
		return countList[i][1].(int) > countList[j][1].(int)
	})
	for _, item := range countList {
		view.addNew(item[0].(string), item[1].(int))
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
	if view.filePath == "" {
		return fmt.Errorf("FrequencyTable: filePath is empty")
	}
	view.saveMutex.Lock()
	defer view.saveMutex.Unlock()
	jsonBytes, err := json.MarshalIndent(view.Counts, "", "\t")
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(view.filePath, jsonBytes, 0o644)
	if err != nil {
		return err
	}
	return nil
}

func (view *FrequencyTable) SaveNoError() {
	err := view.Save()
	if err != nil {
		qerr.Errorf("Error saving frequency: %v", err)
	}
}
