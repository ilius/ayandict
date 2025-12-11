package favoritebutton

import qt "github.com/mappu/miqt/qt6"

func hsvString(h int, s int, v int) string {
	color := qt.NewQColor()
	color.SetHsv(h, s, v)
	return color.ToQVariant().ToString()
}
