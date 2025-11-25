package qsettings

import (
	qt "github.com/mappu/miqt/qt6"
)

func SaveTableColumnsWidth(table *qt.QTableWidget, mainKey string) {
	count := table.ColumnCount()
	widthList := make([]int, count)
	for i := range count {
		widthList[i] = table.ColumnWidth(i)
	}
	saveJson(widthList, mainKey)
}

func RestoreTableColumnsWidth(table *qt.QTableWidget, mainKey string) {
	widthList := loadJsonIntSlice(mainKey)
	header := table.HorizontalHeader()
	for index, width := range widthList {
		header.ResizeSection(index, int(width))
	}
}
