package application

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ilius/ayandict/pkg/config"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

// the current conf.Style value (unchanged config value)
var currentStyle = ""

var definitionStyleString = ""

func readDefinitionStyle(stylePath string) error {
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
	styleBytes, err := ioutil.ReadFile(stylePath)
	if err != nil {
		return err
	}
	definitionStyleString = "<style>" + string(styleBytes) + "</style>"
	return nil
}

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
	{
		err := readDefinitionStyle(conf.DefinitionStyle)
		if err != nil {
			fmt.Println(err)
		}
	}
}
