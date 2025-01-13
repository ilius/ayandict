//go:build nogui

package main

import (
	"flag"
	"log/slog"
	"os"

	"github.com/ilius/ayandict/v2/pkg/config"
	"github.com/ilius/ayandict/v2/pkg/dictmgr"
	"github.com/ilius/ayandict/v2/pkg/logging"
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

	noColor := os.Getenv("NO_COLOLR") != ""
	logging.SetupLogger(noColor, logging.DefaultLevel)

	if *createConfigFlag {
		err := config.EnsureExists(conf)
		if err != nil {
			slog.Error("Failed creating config file", "err", err)
		}
	}
	if !conf.WebEnable {
		slog.Warn("Web is not enabled, set web_enable = true in " + config.Path())
	}

	dictmgr.InitDicts(conf)
	server.StartServer(conf.LocalServerPorts[0])
}
