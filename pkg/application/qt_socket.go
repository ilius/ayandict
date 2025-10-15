package application

import (
	"log/slog"
	"strings"

	"github.com/ilius/ayandict/v3/pkg/appinfo"
	"github.com/mappu/miqt/qt6/network"
)

const (
	s_scanPopup          = "scanpopup:"
	s_statusIconActivate = "statusicon:activate"
)

func (app *Application) startLocalSocketServer() bool {
	server := network.NewQLocalServer()
	if !server.Listen(appinfo.LOCAL_SOCKET_NAME) {
		return false
	}
	server.OnNewConnection(func() {
		conn := server.NextPendingConnection()
		conn.OnReadyRead(func() {
			cmd := strings.TrimSpace(string(conn.ReadAll()))
			slog.Debug("LocalSocketServer: received:", "data", cmd)
			if strings.HasPrefix(cmd, s_scanPopup) {
				if conf.ScanPopupAPI {
					query := cmd[len(s_scanPopup):]
					app.scanPopup(query)
				}
			} else if cmd == s_statusIconActivate {
				app.onStatusIconClick()
			} else {
				slog.Warn("server: unsupported command", "cmd", cmd)
			}
		})
	})
	return true
}
