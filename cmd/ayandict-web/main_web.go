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
	"github.com/ilius/ayandict/v3/pkg/webserver"
)

func main() {
	// slog uses stdout

	portFlag := flag.String(
		"port",
		"",
		"Web port (default read from config)",
	)
	exposeFlag := flag.Bool(
		"expose",
		false,
		"Expose web service & web app to outside (otherwise only available to 127.0.0.1)",
	)
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
		fmt.Printf("%v %v (web build)\n", appinfo.APP_DESC, appinfo.VERSION)
		os.Exit(0)
	}

	conf, err := config.Load()
	if err != nil {
		panic(err)
	}

	noColor := os.Getenv("NO_COLOR") != ""
	handler := logging.NewColoredHandler(noColor, logging.DefaultLevel)
	slog.SetDefault(slog.New(handler))

	if *createConfigFlag {
		err := config.EnsureExists(conf)
		if err != nil {
			slog.Error("Failed creating config file", "err", err)
		}
	}
	conf.WebEnable = true
	conf.WebExpose = *exposeFlag
	if *portFlag != "" {
		conf.LocalServerPorts = []string{*portFlag}
	}
	slog.Info("Web ports", "ports", conf.LocalServerPorts)

	dictmgr.InitDicts(conf)
	webserver.StartServer(conf.LocalServerPorts[0])
}
