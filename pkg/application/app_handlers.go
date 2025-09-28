package application

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/ilius/ayandict/v3/pkg/config"
	"github.com/ilius/ayandict/v3/pkg/dictmgr"
	"github.com/ilius/ayandict/v3/pkg/dictmgr/qdictmgr"
	"github.com/ilius/ayandict/v3/pkg/qtcommon/qsettings"
	common "github.com/ilius/go-dict-commons"
	qt "github.com/mappu/miqt/qt6"
)

func (app *Application) setupHandlers() {
	app.articleView.SetupCustomHandlers()
	app.historyView.SetupCustomHandlers()

	entry := app.entry
	queryArgs := app.queryArgs
	frequencyTable := app.frequencyTable

	frequencyTable.OnItemActivated(func(item *qt.QTableWidgetItem) {
		key := frequencyTable.Keys[item.Row()]
		app.doQuery(key)
		newRow := frequencyTable.KeyMap[key]
		// item.Column() panics!
		frequencyTable.SetCurrentCell(newRow, 0)
	})
	app.favoritesWidget.OnItemActivated(func(item *qt.QListWidgetItem) {
		app.doQuery(item.Text())
	})

	app.reloadDictsButton.OnClicked(func() {
		qdictmgr.InitDicts(conf, true)
		app.dictManager = nil
		onQuery(entry.Text(), queryArgs, false)
	})
	app.closeDictsButton.OnClicked(dictmgr.CloseDicts)
	app.openConfigButton.OnClicked(OpenConfig)
	app.reloadConfigButton.OnClicked(app.ReloadConfig)
	app.reloadStyleButton.OnClicked(func() {
		app.LoadUserStyle()
		onQuery(entry.Text(), queryArgs, false)
	})
	app.saveHistoryButton.OnClicked(func() {
		app.historyView.Save()
		frequencyTable.SaveNoError()
	})
	app.clearHistoryButton.OnClicked(app.clearHistoryClicked)
	app.saveFavoritesButton.OnClicked(app.saveFavoritesClicked)
	app.clearButton.OnClicked(app.resetQuery)
	app.dictsButton.OnClicked(app.dictsButtonClicked)
	app.randomEntryButton.OnClicked(app.randomEntryClicked)
	app.randomFavoriteButton.OnClicked(app.randomFavoriteClicked)
	entry.OnKeyPressEvent(app.onEntryKeyPress)

	if config.PrivateMode {
		app.favoriteButton.SetDisabled(true)
		app.queryFavoriteButton.SetDisabled(true)
	}
	// slog.Error("test error", "s", "hello", "n", 2, "b", true)
}

func (app *Application) setupKeyPressEvent(widget KeyPressIface) {
	widget.OnKeyPressEvent(func(super func(event *qt.QKeyEvent), event *qt.QKeyEvent) {
		switch event.Key() {
		case int(qt.Key_Space): // " "
			app.entry.SetFocusWithReason(qt.ShortcutFocusReason)
			return
		case int(qt.Key_Plus), int(qt.Key_Equal): // "+", "="
			app.articleView.ZoomIn()
			return
		case int(qt.Key_Minus): // "-"
			app.articleView.ZoomOut()
			return
		case escape: // event.Text()="\x1b"
			app.resetQuery()
			return
		case int(qt.Key_F1):
			aboutClicked(app.window.QWidget)
			return
		}
		super(event)
	})
}

func (app *Application) onEntryKeyPress(super func(*qt.QKeyEvent), event *qt.QKeyEvent) {
	// slog.Info(
	// 	"entry: KeyPressEvent",
	// 	"text", fmt.Sprintf("%#v", event.Text()),
	// 	"key", event.Key(),
	// )
	key := event.Key()
	switch key {
	case escape: // event.Text()="\x1b"
		app.window.SetFocus()
		return
	case int(qt.Key_Return), int(qt.Key_Enter): // event.Text()="\r"
		onQuery(app.entry.Text(), app.queryArgs, false)
		return
	}

	super(event)

	// event.Modifiers(): qt.NoModifier, qt.ShiftModifier, KeypadModifier
	if conf.SearchOnType && key < escape {
		if int(event.Modifiers())&shortcutModifierMask == 0 {
			text := app.entry.Text()
			// slog.Debug("checking SearchOnType") // FIXME: panics
			if len(text) >= conf.SearchOnTypeMinLength {
				onQuery(text, app.queryArgs, true)
			}
			return
		}
	}
}

