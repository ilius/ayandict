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
	// menu := qt.NewQMenu2()
	// menu.AddAction2()
	// trayIcon.SetContextMenu(menu)
	trayIcon.Show()
}
