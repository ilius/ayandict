//go:build !windows
// +build !windows

package config

import (
	"os"
	"path/filepath"
)

func GetConfigDir() string {
	parent := os.Getenv("XDG_CONFIG_HOME")
	if parent == "" {
		parent = filepath.Join(os.Getenv("HOME"), ".config")
	}
	return filepath.Join(parent, "ayandict")
}
