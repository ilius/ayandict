package application

import (
	common "codeberg.org/ilius/go-dict-commons"
	qt "github.com/mappu/miqt/qt6"
)

type KeyPressIface interface {
	OnKeyPressEvent(func(func(event *qt.QKeyEvent), *qt.QKeyEvent))
}

type ResultsIface interface {
	Clear()
	QWidget() *qt.QWidget
	SetResults(results []common.SearchResultIface)
	Active() common.SearchResultIface
	ReloadList()
	OnKeyPressEvent(func(func(event *qt.QKeyEvent), *qt.QKeyEvent))
	GoPrevious()
	GoNext()
}

type FavoriteButtonInterface interface {
	SetChecked(bool)
	ToggleChecked() bool
	SetToolTips(string, string)
	QWidget() *qt.QWidget
	Hide()
	Show()
	SetDisabled(bool)
	SetFixedSize2(int, int)
}
