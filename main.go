package main

import (
	"os"

	"github.com/ilius/ayandict/pkg/stardict"
	"github.com/therecipe/qt/widgets"
)

func main() {
	stardict.Init()
	app := widgets.NewQApplication(len(os.Args), os.Args)
	// icon := gui.NewQIcon5("./img/icon.png")

	window := widgets.NewQMainWindow(nil, 0)
	window.SetWindowTitle("AyanDict")
	window.Resize2(600, 400)

	webview := widgets.NewQTextBrowser(nil)
	webview.SetReadOnly(true)
	webview.SetOpenExternalLinks(true)

	entry := widgets.NewQLineEdit(nil)
	entry.SetPlaceholderText("")
	entry.SetFixedHeight(25)
	entry.ConnectReturnPressed(func() {
		onQuery(entry.Text(), webview)
	})

	okButton := widgets.NewQPushButton2("OK", nil)
	okButton.ConnectClicked(func(bool) {
		onQuery(entry.Text(), webview)
	})

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
	window.Show()
	app.Exec()
}
