package main

import (
	"os"
	"runtime"

	"github.com/ilius/ayandict/v3/pkg/activity"
	"github.com/ilius/ayandict/v3/pkg/application/frequency"
	"github.com/ilius/ayandict/v3/pkg/config"
	qt "github.com/mappu/miqt/qt6"
)

func main() {
	runtime.LockOSThread()
	_ = qt.NewQApplication(os.Args)
	window := qt.NewQMainWindow2()
	window.SetWindowTitle("FrequencyView")
	window.Resize(600, 400)

	entry := qt.NewQLineEdit2()
	entry.SetPlaceholderText("")
	entry.SetFixedHeight(25)

	activityStorage := activity.NewActivityStorage(config.Default(), config.GetConfigDir())

	view := frequency.NewFrequencyView(activityStorage, 6)
	view.SetHorizontalHeaderItem(0, qt.NewQTableWidgetItem2("Key"))
	view.SetHorizontalHeaderItem(1, qt.NewQTableWidgetItem2("Count"))

	centralWidget := qt.NewQWidget2()
	mainLayout := qt.NewQVBoxLayout2()
	centralWidget.SetLayout(mainLayout.Layout())
	mainLayout.AddWidget(entry.QWidget)
	mainLayout.AddWidget(view.QWidget)
	window.SetCentralWidget(centralWidget)

	entry.OnReturnPressed(func() {
		view.Add(entry.Text(), 1)
	})
	window.Show()
	_ = qt.QApplication_Exec()
}
