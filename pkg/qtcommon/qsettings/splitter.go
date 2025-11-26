package qsettings

import (
	"time"

	qt "github.com/mappu/miqt/qt6"
)

func SetupSplitterSizesSave(splitter *qt.QSplitter, mainKey string) {
	ch := make(chan time.Time, 100)
	splitter.OnSplitterMoved(func(pos int, index int) {
		ch <- time.Now()
	})
	go ActionSaveLoop(ch, func() {
		saveSplitterSizes(splitter, mainKey)
	})
}

func saveSplitterSizes(splitter *qt.QSplitter, mainKey string) {
	saveJson(splitter.Sizes(), mainKey)
}

func RestoreSplitterSizes(splitter *qt.QSplitter, mainKey string) {
	sizes := loadJsonIntSlice(mainKey)
	if sizes == nil {
		return
	}
	splitter.SetSizes(sizes)
}
