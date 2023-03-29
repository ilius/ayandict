package application

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ilius/ayandict/pkg/config"

	// "github.com/therecipe/qt/webengine"

	"github.com/ilius/ayandict/pkg/frequency"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/multimedia"
	"github.com/therecipe/qt/widgets"
)

var expanding = widgets.QSizePolicy__Expanding

var frequencyTable *frequency.FrequencyTable

func Run() {
	app := widgets.NewQApplication(len(os.Args), os.Args)
	LoadConfig(app)
	initDicts()

	frequencyTable = frequency.NewFrequencyView(conf.MostFrequentMaxSize)

	// icon := gui.NewQIcon5("./img/icon.png")

	window := widgets.NewQMainWindow(nil, 0)
	window.SetWindowTitle("AyanDict")
	window.Resize2(600, 400)

	defiTitleLabel := widgets.NewQLabel(nil, 0)
	defiTitleLabel.SetTextInteractionFlags(core.Qt__TextSelectableByMouse)
	// | core.Qt__TextSelectableByKeyboard
	defiTitleLabel.SetAlignment(core.Qt__AlignVCenter)
	defiTitleLabel.SetContentsMargins(20, 0, 0, 0)
	defiTitleLabel.SetTextFormat(core.Qt__RichText)

	webview := widgets.NewQTextBrowser(nil)
	// webview := webengine.NewQWebEngineView(nil)
	webview.SetReadOnly(true)
	webview.SetOpenExternalLinks(true)
	webview.SetOpenLinks(false)

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

	historyView := widgets.NewQListWidget(nil)

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

	addHistoryGUI = func(query string) {
		historyView.InsertItem2(0, query)
	}
	trimHistoryGUI = func(maxSize int) {
		count := historyView.Count()
		if count <= maxSize {
			return
		}
		for i := maxSize; i < count; i++ {
			historyView.TakeItem(maxSize)
		}
	}

	miscBox := widgets.NewQFrame(nil, 0)
	miscLayout := widgets.NewQVBoxLayout2(miscBox)
	miscLayout.SetContentsMargins(0, 0, 0, 0)
	reloadDictsButton := widgets.NewQPushButton2("Reload Dicts", nil)
	miscLayout.AddWidget(reloadDictsButton, 0, 0)
	openConfigButton := widgets.NewQPushButton2("Open Config", nil)
	miscLayout.AddWidget(openConfigButton, 0, 0)
	reloadConfigButton := widgets.NewQPushButton2("Reload Config", nil)
	miscLayout.AddWidget(reloadConfigButton, 0, 0)
	reloadStyleButton := widgets.NewQPushButton2("Reload Style", nil)
	miscLayout.AddWidget(reloadStyleButton, 0, 0)
	saveHistoryButton := widgets.NewQPushButton2("Save History", nil)
	miscLayout.AddWidget(saveHistoryButton, 0, 0)
	clearHistoryButton := widgets.NewQPushButton2("Clear History", nil)
	miscLayout.AddWidget(clearHistoryButton, 0, 0)

	bottomLayout := widgets.NewQHBoxLayout2(nil)
	bottomLayout.SetContentsMargins(0, 0, 0, 0)
	bottomLayout.SetSpacing(10)
	dictsButton := widgets.NewQPushButton2("Dictionaries", nil)
	bottomLayout.AddWidget(dictsButton, 0, core.Qt__AlignLeft)
	clearButton := widgets.NewQPushButton2("Clear", nil)
	bottomLayout.AddWidget(clearButton, 0, core.Qt__AlignRight)

	leftMainWidget := widgets.NewQWidget(nil, 0)
	leftMainLayout := widgets.NewQVBoxLayout2(leftMainWidget)
	leftMainLayout.SetContentsMargins(0, 0, 0, 0)
	leftMainLayout.SetSpacing(0)
	leftMainLayout.AddWidget(queryBox, 0, 0)
	leftMainLayout.AddSpacing(5)
	leftMainLayout.AddWidget(defiTitleLabel, 0, core.Qt__AlignVCenter)
	leftMainLayout.AddSpacing(5)
	leftMainLayout.AddWidget(webview, 0, 0)
	leftMainLayout.AddSpacing(5)
	leftMainLayout.AddLayout(bottomLayout, 0)

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
	resultList := NewResultListWidget(webview, defiTitleLabel)
	leftPanelLayout.AddWidget(resultList, 0, 0)

	queryWidgets := &QueryWidgets{
		Webview:    webview,
		ResultList: resultList,
	}

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
	mainLayout.AddLayout(bottomLayout, 0)

	centralWidget := widgets.NewQWidget(nil, 0)
	centralWidget.SetLayout(mainLayout)
	window.SetCentralWidget(centralWidget)

	mediaPlayer := multimedia.NewQMediaPlayer(nil, 0)

	doQuery := func(query string) {
		onQuery(query, queryWidgets, false)
		entry.SetText(query)
	}

	resetQuery := func() {
		entry.SetText("")
		resultList.Clear()
		webview.SetHtml("")
		defiTitleLabel.SetText("")
	}

	entry.ConnectReturnPressed(func() {
		onQuery(entry.Text(), queryWidgets, false)
	})
	okButton.ConnectClicked(func(bool) {
		onQuery(entry.Text(), queryWidgets, false)
	})
	resultList.ConnectKeyPressEvent(func(event *gui.QKeyEvent) {
		switch event.Text() {
		case " ":
			entry.SetFocus(core.Qt__ShortcutFocusReason)
			return
		}
		resultList.KeyPressEventDefault(event)
	})
	webview.ConnectAnchorClicked(func(link *core.QUrl) {
		host := link.Host(core.QUrl__FullyDecoded)
		// fmt.Printf(
		// 	"AnchorClicked: %#v, host=%#v = %#v\n",
		// 	link.ToString(core.QUrl__None),
		// 	host,
		// 	link.Host(core.QUrl__FullyEncoded),
		// )
		if link.Scheme() == "bword" {
			if host != "" {
				doQuery(host)
			} else {
				fmt.Printf("AnchorClicked: %#v\n", link.ToString(core.QUrl__None))
			}
			return
		}
		path := link.Path(core.QUrl__FullyDecoded)
		// fmt.Printf("scheme=%#v, host=%#v, path=%#v", link.Scheme(), host, path)
		switch link.Scheme() {
		case "":
			doQuery(path)
			return
		case "file", "http", "https":
			// fmt.Printf("host=%#v, ext=%#v", host, ext)
			switch filepath.Ext(path) {
			case ".wav", ".mp3", ".ogg":
				fmt.Println("Playing audio", link.ToString(core.QUrl__None))
				mediaPlayer.SetMedia(multimedia.NewQMediaContent2(link), nil)
				mediaPlayer.Play()
				return
			}
		}
		gui.QDesktopServices_OpenUrl(link)
	})

	// menuStyleOpt := widgets.NewQStyleOptionMenuItem()
	// style := app.Style()
	// queryMenuIcon := style.StandardIcon(widgets.QStyle__SP_ArrowUp, menuStyleOpt, nil)

	webview.ConnectContextMenuEvent(func(event *gui.QContextMenuEvent) {
		event.Ignore()
		// menu := webview.CreateStandardContextMenu2(event.GlobalPos())
		menu := webview.CreateStandardContextMenu()
		// actions := menu.Actions()
		// fmt.Println("actions", actions)
		// menu.Actions() panic
		// https://github.com/therecipe/qt/issues/1286
		// firstAction := menu.ActiveAction()
		action := widgets.NewQAction2("Query", webview)
		action.ConnectTriggered(func(checked bool) {
			text := webview.TextCursor().SelectedText()
			if text != "" {
				doQuery(text)
			}
		})
		menu.InsertAction(nil, action)
		menu.Popup(event.GlobalPos(), nil)
	})
	webview.ConnectMouseReleaseEvent(func(event *gui.QMouseEvent) {
		switch event.Button() {
		case core.Qt__MiddleButton:
			text := webview.TextCursor().SelectedText()
			if text != "" {
				doQuery(text)
			}
			return
		}
		webview.MouseReleaseEventDefault(event)
	})
	webview.ConnectKeyPressEvent(func(event *gui.QKeyEvent) {
		switch event.Text() {
		case " ":
			entry.SetFocus(core.Qt__ShortcutFocusReason)
		}
	})

	// historyView.SelectedItems() panics
	// and even after fixing panic, doesn't return anything
	// you have to use historyView.CurrentIndex()
	historyView.ConnectMousePressEvent(func(event *gui.QMouseEvent) {
		historyView.MousePressEventDefault(event)
		index := historyView.CurrentIndex()
		if index == nil {
			return
		}
		historyView.Activated(index)
	})

	historyView.ConnectItemActivated(func(item *widgets.QListWidgetItem) {
		doQuery(item.Text())
	})

	historyView.ConnectKeyPressEvent(func(event *gui.QKeyEvent) {
		switch event.Text() {
		case " ":
			entry.SetFocus(core.Qt__ShortcutFocusReason)
			return
		}
		historyView.KeyPressEventDefault(event)
	})
	historyView.ConnectItemClicked(func(item *widgets.QListWidgetItem) {
		doQuery(item.Text())
	})
	frequencyTable.ConnectItemClicked(func(item *widgets.QTableWidgetItem) {
		doQuery(frequencyTable.Keys[item.Row()])
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
		LoadConfig(app)
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
		clearHistory()
		historyView.Clear()
		// frequencyTable.Clear()
	})
	clearButton.ConnectClicked(func(checked bool) {
		resetQuery()
	})

	const dialogAccepted = int(widgets.QDialog__Accepted)

	var dictManager *DictManager
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
			for _, query := range history {
				addHistoryGUI(query)
			}
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
	reastoreMainWinGeometry(qs, window)
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

	window.Show()
	app.Exec()
}
