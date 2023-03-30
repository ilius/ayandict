package application

import (
	"fmt"

	"github.com/ilius/ayandict/pkg/stardict"
)

var dictsOrder map[string]int

func initDicts() {
	var err error
	dictSettingsMap, dictsOrder, err = loadDictsSettings()
	if err != nil {
		fmt.Println(err)
	}
	stardict.Init(conf.DirectoryList, dictsOrder)
}

func reloadDicts() {
	// do we need mutex for this?
	stardict.Init(conf.DirectoryList, dictsOrder)
}

func closeDicts() {
	stardict.CloseDictFiles()
}
