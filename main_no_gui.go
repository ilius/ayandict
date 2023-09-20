//go:build nogui

package main

import (
	"log"
	"os"

	"github.com/ilius/ayandict/v2/pkg/config"
	"github.com/ilius/ayandict/v2/pkg/dictmgr"
	"github.com/ilius/ayandict/v2/pkg/server"
)

func main() {
	log.SetOutput(os.Stdout)
	conf, err := config.Load()
	if err != nil {
		panic(err)
	}
	dictmgr.InitDicts(conf, false)
	server.StartServer(conf.LocalServerPorts[0])
}
