package qdictmgr

import (
	"log/slog"
	"os"
	"path/filepath"
	"strconv"

	"github.com/ilius/ayandict/v3/pkg/config"
	"github.com/ilius/ayandict/v3/pkg/dictmgr/internal/dicts"
	"github.com/ilius/ayandict/v3/pkg/qtcommon"
	"github.com/ilius/ayandict/v3/pkg/qtcommon/qsettings"
	common "github.com/ilius/go-dict-commons"
	qt "github.com/mappu/miqt/qt6"
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

	columns = 5
)

type DictManager struct {
	Dialog      *qt.QDialog
	TableWidget *qt.QTableWidget
	VolumeInput *qt.QSpinBox
	TextWidgets []qtcommon.HasSetFont

	infoMap map[string]common.Dictionary

	app *qt.QApplication

	toolbar   *qt.QToolBar
	buttonBox *qt.QDialogButtonBox

	settings *qt.QSettings
}

func makeDictInfoMap(infos []common.Dictionary) map[string]common.Dictionary {
	infoMap := make(map[string]common.Dictionary, len(infos))
	for _, info := range infos {
		infoMap[info.DictName()] = info
	}
	return infoMap
}

// TODO: break down
func NewDictManager(
	app *qt.QApplication,
	parent *qt.QWidget,
	conf *config.Config,
) *DictManager {
	window := qt.NewQDialog(parent)
	window.SetWindowTitle("Dictionaries")
	window.Resize(900, 800)

	infoMap := makeDictInfoMap(dicts.DictList)

	qs := qsettings.GetQSettings(window.QObject)
	qsettings.RestoreWinGeometry(qs, window.QWidget, QS_dictManager)

	table := qt.NewQTableWidget2()
	volumeInput := qt.NewQSpinBox2()
	toolbar := qt.NewQToolBar3()

	buttonBox := qt.NewQDialogButtonBox2()
	okButton := buttonBox.AddButton2("OK", qt.QDialogButtonBox__AcceptRole)
	cancelButton := buttonBox.AddButton2("Cancel", qt.QDialogButtonBox__RejectRole)

	okButton.OnClicked(func() {
		window.Accept()
	})
	cancelButton.OnClicked(func() {
		window.Reject()
	})

	dictMgr := &DictManager{
		Dialog:      window,
		TableWidget: table,
		VolumeInput: volumeInput,
		TextWidgets: []qtcommon.HasSetFont{
			table,
			toolbar,
			okButton,
			cancelButton,
		},
		infoMap:   infoMap,
		app:       app,
		toolbar:   toolbar,
		buttonBox: buttonBox,
		settings:  qs,
	}
	dictMgr.prepareWidgets(conf)
	return dictMgr
}

func (dm *DictManager) newItem(text string) *qt.QTableWidgetItem {
	item := qt.NewQTableWidgetItem2(text)
	item.SetFlags(qt.ItemIsSelectable | qt.ItemIsEnabled)
	return item
}

func (dm *DictManager) setItem(
	index int,
	dictName string,
	ds *dicts.DictionarySettings,
) {
	table := dm.TableWidget
	info, ok := dm.infoMap[dictName]
	if !ok {
		slog.Error("dictName not in infoMap", "dictName", dictName)
		return
	}
	enabledItem := qt.NewQTableWidgetItem()
	if ds.Order < 0 {
		enabledItem.SetCheckState(qt.Unchecked)
	} else {
		enabledItem.SetCheckState(qt.Checked)
	}
	table.SetItem(index, dm_col_enable, enabledItem)

	headerItem := qt.NewQTableWidgetItem()
	if ds.HideTermsHeader {
		headerItem.SetCheckState(qt.Unchecked)
	} else {
		headerItem.SetCheckState(qt.Checked)
	}
	table.SetItem(index, dm_col_header, headerItem)

	symbolItem := dm.newItem(ds.Symbol)
	symbolItem.SetFlags(qt.ItemIsEnabled |
		qt.ItemIsSelectable |
		qt.ItemIsEditable)
	table.SetItem(index, dm_col_symbol, symbolItem)

	entries, err := info.EntryCount()
	if err != nil {
		slog.Error("error from info.EntryCount: " + err.Error())
		return
	}
	table.SetItem(
		index, dm_col_entries,
		dm.newItem(strconv.FormatInt(int64(entries), 10)),
	)
	table.SetItem(index, dm_col_dictName, dm.newItem(dictName))
}

// table.SelectedIndexes() panics/crashes
// so do methods in table.SelectionModel()
// you have to use table.CurrentRow(), table.CurrentIndex()
// or table.CurrentItem()
func (dm *DictManager) toolbarUp() {
	table := dm.TableWidget
	row := table.CurrentRow()
	if row < 1 {
		return
	}
	for col := range columns {
		item1 := table.TakeItem(row, col)
		item2 := table.TakeItem(row-1, col)
		table.SetItem(row-1, col, item1)
		table.SetItem(row, col, item2)
	}
	table.SetCurrentCell(row-1, table.CurrentColumn())
}

