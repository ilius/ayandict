package server

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"text/template"
	"time"

	"github.com/ilius/ayandict/v2/pkg/appinfo"
	"github.com/ilius/ayandict/v2/pkg/config"
	"github.com/ilius/ayandict/v2/pkg/dictmgr"
	"github.com/ilius/ayandict/v2/pkg/qtcommon/qerr"
	"github.com/ilius/ayandict/v2/web"
	common "github.com/ilius/go-dict-commons"
)

const (
	localhost    = "127.0.0.1"
	path_appName = "app-name"
	path_query   = "query"
	path_random  = "random"
)

var (
	conf = config.MustLoad()

	homeTpl *template.Template
)

const resultFlags = uint32(common.ResultFlag_FixAudio |
	common.ResultFlag_FixFileSrc |
	common.ResultFlag_Web)

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
		slog.Error("error in Write", "err", err)
	}
}

func getAppName(w http.ResponseWriter, _ *http.Request) {
	writeMsg(w, appinfo.APP_NAME)
}

func query(w http.ResponseWriter, r *http.Request) {
	t := time.Now()

	jsonEncoder := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	query := r.FormValue("query")
	if query == "" {
		err := jsonEncoder.Encode(ErrorResponse{Error: "missing query"})
		if err != nil {
			slog.Error("error in jsonEncoder.Encode", "err", err)
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	mode := dictmgr.QueryModeFuzzy
	switch r.FormValue("mode") {
	case "", "fuzzy":
	case "startWith":
		mode = dictmgr.QueryModeStartWith
	case "regex":
		mode = dictmgr.QueryModeRegex
	case "glob":
		mode = dictmgr.QueryModeGlob
	default:
		err := jsonEncoder.Encode(ErrorResponse{Error: "invalid mode"})
		if err != nil {
			slog.Error("error in jsonEncoder.Encode", "err", err)
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	flags := resultFlags
	switch r.FormValue("qt") {
	case "":
	case "5", "6":
		flags = flags | common.ResultFlag_FixWordLink | common.ResultFlag_ColorMapping
	default:
		err := jsonEncoder.Encode(ErrorResponse{Error: "invalid qt version, must be 5 or 6"})
		if err != nil {
			slog.Error("error in jsonEncoder.Encode", "err", err)
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	limit := 0
	limitStr := r.FormValue("limit")
	if limitStr != "" {
		limitI64, err := strconv.ParseUint(limitStr, 10, 0)
		if err != nil {
			err := jsonEncoder.Encode(ErrorResponse{Error: "invalid limit"})
			if err != nil {
				slog.Error("error in jsonEncoder.Encode", "err", err)
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		limit = int(limitI64)
	}

	raw_results := dictmgr.LookupHTML(query, conf, mode, flags, limit)
	// pass resultFlags to LookupHTML
	results := make([]Result, len(raw_results))
	for i, entry := range raw_results {
		results[i] = Result{
			DictName:        entry.DictName(),
			Terms:           entry.Terms(),
			DefinitionsHTML: entry.DefinitionsHTML(),
			EntryIndex:      entry.EntryIndex(),
			Score:           entry.Score(),
		}
		// entry.ResourceDir()
	}
	slog.Info("LookupHTML running time", "dt", time.Since(t), "query", query)
	err := jsonEncoder.Encode(results)
	if err != nil {
		slog.Error("error in jsonEncoder.Encode", "err", err)
		err2 := jsonEncoder.Encode(ErrorResponse{Error: err.Error()})
		if err2 != nil {
			slog.Error("error in jsonEncoder.Encode", "err2", err2)
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func random(w http.ResponseWriter, _ *http.Request) {
	jsonEncoder := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	entry := dictmgr.RandomEntry(conf, resultFlags)
	err := jsonEncoder.Encode(Result{
		DictName:        entry.DictName(),
		Terms:           entry.Terms(),
		DefinitionsHTML: entry.DefinitionsHTML(),
		EntryIndex:      entry.EntryIndex(),
		Score:           entry.Score(),
	})
	if err != nil {
		slog.Error("error in jsonEncoder.Encode", "err", err)
		err2 := jsonEncoder.Encode(ErrorResponse{Error: err.Error()})
		if err2 != nil {
			slog.Error("error in jsonEncoder.Encode", "err2", err2)
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

type homeTemplateParams struct {
	Config *config.Config
}

func home(w http.ResponseWriter, _ *http.Request) {
	err := homeTpl.Execute(w, homeTemplateParams{
		Config: conf,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
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
	http.HandleFunc("/"+path_random, random)
	http.HandleFunc("/", home)
	http.HandleFunc(dictmgr.DictResPathBase, dictRes)

	http.Handle("/web/", http.FileServer(&httpFileSystem{
		fs:     web.FS,
		prefix: "web",
	}))
}

func loadWebTemplates() error {
	file, err := web.FS.Open("web/index.html")
	if err != nil {
		return err
	}
	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	tpl, err := template.New("index").Parse(string(data))
	if err != nil {
		return err
	}
	homeTpl = tpl
	return nil
}

func StartServer(port string) {
	http.HandleFunc("/"+path_appName, getAppName)

	if conf.WebEnable {
		err := loadWebTemplates()
		if err != nil {
			panic(err)
		}
		addWebHandlers()
	}

	slog.Info("Starting local server", "port", port)
	addr := "127.0.0.1:" + port
	if conf.WebExpose {
		addr = ":" + port
	}
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		qerr.Error(err)
	}
}
