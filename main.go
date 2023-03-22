package main

import (
	"fmt"
	"os"

	"github.com/ilius/ayandict/pkg/config"
	"github.com/ilius/ayandict/pkg/stardict"

	// "github.com/therecipe/qt/webengine"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

var expanding = widgets.QSizePolicy__Expanding

func main() {
	stardict.Init()
	app := widgets.NewQApplication(len(os.Args), os.Args)
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
	reloadDictsButton := widgets.NewQPushButton2("Reload Dictionaries", nil)
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
		fmt.Printf("scheme=%#v, host=%#v, path=%#v", link.Scheme(), host, path)
		// if path == "" {
		// 	ext := filepath.Ext(host)
		// 	// fmt.Printf("host=%#v, ext=%#v", host, ext)
		// 	switch ext {
		// 	case ".wav", ".mp3", ".ogg":

		// 	}
		// }
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

	LoadConfig(app)
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

	window.Show()
	app.Exec()
}
