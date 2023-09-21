package server

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
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

var conf *config.Config

const resultFlags = common.ResultFlag_Web |
	common.ResultFlag_FixAudio |
	common.ResultFlag_FixFileSrc

// 	common.ResultFlag_ColorMapping)

func init() {
	var err error
	conf, err = config.Load()
	if err != nil {
		panic(err)
	}
}

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

func getAppName(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte(appinfo.APP_NAME))
	if err != nil {
		log.Println(err)
	}
}

func query(w http.ResponseWriter, r *http.Request) {
	t := time.Now()
	mode := dictmgr.QueryModeFuzzy

	jsonEncoder := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	query := r.FormValue("query")
	if query == "" {
		jsonEncoder.Encode(ErrorResponse{Error: "missing query"})
		w.WriteHeader(400)
		return
	}

	// mode = dictmgr.QueryModeStartWith
	// mode = dictmgr.QueryModeRegex
	// mode = dictmgr.QueryModeGlob

	raw_results := dictmgr.LookupHTML(query, conf, mode)
	results := make([]Result, len(raw_results))
	for i, rr := range raw_results {
		results[i] = Result{
			DictName:        rr.DictName(),
			Terms:           rr.Terms(),
			DefinitionsHTML: rr.DefinitionsHTML(resultFlags),
			EntryIndex:      rr.EntryIndex(),
			Score:           rr.Score(),
		}
		// rr.ResourceDir()
	}
	log.Printf("LookupHTML took %v for %#v", time.Since(t), query)
	err := jsonEncoder.Encode(results)
	if err != nil {
		log.Println(err)
		jsonEncoder.Encode(ErrorResponse{Error: err.Error()})
		w.WriteHeader(500)
		return
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	file, err := web.FS.Open("web/index.html")
	if err != nil {
		w.WriteHeader(404)
		return
	}
	content := file.(io.ReadSeeker)
	http.ServeContent(w, r, "", time.Now(), content)
}

func addWebHandlers() {
	http.HandleFunc("/"+path_query, query)
	http.HandleFunc("/", home)

	fs := &httpFileSystem{
		fs:     web.FS,
		prefix: "web",
	}
	// http.Handle("/web", http.StripPrefix("/web", http.FileServer(fs)))
	http.Handle("/web/", http.FileServer(fs))
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
