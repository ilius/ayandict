package application

import (
	"fmt"
	"sync"

	"github.com/ilius/ayandict/pkg/config"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

var (
	conf      = &config.Config{}
	confMutex sync.Mutex
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
	if conf.Style != currentStyle {
		LoadUserStyle(app)
	}

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
}
