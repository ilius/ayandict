package frequency

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

func NewFrequencyView(maxSize int) *FrequencyTable {
	widget := widgets.NewQTableWidget(nil)
	widget.SetColumnCount(2)

	widget.ConnectItemClicked(func(item *widgets.QTableWidgetItem) {
		widget.ItemActivated(item)
	})

	return &FrequencyTable{
		QTableWidget: widget,
		maxSize:      maxSize,
		KeyMap:       map[string]int{},
		Counts:       map[string]int{},
	}
}

type FrequencyTable struct {
	*widgets.QTableWidget

	maxSize int

	Counts map[string]int
	Keys   []string
	KeyMap map[string]int
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
	jsonBytes, err := ioutil.ReadFile(pathStr)
	if err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("Error loading history: %v\n", err)
		}
		return nil
	}
	countMap := map[string]int{}
	err = json.Unmarshal(jsonBytes, &countMap)
	if err != nil {
		return fmt.Errorf("Bad history file %#v: %v\n", pathStr, err)
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
	// log.Printf("Triming %d items to %d\n", len(view.Counts), maxSize)
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

func (view *FrequencyTable) SaveToFile(pathStr string) error {
	jsonBytes, err := json.MarshalIndent(view.Counts, "", "\t")
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(pathStr, jsonBytes, 0o644)
	if err != nil {
		return err
	}
	return nil
}
