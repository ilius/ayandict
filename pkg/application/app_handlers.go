package application

import (
	"log/slog"
	"time"

	common "codeberg.org/ilius/go-dict-commons"
	"github.com/ilius/ayandict/v3/pkg/dictmgr"
	"github.com/ilius/ayandict/v3/pkg/dictmgr/qdictmgr"
	qt "github.com/mappu/miqt/qt6"
)

func (app *Application) addShortcut(keyDesc string, slot func()) {
	qt.NewQShortcut2(qt.NewQKeySequence2(keyDesc), app.window.QObject).OnActivated(slot)
}

func (app *Application) setupHandlers() {
	// MUST not call OnKeyPressEvent multiple times on the same widget

	app.addShortcut("Escape", app.onEscape)
	app.addShortcut("Space", app.onSpaceKey)

	app.addShortcut("+", app.articleView.ZoomIn)
	app.addShortcut("=", app.articleView.ZoomIn)
	app.addShortcut("-", app.articleView.ZoomOut)

	app.addShortcut("F", app.onFKey)
	app.addShortcut("F1", app.ShowAbout)
	app.addShortcut("Ctrl+Q", app.Exit)
	app.addShortcut("Ctrl+D", app.dictsButtonClicked)
	app.addShortcut("Ctrl+R", app.reloadConfigByShortcut)
	app.addShortcut("Ctrl+Del", app.clearHistoryClicked)

	app.addShortcut("Alt+Left", app.goBackInHistory)
	app.addShortcut("Ctrl+Left", app.goBackInHistory)

	app.addShortcut("Alt+Right", app.goForwardInHistory)
	app.addShortcut("Ctrl+Right", app.goForwardInHistory)

	app.addShortcut("Alt+Up", app.resultList.GoPrevious)
	app.addShortcut("Alt+Down", app.resultList.GoNext)

	app.setupKeyPressEvent(app.window)
	app.setupKeyPressEvent(app.resultList)
	app.setupKeyPressEvent(app.historyView.QListWidget)

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

	app.openConfigButton.OnClicked(OpenConfig)
	app.reloadButton.OnClicked(app.reloadButtonClicked)
	app.saveHistoryButton.OnClicked(func() {
		app.historyView.Save()
		frequencyTable.SaveNoError()
	})
	app.clearHistoryButton.OnClicked(app.clearHistoryClicked)
	app.clearButton.OnClicked(app.clearButtonClicked)
	app.dictsButton.OnClicked(app.dictsButtonClicked)
	app.randomEntryButton.OnClicked(app.randomEntryClicked)
	app.randomFavoriteButton.OnClicked(app.randomFavoriteClicked)
	entry.OnKeyPressEvent(app.onEntryKeyPress)

	// slog.Error("test error", "s", "hello", "n", 2, "b", true)
}

func (app *Application) onEscape() {
	if app.articleView.OnEscape() {
		return
	}
	if app.queryArgs.Entry.HasFocus() {
		app.window.SetFocus()
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
	widget.OnKeyPressEvent(func(super func(*qt.QKeyEvent), event *qt.QKeyEvent) {
		switch event.Key() {
		case int(qt.Key_PageUp), int(qt.Key_PageDown):
			if event.Modifiers() == 0 {
				app.sendKeyEventToArticleView(event)
			} else {
				super(event)
			}
		default:
			super(event)
		}
	})
}

func (app *Application) onSpaceKey() {
	app.entry.SetFocusWithReason(qt.ShortcutFocusReason)
}

func (app *Application) onFKey() {
	app.favoriteButtonClicked(app.favoriteButton.ToggleChecked())
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

func (app *Application) clearButtonClicked() {
	if qt.QGuiApplication_KeyboardModifiers()&qt.ControlModifier != 0 {
		app.clearHistoryClicked()
	}
	app.resetQuery()
}

func (app *Application) reloadConfigByShortcut() {
	app.reloadConfig()
	app.onQuery(app.entry.Text(), false)
}

func (app *Application) reloadButtonClicked() {
	app.reloadConfig()
	modifiers := qt.QGuiApplication_KeyboardModifiers()
	if modifiers == 0 {
		return
	}
	if modifiers&qt.ControlModifier != 0 {
		slog.Info("Reloading dictionaries")
		qdictmgr.InitDicts(conf, true)
		app.dictManager = nil
	}
	if modifiers&(qt.ControlModifier|qt.ShiftModifier) != 0 {
		slog.Info("Reloading user style")
		app.LoadUserStyle()
	}
	app.onQuery(app.entry.Text(), false)
}
