//go:build !windows && !darwin
// +build !windows,!darwin

package config

import (
	"os"
	"path/filepath"

	"github.com/ilius/ayandict/v3/pkg/appinfo"
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

func GetStateDir() string {
	// $XDG_STATE_HOME contains state data that should persist between (app)
	// restarts, but that is not important or portable enough to the user that
	// it should be stored in $XDG_DATA_HOME
	// If either not set or empty, `$HOME/.local/state` should be used.
	parent := os.Getenv("XDG_STATE_HOME")
	if parent != "" {
		return filepath.Join(parent, appinfo.APP_NAME)
	}
	return filepath.Join(
		os.Getenv(S_HOME),
		".local",
		"state",
		appinfo.APP_NAME,
	)
}
