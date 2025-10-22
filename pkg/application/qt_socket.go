package application

import (
	"encoding/json"
	"log/slog"
	"strings"

	"github.com/ilius/ayandict/v3/pkg/appinfo"
	"github.com/ilius/ayandict/v3/pkg/dictmgr"
	"github.com/ilius/ayandict/v3/pkg/headerlib"
	"github.com/ilius/ayandict/v3/pkg/jsonapi"
	common "github.com/ilius/go-dict-commons"
	"github.com/mappu/miqt/qt6/network"
)

const (
	s_scanPopup          = "scanpopup:"
	s_statusIconActivate = "statusicon:activate"
	s_query              = "query:"
	s_limit              = 10
)

const sockerResultFlags = uint32(
	common.ResultFlag_FixAudio |
		common.ResultFlag_FixFileSrc |
		common.ResultFlag_Web,
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
			} else if strings.HasPrefix(cmd, s_query) {
				socketQuery(conn, cmd[len(s_query):])
			} else {
				slog.Warn("server: unsupported command", "cmd", cmd)
			}
		})
	})
	return true
}

func encodeSocketResults(raw_results []common.SearchResultIface) ([]jsonapi.Result, error) {
	results := make([]jsonapi.Result, len(raw_results))
	for i, res := range raw_results {
		header, err := headerlib.GetHeader(headerTpl, res)
		if err != nil {
			return nil, err
		}
		results[i] = jsonapi.Result{
			DictName:        res.DictName(),
			Terms:           res.Terms(),
			DefinitionsHTML: res.DefinitionsHTML(),
			EntryIndex:      res.EntryIndex(),
			Score:           res.Score(),
			HeaderHTML:      header,
		}
		// entry.ResourceDir()
	}
	return results, nil
}

func socketQuery(conn *network.QLocalSocket, cmd string) {
	defer conn.Flush()
	sendError := func(msg string) {
		slog.Error("error in socketQuery", "err", msg)
		data, err := json.Marshal(jsonapi.ErrorResponse{Error: msg})
		if err != nil {
			slog.Error("error in json.Marshal", "err", err)
		} else {
			conn.Write2(data)
		}
	}
	parts := strings.Split(cmd, ":")
	if len(parts) != 2 {
		sendError("invalid number of parts")
		return
	}
	mode, ok := dictmgr.SearchModeByName(parts[0])
	if !ok {
		sendError("invalid mode")
		return
	}
	query := parts[1]
	raw_results := dictmgr.LookupHTML(query, conf, mode, sockerResultFlags, s_limit)
	// pass resultFlags to LookupHTML
	results, err := encodeSocketResults(raw_results)
	if err != nil {
		sendError(err.Error())
		return
	}
	data, err := json.Marshal(results)
	if err != nil {
		sendError(err.Error())
		return
	}
	conn.Write2(data)
}
