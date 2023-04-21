package application

import (
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/ilius/ayandict/pkg/common"
	"github.com/ilius/ayandict/pkg/qerr"
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
		log.Println("Trying", _urlStr)
		t := time.Now()
		res, err := client.Get(_urlStr)
		// fmt.Printf("%T, %v", err, err)
		if err != nil {
			continue
		}
		if res.Body == nil {
			continue
		}
		log.Printf("%s responded in %v", _urlStr, time.Now().Sub(t))
		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			qerr.Error(err)
			continue
		}
		res.Body.Close()
		if string(data) == common.APP_NAME {
			return true, port
		}
	}
	return false, ""
}
