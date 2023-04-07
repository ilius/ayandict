package favorites

import (
	"path/filepath"

	"github.com/ilius/ayandict/pkg/config"
)

const favoritesFilename = "favorites.json"

func Path() string {
	return filepath.Join(config.GetConfigDir(), favoritesFilename)
}
