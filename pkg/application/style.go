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
	configDir := config.GetConfigDir()
	stylePath := conf.Style
	if stylePath == "" {
		return
	}
	stylePath = PathFromUnix(stylePath)
	if !filepath.IsAbs(stylePath) {
		stylePath = filepath.Join(configDir, stylePath)
	}
	// _, err := os.Stat(stylePath)
	// if err != nil {
	// 	qerr.Errorf("Error loading style file %#v: %v\n", stylePath, err)
	// 	return
	// }
	slog.Info("Loading user style", "stylePath", stylePath)
	// file := qt.NewQFile2(stylePath)
	// file.Open(qt.QIODeviceBase__ReadOnly | qt.QIODeviceBase__Text)
	// stream := qt.NewQTextStream2(file)
	// app.SetStyleSheet(stream.ReadAll())
	styleBytes, err := os.ReadFile(stylePath)
	if err != nil {
		slog.Error("Error loading style file: "+err.Error(), "stylePath", stylePath)
		return
	}
	_ = qt.QApplication_SetStyleWithStyle(string(styleBytes))
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
