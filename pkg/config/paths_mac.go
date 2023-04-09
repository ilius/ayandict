//go:build darwin
// +build darwin

package config

import (
	"os"
	"path/filepath"
)

func GetConfigDir() string {
	return filepath.Join(
		os.Getenv("HOME"),
		"Library/Preferences/AyanDict",
	)
}
