package application

import (
	"log/slog"
	"strings"

	"github.com/ilius/ayandict/v3/pkg/dictmgr"
	"github.com/ilius/ayandict/v3/pkg/qplatform"
	qt "github.com/mappu/miqt/qt6"
)

func (app *Application) setupScanPopup() {
	clipboard := qt.QGuiApplication_Clipboard()
	clipboard.OnChanged(func(mode qt.QClipboard__Mode) {
		if mode == qt.QClipboard__Clipboard && !conf.ScanPopupClipboard {
			return
		}
		if mode == qt.QClipboard__Selection && !conf.ScanPopupSelection {
			return
		}
		app.scanPopup(clipboard.TextWithMode(mode))
	})
}

func (app *Application) onScanPopupCloseEvent(super func(*qt.QCloseEvent), event *qt.QCloseEvent) {
	app.scanPopupCount.Add(-1)
}

func (app *Application) OnScanPopupShow() {
	app.scanPopupCount.Add(1)
}

func (app *Application) PreparePopup(popup *qt.QWidget) {
	if qplatform.CanMoveWindow() {
		return
	}
	if app.desktopWidget != nil {
		app.desktopWidget.Show()
		app.desktopWidget.ActivateWindow()
		app.desktopWidget.Raise()
		popup.CreateWinId() // forces creation of QWindow
		handle := popup.WindowHandle()
		if handle == nil {
			slog.Warn("popup handle is nil")
		} else {
			handle.SetParent(app.desktopWidget.WindowHandle())
		}
	}
}

func (app *Application) ShowWindowAndQuery(query string) {
	app.window.ShowNormal()
	app.window.ActivateWindow()
	app.doQuery(query)
}

func (app *Application) AddHistoryAndFrequency(query string) {
	app.queryArgs.AddHistoryAndFrequency(query)
}

func (app *Application) scanPopup(query string) {
	slog.Debug("app.scanPopup", "query", query, "count", app.scanPopupCount.Load())
	if conf.ScanPopupMaxCount > 0 && app.scanPopupCount.Load() >= conf.ScanPopupMaxCount {
		return
	}
	query = strings.TrimSpace(query)
	query = strings.Trim(query, punctuation)
	if query == "" {
		return
	}
	mode, valid := dictmgr.SearchModeByName(conf.ScanPopupMode)
	if !valid {
		slog.Error("invalid scan_popup_mode", "value", conf.ScanPopupMode)
	}

	p := NewScanPopup(
		query,
		mode,
		app.onScanPopupCloseEvent,
		app,
	)
	p.Run(qt.QCursor_Pos(), app.icon)
}

func (app *Application) randomFavoritePopup(onClose func()) {
	term := app.favoritesWidget.Data.Random()
	if term == "" {
		// show "No Favorites" error?
		return
	}
	if conf.ScanPopupMaxCount > 0 && app.scanPopupCount.Load() >= conf.ScanPopupMaxCount {
		return
	}

	onCloseNew := func(super func(event *qt.QCloseEvent), event *qt.QCloseEvent) {
		app.scanPopupCount.Add(-1)
		onClose()
	}

	p := NewScanPopup(
		term,
		dictmgr.SearchModeStartWith,
		onCloseNew,
		app,
	)
	p.Run(nil, app.icon) // on center of primary screen
}
