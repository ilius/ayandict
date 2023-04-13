package application

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/ilius/ayandict/pkg/common"
	"github.com/ilius/ayandict/pkg/qerr"
	"github.com/ilius/ayandict/pkg/stardict"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

const (
	dictManager_up       = "Up"
	dictManager_down     = "Down"
	dictManager_openDirs = "Open Directories"
)

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
	app *Application,
	parent widgets.QWidget_ITF,
) *DictManager {
	infoList := stardict.GetInfoList()
	infoMap := makeDictInfoMap(infoList)

	window := widgets.NewQDialog(parent, core.Qt__Dialog)
	window.SetWindowTitle("Dictionaries")
	window.Resize2(800, 600)

	qs := getQSettings(window)
	restoreWinGeometry(app, qs, &window.QWidget, QS_dictManager)

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
		toolbar.AddAction2(icon, dictManager_up)
	}
	toolbar.AddSeparator()
	{
		icon := style.StandardIcon(widgets.QStyle__SP_ArrowDown, tbOpt, nil)
		toolbar.AddAction2(icon, dictManager_down)
	}
	toolbar.AddSeparator()
	{
		icon := style.StandardIcon(widgets.QStyle__SP_DirOpenIcon, tbOpt, nil)
		toolbar.AddAction2(icon, dictManager_openDirs)
	}
	newItem := func(text string) *widgets.QTableWidgetItem {
		item := widgets.NewQTableWidgetItem2(text, 0)
		item.SetFlags(core.Qt__ItemIsSelectable | core.Qt__ItemIsEnabled)
		return item
	}
	setItem := func(index int, dictName string, ds *common.DictSettings) {
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
			qerr.Error(err)
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
	openFolder := func() {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			qerr.Error(err)
			return
		}
		for _, p := range conf.DirectoryList {
			if !filepath.IsAbs(p) {
				p = filepath.Join(homeDir, p)
			}
			url := core.NewQUrl()
			url.SetScheme("file")
			url.SetPath(p, core.QUrl__TolerantMode)
			gui.QDesktopServices_OpenUrl(url)
		}
	}
	toolbar.ConnectActionTriggered(func(action *widgets.QAction) {
		switch action.Text() {
		case dictManager_up:
			toolbarUp()
		case dictManager_down:
			toolbarDown()
		case dictManager_openDirs:
			openFolder()
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

	table.SetRowCount(len(infoList))
	for index, info := range infoList {
		dictName := info.DictName()
		ds := dictSettingsMap[dictName]
		if ds == nil {
			log.Printf("dict manager: found new dict: %v\n", dictName)
			ds = common.NewDictSettings(info, index)
			dictSettingsMap[dictName] = ds
		}
		setItem(index, dictName, ds)
	}

	restoreTableColumnsWidth(
		qs,
		table,
		QS_dictManager,
	)
	table.HorizontalHeader().ConnectSectionResized(func(logicalIndex int, oldSize int, newSize int) {
		saveTableColumnsWidth(qs, table, QS_dictManager)
	})

	{
		ch := make(chan time.Time, 100)
		window.ConnectMoveEvent(func(event *gui.QMoveEvent) {
			ch <- time.Now()
		})
		window.ConnectResizeEvent(func(event *gui.QResizeEvent) {
			ch <- time.Now()
		})
		go actionSaveLoop(ch, func() {
			saveWinGeometry(qs, &window.QWidget, QS_dictManager)
		})
	}

	app.allTextWidgets = append(
		app.allTextWidgets,
		table,
		toolbar,
		okButton,
		cancelButton,
	)

	return &DictManager{
		Dialog:      window,
		TableWidget: table,
	}
}

// updates global var dictSettingsMap
// and returns dicts order
func (dm *DictManager) updateMap() map[string]int {
	table := dm.TableWidget
	order := map[string]int{}
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
		ds := dictSettingsMap[dictName]
		if ds == nil {
			ds = &common.DictSettings{}
			dictSettingsMap[dictName] = ds
		}
		ds.Symbol = symbol
		ds.Order = value
	}
	return order
}

// Run shows the dialog, if it Cancel was clicked it returns false
// if OK was clicked, then applies and saves changes
// and returs true
func (dm *DictManager) Run() bool {
	if dm.Dialog.Exec() != dialogAccepted {
		return false
	}
	dictsOrder = dm.updateMap()

	stardict.ApplyDictsOrder(dictsOrder)
	err := saveDictsSettings(dictSettingsMap)
	if err != nil {
		qerr.Error(err)
	}
	return true
}
