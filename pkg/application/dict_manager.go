package application

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"unicode/utf8"

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
	for bookName, ds := range settingsMap {
		order[bookName] = ds.Order
		if ds.Symbol == "" {
			symbol, _ := utf8.DecodeRune([]byte(bookName))
			ds.Symbol = fmt.Sprintf("[%s]", string(symbol))
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

func NewDictManager(
	app *widgets.QApplication,
	parent widgets.QWidget_ITF,
) *DictManager {
	infos := stardict.GetInfoList()

	window := widgets.NewQDialog(parent, core.Qt__Dialog)
	window.SetWindowTitle("Dictionaries")
	window.Resize2(400, 400)
	// window.SetSizePolicy2(expanding, expanding)

	table := widgets.NewQTableWidget(nil)
	table.SetColumnCount(3)
	header := table.HorizontalHeader()
	header.ResizeSection(0, 10)
	header.ResizeSection(1, 20)
	header.ResizeSection(2, 300)

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
		widgets.NewQTableWidgetItem2("Name", 0),
	)

	// table.SetSizePolicy2(expanding, expanding)
	// table.SetHorizontalScrollBarPolicy(core.Qt__ScrollBarAlwaysOff)
	// table.SetResizeMode(widgets.QListView__Adjust)
	// table.SetSizeAdjustPolicy(widgets.QAbstractScrollArea__AdjustToContents)

	// fmt.Println("Layout", window.Layout())
	// window.Layout().DestroyQLayout()

	// mainWidget := widgets.NewQWidget(nil, 0)
	// mainWidget.SetLayout(mainHBox)
	// window.SetLayout(mainHBox)
	mainHBox := widgets.NewQHBoxLayout2(nil)
	mainHBox.AddWidget(table, 10, core.Qt__AlignJustify)
	// window.Layout().AddChildWidget(mainWidget)
	// layout.SetStretch(0, 10)
	// layout.SetSizeConstraint(widgets.QLayout__SetMaximumSize)
	// layout.SetStretchFactor(table, 10)

	toolbar := widgets.NewQToolBar2(nil)
	mainHBox.AddWidget(toolbar, 0, 0)
	toolbar.SetOrientation(core.Qt__Vertical)

	style := app.Style()
	tbOpt := widgets.NewQStyleOptionToolBar()
	{
		icon := style.StandardIcon(widgets.QStyle__SP_ArrowUp, tbOpt, nil)
		toolbar.AddAction2(icon, "Up")
	}
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
		// fmt.Printf("dictName=%#v, ds=%#v\n", dictName, ds)
		checkItem := widgets.NewQTableWidgetItem(0)
		if ds.Order >= 0 {
			checkItem.SetCheckState(core.Qt__Checked)
		} else {
			checkItem.SetCheckState(core.Qt__Unchecked)
		}
		table.SetItem(index, 0, checkItem)

		symbolItem := newItem(ds.Symbol)
		symbolItem.SetFlags(core.Qt__ItemIsEnabled |
			core.Qt__ItemIsSelectable |
			core.Qt__ItemIsEditable)

		table.SetItem(index, 1, symbolItem)
		table.SetItem(index, 2, newItem(dictName))
	}

	// table.SelectedIndexes() panics/crashes
	// so do methods in table.SelectionModel()
	// you have to use table.CurrentIndex() or table.CurrentItem()
	toolbarUp := func() {
		// qIndex := table.CurrentIndex()
		// index := qIndex.Row()
		// if index < 1 {
		// 	return
		// }
		// item := table.TakeItem(index)
		// table.InsertRow(index-1)
		// setItem(index-1)
		// table.Item()
		// table.InsertItem(index-1, item)
		// table.SetCurrentRow(index - 1)
	}
	toolbarDown := func() {
		// qIndex := table.CurrentIndex()
		// index := qIndex.Row()
		// if index > table.Count()-2 {
		// 	return
		// }
		// item := table.TakeItem(index)
		// table.InsertItem(index+1, item)
		// table.SetCurrentRow(index + 1)
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
		setItem(index, info.BookName())
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
		dictName := table.Item(index, 2).Text()
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
		fmt.Println(err)
	}
}
