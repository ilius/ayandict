package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ilius/ayandict/pkg/config"
	"github.com/ilius/ayandict/pkg/stardict"

	// "github.com/therecipe/qt/webengine"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/multimedia"
	"github.com/therecipe/qt/widgets"
)

var expanding = widgets.QSizePolicy__Expanding

const (
	QS_mainwindow = "mainwindow"
	QS_geometry   = "geometry"
	QS_savestate  = "savestate"
	QS_maximized  = "maximized"
	QS_pos        = "pos"
	QS_size       = "size"
)

func main() {
	app := widgets.NewQApplication(len(os.Args), os.Args)
	LoadConfig(app)
	stardict.Init(conf.DirectoryList)

	// icon := gui.NewQIcon5("./img/icon.png")

	window := widgets.NewQMainWindow(nil, 0)
	window.SetWindowTitle("AyanDict")
	window.Resize2(600, 400)

	webview := widgets.NewQTextBrowser(nil)
	// webview := webengine.NewQWebEngineView(nil)
	webview.SetReadOnly(true)
	webview.SetOpenExternalLinks(true)
	webview.SetOpenLinks(false)

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
	queryBoxLayout.AddWidget(widgets.NewQLabel2("Query:", nil, 0), 0, 0)
	// queryBoxLayout.AddSpacing(10)
	queryBoxLayout.AddWidget(entry, 0, 0)
	// queryBoxLayout.AddSpacing(10)
	queryBoxLayout.AddWidget(okButton, 0, 0)

	historyView := widgets.NewQListWidget(nil)

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

	sideBar := widgets.NewQTabWidget(nil)
	sideBar.AddTab(historyView, "History")
	sideBar.AddTab(miscBox, "Misc")

	mainSplitter := widgets.NewQSplitter(nil)
	mainSplitter.SetSizePolicy2(expanding, expanding)
	mainSplitter.AddWidget(webview)
	mainSplitter.AddWidget(sideBar)
	mainSplitter.SetStretchFactor(0, 5)
	mainSplitter.SetStretchFactor(1, 1)

	mainLayout := widgets.NewQVBoxLayout()
	mainLayout.AddWidget(queryBox, 0, 0)
	mainLayout.AddWidget(mainSplitter, 0, 0)

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
	historyView.ConnectItemClicked(func(item *widgets.QListWidgetItem) {
		doQuery(item.Text())
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
	})
	reloadStyleButton.ConnectClicked(func(checked bool) {
		LoadUserStyle(app)
	})
	saveHistoryButton.ConnectClicked(func(checked bool) {
		SaveHistory()
	})
	clearHistoryButton.ConnectClicked(func(checked bool) {
		clearHistory()
		historyView.Clear()
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
