package application

import (
	"log/slog"

	qt "github.com/mappu/miqt/qt6"
)

func (app *Application) setupTrayIcon(icon *qt.QIcon) {
	trayIcon := qt.NewQSystemTrayIcon2(icon)
	app.trayIcon = trayIcon
	trayIcon.OnActivated(func(reason qt.QSystemTrayIcon__ActivationReason) {
		app.onStatusIconClick()
	})
	trayIcon.OnMessageClicked(func() {
		slog.Info("trayIcon.OnMessageClicked")
	})
	app.statusIconActions = app.createStatusIconActions()
	app.setTrayMenu()

	trayIcon.Show()
}

func (app *Application) createStatusIconActions() []*qt.QAction {
	actions := []*qt.QAction{}
	{
		action := qt.NewQAction2("Quit")
		action.OnTriggered(app.Exit)
		actions = append(actions, action)
	}
	{
		action := qt.NewQAction2("About")
		action.OnTriggered(func() {
			widget := qt.NewQDialog(app.window.QWidget)
			aboutClickedWidget(widget.QWidget, app.icon)
			widget.Show()
		})
		actions = append(actions, action)
	}
	{
		action := qt.NewQAction2("Show Window")
		action.OnTriggered(func() {
			app.window.ShowNormal()
			app.window.ActivateWindow()
		})
		actions = append(actions, action)
	}
	if qt.QGuiApplication_Clipboard().SupportsSelection() {
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
	return actions
}

func (app *Application) setTrayMenu() {
	// icon will not align with checkboxes! so forget it!
	menu := qt.NewQMenu2()
	// TODO: reverse order of actions if tray icon is closer to top of screen
	for _, action := range app.statusIconActions {
		menu.AddAction(action)
	}
	app.trayIcon.SetContextMenu(menu)
}

func (app *Application) updateTrayMenu() {
	if app.trayScanSelection != nil {
		app.trayScanSelection.SetChecked(conf.ScanPopupSelection)
	}
	app.trayScanClipboard.SetChecked(conf.ScanPopupClipboard)
	app.trayScanAPI.SetChecked(conf.ScanPopupAPI)
}
