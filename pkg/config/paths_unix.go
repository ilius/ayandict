//go:build !windows && !darwin
// +build !windows,!darwin

package config

import (
	"os"
	"path/filepath"
)

func GetConfigDir() string {
	if os.Getenv("CONFIG_FILE") != "" {
		return filepath.Dir(Path())
	}
	parent := os.Getenv("XDG_CONFIG_HOME")
	if parent == "" {
		parent = filepath.Join(os.Getenv("HOME"), ".config")
	}
	return filepath.Join(parent, "ayandict")
}

func GetCacheDir() string {
	return filepath.Join(os.Getenv("HOME"), ".cache", "ayandict")
}
