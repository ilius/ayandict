package qsettings

import (
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/ilius/ayandict/v2/pkg/appinfo"
	"github.com/ilius/ayandict/v2/pkg/qtcommon/qerr"
	"github.com/ilius/qt/core"
	"github.com/ilius/qt/gui"
	"github.com/ilius/qt/widgets"
)

const (
	QS_mainwindow = "mainwindow"
	QS_geometry   = "geometry"
	QS_savestate  = "savestate"
	QS_maximized  = "maximized"
	QS_pos        = "pos"
	QS_size       = "size"

	QS_columnwidth = "columnwidth"

	QS_sizes = "sizes"

	QS_search   = "search"
	QS_activity = "activity"
	QS_mode     = "mode"
)

func joinIntList(nums []int) string {
	strs := make([]string, len(nums))
	for i, num := range nums {
		strs[i] = strconv.FormatInt(int64(num), 10)
	}
	return strings.Join(strs, ",")
}

func splitIntList(st string) ([]int, error) {
	strs := strings.Split(st, ",")
	nums := make([]int, len(strs))
	for i, st := range strs {
		n, err := strconv.ParseInt(st, 10, 64)
		if err != nil {
			return nil, err
		}
		nums[i] = int(n)
	}
	return nums, nil
}

func splitterSizes(splitter *widgets.QSplitter) []int {
	itemCount := splitter.Count()
	widthList := make([]int, itemCount)
	for i := 0; i < itemCount; i++ {
		widthList[i] = splitter.Widget(i).Geometry().Width()
	}
	return widthList
}

func GetQSettings(parent core.QObject_ITF) *core.QSettings {
	return core.NewQSettings("ilius", appinfo.APP_NAME, parent)
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

func restoreIntSetting(
	qs *core.QSettings,
	key string,
	apply func(int),
) {
	if !qs.Contains(key) {
		return
	}
	value := qs.Value(key, core.NewQVariant1(nil))
	valueInt, err := strconv.ParseInt(value.ToString(), 10, 64)
	if err != nil {
		qerr.Error(err)
		return
	}
	apply(int(valueInt))
}

func saveMainWinGeometry(qs *core.QSettings, window *widgets.QMainWindow) {
	// slog.Info("Saving main window geometry")
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

func SaveWinGeometry(qs *core.QSettings, window *widgets.QWidget, mainKey string) {
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
	app *widgets.QApplication,
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
		pos.SetX(screenSize.Width() >> 1)
	}
	switch {
	case y < 0:
		pos.SetY(0)
	case y > screenSize.Height():
		pos.SetY(screenSize.Height() >> 1)
	}
	window.Move(pos)
}

func setWinSize(
	app *widgets.QApplication,
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

func RestoreMainWinGeometry(
	app *widgets.QApplication,
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

func RestoreWinGeometry(
	app *widgets.QApplication,
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

func SaveTableColumnsWidth(qs *core.QSettings, table *widgets.QTableWidget, mainKey string) {
	qs.BeginGroup(mainKey)
	defer qs.EndGroup()
	count := table.ColumnCount()
	widths := make([]int, count)
	for i := 0; i < count; i++ {
		widths[i] = table.ColumnWidth(i)
	}
	qs.SetValue(QS_columnwidth, core.NewQVariant1(joinIntList(widths)))
}

func RestoreTableColumnsWidth(qs *core.QSettings, table *widgets.QTableWidget, mainKey string) {
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
	// slog.Info("Saving splitter sizes")
	qs.BeginGroup(mainKey)
	defer qs.EndGroup()
	sizes := splitterSizes(splitter)
	qs.SetValue(QS_sizes, core.NewQVariant1(joinIntList(sizes)))
}

func RestoreSplitterSizes(qs *core.QSettings, splitter *widgets.QSplitter, mainKey string) {
	qs.BeginGroup(mainKey)
	defer qs.EndGroup()
	if !qs.Contains(QS_sizes) {
		return
	}
	sizesStr := qs.Value(QS_sizes, core.NewQVariant1("")).ToString()
	sizes, err := splitIntList(sizesStr)
	if err != nil {
		slog.Error("error", "err", err)
		return
	}
	splitter.SetSizes(sizes)
}

// QSplitter.Sizes() panics:
// interface conversion: interface {} is []interface {}, not []int

func actionSaveLoop(ch <-chan time.Time, callable func()) {
	var lastSave time.Time
	for {
		lastEvent := <-ch
	Loop1:
		for {
			select {
			case t := <-ch:
				lastEvent = t
			case <-time.After(500 * time.Millisecond):
				break Loop1
			}
		}
		if lastEvent.After(lastSave) {
			callable()
			lastSave = time.Now()
		}
	}
}

func SetupSplitterSizesSave(qs *core.QSettings, splitter *widgets.QSplitter, mainKey string) {
	ch := make(chan time.Time, 100)
	splitter.ConnectSplitterMoved(func(pos int, index int) {
		ch <- time.Now()
	})
	go actionSaveLoop(ch, func() {
		saveSplitterSizes(qs, splitter, mainKey)
	})
}

func SetupWinGeometrySave(
	qs *core.QSettings,
	window *widgets.QWidget,
	mainKey string,
) {
	ch := make(chan time.Time, 100)

	window.ConnectMoveEvent(func(event *gui.QMoveEvent) {
		ch <- time.Now()
	})
	window.ConnectResizeEvent(func(event *gui.QResizeEvent) {
		ch <- time.Now()
	})
	go actionSaveLoop(ch, func() {
		SaveWinGeometry(qs, window, mainKey)
	})
}

func SetupMainWinGeometrySave(qs *core.QSettings, window *widgets.QMainWindow) {
	ch := make(chan time.Time, 100)

	window.ConnectMoveEvent(func(event *gui.QMoveEvent) {
		ch <- time.Now()
	})
	window.ConnectResizeEvent(func(event *gui.QResizeEvent) {
		ch <- time.Now()
	})
	go actionSaveLoop(ch, func() {
		saveMainWinGeometry(qs, window)
	})
}

func SaveSearchSettings(qs *core.QSettings, combo *widgets.QComboBox) {
	qs.BeginGroup(QS_search)
	defer qs.EndGroup()

	qs.SetValue(QS_mode, core.NewQVariant1(
		strconv.FormatInt(int64(combo.CurrentIndex()), 10),
	))
}

func RestoreSearchSettings(qs *core.QSettings, combo *widgets.QComboBox) {
	qs.BeginGroup(QS_search)
	defer qs.EndGroup()

	restoreIntSetting(qs, QS_mode, combo.SetCurrentIndex)
}

func SaveActivityMode(qs *core.QSettings, combo *widgets.QComboBox) {
	qs.BeginGroup(QS_activity)
	defer qs.EndGroup()

	qs.SetValue(QS_mode, core.NewQVariant1(
		strconv.FormatInt(int64(combo.CurrentIndex()), 10),
	))
}

func RestoreActivityMode(qs *core.QSettings, combo *widgets.QComboBox) {
	qs.BeginGroup(QS_activity)
	defer qs.EndGroup()

	restoreIntSetting(qs, QS_mode, combo.SetCurrentIndex)
}
