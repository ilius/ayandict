//go:build windows
// +build windows

package config

import (
	"log"
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

func GetCacheDir() string {
	localAppData := os.Getenv("LOCALAPPDATA")
	if localAppData == "" {
		// Windows Vista or older
		appData := os.Getenv("APPDATA")
		var err error
		localAppData, err = filepath.Abs(filepath.Join(appData, "..", "Local"))
		if err != nil {
			log.Println(err)
			return ""
		}
	}
	return filepath.Join(localAppData, "AyanDict", "Cache")
}
