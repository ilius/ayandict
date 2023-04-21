package main

import (
	"os"

	"github.com/ilius/ayandict/pkg/frequency"
	"github.com/therecipe/qt/widgets"
)

func main() {
	app := widgets.NewQApplication(len(os.Args), os.Args)
	window := widgets.NewQMainWindow(nil, 0)
	window.SetWindowTitle("FrequencyView")
	window.Resize2(600, 400)

	entry := widgets.NewQLineEdit(nil)
	entry.SetPlaceholderText("")
	entry.SetFixedHeight(25)

	view := frequency.NewFrequencyView("", 6)
	view.SetHorizontalHeaderItem(0, widgets.NewQTableWidgetItem2("Key", 0))
	view.SetHorizontalHeaderItem(1, widgets.NewQTableWidgetItem2("Count", 0))

	centralWidget := widgets.NewQWidget(nil, 0)
	mainLayout := widgets.NewQVBoxLayout()
	centralWidget.SetLayout(mainLayout)
	mainLayout.AddWidget(entry, 0, 0)
	mainLayout.AddWidget(view, 0, 0)
	window.SetCentralWidget(centralWidget)

	entry.ConnectReturnPressed(func() {
		view.Add(entry.Text(), 1)
	})
	window.Show()
	app.Exec()
}
