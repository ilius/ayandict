package qtutils

import (
	"log/slog"

	qt "github.com/mappu/miqt/qt6"
)

func FontPixelSize(font *qt.QFont, screen *qt.QScreen) float64 {
	pixels := font.PixelSize()
	if pixels > 0 {
		return float64(pixels)
	}

	points := font.PointSizeF()
	dpi := screen.PhysicalDotsPerInch()
	if dpi <= 0 {
		panic("failed to get DPI")
	}
	return points * dpi / 72.0
}

func FontPointSize(font *qt.QFont, dpi float64) float64 {
	points := font.PointSizeF()
	if points > 0 {
		return points
	}
	pixels := font.PixelSize()
	if pixels <= 0 {
		slog.Error("bad font size", "points", font.PointSizeF(), "pixels", pixels)
	}
	return float64(pixels) * 72.0 / dpi
}
