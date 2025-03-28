package qsettings

import (
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/ilius/ayandict/v3/pkg/appinfo"
	qt "github.com/mappu/miqt/qt6"
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

var (
	qs_mainwindow = *qt.NewQAnyStringView3(QS_mainwindow)
	qs_geometry   = *qt.NewQAnyStringView3(QS_geometry)
	qs_savestate  = *qt.NewQAnyStringView3(QS_savestate)
	qs_maximized  = *qt.NewQAnyStringView3(QS_maximized)
	qs_pos        = *qt.NewQAnyStringView3(QS_pos)
	qs_size       = *qt.NewQAnyStringView3(QS_size)

	qs_columnwidth = *qt.NewQAnyStringView3(QS_columnwidth)

	qs_sizes = *qt.NewQAnyStringView3(QS_sizes)

	qs_search   = *qt.NewQAnyStringView3(QS_search)
	qs_activity = *qt.NewQAnyStringView3(QS_activity)
	qs_mode     = *qt.NewQAnyStringView3(QS_mode)
)

var (
	newQVarBool  = qt.NewQVariant8
	newQVarBytes = qt.NewQVariant12
	newQVarStr   = qt.NewQVariant14
	newQVarSize  = qt.NewQVariant22
	newQVarPoint = qt.NewQVariant24
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

func splitterSizes(splitter *qt.QSplitter) []int {
	itemCount := splitter.Count()
	widthList := make([]int, itemCount)
	for i := range itemCount {
		widthList[i] = splitter.Widget(i).Geometry().Width()
	}
	return widthList
}

func GetQSettings(parent *qt.QObject) *qt.QSettings {
	return qt.NewQSettings8("ilius", appinfo.APP_NAME, parent)
}

func restoreSetting(qs *qt.QSettings, qKey qt.QAnyStringView, apply func(*qt.QVariant)) {
	if !qs.Contains(qKey) {
		return
	}
	apply(qs.ValueWithKey(qKey))
}

func restoreBoolSetting(
	qs *qt.QSettings,
	qKey qt.QAnyStringView,
	_default bool,
	apply func(bool),
) {
	if !qs.Contains(qKey) {
		apply(_default)
		return
	}
	apply(qs.ValueWithKey(qKey).ToBool())
}

func restoreIntSetting(
	qs *qt.QSettings,
	qKey qt.QAnyStringView,
	apply func(int),
) {
	if !qs.Contains(qKey) {
		return
	}
	value := qs.ValueWithKey(qKey)
	valueInt, err := strconv.ParseInt(value.ToString(), 10, 64)
	if err != nil {
		slog.Error("error in restoreIntSetting: bad int value: "+err.Error(), "value", value.ToString())
		return
	}
	apply(int(valueInt))
}

func saveMainWinGeometry(qs *qt.QSettings, window *qt.QMainWindow) {
	// slog.Info("Saving main window geometry")
	qs.BeginGroup(qs_mainwindow)
	defer qs.EndGroup()

	qs.SetValue(qs_geometry, newQVarBytes(window.SaveGeometry()))
	qs.SetValue(qs_savestate, newQVarBytes(window.SaveState()))
	qs.SetValue(qs_maximized, newQVarBool(window.IsMaximized()))
	if !window.IsMaximized() {
		qs.SetValue(qs_pos, newQVarPoint(window.Pos()))
		qs.SetValue(qs_size, newQVarSize(window.Size()))
	}
}

func SaveDialogGeometry(qs *qt.QSettings, dialog *qt.QDialog, mainKey string) {
	qs.BeginGroup(*qt.NewQAnyStringView3(mainKey))
	defer qs.EndGroup()

	qs.SetValue(qs_geometry, newQVarBytes(dialog.SaveGeometry()))
	qs.SetValue(qs_maximized, newQVarBool(dialog.IsMaximized()))
	if !dialog.IsMaximized() {
		qs.SetValue(qs_pos, newQVarPoint(dialog.Pos()))
		qs.SetValue(qs_size, newQVarSize(dialog.Size()))
	}
}

func setWinPosition(
	window *qt.QWidget,
	pos *qt.QPoint,
) {
	screenSize := qt.QGuiApplication_PrimaryScreen().AvailableGeometry()
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
	window.MoveWithQPoint(pos)
}

func setWinSize(
	window *qt.QWidget,
	size *qt.QSize,
) {
	screenSize := qt.QGuiApplication_PrimaryScreen().AvailableGeometry()
	if size.Width() > screenSize.Width() {
		size.SetWidth(screenSize.Width())
	}
	if size.Height() > screenSize.Height() {
		size.SetHeight(screenSize.Height())
	}
	window.ResizeWithQSize(size)
}

func RestoreMainWinGeometry(
	qs *qt.QSettings,
	window *qt.QMainWindow,
) {
	qs.BeginGroup(qs_mainwindow)
	defer qs.EndGroup()

	restoreSetting(qs, qs_geometry, func(value *qt.QVariant) {
		window.RestoreGeometry(value.ToByteArray())
	})
	restoreSetting(qs, qs_savestate, func(value *qt.QVariant) {
		window.RestoreState(value.ToByteArray())
	})
	restoreBoolSetting(qs, qs_maximized, false, func(maximized bool) {
		if maximized {
			window.ShowMaximized()
			return
		}
		restoreSetting(qs, qs_pos, func(value *qt.QVariant) {
			setWinPosition(window.QWidget, value.ToPoint())
		})
		restoreSetting(qs, qs_size, func(value *qt.QVariant) {
			setWinSize(window.QWidget, value.ToSize())
		})
	})
}

func RestoreWinGeometry(
	qs *qt.QSettings,
	window *qt.QWidget,
	mainKey string,
) {
	qs.BeginGroup(*qt.NewQAnyStringView3(mainKey))
	defer qs.EndGroup()

	restoreSetting(qs, qs_geometry, func(value *qt.QVariant) {
		window.RestoreGeometry(value.ToByteArray())
	})
	restoreBoolSetting(qs, qs_maximized, false, func(maximized bool) {
		if maximized {
			window.ShowMaximized()
			return
		}
		restoreSetting(qs, qs_pos, func(value *qt.QVariant) {
			setWinPosition(window, value.ToPoint())
		})
		restoreSetting(qs, qs_size, func(value *qt.QVariant) {
			setWinSize(window, value.ToSize())
		})
	})
}

func SaveTableColumnsWidth(qs *qt.QSettings, table *qt.QTableWidget, mainKey string) {
	qs.BeginGroup(*qt.NewQAnyStringView3(mainKey))
	defer qs.EndGroup()
	count := table.ColumnCount()
	widths := make([]int, count)
	for i := range count {
		widths[i] = table.ColumnWidth(i)
	}
	qs.SetValue(qs_columnwidth, newQVarStr(joinIntList(widths)))
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

func saveSplitterSizes(qs *qt.QSettings, splitter *qt.QSplitter, mainKey string) {
	// slog.Info("Saving splitter sizes")
	qs.BeginGroup(*qt.NewQAnyStringView3(mainKey))
	defer qs.EndGroup()
	sizes := splitterSizes(splitter)
	qs.SetValue(qs_sizes, newQVarStr(joinIntList(sizes)))
}

func RestoreSplitterSizes(qs *qt.QSettings, splitter *qt.QSplitter, mainKey string) {
	qs.BeginGroup(*qt.NewQAnyStringView3(mainKey))
	defer qs.EndGroup()
	if !qs.Contains(qs_sizes) {
		return
	}
	sizesStr := qs.ValueWithKey(qs_sizes).ToString()
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

func SetupSplitterSizesSave(qs *qt.QSettings, splitter *qt.QSplitter, mainKey string) {
	ch := make(chan time.Time, 100)
	splitter.OnSplitterMoved(func(pos int, index int) {
		ch <- time.Now()
	})
	go actionSaveLoop(ch, func() {
		saveSplitterSizes(qs, splitter, mainKey)
	})
}

func SetupDialogGeometrySave(
	qs *qt.QSettings,
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
	go actionSaveLoop(ch, func() {
		SaveDialogGeometry(qs, dialog, mainKey)
	})
}

func SetupMainWinGeometrySave(qs *qt.QSettings, window *qt.QMainWindow) {
	ch := make(chan time.Time, 100)

	window.OnMoveEvent(func(super func(*qt.QMoveEvent), event *qt.QMoveEvent) {
		ch <- time.Now()
	})
	window.OnResizeEvent(func(super func(*qt.QResizeEvent), event *qt.QResizeEvent) {
		ch <- time.Now()
	})
	go actionSaveLoop(ch, func() {
		saveMainWinGeometry(qs, window)
	})
}

func SaveSearchSettings(qs *qt.QSettings, combo *qt.QComboBox) {
	qs.BeginGroup(qs_search)
	defer qs.EndGroup()

	qs.SetValue(qs_mode, newQVarStr(
		strconv.FormatInt(int64(combo.CurrentIndex()), 10),
	))
}

func RestoreSearchSettings(qs *qt.QSettings, combo *qt.QComboBox) {
	qs.BeginGroup(qs_search)
	defer qs.EndGroup()

	restoreIntSetting(qs, qs_mode, combo.SetCurrentIndex)
}

func SaveActivityMode(qs *qt.QSettings, combo *qt.QComboBox) {
	qs.BeginGroup(qs_activity)
	defer qs.EndGroup()

	qs.SetValue(qs_mode, newQVarStr(
		strconv.FormatInt(int64(combo.CurrentIndex()), 10),
	))
}

func RestoreActivityMode(qs *qt.QSettings, combo *qt.QComboBox) {
	qs.BeginGroup(qs_activity)
	defer qs.EndGroup()

	restoreIntSetting(qs, qs_mode, combo.SetCurrentIndex)
}
