package application

import (
	"log/slog"
	"strings"

	"github.com/ilius/ayandict/v3/pkg/dictmgr"
	"github.com/ilius/ayandict/v3/pkg/scanpopup"
	"github.com/ilius/ayandict/v3/pkg/utils"
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
		app.QueryPopup(clipboard.TextWithMode(mode))
	})
}

func (app *Application) OnScanPopupClose(super func(*qt.QCloseEvent), event *qt.QCloseEvent) {
	app.scanPopupCount.Add(-1)
}

func (app *Application) OnScanPopupShow() {
	app.scanPopupCount.Add(1)
}

func (app *Application) ShowWindowAndQuery(query string) {
	app.window.ShowNormal()
	app.window.ActivateWindow()
	app.Query(query)
}

func (app *Application) AddHistoryAndFrequency(query string) {
	app.queryArgs.AddHistoryAndFrequency(query)
}

func (app *Application) QueryPopup(query string) {
	slog.Debug("app.QueryPopup", "query", query, "count", app.scanPopupCount.Load())
	if conf.ScanPopupMaxCount > 0 && app.scanPopupCount.Load() >= conf.ScanPopupMaxCount {
		return
	}
	query = strings.TrimSpace(query)
	query = strings.Trim(query, utils.Punctuation)
	if query == "" {
		return
	}
	mode, valid := dictmgr.SearchModeByName(conf.ScanPopupMode)
	if !valid {
		slog.Error("invalid scan_popup_mode", "value", conf.ScanPopupMode)
	}

	p := scanpopup.NewScanPopup(
		conf,
		query,
		mode,
		app.OnScanPopupClose,
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

	p := scanpopup.NewScanPopup(
		conf,
		term,
		dictmgr.SearchModeStartWith,
		onCloseNew,
		app,
	)
	p.Run(nil, app.icon) // on center of primary screen
}
