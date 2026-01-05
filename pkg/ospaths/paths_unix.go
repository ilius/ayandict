//go:build !windows && !darwin

package ospaths

import (
	"os"
	"path/filepath"
)

func (p *Paths) ConfigDir() string {
	parent := os.Getenv("XDG_CONFIG_HOME")
	if parent == "" {
		parent = filepath.Join(p.Home, ".config")
	}
	return filepath.Join(parent, p.AppName)
}

func (p *Paths) CacheDir() string {
	parent := os.Getenv("XDG_CACHE_HOME")
	if parent != "" {
		return filepath.Join(parent, p.AppName)
	}
	return filepath.Join(p.Home, ".cache", p.AppName)
}

func (p *Paths) StateDir() string {
	// $XDG_STATE_HOME contains state data that should persist between (app)
	// restarts, but that is not important or portable enough to the user that
	// it should be stored in $XDG_DATA_HOME
	// If either not set or empty, `$HOME/.local/state` should be used.
	parent := os.Getenv("XDG_STATE_HOME")
	if parent != "" {
		return filepath.Join(parent, p.AppName)
	}
	return filepath.Join(
		p.Home,
		".local",
		"state",
		p.AppName,
	)
}
