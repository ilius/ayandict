package qsettings

import (
	"time"

	"github.com/ilius/ayandict/v3/pkg/qtutils"
	qt "github.com/mappu/miqt/qt6"
)

// we need this because passing dialog.QWidget to SetupWindowGeometrySave causes:
// panic: miqt: can only override virtual methods for directly constructed types [recovered]
type windowSaveInterface interface {
	OnMoveEvent(func(func(*qt.QMoveEvent), *qt.QMoveEvent))
	OnResizeEvent(func(func(*qt.QResizeEvent), *qt.QResizeEvent))
	Pos() *qt.QPoint
	Size() *qt.QSize
}

func SaveWindowGeometry(window windowSaveInterface, mainKey string) {
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
	qtutils.SetWinPosition(window, s.X, s.Y)
	qtutils.SetWinSize(window, qt.NewQSize2(s.Width, s.Height))
	if s.Maximized {
		window.ShowMaximized()
	}
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
