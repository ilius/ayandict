package qtutils

import qt "github.com/mappu/miqt/qt6"

func SetWinPosition(window *qt.QWidget, pos *qt.QPoint) {
	screenSize := qt.QGuiApplication_PrimaryScreen().AvailableGeometry()
	x := pos.X()
	y := pos.Y()
	switch {
	case x < 0:
		pos.SetX(0)
	case x > screenSize.Width():
		pos.SetX(screenSize.Width() >> 1)
	}
	switch {
	case y < 0:
		pos.SetY(0)
	case y > screenSize.Height():
		pos.SetY(screenSize.Height() >> 1)
	}
	window.MoveWithQPoint(pos)
}

func SetWinSize(window *qt.QWidget, size *qt.QSize) {
	screenSize := qt.QGuiApplication_PrimaryScreen().AvailableGeometry()
	if size.Width() > screenSize.Width() {
		size.SetWidth(screenSize.Width())
	}
	if size.Height() > screenSize.Height() {
		size.SetHeight(screenSize.Height())
	}
	window.ResizeWithQSize(size)
}
