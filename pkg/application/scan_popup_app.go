package application

import (
	"log/slog"
	"strings"

	"github.com/ilius/ayandict/v3/pkg/dictmgr"
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

func (app *Application) scanPopup(query string) *ScanPopup {
	if conf.ScanPopupMaxCount > 0 && app.scanPopupCount.Load() >= conf.ScanPopupMaxCount {
		return nil
	}
	query = strings.TrimSpace(query)
	query = strings.Trim(query, punctuation)
	if query == "" {
		return nil
	}
	mode, valid := dictmgr.SearchModeByName(conf.ScanPopupMode)
	if !valid {
		slog.Error("invalid scan_popup_mode", "value", conf.ScanPopupMode)
	}

	app.scanPopupCount.Add(1)

	p := NewScanPopup(
		query,
		mode,
		qt.QCursor_Pos(),
		app.icon,
		app.showWindowAndQuery,
		app.queryArgs.AddHistoryAndFrequency,
		app.onScanPopupCloseEvent,
	)
	p.Run()
	return p
}

func (app *Application) showWindowAndQuery(query string) {
	app.window.ShowNormal()
	app.window.ActivateWindow()
	app.doQuery(query)
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
	app.scanPopupCount.Add(1)

	onCloseNew := func(super func(event *qt.QCloseEvent), event *qt.QCloseEvent) {
		app.scanPopupCount.Add(-1)
		onClose()
	}

	p := NewScanPopup(
		term,
		dictmgr.SearchModeStartWith,
		nil, // on center of primary screen
		app.icon,
		app.showWindowAndQuery,
		app.queryArgs.AddHistoryAndFrequency,
		onCloseNew,
	)
	p.Run()
}
