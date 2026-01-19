package qtutils

import (
	"log/slog"

	qt "github.com/mappu/miqt/qt6"
)

func SetWinPosition(window *qt.QWidget, x int, y int) {
	screenSize := qt.QGuiApplication_PrimaryScreen().AvailableGeometry()
	switch {
	case x < 0:
		x = 0
	case x > screenSize.Width():
		slog.Warn("SetWinPosition: exceeds screen", "x", x, "width", screenSize.Width())
		x = screenSize.Width() - 100
	}
	switch {
	case y < 0:
		y = 0
	case y > screenSize.Height():
		slog.Warn("SetWinPosition: exceeds screen", "y", y, "height", screenSize.Height())
		y = screenSize.Height() - 100
	}
	window.Move(x, y)
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
