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

var frequencyView *frequency.FrequencyView

const (
	QS_mainwindow = "mainwindow"
	QS_geometry   = "geometry"
	QS_savestate  = "savestate"
	QS_maximized  = "maximized"
	QS_pos        = "pos"
	QS_size       = "size"
)

func Run() {
	app := widgets.NewQApplication(len(os.Args), os.Args)
	LoadConfig(app)
	initDicts()

	frequencyView = frequency.NewFrequencyView(conf.MostFrequentMaxSize)

	// icon := gui.NewQIcon5("./img/icon.png")

	window := widgets.NewQMainWindow(nil, 0)
	window.SetWindowTitle("AyanDict")
	window.Resize2(600, 400)

	webview := widgets.NewQTextBrowser(nil)
	// webview := webengine.NewQWebEngineView(nil)
	webview.SetReadOnly(true)
	webview.SetOpenExternalLinks(true)
	webview.SetOpenLinks(false)
	// webview.SetContentsMargins(0, 0, 0, 0)

	updateWebView := func(s string) {
		// webview.SetHtml(s, core.NewQUrl())
		webview.SetHtml(s)
	}

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

	frequencyView.SetHorizontalHeaderItem(
		0,
		widgets.NewQTableWidgetItem2("Query", 0),
	)
	frequencyView.SetHorizontalHeaderItem(
		1,
		widgets.NewQTableWidgetItem2("Count", 0),
	)
	if !conf.MostFrequentDisable {
		frequencyView.LoadFromFile(frequencyFilePath())
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
	leftMainLayout.AddWidget(queryBox, 0, 0)
	leftMainLayout.AddWidget(webview, 0, 0)
	leftMainLayout.AddLayout(bottomLayout, 0)

	activityTypeCombo := widgets.NewQComboBox(nil)
	activityTypeCombo.AddItems([]string{
		"Recent",
		"Most Frequent",
	})

	frequencyView.Hide()

	activityWidget := widgets.NewQWidget(nil, 0)
	activityLayout := widgets.NewQVBoxLayout2(activityWidget)
	activityLayout.SetContentsMargins(5, 5, 5, 5)
	activityLayout.AddWidget(activityTypeCombo, 0, 0)
	activityLayout.AddWidget(historyView, 0, 0)
	activityLayout.AddWidget(frequencyView, 0, 0)

	activityTypeCombo.ConnectCurrentIndexChanged(func(index int) {
		switch index {
		case 0:
			historyView.Show()
			frequencyView.Hide()
		case 1:
			historyView.Hide()
			frequencyView.Show()
		}
	})

	sideBar := widgets.NewQTabWidget(nil)
	sideBar.AddTab(activityWidget, "Activity")
	sideBar.AddTab(miscBox, "Misc")

	mainSplitter := widgets.NewQSplitter(nil)
	mainSplitter.SetSizePolicy2(expanding, expanding)
	mainSplitter.AddWidget(leftMainWidget)
	mainSplitter.AddWidget(sideBar)
	mainSplitter.SetStretchFactor(0, 5)
	mainSplitter.SetStretchFactor(1, 1)

	mainLayout := widgets.NewQVBoxLayout()
	mainLayout.SetContentsMargins(5, 5, 5, 5)
	mainLayout.AddWidget(mainSplitter, 0, 0)
	mainLayout.AddLayout(bottomLayout, 0)

	centralWidget := widgets.NewQWidget(nil, 0)
	centralWidget.SetLayout(mainLayout)
	window.SetCentralWidget(centralWidget)

	mediaPlayer := multimedia.NewQMediaPlayer(nil, 0)

	doQuery := func(query string) {
		onQuery(query, updateWebView, false)
		entry.SetText(query)
	}

	resetQuery := func() {
		entry.SetText("")
		updateWebView("")
	}

	entry.ConnectReturnPressed(func() {
		onQuery(entry.Text(), updateWebView, false)
	})
	okButton.ConnectClicked(func(bool) {
		onQuery(entry.Text(), updateWebView, false)
	})
	webview.ConnectAnchorClicked(func(link *core.QUrl) {
		host := link.Host(core.QUrl__FullyDecoded)
		if link.Scheme() == "bword" {
			doQuery(host)
			return
		}
		path := link.Path(core.QUrl__FullyDecoded)
		// fmt.Printf("scheme=%#v, host=%#v, path=%#v", link.Scheme(), host, path)
		switch link.Scheme() {
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
	webview.ConnectMouseReleaseEvent(func(event *gui.QMouseEvent) {
		switch event.Button() {
		case core.Qt__MiddleButton:
			doQuery(webview.TextCursor().SelectedText())
			return
		}
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
	frequencyView.ConnectItemClicked(func(item *widgets.QTableWidgetItem) {
		index := item.Row()
		key := frequencyView.Keys[index]
		doQuery(key)
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
		onQuery(entry.Text(), updateWebView, false)
	})
	reloadStyleButton.ConnectClicked(func(checked bool) {
		LoadUserStyle(app)
	})
	saveHistoryButton.ConnectClicked(func(checked bool) {
		SaveHistory()
		SaveFrequency()
	})
	clearHistoryButton.ConnectClicked(func(checked bool) {
		clearHistory()
		historyView.Clear()
		// frequencyView.Clear()
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
			onQuery(entry.Text(), updateWebView, false)
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
			onQuery(text, updateWebView, true)
		}
	})

	qsettings := core.NewQSettings("ilius", "ayandict", window)
	reastoreMainWinGeometry(qsettings, window)
	window.ConnectResizeEvent(func(event *gui.QResizeEvent) {
		saveMainWinGeometry(qsettings, window)
	})
	window.ConnectMoveEvent(func(event *gui.QMoveEvent) {
		saveMainWinGeometry(qsettings, window)
	})

	window.Show()
	app.Exec()
}

func reastoreSetting(qsettings *core.QSettings, key string, apply func(*core.QVariant)) {
	if !qsettings.Contains(key) {
		return
	}
	apply(qsettings.Value(key, core.NewQVariant1(nil)))
}

func reastoreBoolSetting(
	qsettings *core.QSettings,
	key string, _default bool,
	apply func(bool),
) {
	if !qsettings.Contains(key) {
		apply(_default)
		return
	}
	apply(qsettings.Value(key, core.NewQVariant1(nil)).ToBool())
}

func saveMainWinGeometry(qsettings *core.QSettings, window *widgets.QMainWindow) {
	qsettings.BeginGroup(QS_mainwindow)

	qsettings.SetValue(QS_geometry, core.NewQVariant13(window.SaveGeometry()))
	qsettings.SetValue(QS_savestate, core.NewQVariant13(window.SaveState(0)))
	qsettings.SetValue(QS_maximized, core.NewQVariant9(window.IsMaximized()))
	if !window.IsMaximized() {
		qsettings.SetValue(QS_pos, core.NewQVariant27(window.Pos()))
		qsettings.SetValue(QS_size, core.NewQVariant25(window.Size()))
	}

	qsettings.EndGroup()
}

func reastoreMainWinGeometry(qsettings *core.QSettings, window *widgets.QMainWindow) {
	qsettings.BeginGroup(QS_mainwindow)

	reastoreSetting(qsettings, QS_geometry, func(value *core.QVariant) {
		window.RestoreGeometry(value.ToByteArray())
	})
	reastoreSetting(qsettings, QS_savestate, func(value *core.QVariant) {
		window.RestoreState(value.ToByteArray(), 0)
	})
	reastoreBoolSetting(qsettings, QS_maximized, false, func(maximized bool) {
		if maximized {
			window.ShowMaximized()
			return
		}
		reastoreSetting(qsettings, QS_pos, func(value *core.QVariant) {
			window.Move(value.ToPoint())
		})
		reastoreSetting(qsettings, QS_size, func(value *core.QVariant) {
			window.Resize(value.ToSize())
		})
	})

	qsettings.EndGroup()
}
