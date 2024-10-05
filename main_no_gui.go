//go:build nogui

package main

import (
	"flag"
	"log"
	"os"

	"github.com/ilius/ayandict/v2/pkg/config"
	"github.com/ilius/ayandict/v2/pkg/dictmgr"
	"github.com/ilius/ayandict/v2/pkg/server"
)

func main() {
	log.SetOutput(os.Stdout)

	createConfigFlag := flag.Bool(
		"create-config",
		false,
		"Create config file (with defaults) if it does not exist",
	)
	flag.Parse()

	conf, err := config.Load()
	if err != nil {
		panic(err)
	}

	if *createConfigFlag {
		err := config.EnsureExists(conf)
		if err != nil {
			log.Printf("Failed creating config file: %v", err)
		}
	}

	dictmgr.InitDicts(conf)
	server.StartServer(conf.LocalServerPorts[0])
}
