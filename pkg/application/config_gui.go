package application

import (
	"html/template"
	"log/slog"
	"reflect"
	"sync"

	"github.com/ilius/ayandict/v3/pkg/config"
	"github.com/ilius/ayandict/v3/pkg/dictmgr/qdictmgr"
	"github.com/ilius/ayandict/v3/pkg/headerlib"
	"github.com/ilius/ayandict/v3/pkg/qlocalserver"
	"github.com/ilius/ayandict/v3/pkg/qtcommon"
	qt "github.com/mappu/miqt/qt6"
)

var (
	conf      = &config.Config{}
	confMutex sync.Mutex

	headerTpl *template.Template
)

func LoadConfig() bool {
	confMutex.Lock()
	defer confMutex.Unlock()
	newConf, err := config.Load()
	if err != nil {
		slog.Error("failed to load config: " + err.Error())
		return false
	}
	conf = newConf

	{
		tpl, err := headerlib.LoadHeaderTemplate(conf)
		if err != nil {
			slog.Error("error loading header template: " + err.Error())
		} else {
			headerTpl = tpl
			qlocalserver.SetHeaderTemplate(tpl)
		}
	}
	return true
}

func shouldReloadDicts(currentList []string, newList []string) bool {
	if len(currentList) != len(newList) {
		return true
	}
	if len(newList) == 0 {
		return false
	}
	return !reflect.DeepEqual(newList, currentList)
}

func (app *Application) ReloadFont() {
	font := qtcommon.ConfigFont(conf)
	// app.SetFont only applies to future widgets (DictManager for example)
	qt.QApplication_SetFont(font)
	for _, w := range qt.QApplication_AllWidgets() {
		w.SetFont(font)
	}
}

func (app *Application) ReloadConfig() {
	slog.Info("Reloading config")
	currentDirList := conf.DirectoryList
	fontFamily := conf.FontFamily
	fontSize := conf.FontSize

	if !LoadConfig() {
		return
	}

	if conf.FontFamily != fontFamily || conf.FontSize != fontSize {
		app.ReloadFont()
	}
	app.articleView.LoadUserStyle()

	if conf.Style != currentStyle {
		app.ReloadUserStyle()
	}
	if shouldReloadDicts(currentDirList, conf.DirectoryList) {
		qdictmgr.InitDicts(conf, true)
		app.dictManager = nil
	}
	app.headerLabel.SetWordWrap(conf.HeaderWordWrap)
	app.audioCache.ReloadConfig()

	app.onQuery(app.entry.Text(), false)

	app.updateMiscButtonsVisibility()
	app.updateMiscButtonsPadding()
	app.updateTrayMenu()
}

func OpenConfig() {
	err := config.EnsureExists(conf)
	if err != nil {
		slog.Error("error checking/creating config file: "+err.Error(), "path", config.Path())
	}
	url := qt.NewQUrl()
	url.SetScheme("file")
	url.SetPath2(config.Path(), qt.QUrl__TolerantMode)
	_ = qt.QDesktopServices_OpenUrl(url)
}
