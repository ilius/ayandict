package application

import (
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

func splitterSizes(splitter *widgets.QSplitter) []int {
	itemCount := splitter.Count()
	widthList := make([]int, itemCount)
	for i := 0; i < itemCount; i++ {
		widthList[i] = splitter.Widget(i).Geometry().Width()
	}
	return widthList
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