func (app *Application) activityComboChanged(index int) {
	switch index {
	case 0:
		app.historyView.Show()
		app.frequencyTable.Hide()
		app.favoritesWidget.Hide()
	case 1:
		app.historyView.Hide()
		app.frequencyTable.Show()
		app.favoritesWidget.Hide()
	case 2:
		app.historyView.Hide()
		app.frequencyTable.Hide()
		app.favoritesWidget.Show()
	}
	qsettings.SaveActivityMode(app.qs, app.activityTypeCombo)
}

func (app *Application) okButtonResized(
	_ func(*qt.QResizeEvent),
	event *qt.QResizeEvent,
) {
	h := event.Size().Height()
	if h > 100 {
		return
	}
	app.queryFavoriteButton.SetFixedSize2(h, h)
	app.favoriteButton.SetFixedSize2(h, h)
}

func (app *Application) queryFavoriteButtonClicked(checked bool) {
	term := app.entry.Text()
	if term == "" {
		app.queryFavoriteButton.SetChecked(false)
		return
	}
	app.favoritesWidget.SetFavorite(term, checked)
	if app.resultList.Active != nil && term == app.resultList.Active.Terms()[0] {
		app.favoriteButton.SetChecked(checked)
	}
}

func (app *Application) favoriteButtonClicked(checked bool) {
	if app.resultList.Active == nil {
		app.favoriteButton.SetChecked(false)
		return
	}
	term := app.resultList.Active.Terms()[0]
	app.favoritesWidget.SetFavorite(term, checked)
	if term == app.entry.Text() {
		app.queryFavoriteButton.SetChecked(checked)
	}
}

func (app *Application) clearHistoryClicked() {
	app.historyView.ClearHistory()
	app.frequencyTable.Clear()
	app.frequencyTable.SaveNoError()
}

func (app *Application) saveFavoritesClicked() {
	err := app.favoritesWidget.Save()
	if err != nil {
		slog.Error("error saving favorites: " + err.Error())
	}
}

func (app *Application) dictsButtonClicked() {
	if app.runDictManager() {
		onQuery(app.entry.Text(), app.queryArgs, false)
	}
}

func (app *Application) randomEntryClicked() {
	res := dictmgr.RandomEntry(conf, resultFlags)
	if res == nil {
		return
	}
	query := res.F_Terms[0]
	app.entry.SetText(query)
	app.queryArgs.ResultsLabel.SetText("Results: 1")
	app.queryArgs.ResultList.SetResults([]common.SearchResultIface{res})
	app.queryArgs.AddHistoryAndFrequency(query)
	app.postQuery(query)
}

func (app *Application) randomFavoriteClicked() {
	term := app.favoritesWidget.Data.Random()
	if term == "" {
		// show "No Favorites" error?
		return
	}
	app.entry.SetText(term)
	queryArgs := app.queryArgs

	t := time.Now()
	mode, ok := dictmgr.QueryModeByName(conf.RandomFavoriteSearchMode)
	if !ok {
		slog.Error("invalid random_favorite_search_mode", "value", conf.RandomFavoriteSearchMode)
		mode = dictmgr.QueryModeWordMatch
	}
	results := dictmgr.LookupHTML(term, conf, mode, resultFlags, 0)
	slog.Debug("LookupHTML running time", "dt", time.Since(t), "query", term)
	queryArgs.ResultList.SetResults(results)
	queryArgs.ResultsLabel.SetText(fmt.Sprintf("Results: %d", len(results)))
	if len(results) == 0 {
		queryArgs.SetNoResult(term)
	}

	queryArgs.AddHistoryAndFrequency(term)
	queryArgs.PostQuery(term)
}
