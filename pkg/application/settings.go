package application

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

func reastoreSetting(qsettings *core.QSettings, key string, apply func(*core.QVariant)) {
	if !qsettings.Contains(key) {
		return
	}
	apply(qsettings.Value(key, core.NewQVariant1(nil)))
}

func reastoreBoolSetting(
	qsettings *core.QSettings,
	key string, _default bool,
	apply func(bool),
) {
	if !qsettings.Contains(key) {
		apply(_default)
		return
	}
	apply(qsettings.Value(key, core.NewQVariant1(nil)).ToBool())
}

func saveMainWinGeometry(qsettings *core.QSettings, window *widgets.QMainWindow) {
	qsettings.BeginGroup(QS_mainwindow)

	qsettings.SetValue(QS_geometry, core.NewQVariant13(window.SaveGeometry()))
	qsettings.SetValue(QS_savestate, core.NewQVariant13(window.SaveState(0)))
	qsettings.SetValue(QS_maximized, core.NewQVariant9(window.IsMaximized()))
	if !window.IsMaximized() {
		qsettings.SetValue(QS_pos, core.NewQVariant27(window.Pos()))
		qsettings.SetValue(QS_size, core.NewQVariant25(window.Size()))
	}

	qsettings.EndGroup()
}

func reastoreMainWinGeometry(qsettings *core.QSettings, window *widgets.QMainWindow) {
	qsettings.BeginGroup(QS_mainwindow)

	reastoreSetting(qsettings, QS_geometry, func(value *core.QVariant) {
		window.RestoreGeometry(value.ToByteArray())
	})
	reastoreSetting(qsettings, QS_savestate, func(value *core.QVariant) {
		window.RestoreState(value.ToByteArray(), 0)
	})
	reastoreBoolSetting(qsettings, QS_maximized, false, func(maximized bool) {
		if maximized {
			window.ShowMaximized()
			return
		}
		reastoreSetting(qsettings, QS_pos, func(value *core.QVariant) {
			window.Move(value.ToPoint())
		})
		reastoreSetting(qsettings, QS_size, func(value *core.QVariant) {
			window.Resize(value.ToSize())
		})
	})

	qsettings.EndGroup()
}
