package application

import (
	"encoding/json"
	"log/slog"
	"os"
	"path"
	"time"

	"github.com/ilius/ayandict/v3/pkg/qtcommon/qsettings"
	"github.com/ilius/ayandict/v3/pkg/qtutils"
	qt "github.com/mappu/miqt/qt6"
)

var mainWindowPath = path.Join(stateDir, qsettings.QS_mainwindow+".json")

type MainWindowSettings struct {
	X            int  `json:"x"`
	Y            int  `json:"y"`
	Width        int  `json:"w"`
	Height       int  `json:"h"`
	Maximized    bool `json:"max"`
	SearchMode   int  `json:"searchmode"`
	ActivityType int  `json:"activitytype"`
}

func (s *MainWindowSettings) Save() {
	b, err := json.Marshal(s)
	if err != nil {
		slog.Error("error encoding main window settings", "err", err, "path", mainWindowPath)
		return
	}
	err = os.WriteFile(mainWindowPath, b, 0o644)
	if err != nil {
		slog.Error("error saving main window settings", "err", err, "path", mainWindowPath)
	}
}

func (s *MainWindowSettings) Load() {
	b, err := os.ReadFile(mainWindowPath)
	if err != nil {
		if !os.IsNotExist(err) {
			slog.Error("error reading main window settings", "err", err, "path", mainWindowPath)
		}
		return
	}
	err = json.Unmarshal(b, s)
	if err != nil {
		slog.Error("error decoding main window settings", "err", err, "path", mainWindowPath)
	}
}

func (app *Application) saveMainWindowSettings() {
	// slog.Info("Saving main window geometry")
	pos := app.window.Pos()
	size := app.window.Size()
	// what is window.SaveState()
	s := &MainWindowSettings{
		X:            pos.X(),
		Y:            pos.Y(),
		Width:        size.Width(),
		Height:       size.Height(),
		SearchMode:   app.searchModeCombo.CurrentIndex(),
		ActivityType: app.activityTypeCombo.CurrentIndex(),
	}
	s.Save()
}

func (app *Application) restoreMainWindowSettings() {
	s := &MainWindowSettings{}
	s.Load()
	window := app.window
	qtutils.SetWinPosition(window.QWidget, qt.NewQPoint2(s.X, s.Y))
	qtutils.SetWinSize(window.QWidget, qt.NewQSize2(s.Width, s.Height))
	if s.Maximized {
		window.ShowMaximized()
	}
	app.searchModeCombo.SetCurrentIndex(s.SearchMode)
	app.activityTypeCombo.SetCurrentIndex(s.ActivityType)
}

func (app *Application) setupMainWindowSettings() {
	ch := app.mainWindowSettingsChan
	app.window.OnMoveEvent(func(super func(*qt.QMoveEvent), event *qt.QMoveEvent) {
		ch <- time.Now()
	})
	app.window.OnResizeEvent(func(super func(*qt.QResizeEvent), event *qt.QResizeEvent) {
		ch <- time.Now()
	})
	// DO NO CALL app.searchModeCombo.OnCurrentIndexChanged
	// OR app.activityTypeCombo.OnCurrentIndexChanged
	go qsettings.ActionSaveLoop(ch, app.saveMainWindowSettings)
}
