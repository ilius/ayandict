package application

import (
	"html/template"
	"reflect"
	"sync"

	"github.com/ilius/ayandict/v2/pkg/config"
	"github.com/ilius/ayandict/v2/pkg/dictmgr/qdictmgr"
	"github.com/ilius/ayandict/v2/pkg/qerr"
	"github.com/ilius/ayandict/v2/pkg/wordwrap"
	"github.com/ilius/qt/core"
	"github.com/ilius/qt/gui"
)

var (
	conf      = &config.Config{}
	confMutex sync.Mutex

	headerTpl *template.Template
)

func ConfigFont() *gui.QFont {
	font := gui.NewQFont()
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
		qerr.Errorf("Failed to load config: %v", err)
		return false
	}
	conf = newConf

	if conf.HistoryMaxSize > 0 {
		historyMaxSize = conf.HistoryMaxSize
	}
	{
		err := readArticleStyle(conf.ArticleStyle)
		if err != nil {
			qerr.Error(err)
		}
	}
	{
		// log.Println("Parsing:", conf.HeaderTemplate)
		headerTplNew := template.New("header").Funcs(template.FuncMap{
			"wrapterms": func(terms []string, limit int) [][]string {
				return wordwrap.WordWrapByWords(terms, limit, " ", " ")
			},
		})
		headerTplNew, err := headerTplNew.Parse(conf.HeaderTemplate)
		if err != nil {
			qerr.Error(err)
		} else {
			headerTpl = headerTplNew
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
	app.SetFont(font, "")
	// widgets.QApplication_AllWidgets panics
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
}

func OpenConfig() {
	err := config.EnsureExists(conf)
	if err != nil {
		qerr.Error(err)
	}
	url := core.NewQUrl()
	url.SetScheme("file")
	url.SetPath(config.Path(), core.QUrl__TolerantMode)
	gui.QDesktopServices_OpenUrl(url)
}
