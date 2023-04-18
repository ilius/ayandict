package application

import (
	"fmt"
	"path/filepath"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
)

func filePathFromQUrl(qUrl *core.QUrl) string {
	fpath := qUrl.Path(core.QUrl__FullyEncoded)
	if fpath == "" {
		return ""
	}
	if filepath.Separator == '\\' {
		fpath = fpath[1:]
	}
	return fpath
}

func plaintextFromHTML(htext string) string {
	doc := gui.NewQTextDocument(nil)
	doc.SetHtml(htext)
	return doc.ToPlainText()
}

func fontPointSize(font *gui.QFont, dpi float64) float64 {
	points := font.PointSizeF()
	if points > 0 {
		return points
	}
	pixels := font.PixelSize()
	return float64(pixels) * 72.0 / dpi
}

func posStr(pos *core.QPoint) string {
	if pos == nil {
		return "nil"
	}
	return fmt.Sprintf("(%v, %v)", pos.X(), pos.Y())
}

func sizeStr(size *core.QSize) string {
	if size == nil {
		return "nil"
	}
	return fmt.Sprintf("(%v, %v)", size.Width(), size.Height())
}
