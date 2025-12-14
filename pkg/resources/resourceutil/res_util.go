package resourceutil

import (
	"embed"
	"log/slog"

	qt "github.com/mappu/miqt/qt6"
)

func LoadPNGIcon(res embed.FS, filename string) (*qt.QIcon, error) {
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

func LoadPNGPixmap(res embed.FS, filename string) (*qt.QPixmap, error) {
	data, err := res.ReadFile("res/" + filename)
	if err != nil {
		return nil, err
	}
	pixmap := qt.NewQPixmap()
	pixmap.LoadFromDataWithData(data)
	return pixmap, nil
}
