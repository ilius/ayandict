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
	DirectoryList []string `toml:"directory_list"`

	Style string `toml:"style"`

	ArticleStyle string `toml:"article_style"`

	FontFamily string `toml:"font_family"`
	FontSize   int    `toml:"font_size"`

	SearchOnType          bool `toml:"search_on_type"`
	SearchOnTypeMinLength int  `toml:"search_on_type_min_length"`

	HeaderTemplate string `toml:"header_template"`

	HistoryDisable  bool `toml:"history_disable"`
	HistoryAutoSave bool `toml:"history_auto_save"`
	HistoryMaxSize  int  `toml:"history_max_size"`

	MostFrequentDisable  bool `toml:"most_frequent_disable"`
	MostFrequentAutoSave bool `toml:"most_frequent_auto_save"`
	MostFrequentMaxSize  int  `toml:"most_frequent_max_size"`

	MaxResultsPerDict int `toml:"max_results_per_dict"`
	MaxResultsTotal   int `toml:"max_results_total"`

	Audio bool `toml:"audio"`

	EmbedExternalStylesheet bool `toml:"embed_external_stylesheet"`

	ColorMapping map[string]string `toml:"color_mapping"`
}

func Default() *Config {
	return &Config{
		DirectoryList: []string{
			".stardict/dic",
		},

		Style: "",

		ArticleStyle: "",

		FontFamily: "",
		FontSize:   0,

		SearchOnType:          false,
		SearchOnTypeMinLength: 3,

		HeaderTemplate: `<h1>{{.Term}} (from {{.DictName}})</h1>`,

		HistoryDisable:  false,
		HistoryAutoSave: true,
		HistoryMaxSize:  100,

		MostFrequentDisable:  false,
		MostFrequentAutoSave: true,
		MostFrequentMaxSize:  100,

		MaxResultsPerDict: 20,
		MaxResultsTotal:   40,

		Audio: true,

		EmbedExternalStylesheet: false,

		ColorMapping: map[string]string{},
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
