package application

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"unicode/utf8"

	"github.com/ilius/ayandict/pkg/config"
	"github.com/ilius/ayandict/pkg/qerr"
	"github.com/ilius/ayandict/pkg/stardict"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

const dictsJsonFilename = "dicts.json"

var dictsOrder map[string]int
var dictSettingsMap = map[string]*DictSettings{}

type DictSettings struct {
	Symbol string `json:"symbol"`
	Order  int    `json:"order"`
	Hash   string `json:"hash"`
}

func defaultDictSymbol(dictName string) string {
	symbol, _ := utf8.DecodeRune([]byte(dictName))
	return fmt.Sprintf("[%s]", string(symbol))
}

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

func loadDictsSettings() (map[string]*DictSettings, map[string]int, error) {
	order := map[string]int{}
	settingsMap := map[string]*DictSettings{}
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
			ds.Symbol = defaultDictSymbol(dictName)
		}
	}
	return settingsMap, order, nil
}

func saveDictsSettings(settingsMap map[string]*DictSettings) error {
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

func calcHashForDictName(dictName string) string {
	dic := stardict.ByDictName(dictName)
	if dic == nil {
		qerr.Errorf("could not find dictionary name %#v", dictName)
		return ""
	}
	log.Println("Calculating hash for", dictName)
	b_hash, err := dic.CalcHash()
	if err != nil {
		qerr.Error(err)
	}
	return fmt.Sprintf("%x", b_hash)
}

func setDictHash() bool {
	modified := false
	for dictName, ds := range dictSettingsMap {
		if ds.Hash != "" {
			continue
		}
		hash := calcHashForDictName(dictName)
		if hash == "" {
			continue
		}
		ds.Hash = hash
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
