//go:build windows
// +build windows

package config

import (
	"os"
	"path/filepath"
)

func GetConfigDir() string {
	// HOMEDRIVE := os.Getenv("HOMEDRIVE")
	// HOMEPATH := os.Getenv("HOMEPATH")
	// homeDir := filepath.Join(HOMEDRIVE, HOMEPATH)
	// user := os.Getenv("USERNAME")
	// tmpDir := os.Getenv("TEMP")
	appData := os.Getenv("APPDATA")
	confDir := filepath.Join(appData, "AyanDict")
	return confDir
}
