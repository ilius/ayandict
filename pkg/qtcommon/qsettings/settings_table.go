package qsettings

import (
	"log/slog"
	"strconv"
	"strings"

	qt "github.com/mappu/miqt/qt6"
)

func SaveTableColumnsWidth(table *qt.QTableWidget, mainKey string) {
	count := table.ColumnCount()
	widths := make([]int, count)
	for i := range count {
		widths[i] = table.ColumnWidth(i)
	}
	saveJson(widths, mainKey)
}

func RestoreTableColumnsWidth(qs *qt.QSettings, table *qt.QTableWidget, mainKey string) {
	qs.BeginGroup(*qt.NewQAnyStringView3(mainKey))
	defer qs.EndGroup()
	if !qs.Contains(qs_columnwidth) {
		return
	}
	header := table.HorizontalHeader()
	// even []string does not work, let alone []int
	widthListStr := qs.ValueWithKey(qs_columnwidth).ToString()
	widthList := strings.Split(widthListStr, ",")
	for index, widthStr := range widthList {
		width, err := strconv.ParseInt(widthStr, 10, 64)
		if err != nil {
			slog.Error("invalid column width=" + widthStr)
			continue
		}
		header.ResizeSection(index, int(width))
	}
}
