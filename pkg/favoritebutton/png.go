package favoritebutton

import (
	"github.com/ilius/ayandict/v3/pkg/resources"
	"github.com/ilius/ayandict/v3/pkg/resources/resourceutil"
	qt "github.com/mappu/miqt/qt6"
)

func loadPNGIcon(filename string) (*qt.QIcon, error) {
	return resourceutil.LoadPNGIcon(resources.Res, filename)
}
