//go:build !nogui

package main

import (
	"flag"
	"log/slog"
	"os"

	"github.com/ilius/ayandict/v2/pkg/application"
	"github.com/ilius/ayandict/v2/pkg/config"
	"github.com/ilius/ayandict/v2/pkg/dictmgr"
	"github.com/ilius/ayandict/v2/pkg/logging"
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
			slog.Error("Failed creating config file", "err", err)
		}
	}
	if !conf.WebEnable {
		slog.Warn("Web is not enabled, set web_enable = true in config.toml file")
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

	// slog uses stdout
	noColor := os.Getenv("NO_COLOLR") != ""
	logging.SetupLogger(noColor, logging.DefaultLevel)

	if *noGuiFlag {
		runServerOnly(*createConfigFlag)
		return
	}

	application.Run()
}
