package application

import (
	"log/slog"

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
	actions := []*qt.QAction{}
	// icon will not align with checkboxes! so forget it!
	{
		action := qt.NewQAction2("Quit")
		action.OnTriggered(app.Exit)
		actions = append(actions, action)
	}
	{
		action := qt.NewQAction2("About")
		action.OnTriggered(func() {
			aboutClicked(menu.QWidget, app.icon)
		})
		actions = append(actions, action)
	}
	{
		action := qt.NewQAction2("Show Window")
		action.OnTriggered(func() {
			app.window.Show()
			app.window.Raise()
		})
		actions = append(actions, action)
	}
	{
		action := qt.NewQAction2("Scan Selection")
		action.SetCheckable(true)
		action.SetChecked(conf.ScanPopupSelection)
		action.OnTriggeredWithChecked(func(checked bool) {
			conf.ScanPopupSelection = checked
		})
		actions = append(actions, action)
		app.trayScanSelection = action
	}
	{
		action := qt.NewQAction2("Scan Clipboard")
		action.SetCheckable(true)
		action.SetChecked(conf.ScanPopupClipboard)
		action.OnTriggeredWithChecked(func(checked bool) {
			conf.ScanPopupClipboard = checked
		})
		actions = append(actions, action)
		app.trayScanClipboard = action
	}
	{
		action := qt.NewQAction2("Scan via API")
		action.SetCheckable(true)
		action.SetChecked(conf.ScanPopupAPI)
		action.OnTriggeredWithChecked(func(checked bool) {
			conf.ScanPopupAPI = checked
		})
		actions = append(actions, action)
		app.trayScanAPI = action
	}
	// TODO: reverse order of actions if tray icon is closer to top of screen
	for _, action := range actions {
		menu.AddAction(action)
	}
	app.trayIcon.SetContextMenu(menu)
}

func (app *Application) updateTrayMenu() {
	app.trayScanSelection.SetChecked(conf.ScanPopupSelection)
	app.trayScanClipboard.SetChecked(conf.ScanPopupClipboard)
	app.trayScanAPI.SetChecked(conf.ScanPopupAPI)
}
