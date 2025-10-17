package application

import (
	"github.com/ilius/ayandict/v3/pkg/appinfo"
	qt "github.com/mappu/miqt/qt6"
)

func showKeyBindings(parent *qt.QWidget, icon *qt.QIcon) {
	dialog := qt.NewQDialog(parent)
	dialog.SetWindowIcon(icon)
	dialog.SetWindowTitle("Keyboard Shortcuts")
	dialog.Resize(800, 600)

	layout := qt.NewQGridLayout2()

	layout.SetVerticalSpacing(0)
	layout.SetHorizontalSpacing(0)

	layout.SetColumnStretch(0, 0)
	layout.SetColumnStretch(1, 0)
	layout.SetColumnStretch(2, 1)

	layout.AddWidget4(qt.NewQLabel3("Key").QWidget, 0, 0, qt.AlignCenter)
	layout.AddWidget4(qt.NewQLabel3("while").QWidget, 0, 1, qt.AlignCenter)
	layout.AddWidget4(qt.NewQLabel3("Action").QWidget, 0, 2, qt.AlignCenter)

	for rowI, data := range appinfo.KeyBindings {
		bottom := 0
		if rowI == len(appinfo.KeyBindings)-1 {
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
	).OnClicked(func() { dialog.Close() })

	mainBox := qt.NewQVBoxLayout2()
	mainBox.AddLayout2(layout.Layout(), 1)
	mainBox.AddSpacing(10)
	mainBox.AddWidget(buttonBox.QWidget)
	dialog.SetLayout(mainBox.Layout())

	dialog.Exec()
}
