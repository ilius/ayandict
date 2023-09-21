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

func runServerOnly() {
	conf, err := config.Load()
	if err != nil {
		panic(err)
	}
	dictmgr.InitDicts(conf, false)
	server.StartServer(conf.LocalServerPorts[0])
}

func main() {
	noGuiFlag := flag.Bool(
		"no-gui",
		false,
		"Do not launch GUI",
	)
	flag.Parse()

	log.SetOutput(os.Stdout)

	if *noGuiFlag {
		runServerOnly()
		return
	}

	application.Run()
}
