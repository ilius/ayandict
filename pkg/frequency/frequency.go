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

func NewFrequencyView(maxSize int) *FrequencyView {
	widget := widgets.NewQTableWidget(nil)
	widget.SetColumnCount(2)
	return &FrequencyView{
		QTableWidget: widget,
		maxSize:      maxSize,
		KeyMap:       map[string]int{},
		Counts:       map[string]int{},
	}
}

type FrequencyView struct {
	*widgets.QTableWidget

	maxSize int

	Counts map[string]int
	Keys   []string
	KeyMap map[string]int
}

func (view *FrequencyView) Clear() {
	view.Counts = map[string]int{}
	view.Keys = []string{}
	view.KeyMap = map[string]int{}
	view.QTableWidget.SetRowCount(0)
}

func (view *FrequencyView) newItem(text string) *widgets.QTableWidgetItem {
	item := widgets.NewQTableWidgetItem2(text, 0)
	item.SetFlags(core.Qt__ItemIsSelectable | core.Qt__ItemIsEnabled)
	return item
}

func (view *FrequencyView) addNew(key string, count int) {
	index := len(view.Keys)
	view.KeyMap[key] = index
	view.Keys = append(view.Keys, key)
	view.Counts[key] = count
	view.InsertRow(index)
	qKeyItem := view.newItem(key)
	view.SetItem(index, 0, qKeyItem)
	qCountItem := view.newItem(
		strconv.FormatInt(int64(count), 10),
	)
	view.SetItem(index, 1, qCountItem)
	view.Trim()
}

func (view *FrequencyView) moveUp(key string) int {
	index := view.KeyMap[key]

	prevKey := view.Keys[index-1]

	view.Keys[index-1] = key
	view.Keys[index] = prevKey
	view.KeyMap[key] = index - 1
	view.KeyMap[prevKey] = index

	view.SetItem(index-1, 0, view.newItem(key))
	view.SetItem(index-1, 1, view.newItem(
		strconv.FormatInt(int64(view.Counts[key]), 10),
	))
	view.SetItem(index, 0, view.newItem(prevKey))
	view.SetItem(index, 1, view.newItem(
		strconv.FormatInt(int64(view.Counts[prevKey]), 10),
	))

	return index - 1
}

func (view *FrequencyView) Add(key string, plus int) {
	index, ok := view.KeyMap[key]
	if !ok {
		view.addNew(key, plus)
		return
	}
	count := view.Counts[key] + plus
	view.Counts[key] = count

	if index < 1 || count <= view.Counts[view.Keys[index-1]] {
		view.Item(index, 1).SetText(strconv.FormatInt(int64(count), 10))
		return
	}
	index = view.moveUp(key)
	for index > 0 && count > view.Counts[view.Keys[index-1]] {
		index = view.moveUp(key)
	}
}

func (view *FrequencyView) LoadFromFile(pathStr string) error {
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

func (view *FrequencyView) Trim() {
	if len(view.Counts) <= view.maxSize {
		return
	}
	maxSize := view.maxSize
	// to avoid trimming on every new item
	if maxSize > 20 {
		maxSize -= 10
	} else {
		maxSize = maxSize * 2 / 3
	}
	// fmt.Printf("Triming %d items to %d\n", len(view.Counts), maxSize)
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

func (view *FrequencyView) SaveToFile(pathStr string) error {
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
