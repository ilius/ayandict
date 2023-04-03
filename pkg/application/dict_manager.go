package application

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"unicode/utf8"

	"github.com/ilius/ayandict/pkg/common"
	"github.com/ilius/ayandict/pkg/config"
	"github.com/ilius/ayandict/pkg/stardict"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

const dictsJsonFilename = "dicts.json"

type DictSettings struct {
	Symbol string `json:"symbol"`
	Order  int    `json:"order"`
}

var dictSettingsMap = map[string]*DictSettings{}

func defaultDictSymbol(dictName string) string {
	symbol, _ := utf8.DecodeRune([]byte(dictName))
	return fmt.Sprintf("[%s]", string(symbol))
}

func loadDictsSettings() (map[string]*DictSettings, map[string]int, error) {
	order := map[string]int{}
	settingsMap := map[string]*DictSettings{}
	fpath := filepath.Join(config.GetConfigDir(), dictsJsonFilename)
	jsonBytes, err := ioutil.ReadFile(fpath)
	if err != nil {
		if os.IsNotExist(err) {
			return settingsMap, order, nil
		}
		return settingsMap, order, err
	}
	err = json.Unmarshal(jsonBytes, &settingsMap)
	if err != nil {
		return settingsMap, order, err
	}
	for dictName, ds := range settingsMap {
		order[dictName] = ds.Order
		if ds.Symbol == "" {
			ds.Symbol = defaultDictSymbol(dictName)
		}
	}
	return settingsMap, order, nil
}

func saveDictsSettings(settingsMap map[string]*DictSettings) error {
	jsonBytes, err := json.MarshalIndent(settingsMap, "", "\t")
	if err != nil {
		return err
	}
	fpath := filepath.Join(config.GetConfigDir(), dictsJsonFilename)
	err = ioutil.WriteFile(fpath, jsonBytes, 0o644)
	if err != nil {
		return err
	}
	return nil
}

type DictManager struct {
	Dialog      *widgets.QDialog
	TableWidget *widgets.QTableWidget
}

func makeDictInfoMap(infos []common.Info) map[string]common.Info {
	infoMap := make(map[string]common.Info, len(infos))
	for _, info := range infos {
		infoMap[info.DictName()] = info
	}
	return infoMap
}

