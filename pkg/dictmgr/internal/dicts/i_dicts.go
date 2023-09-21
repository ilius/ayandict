package dicts

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/ilius/ayandict/v2/pkg/config"
	"github.com/ilius/ayandict/v2/pkg/qerr"
	common "github.com/ilius/go-dict-commons"
	"github.com/ilius/go-stardict/v2"
)

const dictsJsonFilename = "dicts.json"

var (
	DictList        []common.Dictionary
	DictByName      = map[string]common.Dictionary{}
	DictsOrder      map[string]int
	DictSettingsMap = map[string]*DictionarySettings{}
)

var sqldictOpen = func([]string, map[string]int) []common.Dictionary {
	return nil
}

func init() {
	stardict.ErrorHandler = func(err error) {
		qerr.Error(err)
	}
}

func absInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

type DictionaryListSorter struct {
	Order map[string]int
	List  []common.Dictionary
}

func (s DictionaryListSorter) Len() int {
	return len(s.List)
}

func (s DictionaryListSorter) Swap(i, j int) {
	s.List[i], s.List[j] = s.List[j], s.List[i]
}

func (s DictionaryListSorter) Less(i, j int) bool {
	return absInt(s.Order[s.List[i].DictName()]) < absInt(s.Order[s.List[j].DictName()])
}

func Reorder(order map[string]int) {
	sort.Sort(DictionaryListSorter{
		List:  DictList,
		Order: order,
	})
}

func loadDictsSettings() (
	map[string]*DictionarySettings,
	map[string]int,
	error,
) {
	order := map[string]int{}
	settingsMap := map[string]*DictionarySettings{}
	fpath := filepath.Join(config.GetConfigDir(), dictsJsonFilename)
	jsonBytes, err := os.ReadFile(fpath)
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

func SaveDictsSettings(settingsMap map[string]*DictionarySettings) error {
	jsonBytes, err := json.MarshalIndent(settingsMap, "", "\t")
	if err != nil {
		return err
	}
	fpath := filepath.Join(config.GetConfigDir(), dictsJsonFilename)
	err = os.WriteFile(fpath, jsonBytes, 0o644)
	if err != nil {
		return err
	}
	return nil
}

func getDictNameByHashMap() map[string][]string {
	byHash := map[string][]string{}
	for dictName, ds := range DictSettingsMap {
		if ds.Hash == "" {
			continue
		}
		byHash[ds.Hash] = append(byHash[ds.Hash], dictName)
	}
	return byHash
}

func InitDicts(conf *config.Config) {
	var err error
	DictSettingsMap, DictsOrder, err = loadDictsSettings()
	if err != nil {
		qerr.Errorf("Error reading dicts.json: %v", err)
	}

	t := time.Now()
	DictList, err = stardict.Open(conf.DirectoryList, DictsOrder)
	if err != nil {
		panic(err)
	}
	if len(conf.SqlDictList) > 0 {
		DictList = append(DictList, sqldictOpen(conf.SqlDictList, DictsOrder)...)
	}

	// to support another format, you can call pkg.Open just like stardict
	// and append them new dicList to this dicList. since we are sorting them
	// here in Reorder after loading all dictionaries

	for _, dic := range DictList {
		DictByName[dic.DictName()] = dic
	}

	log.Println("Loading dictionaries took", time.Since(t))

	nameByHash := getDictNameByHashMap()

	newDictSettings := func(dic common.Dictionary, index int) *DictionarySettings {
		hash := Hash(dic)
		if hash != "" {
			prevNames := nameByHash[hash]
			if len(prevNames) > 0 {
				log.Println("init: found renamed dicts:", prevNames)
				prevName := prevNames[0]
				ds := DictSettingsMap[prevName]
				delete(DictSettingsMap, prevName)
				return ds
			}
		}
		return &DictionarySettings{
			Symbol: common.DefaultSymbol(dic.DictName()),
			Order:  index,
			Hash:   hash,
		}
	}

	modified := false
	for index, dic := range DictList {
		dictName := dic.DictName()
		ds := DictSettingsMap[dictName]
		if ds == nil {
			log.Printf("init: found new dict: %v\n", dictName)
			ds = newDictSettings(dic, index)
			DictSettingsMap[dictName] = ds
			DictsOrder[dictName] = ds.Order
			modified = true
			continue
		}
		if ds.Hash == "" {
			hash := Hash(dic)
			if hash != "" {
				ds.Hash = hash
				modified = true
			}
		}
	}
	if modified {
		err := SaveDictsSettings(DictSettingsMap)
		if err != nil {
			qerr.Error(err)
		}
	}

	Reorder(DictsOrder)
}
