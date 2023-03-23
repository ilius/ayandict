package config

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/BurntSushi/toml"
)

const fileName = "config.toml"

var mutex sync.Mutex

// if you set FontSize, you can not change font size of
// html definition view using mouse scroll

type Config struct {
	Style string `toml:"style"`

	FontFamily string `toml:"font_family"`
	FontSize   int    `toml:"font_size"`

	SearchOnType          bool `toml:"search_on_type"`
	SearchOnTypeMinLength int  `toml:"search_on_type_min_length"`

	DictHeaderTag string `toml:"dict_header_tag"`
	WordHeaderTag string `toml:"word_header_tag"`

	HistoryDisable  bool `toml:"history_disable"`
	HistoryAutoSave bool `toml:"history_auto_save"`
	HistoryMaxSize  int  `toml:"history_max_size"`

	Audio bool `toml:"audio"`
}

func Default() *Config {
	return &Config{
		Style:      "",
		FontFamily: "",
		FontSize:   0,

		SearchOnType:          false,
		SearchOnTypeMinLength: 3,

		DictHeaderTag: "h1",
		WordHeaderTag: "h2",

		HistoryDisable:  false,
		HistoryAutoSave: false,
		HistoryMaxSize:  100,

		Audio: true,
	}
}

func Path() string {
	return filepath.Join(GetConfigDir(), fileName)
}

func loadFile() ([]byte, error) {
	mutex.Lock()
	defer mutex.Unlock()

	pathStr := filepath.Join(GetConfigDir(), fileName)
	tomlBytes, err := ioutil.ReadFile(pathStr)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	return tomlBytes, nil
}

func Load() (*Config, error) {
	tomlBytes, err := loadFile()
	if err != nil {
		return nil, err
	}
	if tomlBytes == nil {
		return Default(), nil
	}
	conf := Default()
	_, err = toml.Decode(string(tomlBytes), conf)
	if err != nil {
		return nil, err
	}
	return conf, nil
}

func MustLoad() *Config {
	conf, err := Load()
	if err != nil {
		panic(err)
	}
	return conf
}

func EnsureExists(conf *Config) error {
	_, err := os.Stat(filepath.Join(GetConfigDir(), fileName))
	if err == nil {
		return nil
	}
	if os.IsNotExist(err) {
		return Save(conf)
	}
	return err
}

func Save(conf *Config) error {
	configDir := GetConfigDir()
	err := os.MkdirAll(configDir, 0o755)
	if err != nil {
		return err
	}
	pathStr := filepath.Join(configDir, fileName)
	buf := bytes.NewBuffer(nil)
	encoder := toml.NewEncoder(buf)
	err = encoder.Encode(conf)
	if err != nil {
		return err
	}

	mutex.Lock()
	defer mutex.Unlock()

	err = ioutil.WriteFile(pathStr, buf.Bytes(), 0o644)
	if err != nil {
		return err
	}
	if runtime.GOOS == "windows" {
		time.Sleep(100 * time.Millisecond)
	}
	return nil
}
