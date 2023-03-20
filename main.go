package main

import (
	"fmt"
	"os"
	"time"

	"github.com/ilius/ayandict/pkg/config"
	"github.com/ilius/ayandict/pkg/stardict"

	// "github.com/therecipe/qt/webengine"

	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

var conf = config.MustLoad()

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

	updateWebView := func(s string) {
		// webview.SetHtml(s, core.NewQUrl())
		webview.SetHtml(s)
	}

	entry := widgets.NewQLineEdit(nil)
	entry.SetPlaceholderText("")
	entry.SetFixedHeight(25)

	okButton := widgets.NewQPushButton2("OK", nil)

	frame1 := widgets.NewQFrame(nil, 0)
	frame1Layout := widgets.NewQHBoxLayout2(frame1)
	frame1Layout.AddWidget(widgets.NewQLabel2("Query:", nil, 0), 0, 0)
	frame1Layout.AddSpacing(10)
	frame1Layout.AddWidget(entry, 0, 0)
	frame1Layout.AddSpacing(10)
	frame1Layout.AddWidget(okButton, 0, 0)

	frame2 := widgets.NewQFrame(nil, 0)
	frame2Layout := widgets.NewQHBoxLayout2(frame2)
	frame2Layout.AddWidget(webview, 0, 0)

	mainLayout := widgets.NewQVBoxLayout()
	mainLayout.AddWidget(frame1, 0, 0)
	mainLayout.AddWidget(frame2, 0, 0)

	centralWidget := widgets.NewQWidget(nil, 0)
	centralWidget.SetLayout(mainLayout)
	window.SetCentralWidget(centralWidget)

	entry.ConnectReturnPressed(func() {
		onQuery(entry.Text(), updateWebView)
	})
	okButton.ConnectClicked(func(bool) {
		onQuery(entry.Text(), updateWebView)
	})

	font := gui.NewQFont()
	if conf.FontFamily != "" {
		font.SetFamily(conf.FontFamily)
	}
	if conf.FontSize > 0 {
		font.SetPixelSize(conf.FontSize)
	}
	app.SetFont(font, "")

	LoadUserStyle(app)

	if conf.SearchOnType {
		minLength := conf.SearchOnTypeMinLength
		entry.ConnectKeyPressEvent(func(event *gui.QKeyEvent) {
			entry.KeyPressEventDefault(event)
			if event.Text() == "" {
				return
			}
			text := entry.Text()
			if len(text) < minLength {
				return
			}
			t := time.Now()
			onQuery(text, updateWebView)
			fmt.Printf("Query %#v took %v\n", text, time.Now().Sub(t))
		})
	}

	window.Show()
	app.Exec()
}
