package dictmgr

import (
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/ilius/ayandict/v2/pkg/config"
	"github.com/ilius/ayandict/v2/pkg/iface"
	"github.com/ilius/ayandict/v2/pkg/qerr"
	"github.com/ilius/ayandict/v2/pkg/settings"
	common "github.com/ilius/go-dict-commons"
	"github.com/ilius/qt/core"
	"github.com/ilius/qt/gui"
	"github.com/ilius/qt/widgets"
)

const (
	QS_dictManager = "dict_manager"

	dm_col_enable   = 0
	dm_col_header   = 1
	dm_col_symbol   = 2
	dm_col_entries  = 3
	dm_col_dictName = 4

	dictManager_up       = "Up"
	dictManager_down     = "Down"
	dictManager_openInfo = "Open Info File"
	dictManager_openDirs = "Open Directories"
)

type DictManager struct {
	Dialog      *widgets.QDialog
	TableWidget *widgets.QTableWidget
	TextWidgets []iface.HasSetFont
}

func makeDictInfoMap(infos []common.Dictionary) map[string]common.Dictionary {
	infoMap := make(map[string]common.Dictionary, len(infos))
	for _, info := range infos {
		infoMap[info.DictName()] = info
	}
	return infoMap
}

func NewDictManager(
	app *widgets.QApplication,
	parent widgets.QWidget_ITF,
	conf *config.Config,
) *DictManager {
	infoMap := makeDictInfoMap(dicList)

	window := widgets.NewQDialog(parent, core.Qt__Dialog)
	window.SetWindowTitle("Dictionaries")
	window.Resize2(900, 800)

	qs := settings.GetQSettings(window)
	settings.RestoreWinGeometry(app, qs, &window.QWidget, QS_dictManager)

	const columns = 5

	table := widgets.NewQTableWidget(nil)
	table.SetColumnCount(columns)
	header := table.HorizontalHeader()
	header.ResizeSection(dm_col_enable, 10)
	header.ResizeSection(dm_col_header, 10)
	header.ResizeSection(dm_col_symbol, 20)
	header.ResizeSection(dm_col_entries, 80)
	header.ResizeSection(dm_col_dictName, 500)

	table.SetHorizontalHeaderItem(
		dm_col_enable,
		widgets.NewQTableWidgetItem2("", 0),
	)
	table.SetHorizontalHeaderItem(
		dm_col_header,
		widgets.NewQTableWidgetItem2("Terms\nHeader", 0),
	)
	table.SetHorizontalHeaderItem(
		dm_col_symbol,
		widgets.NewQTableWidgetItem2("Sym", 0),
	)
	table.SetHorizontalHeaderItem(
		dm_col_entries,
		widgets.NewQTableWidgetItem2("Entries", 0),
	)
	table.SetHorizontalHeaderItem(
		dm_col_dictName,
		widgets.NewQTableWidgetItem2("Name", 0),
	)

	extraOptionsWidget := widgets.NewQWidget(nil, 0)
	extraOptionsVBox := widgets.NewQVBoxLayout2(nil)
	extraOptionsWidget.SetLayout(extraOptionsVBox)
	extraOptionsWidget.Hide()

	flagsCBWidget := NewDictFlagsCheckboxes(func() {
		extraOptionsWidget.Hide()
	})
	extraOptionsVBox.AddWidget(flagsCBWidget, 0, 0)

	mainVBox := widgets.NewQVBoxLayout2(nil)
	mainVBox.AddWidget(table, 3, 0)
	mainVBox.AddWidget(extraOptionsWidget, 1, 0)

	mainHBox := widgets.NewQHBoxLayout2(nil)
	mainHBox.AddLayout(mainVBox, 1)

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
		icon := style.StandardIcon(widgets.QStyle__SP_FileIcon, tbOpt, nil)
		toolbar.AddAction2(icon, dictManager_openInfo)
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
	setItem := func(index int, dictName string, ds *DictSettings) {
		info, ok := infoMap[dictName]
		if !ok {
			log.Printf("dictName=%#v not in infoMap\n", dictName)
			return
		}
		enabledItem := widgets.NewQTableWidgetItem(0)
		if ds.Order < 0 {
			enabledItem.SetCheckState(core.Qt__Unchecked)
		} else {
			enabledItem.SetCheckState(core.Qt__Checked)
		}
		table.SetItem(index, dm_col_enable, enabledItem)

		headerItem := widgets.NewQTableWidgetItem(1)
		if ds.HideTermsHeader {
			headerItem.SetCheckState(core.Qt__Unchecked)
		} else {
			headerItem.SetCheckState(core.Qt__Checked)
		}
		table.SetItem(index, dm_col_header, headerItem)

		symbolItem := newItem(ds.Symbol)
		symbolItem.SetFlags(core.Qt__ItemIsEnabled |
			core.Qt__ItemIsSelectable |
			core.Qt__ItemIsEditable)
		table.SetItem(index, dm_col_symbol, symbolItem)

		entries, err := info.EntryCount()
		if err != nil {
			qerr.Error(err)
			return
		}
		table.SetItem(
			index, dm_col_entries,
			newItem(strconv.FormatInt(int64(entries), 10)),
		)
		table.SetItem(index, dm_col_dictName, newItem(dictName))
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
	openInfoFile := func() {
		row := table.CurrentRow()
		if row < 0 {
			return
		}
		dictName := table.Item(row, dm_col_dictName).Text()
		dic := dicMap[dictName]
		if dic == nil {
			qerr.Errorf("No dictionary %#v found", dictName)
			return
		}
		path := dic.InfoPath()
		if path == "" {
			return
		}
		url := core.NewQUrl()
		url.SetScheme("file")
		url.SetPath(path, core.QUrl__TolerantMode)
		gui.QDesktopServices_OpenUrl(url)
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
		case dictManager_openInfo:
			openInfoFile()
		case dictManager_openDirs:
			openFolder()
		}
	})

	table.ConnectCellClicked(func(row int, column int) {
		if column < 3 {
			return
		}
		dictName := table.Item(row, dm_col_dictName).Text()
		ds := dictSettingsMap[dictName]
		if ds == nil {
			return
		}
		flagsCBWidget.SetActiveDictSetting(ds)
		extraOptionsWidget.Show()
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

	table.SetRowCount(len(dicList))
	for index, dic := range dicList {
		dictName := dic.DictName()
		ds := dictSettingsMap[dictName]
		if ds == nil {
			log.Printf("dict manager: found new dict: %v\n", dictName)
			ds = NewDictSettings(dic, index)
			ds.Hash = Hash(dic)
			dictSettingsMap[dictName] = ds
		}
		setItem(index, dictName, ds)
	}

	settings.RestoreTableColumnsWidth(
		qs,
		table,
		QS_dictManager,
	)
	table.HorizontalHeader().ConnectSectionResized(func(logicalIndex int, oldSize int, newSize int) {
		settings.SaveTableColumnsWidth(qs, table, QS_dictManager)
	})

	settings.SetupWinGeometrySave(qs, &window.QWidget, QS_dictManager)

	return &DictManager{
		Dialog:      window,
		TableWidget: table,
		TextWidgets: []iface.HasSetFont{
			table,
			toolbar,
			okButton,
			cancelButton,
		},
	}
}

// updates global var dictSettingsMap
// and returns dicts order
func (dm *DictManager) updateMap() map[string]int {
	table := dm.TableWidget
	order := map[string]int{}
	count := table.RowCount()
	for index := 0; index < count; index++ {
		disable := table.Item(index, dm_col_enable).CheckState() != core.Qt__Checked
		hideHeader := table.Item(index, dm_col_header).CheckState() != core.Qt__Checked
		symbol := table.Item(index, dm_col_symbol).Text()
		dictName := table.Item(index, dm_col_dictName).Text()
		value := index + 1
		if disable {
			value = -value
		}
		order[dictName] = value
		ds := dictSettingsMap[dictName]
		if ds == nil {
			ds = &DictSettings{}
			dictSettingsMap[dictName] = ds
		}
		ds.HideTermsHeader = hideHeader
		ds.Symbol = symbol
		ds.Order = value
	}
	return order
}

// Run shows the dialog, if it Cancel was clicked it returns false
// if OK was clicked, then applies and saves changes
// and returs true
func (dm *DictManager) Run() bool {
	if dm.Dialog.Exec() != int(widgets.QDialog__Accepted) {
		return false
	}
	dictsOrder = dm.updateMap()

	Reorder(dictsOrder)

	for _, dic := range dicList {
		disabled := dic.Disabled()
		dic.SetDisabled(dictsOrder[dic.DictName()] < 0)
		if disabled && !dic.Disabled() {
			dic.Load()
		}
	}

	err := saveDictsSettings(dictSettingsMap)
	if err != nil {
		qerr.Error(err)
	}
	return true
}
