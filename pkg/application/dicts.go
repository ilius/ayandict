package application

import (
	"fmt"
	"log"

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

func setDictHash() bool {
	modified := false
	for dictName, ds := range dictSettingsMap {
		if ds.Hash != "" {
			continue
		}
		dic := stardict.ByDictName(dictName)
		if dic == nil {
			qerr.Errorf("could not find dictionary name %#v", dictName)
			continue
		}
		log.Println("Calculating hash for", dictName)
		b_hash, err := dic.CalcHash()
		if err != nil {
			qerr.Error(err)
		}
		ds.Hash = fmt.Sprintf("%x", b_hash)
		modified = true
	}
	return modified
}

func initDicts() {
	var err error
	popup := loadingDictsPopup()
	defer popup.Destroy(true, true)

	dictSettingsMap, dictsOrder, err = loadDictsSettings()
	if err != nil {
		qerr.Errorf("Error reading dicts.json: %v", err)
	}
	stardict.Init(conf.DirectoryList, dictsOrder)
	if setDictHash() {
		err := saveDictsSettings(dictSettingsMap)
		if err != nil {
			qerr.Error(err)
		}
	}
}

func reloadDicts() {
	// do we need mutex for this?
	popup := loadingDictsPopup()
	stardict.Init(conf.DirectoryList, dictsOrder)
	popup.Destroy(true, true)
}

func closeDicts() {
	stardict.CloseDictFiles()
}
