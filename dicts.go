package main

import "github.com/ilius/ayandict/pkg/stardict"

func reloadDicts() {
	// do we need mutex for this?
	stardict.Init(conf.DirectoryList)
}
