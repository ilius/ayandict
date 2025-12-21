package qdictmgr

import (
	qt "github.com/mappu/miqt/qt6"
)

func makeCheckMarkBig(check *qt.QCheckBox) {
	check.SetStyleSheet(`
	QCheckBox {
		spacing: 0.2em;
	}
	QCheckBox::indicator {
		width: 1em;
		height: 1em;
	}
	`)
}
