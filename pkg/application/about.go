package application

import (
	"fmt"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

func addTabWithIcon(
	tabWidget *widgets.QTabWidget,
	widget widgets.QWidget_ITF,
	label string,
	filename string,
) {
	icon := loadIcon(filename)
	if icon == nil {
		tabWidget.AddTab(widget, label)
		return
	}
	tabWidget.AddTab2(widget, icon, label)
}

func aboutClicked(
	parent widgets.QWidget_ITF,
) {
	window := widgets.NewQDialog(parent, core.Qt__Dialog)
	window.SetWindowTitle("About AyanDict")
	window.Resize2(800, 400)

	topHBox := widgets.NewQHBoxLayout()
	topLabel := widgets.NewQLabel2(fmt.Sprintf(
		"AyanDict\nVersion %s",
		VERSION,
	), nil, 0)
	topHBox.AddWidget(topLabel, 0, 0)

	tabWidget := widgets.NewQTabWidget(nil)
	tabWidget.SetSizePolicy2(expanding, expanding)
	tabWidget.SetIconSize(core.NewQSize2(22, 22))

	// tabWidget.SetTabPosition(widgets.QTabWidget__West)
	// tabBar := tabWidget.TabBar()
	// tabWidget.SetStyleSheet(`
	// QTabBar::tab {
	// 	direction: ltr;
	// 	padding: 15px;
	// }`)

	aboutLabel := widgets.NewQLabel2(ABOUT, nil, 0)
	aboutLabel.SetTextInteractionFlags(core.Qt__TextSelectableByMouse)
	aboutLabel.SetAlignment(core.Qt__AlignTop)
	addTabWithIcon(tabWidget, aboutLabel, "About", "dialog-information-22.png")

	authorsLabel := widgets.NewQLabel2(AUTHORS, nil, 0)
	authorsLabel.SetTextInteractionFlags(core.Qt__TextSelectableByMouse)
	authorsLabel.SetAlignment(core.Qt__AlignTop)
	addTabWithIcon(tabWidget, authorsLabel, "Authors", "author-22.png")

	licenseWidget := widgets.NewQTextEdit(nil)
	licenseWidget.SetReadOnly(true)
	licenseWidget.SetPlainText(LICENSE)
	addTabWithIcon(tabWidget, licenseWidget, "License", "license-22.png")

	buttonBox := widgets.NewQDialogButtonBox(nil)
	okButton := buttonBox.AddButton2("OK", widgets.QDialogButtonBox__AcceptRole)
	okButton.ConnectClicked(func(checked bool) {
		window.Accept()
	})

	mainBox := widgets.NewQVBoxLayout2(window)
	mainBox.AddLayout(topHBox, 0)
	mainBox.AddWidget(tabWidget, 0, 0)
	mainBox.AddWidget(buttonBox, 0, 0)

	window.Exec()
}