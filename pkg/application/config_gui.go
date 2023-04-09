package application

import (
	"html/template"
	"log"
	"reflect"
	"sync"

	"github.com/ilius/ayandict/pkg/config"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
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

func LoadConfig(app *widgets.QApplication) {
	confMutex.Lock()
	defer confMutex.Unlock()
	newConf, err := config.Load()
	if err != nil {
		log.Printf("Failed to load config: %v\n", err)
		return
	}
	conf = newConf

	if conf.HistoryMaxSize > 0 {
		historyMaxSize = conf.HistoryMaxSize
	}
	{
		err := readArticleStyle(conf.ArticleStyle)
		if err != nil {
			log.Println(err)
		}
	}
	{
		// log.Println("Parsing:", conf.HeaderTemplate)
		headerTplNew, err := template.New("header").Parse(conf.HeaderTemplate)
		if err != nil {
			log.Println(err)
		} else {
			headerTpl = headerTplNew
		}
	}
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

func ReloadConfig(app *widgets.QApplication) {
	currentDirList := conf.DirectoryList

	LoadConfig(app)
	app.SetFont(ConfigFont(), "")

	if conf.Style != currentStyle {
		ReloadUserStyle(app)
	}
	if shouldReloadDicts(currentDirList, conf.DirectoryList) {
		reloadDicts()
	}
}

func OpenConfig() {
	err := config.EnsureExists(conf)
	if err != nil {
		log.Println(err)
	}
	url := core.NewQUrl()
	url.SetScheme("file")
	url.SetPath(config.Path(), core.QUrl__TolerantMode)
	gui.QDesktopServices_OpenUrl(url)
}
