package main

import (
	"fmt"
	"sync"

	"github.com/ilius/ayandict/pkg/config"
	"github.com/therecipe/qt/widgets"
)

var (
	conf      = config.MustLoad()
	confMutex sync.Mutex
)

func ReloadConfig(app *widgets.QApplication) {
	confMutex.Lock()
	defer confMutex.Unlock()
	newConf, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
	}
	conf = newConf
	if conf.Style != currentStyle {
		LoadUserStyle(app)
	}
}
