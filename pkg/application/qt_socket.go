package application

import (
	"log/slog"
	"strings"

	"github.com/ilius/ayandict/v3/pkg/appinfo"
	"github.com/mappu/miqt/qt6/network"
)

const (
	s_scanpopup = "scanpopup:"
)

func (app *Application) startLocalSocketServer() bool {
	server := network.NewQLocalServer()
	if !server.Listen(appinfo.APP_NAME) {
		return false
	}
	server.OnNewConnection(func() {
		conn := server.NextPendingConnection()
		conn.OnReadyRead(func() {
			data := string(conn.ReadAll())
			slog.Debug("LocalSocketServer: received:", "data", data)
			if strings.HasPrefix(data, s_scanpopup) {
				if conf.ScanPopupAPI {
					query := data[len(s_scanpopup):]
					app.scanPopup(query)
				}
			}
		})
	})
	return true
}
