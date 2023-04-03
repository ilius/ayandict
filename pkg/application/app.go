package application

import (
	"fmt"
	"os"

	"github.com/ilius/ayandict/pkg/config"

	// "github.com/therecipe/qt/webengine"

	"github.com/ilius/ayandict/pkg/frequency"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

var dictManager *DictManager

func Run() {
	app := widgets.NewQApplication(len(os.Args), os.Args)
	LoadConfig(app)
	LoadUserStyle(app)
	initDicts()

	frequencyTable = frequency.NewFrequencyView(conf.MostFrequentMaxSize)

	// icon := gui.NewQIcon5("./img/icon.png")

	window := widgets.NewQMainWindow(nil, 0)
	window.SetWindowTitle("AyanDict")
	window.Resize2(600, 400)

	headerLabel := widgets.NewQLabel(nil, 0)
	headerLabel.SetTextInteractionFlags(core.Qt__TextSelectableByMouse)
	// | core.Qt__TextSelectableByKeyboard
	headerLabel.SetAlignment(core.Qt__AlignVCenter)
	headerLabel.SetContentsMargins(20, 0, 0, 0)
	headerLabel.SetTextFormat(core.Qt__RichText)
	headerLabel.SetWordWrap(true)

	articleView := NewArticleView()

	entry := widgets.NewQLineEdit(nil)
	entry.SetPlaceholderText("")
	entry.SetFixedHeight(25)

	okButton := widgets.NewQPushButton2("OK", nil)

	queryBox := widgets.NewQFrame(nil, 0)
	queryBoxLayout := widgets.NewQHBoxLayout2(queryBox)
	queryBoxLayout.SetContentsMargins(5, 5, 5, 0)
	queryBoxLayout.SetSpacing(10)
	queryBoxLayout.AddWidget(widgets.NewQLabel2("Query:", nil, 0), 0, 0)
	queryBoxLayout.AddWidget(entry, 0, 0)
	queryBoxLayout.AddWidget(okButton, 0, 0)
	// queryBoxLayout.SetSpacing(10)

	historyView := NewHistoryView()

	frequencyTable.SetHorizontalHeaderItem(
		0,
		widgets.NewQTableWidgetItem2("Query", 0),
	)
	frequencyTable.SetHorizontalHeaderItem(
		1,
		widgets.NewQTableWidgetItem2("Count", 0),
	)
	if !conf.MostFrequentDisable {
		frequencyTable.LoadFromFile(frequencyFilePath())
	}
	// TODO: save the width of 2 columns

	miscBox := widgets.NewQFrame(nil, 0)
	miscLayout := widgets.NewQVBoxLayout2(miscBox)
	miscLayout.SetContentsMargins(0, 0, 0, 0)

	saveHistoryButton := widgets.NewQPushButton2("Save History", nil)
	miscLayout.AddWidget(saveHistoryButton, 0, 0)
	clearHistoryButton := widgets.NewQPushButton2("Clear History", nil)
	miscLayout.AddWidget(clearHistoryButton, 0, 0)

	reloadDictsButton := widgets.NewQPushButton2("Reload Dicts", nil)
	miscLayout.AddWidget(reloadDictsButton, 0, 0)
	closeDictsButton := widgets.NewQPushButton2("Close Dicts", nil)
	miscLayout.AddWidget(closeDictsButton, 0, 0)
	reloadStyleButton := widgets.NewQPushButton2("Reload Style", nil)
	miscLayout.AddWidget(reloadStyleButton, 0, 0)

	bottomBox := widgets.NewQHBoxLayout2(nil)
	bottomBox.SetContentsMargins(0, 0, 0, 0)
	bottomBox.SetSpacing(10)

	bottomBoxStyleOpt := widgets.NewQStyleOptionButton()
	style := app.Style()

	newIconTextButton := func(label string, pix widgets.QStyle__StandardPixmap) *widgets.QPushButton {
		return widgets.NewQPushButton3(
			style.StandardIcon(
				pix, bottomBoxStyleOpt, nil,
			),
			label, nil,
		)
	}

	dictsButton := newIconTextButton("Dictionaries", widgets.QStyle__SP_FileDialogDetailedView)
	bottomBox.AddWidget(dictsButton, 0, core.Qt__AlignLeft)

	aboutButton := newIconTextButton("About", widgets.QStyle__SP_MessageBoxInformation)
	bottomBox.AddWidget(aboutButton, 0, core.Qt__AlignLeft)

	bottomBox.AddStretch(1)

	openConfigButton := newIconTextButton("Config", widgets.QStyle__SP_DialogOpenButton)
	bottomBox.AddWidget(openConfigButton, 0, 0)
	reloadConfigButton := newIconTextButton("Reload", widgets.QStyle__SP_BrowserReload)
	bottomBox.AddWidget(reloadConfigButton, 0, 0)

	bottomBox.AddStretch(1)

	clearButton := widgets.NewQPushButton2("Clear", nil)
	bottomBox.AddWidget(clearButton, 0, core.Qt__AlignRight)

	leftMainWidget := widgets.NewQWidget(nil, 0)
	leftMainLayout := widgets.NewQVBoxLayout2(leftMainWidget)
	leftMainLayout.SetContentsMargins(0, 0, 0, 0)
	leftMainLayout.SetSpacing(0)
	leftMainLayout.AddWidget(queryBox, 0, 0)
	leftMainLayout.AddSpacing(5)
	leftMainLayout.AddWidget(headerLabel, 0, core.Qt__AlignVCenter)
	leftMainLayout.AddSpacing(5)
	leftMainLayout.AddWidget(articleView, 0, 0)
	leftMainLayout.AddSpacing(5)
	leftMainLayout.AddLayout(bottomBox, 0)

	activityTypeCombo := widgets.NewQComboBox(nil)
	activityTypeCombo.AddItems([]string{
		"Recent",
		"Most Frequent",
	})

	frequencyTable.Hide()

	activityWidget := widgets.NewQWidget(nil, 0)
	activityLayout := widgets.NewQVBoxLayout2(activityWidget)
	activityLayout.SetContentsMargins(5, 5, 5, 5)
	activityLayout.AddWidget(activityTypeCombo, 0, 0)
	activityLayout.AddWidget(historyView, 0, 0)
	activityLayout.AddWidget(frequencyTable, 0, 0)

	activityTypeCombo.ConnectCurrentIndexChanged(func(index int) {
		switch index {
		case 0:
			historyView.Show()
			frequencyTable.Hide()
		case 1:
			historyView.Hide()
			frequencyTable.Show()
		}
	})

	leftPanel := widgets.NewQWidget(nil, 0)
	leftPanelLayout := widgets.NewQVBoxLayout2(leftPanel)
	leftPanelLayout.AddWidget(widgets.NewQLabel2("Results", nil, 0), 0, 0)
	resultList := NewResultListWidget(articleView, headerLabel)
	leftPanelLayout.AddWidget(resultList, 0, 0)

	queryWidgets := &QueryWidgets{
		ArticleView: articleView,
		ResultList:  resultList,
		HeaderLabel: headerLabel,
		HistoryView: historyView,
	}

	doQuery := func(query string) {
		onQuery(query, queryWidgets, false)
		entry.SetText(query)
	}
	articleView.doQuery = doQuery
	historyView.doQuery = doQuery

	rightPanel := widgets.NewQTabWidget(nil)
	rightPanel.AddTab(activityWidget, "Activity")
	rightPanel.AddTab(miscBox, "Misc")

	mainSplitter := widgets.NewQSplitter(nil)
	mainSplitter.SetSizePolicy2(expanding, expanding)
	mainSplitter.AddWidget(leftPanel)
	mainSplitter.AddWidget(leftMainWidget)
	mainSplitter.AddWidget(rightPanel)
	mainSplitter.SetStretchFactor(0, 1)
	mainSplitter.SetStretchFactor(1, 5)
	mainSplitter.SetStretchFactor(2, 1)

	mainLayout := widgets.NewQVBoxLayout()
	mainLayout.SetContentsMargins(5, 5, 5, 5)
	mainLayout.AddWidget(mainSplitter, 0, 0)
	mainLayout.AddLayout(bottomBox, 0)

	centralWidget := widgets.NewQWidget(nil, 0)
	centralWidget.SetLayout(mainLayout)
	window.SetCentralWidget(centralWidget)

	resetQuery := func() {
		entry.SetText("")
		resultList.Clear()
		articleView.SetHtml("")
		headerLabel.SetText("")
	}

	entry.ConnectReturnPressed(func() {
		onQuery(entry.Text(), queryWidgets, false)
	})
	okButton.ConnectClicked(func(bool) {
		onQuery(entry.Text(), queryWidgets, false)
	})
	aboutButton.ConnectClicked(func(bool) {
		aboutClicked(window)
	})
	resultList.ConnectKeyPressEvent(func(event *gui.QKeyEvent) {
		switch event.Text() {
		case " ":
			entry.SetFocus(core.Qt__ShortcutFocusReason)
			return
		}
		resultList.KeyPressEventDefault(event)
	})

	articleView.SetupCustomHandlers()
	articleView.ConnectKeyPressEvent(func(event *gui.QKeyEvent) {
		switch event.Text() {
		case " ":
			entry.SetFocus(core.Qt__ShortcutFocusReason)
		}
	})

	historyView.SetupCustomHandlers()
	historyView.ConnectKeyPressEvent(func(event *gui.QKeyEvent) {
		switch event.Text() {
		case " ":
			entry.SetFocus(core.Qt__ShortcutFocusReason)
			return
		}
		historyView.KeyPressEventDefault(event)
	})

	frequencyTable.ConnectItemClicked(func(item *widgets.QTableWidgetItem) {
		frequencyTable.ItemActivated(item)
	})
	frequencyTable.ConnectItemActivated(func(item *widgets.QTableWidgetItem) {
		key := frequencyTable.Keys[item.Row()]
		doQuery(key)
		newRow := frequencyTable.KeyMap[key]
		// item.Column() panics!
		frequencyTable.SetCurrentCell(newRow, 0)
	})
	reloadDictsButton.ConnectClicked(func(checked bool) {
		reloadDicts()
	})
	closeDictsButton.ConnectClicked(func(checked bool) {
		closeDicts()
	})
	openConfigButton.ConnectClicked(func(checked bool) {
		err := config.EnsureExists(conf)
		if err != nil {
			fmt.Println(err)
		}
		url := core.NewQUrl()
		url.SetScheme("file")
		url.SetPath(config.Path(), core.QUrl__TolerantMode)
		gui.QDesktopServices_OpenUrl(url)
	})
	reloadConfigButton.ConnectClicked(func(checked bool) {
		ReloadConfig(app)
		onQuery(entry.Text(), queryWidgets, false)
	})
	reloadStyleButton.ConnectClicked(func(checked bool) {
		LoadUserStyle(app)
		onQuery(entry.Text(), queryWidgets, false)
	})
	saveHistoryButton.ConnectClicked(func(checked bool) {
		SaveHistory()
		SaveFrequency()
	})
	clearHistoryButton.ConnectClicked(func(checked bool) {
		historyView.ClearHistory()
		frequencyTable.Clear()
		SaveFrequency()
	})
	clearButton.ConnectClicked(func(checked bool) {
		resetQuery()
	})

	const dialogAccepted = int(widgets.QDialog__Accepted)

	dictsButton.ConnectClicked(func(checked bool) {
		if dictManager == nil {
			dictManager = NewDictManager(app, window)
		}
		if dictManager.Dialog.Exec() == dialogAccepted {
			SaveDictManagerDialog(dictManager)
			onQuery(entry.Text(), queryWidgets, false)
		}
	})

	if !conf.HistoryDisable {
		err := LoadHistory()
		if err != nil {
			fmt.Println(err)
		} else {
			historyView.AddHistoryList(history)
		}
	}

	entry.ConnectKeyPressEvent(func(event *gui.QKeyEvent) {
		entry.KeyPressEventDefault(event)
		switch event.Text() {
		case "", "\b":
			return
		case "\x1b":
			// Escape, is there a more elegant way?
			resetQuery()
			return
		}
		if conf.SearchOnType {
			text := entry.Text()
			if len(text) < conf.SearchOnTypeMinLength {
				return
			}
			onQuery(text, queryWidgets, true)
		}
	})

	qs := core.NewQSettings("ilius", "ayandict", window)
	restoreSplitterSizes(qs, mainSplitter, QS_mainSplitter)
	restoreMainWinGeometry(app, qs, window)
	window.ConnectResizeEvent(func(event *gui.QResizeEvent) {
		saveMainWinGeometry(qs, window)
	})
	window.ConnectMoveEvent(func(event *gui.QMoveEvent) {
		saveMainWinGeometry(qs, window)
	})
	restoreTableColumnsWidth(
		qs,
		frequencyTable.QTableWidget,
		QS_frequencyTable,
	)
	// frequencyTable.ConnectColumnResized does not work
	frequencyTable.HorizontalHeader().ConnectSectionResized(func(logicalIndex int, oldSize int, newSize int) {
		saveTableColumnsWidth(qs, frequencyTable.QTableWidget, QS_frequencyTable)
	})

	// QSplitter.Sizes() panics:
	// interface conversion: interface {} is []interface {}, not []int

	setupSplitterSizesSave(qs, mainSplitter, QS_mainSplitter)

	window.Show()
	app.Exec()
}
