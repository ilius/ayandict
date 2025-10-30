//go:build darwin
// +build darwin

package config

import (
	"os"
	"path/filepath"

	"github.com/ilius/ayandict/v3/pkg/appinfo"
)

func platformConfigDir() string {
	return filepath.Join(
		os.Getenv(S_HOME),
		"Library",
		"Preferences",
		appinfo.APP_DESC,
	)
}

func GetCacheDir() string {
	return filepath.Join(os.Getenv(S_HOME), "Library", "Caches", appinfo.APP_DESC)
}
