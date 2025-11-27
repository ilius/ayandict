//go:build windows
// +build windows

package ospaths

import (
	"log/slog"
	"os"
	"path/filepath"
)

func (p *Paths) ConfigDir() string {
	// HOMEDRIVE := os.Getenv("HOMEDRIVE")
	// HOMEPATH := os.Getenv("HOMEPATH")
	// homeDir := filepath.Join(HOMEDRIVE, HOMEPATH)
	// user := os.Getenv("USERNAME")
	// tmpDir := os.Getenv("TEMP")
	appData := os.Getenv("APPDATA")
	confDir := filepath.Join(appData, p.AppDesc)
	return confDir
}

func (p *Paths) CacheDir() string {
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
	return filepath.Join(localAppData, p.AppDesc, "Cache")
}

func (p *Paths) StateDir() string {
	return filepath.Join(p.ConfigDir(), "State")
}
