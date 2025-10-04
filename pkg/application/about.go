package application

import (
	"fmt"
	"log/slog"
	"runtime"

	"github.com/ilius/ayandict/v3/pkg/appinfo"
	qt "github.com/mappu/miqt/qt6"
)

func addTabWithIcon(
	tabWidget *qt.QTabWidget,
	widget *qt.QWidget,
	label string,
	filename string,
) {
	icon, err := loadPNGIcon(filename)
	if err != nil {
		fmt.Println(err)
	}
	if icon == nil {
		_ = tabWidget.AddTab2(widget, nil, label)
		return
	}
	_ = tabWidget.AddTab2(widget, icon, label)
}

func aboutClicked(
	parent *qt.QWidget,
) {
	window := qt.NewQDialog(parent)
	window.SetWindowTitle("About AyanDict")
	window.Resize(700, 500)
	window.SetWindowIcon(parent.WindowIcon())

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
	buttonBox.AddButton2(
		"Close",
		qt.QDialogButtonBox__AcceptRole,
	).OnClicked(func() {
		window.Accept()
	})

	mainBox := qt.NewQVBoxLayout(window.QWidget)
	mainBox.AddWidget(topHBox.QWidget)
	mainBox.AddWidget(tabWidget.QWidget)
	mainBox.AddWidget(buttonBox.QWidget)

	_ = window.Exec()
}
