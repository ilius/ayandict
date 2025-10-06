package application

import (
	"log/slog"
	"os"

	qt "github.com/mappu/miqt/qt6"
)

func (app *Application) setupTrayIcon(icon *qt.QIcon) {
	window := app.window
	trayIcon := qt.NewQSystemTrayIcon2(icon)
	app.trayIcon = trayIcon
	trayIcon.OnActivated(func(reason qt.QSystemTrayIcon__ActivationReason) {
		if window.IsVisible() {
			window.Hide()
		} else {
			window.Show()
		}
	})
	trayIcon.OnMessageClicked(func() {
		slog.Info("trayIcon.OnMessageClicked")
	})
	app.setTrayMenu()

	trayIcon.Show()
}

func (app *Application) setTrayMenu() {
	menu := qt.NewQMenu2()
	{
		action := qt.NewQAction2("Quit")
		// icon will not align with checkboxes! so forget it!
		action.OnTriggered(func() { os.Exit(0) })
		menu.AddAction(action)
	}
	{
		action := qt.NewQAction2("Scan Selection")
		action.SetCheckable(true)
		action.SetChecked(conf.ScanPopupSelection)
		action.OnTriggeredWithChecked(func(checked bool) {
			conf.ScanPopupSelection = checked
		})
		menu.AddAction(action)
		app.trayScanSelection = action
	}
	{
		action := qt.NewQAction2("Scan Clipboard")
		action.SetCheckable(true)
		action.SetChecked(conf.ScanPopupClipboard)
		action.OnTriggeredWithChecked(func(checked bool) {
			conf.ScanPopupClipboard = checked
		})
		menu.AddAction(action)
		app.trayScanClipboard = action
	}

	app.trayIcon.SetContextMenu(menu)
}

func (app *Application) updateTrayMenu() {
	app.trayScanSelection.SetChecked(conf.ScanPopupSelection)
	app.trayScanClipboard.SetChecked(conf.ScanPopupClipboard)
}
