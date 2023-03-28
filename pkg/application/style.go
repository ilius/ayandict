package application

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ilius/ayandict/pkg/config"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

// the current conf.Style value (unchanged config value)
var currentStyle = ""

func LoadUserStyle(app *widgets.QApplication) {
	configDir := config.GetConfigDir()
	stylePath := conf.Style
	if stylePath == "" {
		return
	}
	stylePath = PathFromUnix(stylePath)
	if !filepath.IsAbs(stylePath) {
		stylePath = filepath.Join(configDir, stylePath)
	}
	_, err := os.Stat(stylePath)
	if err != nil {
		fmt.Printf("Error loading style file %#v: %v\n", stylePath, err)
		return
	}
	fmt.Println("Loading", stylePath)
	file := core.NewQFile2(stylePath)
	file.Open(core.QIODevice__ReadOnly | core.QIODevice__Text)
	stream := core.NewQTextStream2(file)
	app.SetStyleSheet(stream.ReadAll())
	currentStyle = conf.Style
}