func NewDictManager(
	app *widgets.QApplication,
	parent widgets.QWidget_ITF,
) *DictManager {
	infos := stardict.GetInfoList()
	infoMap := makeDictInfoMap(infos)

	window := widgets.NewQDialog(parent, core.Qt__Dialog)
	window.SetWindowTitle("Dictionaries")
	window.Resize2(800, 600)

	const columns = 4

	table := widgets.NewQTableWidget(nil)
	table.SetColumnCount(columns)
	header := table.HorizontalHeader()
	header.ResizeSection(0, 10)
	header.ResizeSection(1, 20)
	header.ResizeSection(2, 80)
	header.ResizeSection(3, 500)

	table.SetHorizontalHeaderItem(
		0,
		widgets.NewQTableWidgetItem2("", 0),
	)
	table.SetHorizontalHeaderItem(
		1,
		widgets.NewQTableWidgetItem2("Sym", 0),
	)
	table.SetHorizontalHeaderItem(
		2,
		widgets.NewQTableWidgetItem2("Entries", 0),
	)
	table.SetHorizontalHeaderItem(
		3,
		widgets.NewQTableWidgetItem2("Name", 0),
	)

	mainHBox := widgets.NewQHBoxLayout2(nil)
	mainHBox.AddWidget(table, 0, 0)

	toolbar := widgets.NewQToolBar2(nil)
	toolbarVBox := widgets.NewQVBoxLayout2(nil)
	toolbarVBox.AddSpacing(80)
	toolbarVBox.AddWidget(toolbar, 0, 0)
	toolbarVBox.SetContentsMargins(0, 0, 0, 0)

	mainHBox.AddLayout(toolbarVBox, 0)
	toolbar.SetOrientation(core.Qt__Vertical)

	style := app.Style()
	tbOpt := widgets.NewQStyleOptionToolBar()
	toolbar.SetIconSize(core.NewQSize2(48, 48))
	{
		icon := style.StandardIcon(widgets.QStyle__SP_ArrowUp, tbOpt, nil)
		toolbar.AddAction2(icon, "Up")
	}
	toolbar.AddSeparator()
	{
		icon := style.StandardIcon(widgets.QStyle__SP_ArrowDown, tbOpt, nil)
		toolbar.AddAction2(icon, "Down")
	}
	newItem := func(text string) *widgets.QTableWidgetItem {
		item := widgets.NewQTableWidgetItem2(text, 0)
		item.SetFlags(core.Qt__ItemIsSelectable | core.Qt__ItemIsEnabled)
		return item
	}
	setItem := func(index int, dictName string) {
		ds := dictSettingsMap[dictName]
		if ds == nil {
			log.Printf("dictName=%#v, ds=%v\n", dictName, ds)
			ds = &DictSettings{
				Symbol: defaultDictSymbol(dictName),
				Order:  index,
			}
			dictSettingsMap[dictName] = ds
		}
		info, ok := infoMap[dictName]
		if !ok {
			log.Printf("dictName=%#v not in infoMap\n", dictName)
			return
		}
		checkItem := widgets.NewQTableWidgetItem(0)
		if ds.Order < 0 {
			checkItem.SetCheckState(core.Qt__Unchecked)
		} else {
			checkItem.SetCheckState(core.Qt__Checked)
		}
		table.SetItem(index, 0, checkItem)

		symbolItem := newItem(ds.Symbol)
		symbolItem.SetFlags(core.Qt__ItemIsEnabled |
			core.Qt__ItemIsSelectable |
			core.Qt__ItemIsEditable)

		table.SetItem(index, 1, symbolItem)
		entries, err := info.EntryCount()
		if err != nil {
			log.Println(err)
			return
		}
		table.SetItem(index, 2, newItem(strconv.FormatInt(int64(entries), 10)))
		table.SetItem(index, 3, newItem(dictName))
	}

	// table.SelectedIndexes() panics/crashes
	// so do methods in table.SelectionModel()
	// you have to use table.CurrentRow(), table.CurrentIndex()
	// or table.CurrentItem()
	toolbarUp := func() {
		row := table.CurrentRow()
		if row < 1 {
			return
		}
		for col := 0; col < columns; col++ {
			item1 := table.TakeItem(row, col)
			item2 := table.TakeItem(row-1, col)
			table.SetItem(row-1, col, item1)
			table.SetItem(row, col, item2)
		}
		table.SetCurrentCell(row-1, table.CurrentColumn())
	}
	toolbarDown := func() {
		row := table.CurrentRow()
		if row > table.RowCount()-2 {
			return
		}
		for col := 0; col < columns; col++ {
			item1 := table.TakeItem(row, col)
			item2 := table.TakeItem(row+1, col)
			table.SetItem(row+1, col, item1)
			table.SetItem(row, col, item2)
		}
		table.SetCurrentCell(row+1, table.CurrentColumn())
	}
	toolbar.ConnectActionTriggered(func(action *widgets.QAction) {
		switch action.Text() {
		case "Up":
			toolbarUp()
		case "Down":
			toolbarDown()
		}
	})

	buttonBox := widgets.NewQDialogButtonBox(nil)
	okButton := buttonBox.AddButton2("OK", widgets.QDialogButtonBox__AcceptRole)
	okButton.ConnectClicked(func(checked bool) {
		window.Accept()
	})
	cancelButton := buttonBox.AddButton2("Cancel", widgets.QDialogButtonBox__RejectRole)
	cancelButton.ConnectClicked(func(checked bool) {
		window.Reject()
	})

	mainBox := widgets.NewQVBoxLayout2(window)
	mainBox.AddLayout(mainHBox, 1)
	mainBox.AddWidget(buttonBox, 0, 0)

	table.SetRowCount(len(infos))
	for index, info := range infos {
		setItem(index, info.DictName())
	}

	return &DictManager{
		Dialog:      window,
		TableWidget: table,
	}
}

func dictsSettingsFromListWidget(table *widgets.QTableWidget) (map[string]*DictSettings, map[string]int) {
	order := map[string]int{}
	settingMap := map[string]*DictSettings{}
	count := table.RowCount()
	for index := 0; index < count; index++ {
		disable := table.Item(index, 0).CheckState() != core.Qt__Checked
		symbol := table.Item(index, 1).Text()
		dictName := table.Item(index, 3).Text()
		value := index + 1
		if disable {
			value = -value
		}
		order[dictName] = value
		settingMap[dictName] = &DictSettings{
			Symbol: symbol,
			Order:  value,
		}
	}
	return settingMap, order
}

func SaveDictManagerDialog(manager *DictManager) {
	dictSettingsMap, dictsOrder = dictsSettingsFromListWidget(manager.TableWidget)

	stardict.ApplyDictsOrder(dictsOrder)

	err := saveDictsSettings(dictSettingsMap)
	if err != nil {
		log.Println(err)
	}
}
