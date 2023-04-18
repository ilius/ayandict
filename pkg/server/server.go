package server

import (
	"log"
	"net/http"

	"github.com/ilius/ayandict/pkg/common"
	"github.com/ilius/ayandict/pkg/qerr"
)

const (
	localhost     = "127.0.0.1"
	serverAppName = "app-name"
)

func handleGetAppName(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(common.APP_NAME))
}

func StartServer(port string) {
	http.HandleFunc("/"+serverAppName, handleGetAppName)
	log.Println("Starting local server on port", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		qerr.Error(err)
	}
}
