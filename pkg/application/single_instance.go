package application

import (
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/ilius/ayandict/pkg/qerr"
)

const (
	localhost     = "127.0.0.1"
	serverAppName = "app-name"
)

var client = http.Client{
	Timeout: 100 * time.Millisecond,
}

func isSingleInstanceRunning(appName string, ports []string) bool {
	ok, _ := findLocalServer(ports)
	return ok
}

func handleGetAppName(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(APP_NAME))
}

func startSingleInstanceServer(appName string, port string) {
	http.HandleFunc("/"+serverAppName, handleGetAppName)
	log.Println("Starting local server on port", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		qerr.Error(err)
	}
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
		log.Printf("%s responsed in %v", _urlStr, time.Now().Sub(t))
		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			qerr.Error(err)
			continue
		}
		res.Body.Close()
		if string(data) == APP_NAME {
			return true, port
		}
	}
	return false, ""
}
