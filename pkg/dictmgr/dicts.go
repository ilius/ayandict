package dictmgr

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/ilius/ayandict/pkg/common"
	"github.com/ilius/ayandict/pkg/config"
	"github.com/ilius/ayandict/pkg/qerr"
	"github.com/ilius/ayandict/pkg/stardict"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

const dictsJsonFilename = "dicts.json"

var (
	dicList         []common.Dictionary
	dicMap          = map[string]common.Dictionary{}
	dictsOrder      map[string]int
	dictSettingsMap = map[string]*common.DictSettings{}
)

func absInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

type DicListSorter struct {
	Order map[string]int
	List  []common.Dictionary
}

func (s DicListSorter) Len() int {
	return len(s.List)
}

func (s DicListSorter) Swap(i, j int) {
	s.List[i], s.List[j] = s.List[j], s.List[i]
}

func (s DicListSorter) Less(i, j int) bool {
	return absInt(s.Order[s.List[i].DictName()]) < absInt(s.Order[s.List[j].DictName()])
}

func Reorder(order map[string]int) {
	sort.Sort(DicListSorter{
		List:  dicList,
		Order: order,
	})
}

func loadingDictsPopup(conf *config.Config) *widgets.QLabel {
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

func InitDicts(conf *config.Config) {
	var err error
	popup := loadingDictsPopup(conf)
	defer popup.Destroy(true, true)

	dictSettingsMap, dictsOrder, err = loadDictsSettings()
	if err != nil {
		qerr.Errorf("Error reading dicts.json: %v", err)
	}

	t := time.Now()
	dicList, err = stardict.Open(conf.DirectoryList, dictsOrder)
	if err != nil {
		panic(err)
	}

	// to support another format, you can call pkg.Open just like stardict
	// and append them new dicList to this dicList. since we are sorting them
	// here in Reorder after loading all dictionaries

	log.Println("Loading dictionaries took", time.Now().Sub(t))

	Reorder(dictsOrder)

	modified := false
	for index, info := range dicList {
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

func CloseDicts() {
	for _, dic := range dicList {
		if dic.Disabled() {
			continue
		}
		dic.Close()
	}
}

func DictSymbol(dictName string) string {
	ds := dictSettingsMap[dictName]
	if ds == nil {
		return ""
	}
	return ds.Symbol
}
