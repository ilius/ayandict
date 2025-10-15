package desktopwidget

import (
	qt "github.com/mappu/miqt/qt6"
)

func loadPNGPixmap(filename string) (*qt.QPixmap, error) {
	data, err := res.ReadFile("res/" + filename)
	if err != nil {
		return nil, err
	}
	pixmap := qt.NewQPixmap()
	pixmap.LoadFromDataWithData(data)
	return pixmap, nil
}
