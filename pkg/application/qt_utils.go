package application

import (
	"log/slog"
	"path/filepath"
	"strings"

	qt "github.com/mappu/miqt/qt6"
)

type KeyPressIface interface {
	OnKeyPressEvent(func(func(event *qt.QKeyEvent), *qt.QKeyEvent))
}

func filePathFromQUrl(qUrl *qt.QUrl) string {
	fpath := qUrl.Path()
	if fpath == "" {
		return ""
	}
	if filepath.Separator == '\\' {
		fpath = strings.TrimLeft(fpath, "/")
	}
	return fpath
}

func plaintextFromHTML(htext string) string {
	doc := qt.NewQTextDocument()
	doc.SetHtml(htext)
	return doc.ToPlainText()
}

func fontPointSize(font *qt.QFont, dpi float64) float64 {
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

func fontPixelSize(font *qt.QFont, dpi float64) float64 {
	pixels := font.PixelSize()
	if pixels > 0 {
		return float64(pixels)
	}

	points := font.PointSizeF()
	return points * dpi / 72.0
}

// func posStr(pos *qt.QPoint) string {
// 	if pos == nil {
// 		return "nil"
// 	}
// 	return fmt.Sprintf("(%v, %v)", pos.X(), pos.Y())
// }

// func sizeStr(size *qt.QSize) string {
// 	if size == nil {
// 		return "nil"
// 	}
// 	return fmt.Sprintf("(%v, %v)", size.Width(), size.Height())
// }
