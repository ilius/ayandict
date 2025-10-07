package main

import (
	"log/slog"
	"os"

	"github.com/ilius/ayandict/v3/pkg/appinfo"
	"github.com/mappu/miqt/qt6/network"
)

const (
	s_scanpopup = "scanpopup:"
)

func main() {
	query := os.Args[1]
	client := network.NewQLocalSocket()
	slog.Info("connecting to server", "name", appinfo.LOCAL_SOCKET_NAME)
	client.ConnectToServerWithName(appinfo.LOCAL_SOCKET_NAME)
	slog.Info("waiting for connection")
	if !client.WaitForConnectedWithMsecs(200) {
		slog.Error("time out while waiting for connection")
		os.Exit(1)
	}
	slog.Info("writing data", "query", query)
	_ = client.Write2([]byte(s_scanpopup + query))
	slog.Info("flushing")
	client.Flush()
}
