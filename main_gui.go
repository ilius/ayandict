//go:build !nogui

package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/ilius/ayandict/v2/pkg/appinfo"
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
		slog.Warn("Web is not enabled, set web_enable = true in " + config.Path())
	}
	dictmgr.InitDicts(conf)
	server.StartServer(conf.LocalServerPorts[0])
}

func main() {
	versionFlag := flag.Bool(
		"version",
		false,
		"Show version and exit",
	)
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

	if *versionFlag {
		fmt.Printf("%v %v\n", appinfo.APP_DESC, appinfo.VERSION)
		os.Exit(0)
	}

	// slog uses stdout
	noColor := os.Getenv("NO_COLOLR") != ""
	logging.SetupGUILogger(noColor, logging.DefaultLevel)

	if *noGuiFlag {
		runServerOnly(*createConfigFlag)
		return
	}

	application.Run()
}
