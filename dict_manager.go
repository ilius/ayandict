package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ilius/ayandict/pkg/config"
	"github.com/ilius/ayandict/pkg/stardict"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

const dictsJsonFilename = "dicts.json"

func loadDictsOrder() (map[string]int, error) {
	m := map[string]int{}
	fpath := filepath.Join(config.GetConfigDir(), dictsJsonFilename)
	jsonBytes, err := ioutil.ReadFile(fpath)
	if err != nil {
		if os.IsNotExist(err) {
			return m, nil
		}
		return m, err
	}
	err = json.Unmarshal(jsonBytes, &m)
	if err != nil {
		return m, err
	}
	return m, nil
}

func saveDictsOrder(order map[string]int) error {
	jsonBytes, err := json.MarshalIndent(order, "", "\t")
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
	Dialog     *widgets.QDialog
	ListWidget *widgets.QListWidget
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

	listWidget := widgets.NewQListWidget(nil)
	// listWidget.SetSizePolicy2(expanding, expanding)
	// listWidget.SetHorizontalScrollBarPolicy(core.Qt__ScrollBarAlwaysOff)
	// listWidget.SetResizeMode(widgets.QListView__Adjust)
	// listWidget.SetSizeAdjustPolicy(widgets.QAbstractScrollArea__AdjustToContents)

	// fmt.Println("Layout", window.Layout())
	// window.Layout().DestroyQLayout()

	// mainWidget := widgets.NewQWidget(nil, 0)
	// mainWidget.SetLayout(mainHBox)
	// window.SetLayout(mainHBox)
	mainHBox := widgets.NewQHBoxLayout2(nil)
	mainHBox.AddWidget(listWidget, 10, core.Qt__AlignJustify)
	// window.Layout().AddChildWidget(mainWidget)
	// layout.SetStretch(0, 10)
	// layout.SetSizeConstraint(widgets.QLayout__SetMaximumSize)
	// layout.SetStretchFactor(listWidget, 10)

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
	// listWidget.SelectedIndexes() panics/crashes
	// so do methods in listWidget.SelectionModel()
	// you have to use listWidget.CurrentIndex() or listWidget.CurrentItem()
	toolbarUp := func() {
		qIndex := listWidget.CurrentIndex()
		index := qIndex.Row()
		if index < 1 {
			return
		}
		item := listWidget.TakeItem(index)
		listWidget.InsertItem(index-1, item)
		listWidget.SetCurrentRow(index - 1)
	}
	toolbarDown := func() {
		qIndex := listWidget.CurrentIndex()
		index := qIndex.Row()
		if index > listWidget.Count()-2 {
			return
		}
		item := listWidget.TakeItem(index)
		listWidget.InsertItem(index+1, item)
		listWidget.SetCurrentRow(index + 1)
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

	for _, info := range infos {
		item := widgets.NewQListWidgetItem2(info.BookName(), listWidget, 0)
		if info.Disabled {
			item.SetCheckState(core.Qt__Unchecked)
		} else {
			item.SetCheckState(core.Qt__Checked)
		}
	}

	return &DictManager{
		Dialog:     window,
		ListWidget: listWidget,
	}
}

func dictsOrderFromListWidget(list *widgets.QListWidget) map[string]int {
	order := map[string]int{}
	count := list.Count()
	for index := 0; index < count; index++ {
		item := list.Item(index)
		bookName := item.Text()
		value := index + 1
		if item.CheckState() != core.Qt__Checked {
			value = -value
		}
		order[bookName] = value
	}
	return order
}

func SaveDictManagerDialog(manager *DictManager) {
	dictsOrder = dictsOrderFromListWidget(manager.ListWidget)

	stardict.ApplyDictsOrder(dictsOrder)

	err := saveDictsOrder(dictsOrder)
	if err != nil {
		fmt.Println(err)
	}
}
