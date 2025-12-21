package application

import (
	common "codeberg.org/ilius/go-dict-commons"
	qt "github.com/mappu/miqt/qt6"
)

type KeyPressIface interface {
	OnKeyPressEvent(func(func(*qt.QKeyEvent), *qt.QKeyEvent))
}

type ResultsIface interface {
	Clear()
	QWidget() *qt.QWidget
	SetResults([]common.SearchResultIface)
	CurrentResult() common.SearchResultIface
	SetCurrentResult(int)
	Reload()
	OnKeyPressEvent(func(func(*qt.QKeyEvent), *qt.QKeyEvent))
	GoPrevious()
	GoNext()
}

type FavoriteButtonInterface interface {
	SetChecked(bool)
	SetTerms(terms []string)
	ToggleChecked() bool
	SetToolTips(inactive string, active string)
	QWidget() *qt.QWidget
	Hide()
	Show()
	SetDisabled(bool)
	SetFixedSize2(int, int)
}
