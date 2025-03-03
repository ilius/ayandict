package application

import (
	"fmt"
	"log/slog"
	"net/url"

	common "github.com/ilius/go-dict-commons"
)

func openEntryInWebInterface(result common.SearchResultIface) {
	if len(conf.LocalServerPorts) == 0 {
		slog.Error("openEntryInWeb: LocalServerPorts is empty")
		return
	}
	// result.DictName()
	// result.EntryIndex()
	port := conf.LocalServerPorts[0]
	addr := url.URL{
		Scheme: "http",
		Host:   "127.0.0.1:" + port,
		// Path: ,
	}
	addrStr := addr.String()
	fmt.Println(addrStr)
}

func openResultInWeb(result common.SearchResultIface) {
	if conf.WebEnable {
		openEntryInWebInterface(result)
		return
	}
	// TODO: save article to tmp file and open it in browser
}