func (dm *DictManager) toolbarDown() {
	table := dm.TableWidget
	row := table.CurrentRow()
	if row > table.RowCount()-2 {
		return
	}
	for col := range columns {
		item1 := table.TakeItem(row, col)
		item2 := table.TakeItem(row+1, col)
		table.SetItem(row+1, col, item1)
		table.SetItem(row, col, item2)
	}
	table.SetCurrentCell(row+1, table.CurrentColumn())
}

func (dm *DictManager) openInfoFile() {
	table := dm.TableWidget
	row := table.CurrentRow()
	if row < 0 {
		return
	}
	dictName := table.Item(row, dm_col_dictName).Text()
	dic := dicts.DictByName[dictName]
	if dic == nil {
		slog.Error("no dictionary was found with this name: " + dictName)
		return
	}
	path := dic.InfoPath()
	if path == "" {
		return
	}
	url := qt.NewQUrl()
	url.SetScheme("file")
	url.SetPath2(path, qt.QUrl__TolerantMode)
	_ = qt.QDesktopServices_OpenUrl(url)
}

func (dm *DictManager) openFolder(conf *config.Config) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		slog.Error("error in os.UserHomeDir: " + err.Error())
		return
	}
	for _, p := range conf.DirectoryList {
		if !filepath.IsAbs(p) {
			p = filepath.Join(homeDir, p)
		}
		url := qt.NewQUrl()
		url.SetScheme("file")
		url.SetPath2(p, qt.QUrl__TolerantMode)
		_ = qt.QDesktopServices_OpenUrl(url)
	}
}

