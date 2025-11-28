package application

import (
	"log/slog"

	qt "github.com/mappu/miqt/qt6"
)

// the current conf.Style value (unchanged config value)
var currentStyle = ""

func (app *Application) LoadUserStyle() {
	stylePath := conf.Style
	if stylePath == "" {
		return
	}
	slog.Info("Loading user style", "stylePath", stylePath)
	file := qt.NewQFile2(stylePath)
	if !file.Exists() {
		slog.Error("style file does not exist", "stylePath", stylePath)
		return
	}
	if !file.Open(qt.QIODeviceBase__ReadOnly) {
		slog.Error("failed to open style file", "error", file.ErrorString())
		return
	}
	defer file.Close()
	styleBytes := file.ReadAll()
	slog.Info("Loaded user style", "stylePath", stylePath, "size", len(styleBytes))
	app.SetStyleSheet(string(styleBytes))
	currentStyle = conf.Style
}

func (app *Application) ReloadUserStyle() {
	if conf.Style == "" {
		app.SetStyleSheet("")
		currentStyle = ""
		return
	}
	app.LoadUserStyle()
	app.articleView.LoadUserStyle()
}
