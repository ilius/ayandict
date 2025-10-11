package application

import (
	"log/slog"
	"os"
	"path/filepath"

	"github.com/ilius/ayandict/v3/pkg/config"
	qt "github.com/mappu/miqt/qt6"
)

// the current conf.Style value (unchanged config value)
var currentStyle = ""

var definitionStyleString = ""

func readArticleStyle(stylePath string) error {
	if stylePath == "" {
		return nil
	}
	configDir := config.GetConfigDir()
	stylePath = PathFromUnix(stylePath)
	if !filepath.IsAbs(stylePath) {
		stylePath = filepath.Join(configDir, stylePath)
	}
	_, err := os.Stat(stylePath)
	if err != nil {
		return err
	}
	styleBytes, err := os.ReadFile(stylePath)
	if err != nil {
		return err
	}
	definitionStyleString = "<style>" + string(styleBytes) + "</style>"
	return nil
}

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
	{
		err := readArticleStyle(conf.ArticleStyle)
		if err != nil {
			slog.Error("error in readArticleStyle: " + err.Error())
		}
	}
}

func (app *Application) ReloadUserStyle() {
	if conf.Style == "" {
		app.SetStyleSheet("")
		currentStyle = ""
		return
	}
	app.LoadUserStyle()
}
