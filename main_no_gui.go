//go:build nogui

package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/ilius/ayandict/v3/pkg/appinfo"
	"github.com/ilius/ayandict/v3/pkg/config"
	"github.com/ilius/ayandict/v3/pkg/dictmgr"
	"github.com/ilius/ayandict/v3/pkg/logging"
	"github.com/ilius/ayandict/v3/pkg/server"
)

func main() {
	// slog uses stdout

	versionFlag := flag.Bool(
		"version",
		false,
		"Show version and exit",
	)
	createConfigFlag := flag.Bool(
		"create-config",
		false,
		"Create config file (with defaults) if it does not exist",
	)
	flag.Parse()

	if *versionFlag {
		fmt.Printf("%v %v (non-GUI build)\n", appinfo.APP_DESC, appinfo.VERSION)
		os.Exit(0)
	}

	conf, err := config.Load()
	if err != nil {
		panic(err)
	}

	noColor := os.Getenv("NO_COLOLR") != ""
	handler := logging.NewColoredHandler(noColor, logging.DefaultLevel)
	slog.SetDefault(slog.New(handler))

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
