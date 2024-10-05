package application

import (
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/ilius/ayandict/v2/pkg/appinfo"
	"github.com/ilius/ayandict/v2/pkg/qtcommon/qerr"
)

const (
	localhost     = "127.0.0.1"
	serverAppName = "app-name"
)

var client = &http.Client{
	Timeout: 100 * time.Millisecond,
}

func findLocalServer(ports []string) (bool, string) {
	for _, port := range ports {
		_url := &url.URL{
			Scheme: "http",
			Host:   net.JoinHostPort(localhost, port),
			Path:   serverAppName,
		}
		_urlStr := _url.String()
		slog.Debug("findLocalServer, trying " + _urlStr)
		t := time.Now()
		res, err := client.Get(_urlStr)
		if err != nil {
			continue
		}
		if res.Body == nil {
			continue
		}
		slog.Debug("local server responded", "url", _urlStr, "dt", time.Since(t))
		data, err := io.ReadAll(res.Body)
		if err != nil {
			qerr.Error(err)
			continue
		}
		res.Body.Close()
		if string(data) == appinfo.APP_NAME {
			return true, port
		}
	}
	return false, ""
}
