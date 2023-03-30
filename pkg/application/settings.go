package application

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

const (
	QS_mainwindow = "mainwindow"
	QS_geometry   = "geometry"
	QS_savestate  = "savestate"
	QS_maximized  = "maximized"
	QS_pos        = "pos"
	QS_size       = "size"

	QS_frequencyTable = "frequencytable"
	QS_columnwidth    = "columnwidth"
)

func restoreSetting(qs *core.QSettings, key string, apply func(*core.QVariant)) {
	if !qs.Contains(key) {
		return
	}
	apply(qs.Value(key, core.NewQVariant1(nil)))
}

func restoreBoolSetting(
	qs *core.QSettings,
	key string, _default bool,
	apply func(bool),
) {
	if !qs.Contains(key) {
		apply(_default)
		return
	}
	apply(qs.Value(key, core.NewQVariant1(nil)).ToBool())
}

func saveMainWinGeometry(qs *core.QSettings, window *widgets.QMainWindow) {
	qs.BeginGroup(QS_mainwindow)
	defer qs.EndGroup()

	qs.SetValue(QS_geometry, core.NewQVariant13(window.SaveGeometry()))
	qs.SetValue(QS_savestate, core.NewQVariant13(window.SaveState(0)))
	qs.SetValue(QS_maximized, core.NewQVariant9(window.IsMaximized()))
	if !window.IsMaximized() {
		qs.SetValue(QS_pos, core.NewQVariant27(window.Pos()))
		qs.SetValue(QS_size, core.NewQVariant25(window.Size()))
	}
}

func setWinPosition(window *widgets.QMainWindow, pos *core.QPoint) {
	x := pos.X()
	y := pos.Y()
	// TODO: check with screen size
	switch {
	case x < 0:
		pos.SetX(0)
	case x > 1000:
		pos.SetX(1000)
	}
	switch {
	case y < 0:
		pos.SetY(0)
	case y > 1000:
		pos.SetY(1000)
	}
	window.Move(pos)
}

func restoreMainWinGeometry(qs *core.QSettings, window *widgets.QMainWindow) {
	qs.BeginGroup(QS_mainwindow)
	defer qs.EndGroup()

	restoreSetting(qs, QS_geometry, func(value *core.QVariant) {
		window.RestoreGeometry(value.ToByteArray())
	})
	restoreSetting(qs, QS_savestate, func(value *core.QVariant) {
		window.RestoreState(value.ToByteArray(), 0)
	})
	restoreBoolSetting(qs, QS_maximized, false, func(maximized bool) {
		if maximized {
			window.ShowMaximized()
			return
		}
		restoreSetting(qs, QS_pos, func(value *core.QVariant) {
			setWinPosition(window, value.ToPoint())
		})
		restoreSetting(qs, QS_size, func(value *core.QVariant) {
			window.Resize(value.ToSize())
		})
	})
}

func saveTableColumnsWidth(qs *core.QSettings, table *widgets.QTableWidget, tableKey string) {
	qs.BeginGroup(tableKey)
	defer qs.EndGroup()
	count := table.ColumnCount()
	widths := make([]string, count)
	for i := 0; i < count; i++ {
		widths[i] = strconv.FormatInt(int64(table.ColumnWidth(i)), 10)
	}
	qs.SetValue(QS_columnwidth, core.NewQVariant1(strings.Join(widths, ",")))
}

func restoreTableColumnsWidth(qs *core.QSettings, table *widgets.QTableWidget, tableKey string) {
	qs.BeginGroup(tableKey)
	defer qs.EndGroup()
	if !qs.Contains(QS_columnwidth) {
		return
	}
	header := table.HorizontalHeader()
	// even []string does not work, let alone []int
	widthListStr := qs.Value(QS_columnwidth, core.NewQVariant1("")).ToString()
	widthList := strings.Split(widthListStr, ",")
	for index, widthStr := range widthList {
		width, err := strconv.ParseInt(widthStr, 10, 64)
		if err != nil {
			fmt.Printf("invalid column width=%#v\n", widthStr)
			continue
		}
		header.ResizeSection(index, int(width))
	}
}
