package qsettings

import (
	"time"

	"github.com/ilius/ayandict/v3/pkg/qtutils"
	qt "github.com/mappu/miqt/qt6"
)

var qs_columnwidth = *qt.NewQAnyStringView3(QS_columnwidth)

// we need this because passing dialog.QWidget to SetupWindowGeometrySave causes:
// panic: miqt: can only override virtual methods for directly constructed types [recovered]
type windowSaveInterface interface {
	OnMoveEvent(func(super func(*qt.QMoveEvent), event *qt.QMoveEvent))
	OnResizeEvent(func(super func(*qt.QResizeEvent), event *qt.QResizeEvent))
	Pos() *qt.QPoint
	Size() *qt.QSize
}

func SaveWindowGeometry(window windowSaveInterface, mainKey string) {
	// slog.Info("Saving main window geometry")
	pos := window.Pos()
	size := window.Size()
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

func SetupWindowGeometrySave(
	window windowSaveInterface,
	mainKey string,
) {
	ch := make(chan time.Time, 100)

	window.OnMoveEvent(func(super func(*qt.QMoveEvent), event *qt.QMoveEvent) {
		ch <- time.Now()
	})
	window.OnResizeEvent(func(super func(*qt.QResizeEvent), event *qt.QResizeEvent) {
		ch <- time.Now()
	})
	go ActionSaveLoop(ch, func() {
		SaveWindowGeometry(window, mainKey)
	})
}
