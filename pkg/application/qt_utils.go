package application

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/ilius/qt/core"
	"github.com/ilius/qt/gui"
)

type KeyPressIface interface {
	ConnectKeyPressEvent(func(event *gui.QKeyEvent))
	KeyPressEventDefault(event gui.QKeyEvent_ITF)
}

func filePathFromQUrl(qUrl *core.QUrl) string {
	fpath := qUrl.Path(core.QUrl__FullyDecoded)
	if fpath == "" {
		return ""
	}
	if filepath.Separator == '\\' {
		fpath = strings.TrimLeft(fpath, "/")
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
