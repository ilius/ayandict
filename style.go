package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ilius/ayandict/pkg/config"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

func LoadUserStyle(app *widgets.QApplication) {
	configDir := config.GetConfigDir()
	stylePath := filepath.Join(configDir, "style.qss")
	_, err := os.Stat(stylePath)
	if err != nil {
		if !os.IsNotExist(err) {
			fmt.Printf("Error loading %s: %v\n", stylePath, err)
		}
		return
	}
	fmt.Println("Loading", stylePath)
	file := core.NewQFile2(stylePath)
	file.Open(core.QIODevice__ReadOnly | core.QIODevice__Text)
	stream := core.NewQTextStream2(file)
	app.SetStyleSheet(stream.ReadAll())
}
