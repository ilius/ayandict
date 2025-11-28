package webclient

import (
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/ilius/ayandict/v3/pkg/appinfo"
	"github.com/ilius/ayandict/v3/pkg/config"
)

const (
	localhost        = "127.0.0.1"
	webServerAppName = "app-name"
)

var client = &http.Client{
	Timeout: 100 * time.Millisecond,
}

func Init(conf *config.Config) {
	client.Timeout = conf.LocalClientTimeout
}

func FindLocalWebServer(ports []string) (bool, string) {
	for _, port := range ports {
		_url := &url.URL{
			Scheme: "http",
			Host:   net.JoinHostPort(localhost, port),
			Path:   webServerAppName,
		}
		_urlStr := _url.String()
		slog.Debug("findLocalWebServer, trying " + _urlStr)
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
			slog.Error("error in findLocalWebServer while reading response body: " + err.Error())
			continue
		}
		res.Body.Close()
		if string(data) == appinfo.APP_NAME {
			return true, port
		}
	}
	return false, ""
}
