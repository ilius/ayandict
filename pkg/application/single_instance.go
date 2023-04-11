package application

import (
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"time"
)

const (
	localhost     = "127.0.0.1"
	serverAppName = "app-name"
)

var client = http.Client{
	Timeout: 100 * time.Millisecond,
}

var localPort = "8350"
var localTryPorts = []string{
	"8350", "8351", "8352",
}

func isSingleInstanceRunning(appName string) bool {
	ok, port := findLocalServer(localTryPorts)
	if ok {
		localPort = port
	}
	return ok
}

func handleGetAppName(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(APP_NAME))
}

func startSingleInstanceServer(appName string) {
	http.HandleFunc("/"+serverAppName, handleGetAppName)
	err := http.ListenAndServe(":"+localPort, nil)
	if err != nil {
		log.Println(err)
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
			log.Println(err)
			continue
		}
		res.Body.Close()
		if string(data) == APP_NAME {
			return true, port
		}
	}
	return false, ""
}
