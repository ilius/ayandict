package application

import (
	"fmt"
	"html/template"
	"reflect"
	"sync"

	"github.com/ilius/ayandict/pkg/config"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

var (
	conf      = &config.Config{}
	confMutex sync.Mutex

	headerTpl *template.Template
)

func LoadConfig(app *widgets.QApplication) {
	confMutex.Lock()
	defer confMutex.Unlock()
	newConf, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		return
	}
	conf = newConf

	font := gui.NewQFont()
	if conf.FontFamily != "" {
		font.SetFamily(conf.FontFamily)
	}
	if conf.FontSize > 0 {
		font.SetPixelSize(conf.FontSize)
	}
	app.SetFont(font, "")

	if conf.HistoryMaxSize > 0 {
		historyMaxSize = conf.HistoryMaxSize
	}
	{
		err := readDefinitionStyle(conf.DefinitionStyle)
		if err != nil {
			fmt.Println(err)
		}
	}
	{
		// fmt.Println("Parsing:", conf.HeaderTemplate)
		headerTplNew, err := template.New("header").Parse(conf.HeaderTemplate)
		if err != nil {
			fmt.Println(err)
		} else {
			headerTpl = headerTplNew
		}
	}
}

func ReloadConfig(app *widgets.QApplication) {
	currentDirList := conf.DirectoryList

	LoadConfig(app)

	if conf.Style != currentStyle {
		ReloadUserStyle(app)
	}

	if !reflect.DeepEqual(conf.DirectoryList, currentDirList) {
		reloadDicts()
	}
}
