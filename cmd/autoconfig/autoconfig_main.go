package main

import (
	"bytes"
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/mappu/autoconfig"
	qt "github.com/mappu/miqt/qt6"
)

type LoggingConfig struct {
	NoColor bool   `toml:"no_color" doc:"Disable log colors" ylabel:"Disable colors"`
	Level   string `toml:"level" doc:"Log level" ylabel:"Level"`
}

type MiscButtonsConfig struct {
	SaveHistory   bool `toml:"save_history" ylabel:"Save History"`
	ClearHistory  bool `toml:"clear_history"  ylabel:"Clear History"`
	SaveFavorites bool `toml:"save_favorites"  ylabel:"Save Favorites"`
}

type TestConfig struct {
	Logging LoggingConfig `toml:"logging" ylabel:"Logging"`

	MiscButtons MiscButtonsConfig `toml:"misc_buttons" ylabel:"Misc"`

	DirectoryList []string `toml:"directory_list" ylabel:"Directories"`
}

func main() {
	qt.NewQApplication(os.Args)
	conf := &TestConfig{}
	autoconfig.OpenDialog(
		conf,
		nil,
		"Config",
		func() {},
	)
	qt.QApplication_Exec()
	buf := bytes.NewBuffer(nil)
	encoder := toml.NewEncoder(buf)
	err := encoder.Encode(conf)
	if err != nil {
		panic(err)
	}
	fmt.Println(buf.String())
}
