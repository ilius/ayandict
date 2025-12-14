package application

import (
	"log/slog"
	"time"

	common "codeberg.org/ilius/go-dict-commons"
	"github.com/ilius/ayandict/v3/pkg/config"
	"github.com/ilius/ayandict/v3/pkg/dictmgr"
	"github.com/ilius/ayandict/v3/pkg/dictmgr/qdictmgr"
	qt "github.com/mappu/miqt/qt6"
)

func (app *Application) setupHandlers() {
	// MUST not call OnKeyPressEvent multiple times on the same widget
	// that's why I separated setupArticleViewKeyPressEvent from setupKeyPressEvent

	app.setupKeyPressEvent(app.window)
	app.setupKeyPressEvent(app.resultList.QListWidget)
	app.setupKeyPressEvent(app.historyView.QListWidget)

	app.articleView.OnKeyPressEvent(app.onArticleViewKeyPressEvent)

	entry := app.entry
	frequencyTable := app.frequencyTable

	frequencyTable.OnItemActivated(func(item *qt.QTableWidgetItem) {
		key := frequencyTable.Keys[item.Row()]
		app.Query(key)
		newRow := frequencyTable.KeyMap[key]
		column := item.Column()
		// slog.Info("frequencyTable.OnItemActivated", "newRow", newRow, "column", column)
		if column < 0 {
			column = 1 // for some reason, it's -1 instead of 1
		}
		frequencyTable.SetCurrentCell(newRow, column)
	})
	app.favoritesWidget.OnItemActivated(func(item *qt.QListWidgetItem) {
		app.Query(item.Text())
	})

	app.reloadDictsButton.OnClicked(func() {
		qdictmgr.InitDicts(conf, true)
		app.dictManager = nil
		app.onQuery(entry.Text(), false)
	})
	app.closeDictsButton.OnClicked(dictmgr.CloseDicts)
	app.openConfigButton.OnClicked(OpenConfig)
	app.reloadConfigButton.OnClicked(app.ReloadConfig)
	app.reloadStyleButton.OnClicked(func() {
		app.LoadUserStyle()
		app.onQuery(entry.Text(), false)
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

func (app *Application) onEscape() {
	if app.articleView.OnEscape() {
		return
	}
	app.resetQuery()
}

func (app *Application) sendKeyEventToArticleView(event *qt.QKeyEvent) {
	// despite this, this still crashes if event has modifier
	defer func() {
		r := recover()
		if r != nil {
			slog.Error("error sending event", "r", r)
		}
	}()
	qt.QCoreApplication_SendEvent(app.articleView.Browser.QObject, event.QEvent)
}

func (app *Application) setupKeyPressEvent(widget KeyPressIface) {
	widget.OnKeyPressEvent(func(super func(event *qt.QKeyEvent), event *qt.QKeyEvent) {
		switch event.Key() {
		case int(qt.Key_Space): // " "
			app.entry.SetFocusWithReason(qt.ShortcutFocusReason)
		case int(qt.Key_Plus), int(qt.Key_Equal): // "+", "="
			app.articleView.ZoomIn()
		case int(qt.Key_Minus): // "-"
			app.articleView.ZoomOut()
		case escape: // event.Text()="\x1b"
			app.onEscape()
		case int(qt.Key_F1):
			app.ShowAbout()
		case int(qt.Key_PageUp), int(qt.Key_PageDown):
			if event.Modifiers() == 0 {
				app.sendKeyEventToArticleView(event)
			} else {
				super(event)
			}
		case int(qt.Key_Q):
			if event.Modifiers()&qt.ControlModifier > 0 {
				app.Exit()
			}
		case int(qt.Key_Left):
			if event.Modifiers()&altCtrlModifier > 0 {
				app.goBackInHistory()
			}
		case int(qt.Key_Right):
			if event.Modifiers()&altCtrlModifier > 0 {
				app.goForwardInHistory()
			}
		case int(qt.Key_Up):
			if event.Modifiers()&qt.AltModifier > 0 {
				app.resultList.GoPrevious()
			} else {
				super(event)
			}
		case int(qt.Key_Down):
			if event.Modifiers()&qt.AltModifier > 0 {
				app.resultList.GoNext()
			} else {
				super(event)
			}
		case int(qt.Key_F):
			app.favoriteButtonClicked(app.favoriteButton.ToggleChecked())
		default:
			super(event)
		}
	})
}

func (app *Application) onArticleViewKeyPressEvent(super func(event *qt.QKeyEvent), event *qt.QKeyEvent) {
	switch event.Key() {
	case int(qt.Key_Space): // " "
		app.entry.SetFocusWithReason(qt.ShortcutFocusReason)
	case escape: // event.Text()="\x1b"
		app.onEscape()
	case int(qt.Key_F1):
		app.ShowAbout()
	case int(qt.Key_Q):
		if event.Modifiers()&qt.ControlModifier > 0 {
			app.Exit()
		}
	case int(qt.Key_Left):
		if event.Modifiers()&altCtrlModifier > 0 {
			app.goBackInHistory()
		}
	case int(qt.Key_Right):
		if event.Modifiers()&altCtrlModifier > 0 {
			app.goForwardInHistory()
		}
	case int(qt.Key_Up):
		if event.Modifiers()&qt.AltModifier > 0 {
			app.resultList.GoPrevious()
		} else {
			super(event)
		}
	case int(qt.Key_Down):
		if event.Modifiers()&qt.AltModifier > 0 {
			app.resultList.GoNext()
		} else {
			super(event)
		}
	case int(qt.Key_F):
		app.favoriteButtonClicked(app.favoriteButton.ToggleChecked())
	default:
		super(event)
	}
}

func (app *Application) goBackInHistory() {
	app.queryArgs.DisableHistory = true
	defer func() {
		app.queryArgs.DisableHistory = false
	}()
	app.historyView.GoBack()
}

func (app *Application) goForwardInHistory() {
	app.queryArgs.DisableHistory = true
	defer func() {
		app.queryArgs.DisableHistory = false
	}()
	app.historyView.GoForward()
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
		app.onQuery(app.entry.Text(), false)
		return
	}

	super(event)

	// event.Modifiers(): qt.NoModifier, qt.ShiftModifier, KeypadModifier
	// also Ctrl+V should trigger SearchOnType
	if conf.SearchOnType && key < escape {
		if int(event.Modifiers())&searchOnTypeNotModifierMask == 0 {
			text := app.entry.Text()
			if len(text) >= conf.SearchOnTypeMinLength {
				app.onQuery(text, true)
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
	app.mainWindowSettingsChan <- time.Now()
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
		app.onQuery(app.entry.Text(), false)
	}
}

func (app *Application) randomEntryClicked() {
	res := dictmgr.RandomEntry(conf, resultFlags)
	if res == nil {
		return
	}
	query := res.F_Terms[0]
	app.entry.SetText(query)
	app.SetResults([]common.SearchResultIface{res})
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
	mode, ok := dictmgr.SearchModeByName(conf.RandomFavoriteSearchMode)
	if !ok {
		slog.Error("invalid random_favorite_search_mode", "value", conf.RandomFavoriteSearchMode)
		mode = dictmgr.SearchModeWordMatch
	}
	results := dictmgr.LookupHTML(term, conf, mode, resultFlags, 0)
	slog.Debug("LookupHTML running time", "dt", time.Since(t), "query", term)
	app.SetResults(results)
	if len(results) == 0 {
		app.SetNoResult(term)
	}

	queryArgs.AddHistoryAndFrequency(term)
	app.postQuery(term)
}
