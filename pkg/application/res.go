package application

import (
	"log/slog"

	qt "github.com/mappu/miqt/qt6"
)

func loadPNGIcon(filename string) (*qt.QIcon, error) {
	data, err := res.ReadFile("res/" + filename)
	if err != nil {
		return nil, err
	}
	pixmap := qt.NewQPixmap()
	pixmap.LoadFromDataWithData(data)
	icon := qt.NewQIcon2(pixmap)
	if icon == nil {
		slog.Error("error loading png icon: icon is nil: " + filename)
		panic("error loading png icon: icon is nil")
	}
	return icon, nil
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
