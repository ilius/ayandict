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

func (app *Application) scanPopup(query string) {
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

	app.scanPopupCount.Add(1)

	p := NewScanPopup(
		query,
		mode,
		app.icon,
		app.showWindowAndQuery,
		app.queryArgs.AddHistoryAndFrequency,
		app.onScanPopupCloseEvent,
	)
	p.Run()
}

func (app *Application) showWindowAndQuery(query string) {
	app.window.ShowNormal()
	app.window.ActivateWindow()
	app.doQuery(query)
}
