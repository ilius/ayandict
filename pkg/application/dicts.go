package application

import (
	"log"

	"github.com/ilius/ayandict/pkg/stardict"
)

var dictsOrder map[string]int

func initDicts() {
	var err error
	dictSettingsMap, dictsOrder, err = loadDictsSettings()
	if err != nil {
		log.Println(err)
	}
	stardict.Init(conf.DirectoryList, dictsOrder)
}

func reloadDicts() {
	// do we need mutex for this?
	stardict.Init(conf.DirectoryList, dictsOrder)
	dictManager = nil
}

func closeDicts() {
	stardict.CloseDictFiles()
}
