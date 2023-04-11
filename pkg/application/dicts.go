package application

import (
	"github.com/ilius/ayandict/pkg/qerr"
	"github.com/ilius/ayandict/pkg/stardict"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

var dictsOrder map[string]int

func loadingDictsPopup() *widgets.QLabel {
	popup := widgets.NewQLabel2(
		`<span style="font-size:xx-large;">Loading dictionaries</span>`,
		nil,
		core.Qt__SplashScreen,
	)
	// Qt__SplashScreen makes it centered on screen
	popup.SetFrameStyle(int(widgets.QFrame__Raised | widgets.QFrame__Shadow(widgets.QFrame__Panel)))
	popup.SetAlignment(core.Qt__AlignCenter)
	popup.SetFixedSize2(300, 100)
	popup.SetWindowModality(core.Qt__WindowModal)
	popup.SetStyleSheet(conf.PopupStyleStr)
	popup.Show()
	core.QCoreApplication_ProcessEvents(core.QEventLoop__AllEvents)
	return popup
}

func initDicts() {
	var err error
	popup := loadingDictsPopup()
	dictSettingsMap, dictsOrder, err = loadDictsSettings()
	if err != nil {
		qerr.Errorf("Error reading dicts.json: %v", err)
	}
	stardict.Init(conf.DirectoryList, dictsOrder)
	popup.Destroy(true, true)
}

func reloadDicts() {
	// do we need mutex for this?
	popup := loadingDictsPopup()
	stardict.Init(conf.DirectoryList, dictsOrder)
	popup.Destroy(true, true)
	dictManager = nil
}

func closeDicts() {
	stardict.CloseDictFiles()
}
