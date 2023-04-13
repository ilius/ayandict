package application

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/ilius/ayandict/pkg/common"
	"github.com/ilius/ayandict/pkg/config"
	"github.com/ilius/ayandict/pkg/qerr"
	"github.com/ilius/ayandict/pkg/stardict"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

const dictsJsonFilename = "dicts.json"

var (
	dictsOrder      map[string]int
	dictSettingsMap = map[string]*common.DictSettings{}
)

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

func loadDictsSettings() (map[string]*common.DictSettings, map[string]int, error) {
	order := map[string]int{}
	settingsMap := map[string]*common.DictSettings{}
	fpath := filepath.Join(config.GetConfigDir(), dictsJsonFilename)
	jsonBytes, err := ioutil.ReadFile(fpath)
	if err != nil {
		if os.IsNotExist(err) {
			return settingsMap, order, nil
		}
		return settingsMap, order, err
	}
	err = json.Unmarshal(jsonBytes, &settingsMap)
	if err != nil {
		return settingsMap, order, err
	}
	for dictName, ds := range settingsMap {
		order[dictName] = ds.Order
		if ds.Symbol == "" {
			ds.Symbol = common.DefaultSymbol(dictName)
		}
	}
	return settingsMap, order, nil
}

func saveDictsSettings(settingsMap map[string]*common.DictSettings) error {
	jsonBytes, err := json.MarshalIndent(settingsMap, "", "\t")
	if err != nil {
		return err
	}
	fpath := filepath.Join(config.GetConfigDir(), dictsJsonFilename)
	err = ioutil.WriteFile(fpath, jsonBytes, 0o644)
	if err != nil {
		return err
	}
	return nil
}

func initDicts() {
	var err error
	popup := loadingDictsPopup()
	defer popup.Destroy(true, true)

	dictSettingsMap, dictsOrder, err = loadDictsSettings()
	if err != nil {
		qerr.Errorf("Error reading dicts.json: %v", err)
	}
	infoList := stardict.Init(conf.DirectoryList, dictsOrder)
	modified := false
	for index, info := range infoList {
		dictName := info.DictName()
		ds := dictSettingsMap[dictName]
		if ds == nil {
			log.Printf("init: found new dict: %v\n", dictName)
			dictSettingsMap[dictName] = common.NewDictSettings(info, index)
			modified = true
			continue
		}
		if ds.Hash == "" {
			hash := common.Hash(info)
			if hash != "" {
				ds.Hash = hash
				modified = true
			}
		}
	}
	if modified {
		err := saveDictsSettings(dictSettingsMap)
		if err != nil {
			qerr.Error(err)
		}
	}
}

func closeDicts() {
	stardict.CloseDictFiles()
}
