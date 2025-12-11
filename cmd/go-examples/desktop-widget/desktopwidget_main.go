package main

import (
	"log/slog"
	"os"

	"github.com/ilius/ayandict/v3/pkg/desktopwidget"
	qt "github.com/mappu/miqt/qt6"
)

func main() {
	_ = qt.NewQApplication(os.Args)
	onActivate := func(_ *qt.QMouseEvent) {
		slog.Info("----- desktopwidget.main: onActivate")
	}
	actions := []*qt.QAction{}
	widget := desktopwidget.NewDekstopWidget(onActivate, actions)
	if widget == nil {
		return
	}
	widget.Show()
	qt.QApplication_Exec()
}