func (dm *DictManager) prepareWidgets(conf *config.Config) {
	var selectedDictSettings *dicts.DictionarySettings

	table := dm.TableWidget
	volumeInput := dm.VolumeInput

	table.SetColumnCount(columns)
	header := table.HorizontalHeader()
	header.ResizeSection(dm_col_enable, 10)
	header.ResizeSection(dm_col_header, 10)
	header.ResizeSection(dm_col_symbol, 20)
	header.ResizeSection(dm_col_entries, 80)
	header.ResizeSection(dm_col_dictName, 500)

	table.SetHorizontalHeaderItem(
		dm_col_enable,
		qt.NewQTableWidgetItem2(""),
	)
	table.SetHorizontalHeaderItem(
		dm_col_header,
		qt.NewQTableWidgetItem2("Terms\nHeader"),
	)
	table.SetHorizontalHeaderItem(
		dm_col_symbol,
		qt.NewQTableWidgetItem2("Sym"),
	)
	table.SetHorizontalHeaderItem(
		dm_col_entries,
		qt.NewQTableWidgetItem2("Entries"),
	)
	table.SetHorizontalHeaderItem(
		dm_col_dictName,
		qt.NewQTableWidgetItem2("Name"),
	)

	extraOptionsWidget := qt.NewQWidget2()
	extraOptionsVBox := qt.NewQVBoxLayout2()
	extraOptionsWidget.SetLayout(extraOptionsVBox.Layout())
	extraOptionsWidget.Hide()

	flagsCBWidget := NewDictFlagsCheckboxes(func() {
		extraOptionsWidget.Hide()
	})
	extraOptionsVBox.AddWidget(flagsCBWidget.QWidget)

	volumeHBox := qt.NewQHBoxLayout2()
	volumeHBox.AddWidget(qt.NewQLabel3("Audio Volume:").QWidget)
	volumeInput.SetMinimum(0)
	volumeInput.SetMaximum(999)
	volumeHBox.AddWidget(volumeInput.QWidget)
	volumeHBox.AddWidget(qt.NewQLabel3("").QWidget)
	extraOptionsVBox.AddLayout(volumeHBox.Layout())
	volumeInput.OnValueChanged(func(value int) {
		if selectedDictSettings == nil {
			slog.Info("ConnectValueChanged: selectedDictSettings == nil")
			return
		}
		selectedDictSettings.AudioVolume = value
	})

	mainVBox := qt.NewQVBoxLayout2()
	mainVBox.AddWidget3(table.QWidget, 3, 0)
	mainVBox.AddWidget3(extraOptionsWidget, 1, 0)

	mainHBox := qt.NewQHBoxLayout2()
	mainHBox.AddLayout2(mainVBox.Layout(), 1)

	toolbar := dm.toolbar
	toolbarVBox := qt.NewQVBoxLayout2()
	toolbarVBox.AddSpacing(80)
	toolbarVBox.AddWidget(toolbar.QWidget)
	toolbarVBox.SetContentsMargins(0, 0, 0, 0)

	mainHBox.AddLayout(toolbarVBox.Layout())
	toolbar.SetOrientation(qt.Vertical)

	style := qt.QApplication_Style()
	tbOpt := qt.NewQStyleOptionToolBar()
	toolbar.SetIconSize(qt.NewQSize2(48, 48))
	{
		icon := style.StandardIcon(qt.QStyle__SP_ArrowUp, tbOpt.QStyleOption, nil)
		_ = toolbar.AddAction2(icon, dictManager_up)
	}
	_ = toolbar.AddSeparator()
	{
		icon := style.StandardIcon(qt.QStyle__SP_ArrowDown, tbOpt.QStyleOption, nil)
		_ = toolbar.AddAction2(icon, dictManager_down)
	}
	_ = toolbar.AddSeparator()
	{
		icon := style.StandardIcon(qt.QStyle__SP_FileIcon, tbOpt.QStyleOption, nil)
		_ = toolbar.AddAction2(icon, dictManager_openInfo)
	}
	_ = toolbar.AddSeparator()
	{
		icon := style.StandardIcon(qt.QStyle__SP_DirOpenIcon, tbOpt.QStyleOption, nil)
		_ = toolbar.AddAction2(icon, dictManager_openDirs)
	}

	toolbar.OnActionTriggered(func(action *qt.QAction) {
		switch action.Text() {
		case dictManager_up:
			dm.toolbarUp()
		case dictManager_down:
			dm.toolbarDown()
		case dictManager_openInfo:
			dm.openInfoFile()
		case dictManager_openDirs:
			dm.openFolder(conf)
		}
	})

	table.OnCellClicked(func(row int, column int) {
		if column < 3 {
			extraOptionsWidget.Hide()
			return
		}
		dictName := table.Item(row, dm_col_dictName).Text()
		ds := dicts.DictSettingsMap[dictName]
		if ds == nil {
			extraOptionsWidget.Hide()
			return
		}
		selectedDictSettings = ds
		flagsCBWidget.SetActiveDictSetting(ds)
		volumeInput.SetValue(ds.AudioVolume)
		extraOptionsWidget.Show()
	})

	mainBox := qt.NewQVBoxLayout(dm.Dialog.QWidget)
	mainBox.AddLayout2(mainHBox.Layout(), 1)
	mainBox.AddWidget(dm.buttonBox.QWidget)

	table.SetRowCount(len(dicts.DictList))
	for index, dic := range dicts.DictList {
		dictName := dic.DictName()
		ds := dicts.DictSettingsMap[dictName]
		if ds == nil {
			slog.Info("dict manager: found new dict", "dictName", dictName)
			ds = dicts.NewDictSettings(dic, index)
			ds.Hash = dicts.Hash(dic)
			dicts.DictSettingsMap[dictName] = ds
		}
		dm.setItem(index, dictName, ds)
	}

	qs := dm.settings
	qsettings.RestoreTableColumnsWidth(
		qs,
		table,
		QS_dictManager,
	)
	table.HorizontalHeader().OnSectionResized(func(logicalIndex int, oldSize int, newSize int) {
		qsettings.SaveTableColumnsWidth(qs, table, QS_dictManager)
	})

	qsettings.SetupDialogGeometrySave(qs, dm.Dialog, QS_dictManager)
}

// updates global var dictSettingsMap
// and returns dicts order
func (dm *DictManager) updateMap() map[string]int {
	table := dm.TableWidget
	order := map[string]int{}
	count := table.RowCount()
	for index := range count {
		disable := table.Item(index, dm_col_enable).CheckState() != qt.Checked
		hideHeader := table.Item(index, dm_col_header).CheckState() != qt.Checked
		symbol := table.Item(index, dm_col_symbol).Text()
		dictName := table.Item(index, dm_col_dictName).Text()
		value := index + 1
		if disable {
			value = -value
		}
		order[dictName] = value
		ds := dicts.DictSettingsMap[dictName]
		if ds == nil {
			ds = &dicts.DictionarySettings{}
			dicts.DictSettingsMap[dictName] = ds
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
	if dm.Dialog.Exec() != int(qt.QDialog__Accepted) {
		return false
	}
	dicts.DictsOrder = dm.updateMap()

	dicts.Reorder(dicts.DictsOrder)

	for _, dic := range dicts.DictList {
		disabled := dic.Disabled()
		dic.SetDisabled(dicts.DictsOrder[dic.DictName()] < 0)
		if disabled && !dic.Disabled() {
			err := dic.Load()
			if err != nil {
				slog.Error("error in dic.Load: " + err.Error())
			}
		}
	}

	err := dicts.SaveDictsSettings(dicts.DictSettingsMap)
	if err != nil {
		slog.Error("error in saving dicts settings: " + err.Error())
	}
	return true
}
