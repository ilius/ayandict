package application

import (
	"fmt"
	"log/slog"
	"runtime"

	"github.com/ilius/ayandict/v3/pkg/appinfo"
	qt "github.com/mappu/miqt/qt6"
)

var isMac = runtime.GOOS == "darwin"

func addTabWithIcon(
	tabWidget *qt.QTabWidget,
	widget *qt.QWidget,
	label string,
	filename string,
) {
	if isMac {
		_ = tabWidget.AddTab(widget, label)
		return
	}
	icon, err := loadPNGIcon(filename)
	if err != nil {
		slog.Error("error loading icon", "filename", filename)
	}
	if icon == nil {
		_ = tabWidget.AddTab(widget, label)
		return
	}
	_ = tabWidget.AddTab2(widget, icon, label)
}

func aboutClickedWidget(widget *qt.QWidget, icon *qt.QIcon) {
	widget.SetWindowTitle("About AyanDict")
	widget.Resize(700, 500)
	widget.SetWindowIcon(icon)

	topHBox := qt.NewQFrame(nil)
	topHBoxLayout := qt.NewQHBoxLayout(topHBox.QWidget)

	{
		pixmap, err := loadPNGPixmap("ayandict-64px.png")
		if err != nil {
			slog.Error("failed to load icon image", "err", err)
		} else {
			label := qt.NewQLabel2()
			label.SetPixmap(pixmap)
			label.SetMinimumWidth(80)
			topHBoxLayout.AddWidget3(label.QWidget, 0, qt.AlignCenter)
		}
	}

	topLabel := qt.NewQLabel3(fmt.Sprintf(
		"AyanDict version %s\nUsing Qt %v and Go %v",
		appinfo.VERSION,
		qt.QLibraryInfo_Version().ToString(),
		runtime.Version()[2:],
	))
	topHBoxLayout.AddWidget(topLabel.QWidget)
	topHBoxLayout.AddStretch()

	tabWidget := qt.NewQTabWidget2()
	tabWidget.SetSizePolicy2(expanding, expanding)
	tabWidget.SetIconSize(qt.NewQSize2(22, 22))

	// tabWidget.SetTabPosition(qt.QTabWidget__West)
	// tabBar := tabWidget.TabBar()
	// tabWidget.SetStyleSheet(`
	// QTabBar::tab {
	// 	direction: ltr;
	// 	padding: 15px;
	// }`)

	aboutLabel := qt.NewQLabel3(appinfo.ABOUT)
	aboutLabel.SetTextInteractionFlags(qt.TextSelectableByMouse)
	aboutLabel.SetAlignment(qt.AlignTop)
	aboutLabel.SetWordWrap(true)
	addTabWithIcon(tabWidget, aboutLabel.QWidget, "About", "dialog-information-22.png")

	authorsLabel := qt.NewQLabel3(appinfo.AUTHORS)
	authorsLabel.SetTextInteractionFlags(qt.TextSelectableByMouse)
	authorsLabel.SetAlignment(qt.AlignTop)
	addTabWithIcon(tabWidget, authorsLabel.QWidget, "Authors", "author-22.png")

	licenseWidget := qt.NewQTextEdit2()
	licenseWidget.SetReadOnly(true)
	licenseWidget.SetPlainText(appinfo.LICENSE)
	addTabWithIcon(tabWidget, licenseWidget.QWidget, "License", "license-22.png")

	buttonBox := qt.NewQDialogButtonBox2()
	buttonBox.AddButton2("Keyboard Shortcuts", qt.QDialogButtonBox__AcceptRole).OnClicked(func() {
		showKeyBindings(widget, icon)
	})
	closeButton := buttonBox.AddButton2("  Close  ", qt.QDialogButtonBox__AcceptRole)
	closeButton.OnClicked(func() {
		widget.Close()
	})
	// closeButton.SetDefault(true)

	mainBox := qt.NewQVBoxLayout(widget)
	mainBox.AddWidget(topHBox.QWidget)
	mainBox.AddWidget(tabWidget.QWidget)
	mainBox.AddWidget(buttonBox.QWidget)
}

func labelFrameWithLine(text string) *qt.QFrame {
	label := qt.NewQLabel3(text)
	frame := qt.NewQFrame2()
	// frame.SetStyleSheet("border: 1px solid rgba(128,128,128,0.5);")
	frame.SetFrameShape(qt.QFrame__Box)    // The border shape
	frame.SetFrameShadow(qt.QFrame__Plain) // Flat (not sunken/raised)
	frame.SetLineWidth(1)
	frameLayout := qt.NewQHBoxLayout(frame.QWidget)
	frameLayout.AddWidget3(label.QWidget, 1, qt.AlignLeft)
	return frame
}

func showKeyBindings(parent *qt.QWidget, icon *qt.QIcon) {
	dialog := qt.NewQDialog(parent)
	dialog.SetWindowIcon(icon)
	dialog.SetWindowTitle("Keyboard Shortcuts")
	dialog.Resize(800, 600)
	layout := qt.NewQGridLayout2()

	layout.SetColumnStretch(0, 0)
	layout.SetColumnStretch(1, 1)
	layout.SetColumnStretch(2, 0)

	layout.AddWidget2(qt.NewQLabel3("Key").QWidget, 0, 0)
	layout.AddWidget2(qt.NewQLabel3("while").QWidget, 0, 1)
	layout.AddWidget2(qt.NewQLabel3("Action").QWidget, 0, 2)

	for rowI, data := range appinfo.KeyBindings {
		if data[1] == "" {
			layout.AddWidget3(
				labelFrameWithLine(data[0]).QWidget,
				rowI+1, 0, 1, 2,
			)
		} else {
			layout.AddWidget3(
				labelFrameWithLine(data[0]).QWidget,
				rowI+1, 0, 1, 1,
			)
			layout.AddWidget3(
				labelFrameWithLine(data[1]).QWidget,
				rowI+1, 1, 1, 1,
			)
		}
		layout.AddWidget3(
			labelFrameWithLine(data[2]).QWidget,
			rowI+1, 2, 1, 1,
		)
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
