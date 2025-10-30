package config

import (
	"bytes"
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"time"

	"github.com/BurntSushi/toml"
)

func Save(conf *Config) error {
	configDir := GetConfigDir()
	err := os.MkdirAll(configDir, 0o755)
	if err != nil {
		return err
	}
	pathStr := Path()
	buf := bytes.NewBuffer(nil)
	encoder := toml.NewEncoder(buf)
	err = encoder.Encode(conf)
	if err != nil {
		return err
	}

	mutex.Lock()
	defer mutex.Unlock()

	err = os.WriteFile(pathStr, buf.Bytes(), 0o644)
	if err != nil {
		return err
	}
	if runtime.GOOS == "windows" {
		time.Sleep(100 * time.Millisecond)
	}
	return nil
}

func EnsureExists(conf *Config) error {
	stat, err := os.Stat(Path())
	if err == nil {
		if stat.IsDir() {
			return fmt.Errorf("config file is a directory")
		}
		return nil
	}
	if !os.IsNotExist(err) {
		return err
	}
	err = Save(conf)
	if err != nil {
		return err
	}
	slog.Info("Created config file: " + Path())
	return nil
}
