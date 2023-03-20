package config

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

const fileName = "config.toml"

// if you set FontSize, you can not change font size of
// html definition view using mouse scroll

type Config struct {
	FontFamily   string `toml:"font_family"`
	FontSize     int    `toml:"font_size"`
	SearchOnType bool   `toml:"search_on_type"`
}

func Load() (*Config, error) {
	pathStr := filepath.Join(GetConfigDir(), fileName)
	tomlBytes, err := ioutil.ReadFile(pathStr)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{}, nil
		}
		return nil, err
	}
	conf := &Config{}
	_, err = toml.Decode(string(tomlBytes), &conf)
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
	err = ioutil.WriteFile(pathStr, buf.Bytes(), 0o644)
	if err != nil {
		return err
	}
	return nil
}
