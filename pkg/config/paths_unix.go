//go:build !windows && !darwin
// +build !windows,!darwin

package config

import (
	"os"
	"path/filepath"

	"github.com/ilius/ayandict/v2/pkg/appinfo"
)

func platformConfigDir() string {
	parent := os.Getenv("XDG_CONFIG_HOME")
	if parent == "" {
		parent = filepath.Join(os.Getenv(S_HOME), ".config")
	}
	return filepath.Join(parent, appinfo.APP_NAME)
}

func GetCacheDir() string {
	parent := os.Getenv("XDG_CACHE_HOME")
	if parent != "" {
		return filepath.Join(parent, appinfo.APP_NAME)
	}
	return filepath.Join(os.Getenv(S_HOME), ".cache", appinfo.APP_NAME)
}
