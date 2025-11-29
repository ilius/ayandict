package application

import (
	"github.com/ilius/ayandict/v3/pkg/config"
	qt "github.com/mappu/miqt/qt6"
)

const expanding = qt.QSizePolicy__Expanding

const altCtrlModifier = qt.AltModifier | qt.ControlModifier

const (
	QS_mainSplitter   = "main_splitter"
	QS_frequencyTable = "frequencytable"

	escape = int(qt.Key_Escape)

	searchOnTypeNotModifierMask = int(qt.AltModifier) | int(qt.MetaModifier)
)

var stateDir = config.Paths.StateDir()
