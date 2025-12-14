package application

import (
	"github.com/ilius/ayandict/v3/pkg/resources"
	"github.com/ilius/ayandict/v3/pkg/resources/resourceutil"
	qt "github.com/mappu/miqt/qt6"
)

func loadPNGIcon(filename string) (*qt.QIcon, error) {
	return resourceutil.LoadPNGIcon(resources.Res, filename)
}

func loadPNGPixmap(filename string) (*qt.QPixmap, error) {
	return resourceutil.LoadPNGPixmap(resources.Res, filename)
}

// func loadSVGIcon(filename string) *qt.QIcon {
// 	data, err := res.ReadFile("res/" + filename)
// 	if err != nil {
// 		slog.Error("error", "err", err)
// 		return nil
// 	}
// 	image := qt.QImage_FromData(data, len(data), "SVG")
// 	image.Rect().SetSize(qt.NewQSize2(36, 36))
// 	pixmap := qt.QPixmap_FromImage(image, qt.AutoColor)
// 	return qt.NewQIcon2(pixmap)
// }
