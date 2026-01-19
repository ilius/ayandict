package about

import (
	"github.com/ilius/ayandict/v3/pkg/appinfo"
	"github.com/ilius/ayandict/v3/pkg/qtutils"
	qt "github.com/mappu/miqt/qt6"
)

func linedLabel(text string, right int, bottom int) *qt.QWidget {
	label := qt.NewQLabel3(text)
	frame := qt.NewQFrame2()
	frame.SetFrameShape(qt.QFrame__Box)    // The border shape
	frame.SetFrameShadow(qt.QFrame__Plain) // Flat (not sunken/raised)
	frame.SetContentsMargins(1, 1, right, bottom)
	frameLayout := qt.NewQHBoxLayout(frame.QWidget)
	frameLayout.AddWidget3(label.QWidget, 1, qt.AlignCenter)
	frameLayout.SetContentsMargins(3, 0, 3, 0)
	return frame.QWidget
}

func newKeyBindingsGrid(keyBindings [][3]string, hasWhile bool) *qt.QGridLayout {
	layout := qt.NewQGridLayout2()

	layout.SetVerticalSpacing(0)
	layout.SetHorizontalSpacing(0)

	layout.SetColumnStretch(0, 0)
	layout.SetColumnStretch(1, 0)
	layout.SetColumnStretch(2, 1)

	if hasWhile {
		layout.AddWidget4(qt.NewQLabel3("Key").QWidget, 0, 0, qt.AlignCenter)
		layout.AddWidget4(qt.NewQLabel3("[while]").QWidget, 0, 1, qt.AlignCenter)
	} else {
		layout.AddWidget5(qt.NewQLabel3("Key").QWidget, 0, 0, 1, 2, qt.AlignCenter)
	}
	layout.AddWidget4(qt.NewQLabel3("Action").QWidget, 0, 2, qt.AlignCenter)

	for rowI, data := range keyBindings {
		bottom := 0
		if rowI == len(keyBindings)-1 {
			bottom = 1
		}
		widget1 := linedLabel(data[0], 0, bottom)
		if data[1] == "" {
			layout.AddWidget3(widget1, rowI+1, 0, 1, 2)
		} else {
			widget2 := linedLabel(data[1], 0, bottom)
			layout.AddWidget3(widget1, rowI+1, 0, 1, 1)
			layout.AddWidget3(widget2, rowI+1, 1, 1, 1)
		}
		layout.AddWidget3(linedLabel(data[2], 1, bottom), rowI+1, 2, 1, 1)
	}

	return layout
}

func showKeyBindings(parent *qt.QWidget, icon *qt.QIcon) {
	dialog := qt.NewQDialog(parent)
	dialog.SetWindowIcon(icon)
	dialog.SetWindowTitle("Keyboard Shortcuts")
	qtutils.SetWinSize(dialog.QWidget, 800, 400)

	grid1 := newKeyBindingsGrid(appinfo.KeyBindings1, false)
	grid2 := newKeyBindingsGrid(appinfo.KeyBindings2, true)

	mainHBox := qt.NewQHBoxLayout2()
	mainHBox.AddLayout2(grid1.QLayout, 1)
	mainHBox.AddSpacing(10)
	mainHBox.AddLayout2(grid2.QLayout, 1)

	dialog.OnKeyPressEvent(func(super func(e *qt.QKeyEvent), e *qt.QKeyEvent) {
		if e.Key() == int(qt.Key_Escape) {
			dialog.Close()
			return
		}
		super(e)
	})

	buttonBox := qt.NewQDialogButtonBox2()
	buttonBox.AddButton2(
		"  Close  ",
		qt.QDialogButtonBox__AcceptRole,
	).OnClicked(func() {
		_ = dialog.Close()
	})

	mainBox := qt.NewQVBoxLayout2()
	mainBox.AddLayout2(mainHBox.Layout(), 1)
	mainBox.AddSpacing(10)
	mainBox.AddWidget(buttonBox.QWidget)
	dialog.SetLayout(mainBox.Layout())

	dialog.Exec()
}
