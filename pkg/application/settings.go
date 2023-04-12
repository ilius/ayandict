package application

import (
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ilius/ayandict/pkg/qerr"
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

	QS_mainSplitter = "main_splitter"
	QS_sizes        = "sizes"

	QS_dictManager = "dict_manager"
)

func getQSettings(parent core.QObject_ITF) *core.QSettings {
	return core.NewQSettings("ilius", "ayandict", parent)
}

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

func saveWinGeometry(qs *core.QSettings, window *widgets.QWidget, mainKey string) {
	qs.BeginGroup(mainKey)
	defer qs.EndGroup()

	qs.SetValue(QS_geometry, core.NewQVariant13(window.SaveGeometry()))
	qs.SetValue(QS_maximized, core.NewQVariant9(window.IsMaximized()))
	if !window.IsMaximized() {
		qs.SetValue(QS_pos, core.NewQVariant27(window.Pos()))
		qs.SetValue(QS_size, core.NewQVariant25(window.Size()))
	}
}

func setWinPosition(
	app *Application,
	window *widgets.QWidget,
	pos *core.QPoint,
) {
	screenSize := app.Desktop().AvailableGeometry(0)
	x := pos.X()
	y := pos.Y()
	switch {
	case x < 0:
		pos.SetX(0)
	case x > screenSize.Width():
		pos.SetX(screenSize.Width() / 2)
	}
	switch {
	case y < 0:
		pos.SetY(0)
	case y > screenSize.Height():
		pos.SetY(screenSize.Height() / 2)
	}
	window.Move(pos)
}

func setWinSize(
	app *Application,
	window *widgets.QWidget,
	size *core.QSize,
) {
	screenSize := app.Desktop().AvailableGeometry(0)
	if size.Width() > screenSize.Width() {
		size.SetWidth(screenSize.Width())
	}
	if size.Height() > screenSize.Height() {
		size.SetHeight(screenSize.Height())
	}
	window.Resize(size)
}

func restoreMainWinGeometry(
	app *Application,
	qs *core.QSettings,
	window *widgets.QMainWindow,
) {
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
			setWinPosition(app, &window.QWidget, value.ToPoint())
		})
		restoreSetting(qs, QS_size, func(value *core.QVariant) {
			setWinSize(app, &window.QWidget, value.ToSize())
		})
	})
}

func restoreWinGeometry(
	app *Application,
	qs *core.QSettings,
	window *widgets.QWidget,
	mainKey string,
) {
	qs.BeginGroup(mainKey)
	defer qs.EndGroup()

	restoreSetting(qs, QS_geometry, func(value *core.QVariant) {
		window.RestoreGeometry(value.ToByteArray())
	})
	restoreBoolSetting(qs, QS_maximized, false, func(maximized bool) {
		if maximized {
			window.ShowMaximized()
			return
		}
		restoreSetting(qs, QS_pos, func(value *core.QVariant) {
			setWinPosition(app, window, value.ToPoint())
		})
		restoreSetting(qs, QS_size, func(value *core.QVariant) {
			setWinSize(app, window, value.ToSize())
		})
	})
}

func saveTableColumnsWidth(qs *core.QSettings, table *widgets.QTableWidget, mainKey string) {
	qs.BeginGroup(mainKey)
	defer qs.EndGroup()
	count := table.ColumnCount()
	widths := make([]int, count)
	for i := 0; i < count; i++ {
		widths[i] = table.ColumnWidth(i)
	}
	qs.SetValue(QS_columnwidth, core.NewQVariant1(joinIntList(widths)))
}

func restoreTableColumnsWidth(qs *core.QSettings, table *widgets.QTableWidget, mainKey string) {
	qs.BeginGroup(mainKey)
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
			qerr.Errorf("invalid column width=%#v\n", widthStr)
			continue
		}
		header.ResizeSection(index, int(width))
	}
}

func saveSplitterSizes(qs *core.QSettings, splitter *widgets.QSplitter, mainKey string) {
	qs.BeginGroup(mainKey)
	defer qs.EndGroup()
	sizes := splitterSizes(splitter)
	qs.SetValue(QS_sizes, core.NewQVariant1(joinIntList(sizes)))
}

func restoreSplitterSizes(qs *core.QSettings, splitter *widgets.QSplitter, mainKey string) {
	qs.BeginGroup(mainKey)
	defer qs.EndGroup()
	if !qs.Contains(QS_sizes) {
		return
	}
	sizesStr := qs.Value(QS_sizes, core.NewQVariant1("")).ToString()
	sizes, err := splitIntList(sizesStr)
	if err != nil {
		log.Println(err)
		return
	}
	splitter.SetSizes(sizes)
}

// QSplitter.Sizes() panics:
// interface conversion: interface {} is []interface {}, not []int

func setupSplitterSizesSave(qs *core.QSettings, splitter *widgets.QSplitter, mainKey string) {
	var mutex sync.Mutex
	var lastSave time.Time
	savedPos := make(map[int]int, 3)

	onMove := func(pos int, index int) {
		// log.Printf("---- splitter: moved: index=%v, pos=%v", index, pos)
		eventTime := time.Now()
		if !tryLockAsManyAs(&mutex, 5, 500*time.Millisecond) {
			return
		}
		defer mutex.Unlock()
		if eventTime.Before(lastSave) {
			return
		}
		if savedPos[index] == pos {
			return
		}
		// log.Printf("---- splitter: saving sizes: index=%v, pos=%v", index, pos)
		saveSplitterSizes(qs, splitter, mainKey)
		lastSave = time.Now()
		savedPos[index] = pos
		time.Sleep(500 * time.Millisecond)
	}

	splitter.ConnectSplitterMoved(func(pos int, index int) {
		go onMove(pos, index)
	})
}
