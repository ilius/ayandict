package application

import (
	"html/template"
	"log/slog"
	"reflect"
	"sync"

	"github.com/ilius/ayandict/v3/pkg/config"
	"github.com/ilius/ayandict/v3/pkg/dictmgr/qdictmgr"
	"github.com/ilius/ayandict/v3/pkg/headerlib"
	qt "github.com/mappu/miqt/qt6"
)

var (
	conf      = &config.Config{}
	confMutex sync.Mutex

	headerTpl *template.Template
)

func ConfigFont() *qt.QFont {
	font := qt.NewQFont()
	if conf.FontFamily != "" {
		font.SetFamily(conf.FontFamily)
	}
	if conf.FontSize > 0 {
		font.SetPixelSize(conf.FontSize)
	}
	return font
}

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
		err := readArticleStyle(conf.ArticleStyle)
		if err != nil {
			slog.Error("error reading srticle style: " + err.Error())
		}
	}
	{
		tpl, err := headerlib.LoadHeaderTemplate(conf)
		if err != nil {
			slog.Error("error loading header template: " + err.Error())
		} else {
			headerTpl = tpl
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
	font := ConfigFont()
	// app.SetFont only applies to future widgets (DictManager for example)
	qt.QApplication_SetFont2(font, "")
	// qt.QApplication_AllWidgets panics
	// app.AllWidgets() panics
	// window.Children() panics
	for _, w := range app.allTextWidgets {
		w.SetFont(font)
	}
}

func (app *Application) ReloadConfig() {
	currentDirList := conf.DirectoryList
	fontFamily := conf.FontFamily
	fontSize := conf.FontSize

	if !LoadConfig() {
		return
	}

	if conf.FontFamily != fontFamily || conf.FontSize != fontSize {
		app.ReloadFont()
	}

	if conf.Style != currentStyle {
		app.ReloadUserStyle()
	}
	if shouldReloadDicts(currentDirList, conf.DirectoryList) {
		qdictmgr.InitDicts(conf, true)
		app.dictManager = nil
	}
	app.headerLabel.ReloadConfig()
	audioCache.ReloadConfig()
	onQuery(app.entry.Text(), app.queryArgs, false)
}

func OpenConfig() {
	err := config.EnsureExists(conf)
	if err != nil {
		slog.Error("error checking/creating config file: " + err.Error())
	}
	url := qt.NewQUrl()
	url.SetScheme("file")
	url.SetPath2(config.Path(), qt.QUrl__TolerantMode)
	_ = qt.QDesktopServices_OpenUrl(url)
}
