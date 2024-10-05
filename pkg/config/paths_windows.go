//go:build windows
// +build windows

package config

import (
	"log/slog"
	"os"
	"path/filepath"

	"github.com/ilius/ayandict/v2/pkg/appinfo"
)

func platformConfigDir() string {
	// HOMEDRIVE := os.Getenv("HOMEDRIVE")
	// HOMEPATH := os.Getenv("HOMEPATH")
	// homeDir := filepath.Join(HOMEDRIVE, HOMEPATH)
	// user := os.Getenv("USERNAME")
	// tmpDir := os.Getenv("TEMP")
	appData := os.Getenv("APPDATA")
	confDir := filepath.Join(appData, appinfo.APP_DESC)
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
			slog.Error("error", "err", err)
			return ""
		}
	}
	return filepath.Join(localAppData, appinfo.APP_DESC, "Cache")
}
