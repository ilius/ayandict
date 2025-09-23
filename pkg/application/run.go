package application

import (
	"log/slog"
	"os"

	"github.com/ilius/ayandict/v3/pkg/appinfo"
	"github.com/ilius/ayandict/v3/pkg/qtcommon"
	"github.com/ilius/ayandict/v3/pkg/qtcommon/qerr"
	qt "github.com/mappu/miqt/qt6"
)

func Run() {
	app := &Application{
		QApplication:   qt.NewQApplication(os.Args),
		window:         qt.NewQMainWindow(nil),
		allTextWidgets: []qtcommon.HasSetFont{},
	}
	qerr.ShowMessage = showErrorMessage
	app.style = qt.QApplication_Style()
	app.bottomBoxStyleOpt = qt.NewQStyleOptionButton()
	qt.QCoreApplication_SetApplicationName(appinfo.APP_DESC)

	if cacheDir == "" {
		slog.Error("cacheDir is empty")
	}
	{
		err := os.MkdirAll(cacheDir, 0o755)
		if err != nil {
			slog.Error("error in MkdirAll: " + err.Error())
		}
	}

	app.Run()
}
