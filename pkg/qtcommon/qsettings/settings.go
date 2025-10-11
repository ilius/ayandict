package qsettings

import (
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/ilius/ayandict/v3/pkg/qtutils"
	qt "github.com/mappu/miqt/qt6"
)

var qs_columnwidth = *qt.NewQAnyStringView3(QS_columnwidth)

func SaveDialogGeometry(dialog *qt.QDialog, mainKey string) {
	// slog.Info("Saving main window geometry")
	pos := dialog.Pos()
	size := dialog.Size()
	// what is window.SaveState()
	s := &WindowSettings{
		X:      pos.X(),
		Y:      pos.Y(),
		Width:  size.Width(),
		Height: size.Height(),
	}
	s.Save(mainKey)
}

func RestoreWindowGeometry(window *qt.QWidget, mainKey string) {
	s := &WindowSettings{}
	s.Load(mainKey)
	qtutils.SetWinPosition(window, qt.NewQPoint2(s.X, s.Y))
	qtutils.SetWinSize(window, qt.NewQSize2(s.Width, s.Height))
	if s.Maximized {
		window.ShowMaximized()
	}
}

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

func saveSplitterSizes(splitter *qt.QSplitter, mainKey string) {
	// slog.Info("Saving splitter sizes")
	saveJson(splitter.Sizes(), mainKey)
}

func RestoreSplitterSizes(splitter *qt.QSplitter, mainKey string) {
	sizes := loadJsonIntSlice(mainKey)
	if sizes == nil {
		return
	}
	splitter.SetSizes(sizes)
}

func SetupSplitterSizesSave(splitter *qt.QSplitter, mainKey string) {
	ch := make(chan time.Time, 100)
	splitter.OnSplitterMoved(func(pos int, index int) {
		ch <- time.Now()
	})
	go ActionSaveLoop(ch, func() {
		saveSplitterSizes(splitter, mainKey)
	})
}

func SetupDialogGeometrySave(
	dialog *qt.QDialog,
	mainKey string,
) {
	ch := make(chan time.Time, 100)

	dialog.OnMoveEvent(func(super func(*qt.QMoveEvent), event *qt.QMoveEvent) {
		ch <- time.Now()
	})
	dialog.OnResizeEvent(func(super func(*qt.QResizeEvent), event *qt.QResizeEvent) {
		ch <- time.Now()
	})
	go ActionSaveLoop(ch, func() {
		SaveDialogGeometry(dialog, mainKey)
	})
}
