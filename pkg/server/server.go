package server

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ilius/ayandict/v2/pkg/appinfo"
	"github.com/ilius/ayandict/v2/pkg/config"
	"github.com/ilius/ayandict/v2/pkg/dictmgr"
	"github.com/ilius/ayandict/v2/pkg/qerr"
	"github.com/ilius/ayandict/v2/web"
	common "github.com/ilius/go-dict-commons"
)

const (
	localhost    = "127.0.0.1"
	path_appName = "app-name"
	path_query   = "query"
)

var conf = config.MustLoad()

const resultFlags = common.ResultFlag_Web |
	common.ResultFlag_FixAudio |
	common.ResultFlag_FixFileSrc

// 	common.ResultFlag_ColorMapping)

type ErrorResponse struct {
	Error string `json:"error"`
}

type Result struct {
	DictName        string   `json:"dictName"`
	Terms           []string `json:"terms"`
	DefinitionsHTML []string `json:"definitionsHTML"`
	EntryIndex      uint64   `json:"entryIndex"`
	Score           uint8    `json:"score"`
	// ResourceDir string
}

func writeMsg(w http.ResponseWriter, msg string) {
	_, err := w.Write([]byte(msg))
	if err != nil {
		log.Println("error in Write:", err)
	}
}

func getAppName(w http.ResponseWriter, r *http.Request) {
	writeMsg(w, appinfo.APP_NAME)
}

func query(w http.ResponseWriter, r *http.Request) {
	t := time.Now()
	mode := dictmgr.QueryModeFuzzy

	jsonEncoder := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	query := r.FormValue("query")
	if query == "" {
		err := jsonEncoder.Encode(ErrorResponse{Error: "missing query"})
		if err != nil {
			log.Println(err)
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// mode = dictmgr.QueryModeStartWith
	// mode = dictmgr.QueryModeRegex
	// mode = dictmgr.QueryModeGlob

	raw_results := dictmgr.LookupHTML(query, conf, mode, resultFlags)
	// pass resultFlags to LookupHTML
	results := make([]Result, len(raw_results))
	for i, rr := range raw_results {
		results[i] = Result{
			DictName:        rr.DictName(),
			Terms:           rr.Terms(),
			DefinitionsHTML: rr.DefinitionsHTML(),
			EntryIndex:      rr.EntryIndex(),
			Score:           rr.Score(),
		}
		// rr.ResourceDir()
	}
	log.Printf("LookupHTML took %v for %#v", time.Since(t), query)
	err := jsonEncoder.Encode(results)
	if err != nil {
		log.Println(err)
		err2 := jsonEncoder.Encode(ErrorResponse{Error: err.Error()})
		if err2 != nil {
			log.Println(err2)
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	file, err := web.FS.Open("web/index.html")
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	content := file.(io.ReadSeeker)
	http.ServeContent(w, r, "", time.Now(), content)
}

func dictRes(w http.ResponseWriter, r *http.Request) {
	dictName := r.FormValue("dictName")
	path := r.FormValue("path")
	if dictName == "" {
		writeMsg(w, "missing dictName")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if path == "" {
		writeMsg(w, "missing path")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fpath, ok := dictmgr.DictResFile(dictName, path)
	if !ok {
		writeMsg(w, "file not found")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	file, err := os.Open(fpath)
	if err != nil {
		writeMsg(w, "file not found")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	http.ServeContent(w, r, "", time.Now(), file)
}

func addWebHandlers() {
	http.HandleFunc("/"+path_query, query)
	http.HandleFunc("/", home)
	http.HandleFunc(dictmgr.DictResPathBase, dictRes)

	http.Handle("/web/", http.FileServer(&httpFileSystem{
		fs:     web.FS,
		prefix: "web",
	}))
}

func StartServer(port string) {
	http.HandleFunc("/"+path_appName, getAppName)

	if conf.WebEnable {
		addWebHandlers()
	}

	log.Println("Starting local server on port", port)
	addr := "127.0.0.1:" + port
	if conf.WebExpose {
		addr = ":" + port
	}
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		qerr.Error(err)
	}
}
