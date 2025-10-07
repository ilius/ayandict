package main

import (
	"log/slog"
	"os"

	"github.com/mappu/miqt/qt6/network"
)

const (
	s_scanpopup = "scanpopup:"
	APP_NAME    = "ayandict"
)

func main() {
	query := os.Args[1]
	client := network.NewQLocalSocket()
	slog.Info("connecting to server", "name", APP_NAME)
	client.ConnectToServerWithName(APP_NAME)
	slog.Info("waiting for connection")
	if !client.WaitForConnectedWithMsecs(200) {
		slog.Error("time out while waiting for connection")
		os.Exit(1)
	}
	slog.Info("sending data")
	_ = client.Write2([]byte(s_scanpopup + query))
	client.Flush()
}
