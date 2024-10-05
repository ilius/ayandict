//go:build nogui

package main

import (
	"flag"
	"log/slog"

	"github.com/ilius/ayandict/v2/pkg/config"
	"github.com/ilius/ayandict/v2/pkg/dictmgr"
	"github.com/ilius/ayandict/v2/pkg/server"
)

func main() {
	// slog uses stdout

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
			slog.Error("Failed creating config file", "err", err)
		}
	}

	dictmgr.InitDicts(conf)
	server.StartServer(conf.LocalServerPorts[0])
}
