package application

import (
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
