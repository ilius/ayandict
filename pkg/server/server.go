package server

import (
	"log"
	"net/http"

	"github.com/ilius/ayandict/v2/pkg/appinfo"
	"github.com/ilius/ayandict/v2/pkg/qerr"
)

const (
	localhost     = "127.0.0.1"
	serverAppName = "app-name"
)

func handleGetAppName(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte(appinfo.APP_NAME))
	if err != nil {
		log.Println(err)
	}
}

func StartServer(port string) {
	http.HandleFunc("/"+serverAppName, handleGetAppName)
	log.Println("Starting local server on port", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		qerr.Error(err)
	}
}
