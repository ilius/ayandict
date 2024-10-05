//go:build !nogui

package main

import (
	"flag"
	"log"
	"os"

	"github.com/ilius/ayandict/v2/pkg/application"
	"github.com/ilius/ayandict/v2/pkg/config"
	"github.com/ilius/ayandict/v2/pkg/dictmgr"
	"github.com/ilius/ayandict/v2/pkg/server"
)

func runServerOnly(createConfig bool) {
	conf, err := config.Load()
	if err != nil {
		panic(err)
	}
	if createConfig {
		err := config.EnsureExists(conf)
		if err != nil {
			log.Printf("Failed creating config file: %v", err)
		}
	}
	dictmgr.InitDicts(conf)
	server.StartServer(conf.LocalServerPorts[0])
}

func main() {
	noGuiFlag := flag.Bool(
		"no-gui",
		false,
		"Do not launch GUI",
	)
	createConfigFlag := flag.Bool(
		"create-config",
		false,
		"With --no-gui: create config file (with defaults) if it does not exist",
	)
	flag.Parse()

	log.SetOutput(os.Stdout)

	if *noGuiFlag {
		runServerOnly(*createConfigFlag)
		return
	}

	application.Run()
}
