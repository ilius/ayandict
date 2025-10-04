package server

import (
	"encoding/json"
	html_template "html/template"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	text_template "text/template"
	"time"

	"github.com/ilius/ayandict/v2/pkg/appinfo"
	"github.com/ilius/ayandict/v2/pkg/config"
	"github.com/ilius/ayandict/v2/pkg/dictmgr"
	"github.com/ilius/ayandict/v2/pkg/headerlib"
	"github.com/ilius/ayandict/v2/pkg/logging"
	"github.com/ilius/ayandict/v2/web"
	common "github.com/ilius/go-dict-commons"
)

const (
	localhost       = "127.0.0.1"
	path_appName    = "app-name"
	path_api_query  = "api/query"
	path_api_random = "api/random"
)

var (
	conf = config.MustLoad()

	homeTpl   *text_template.Template
	headerTpl *html_template.Template
)

// using a different logger here, so that it does not show errors in GUI
// because there is a little risk in showing web-user-input values in GUI
var logger = slog.New(logging.NewColoredHandler(
	os.Getenv("NO_COLOLR") != "",
	logging.DefaultLevel,
))

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
	HeaderHTML      string   `json:"header_html"`
	// ResourceDir string
}

func writeMsg(w http.ResponseWriter, msg string) {
	_, err := w.Write([]byte(msg))
	if err != nil {
		logger.Error("error in Write", "err", err)
	}
}

func getAppName(w http.ResponseWriter, _ *http.Request) {
	writeMsg(w, appinfo.APP_NAME)
}

func queryModeParam(r *http.Request) (dictmgr.QueryMode, bool) {
	switch r.FormValue("mode") {
	case "", "fuzzy":
		return dictmgr.QueryModeFuzzy, true
	case "startWith":
		return dictmgr.QueryModeStartWith, true
	case "regex":
		return dictmgr.QueryModeRegex, true
	case "glob":
		return dictmgr.QueryModeGlob, true
	case "wordMatch":
		return dictmgr.QueryModeWordMatch, true
	}
	return dictmgr.QueryMode(0), false
}

func api_query(w http.ResponseWriter, r *http.Request) {
	t := time.Now()

	jsonEncoder := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	query := r.FormValue("query")
	if query == "" {
		err := jsonEncoder.Encode(ErrorResponse{Error: "missing query"})
		if err != nil {
			logger.Error("error in jsonEncoder.Encode", "err", err)
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	mode, ok := queryModeParam(r)
	if !ok {
		err := jsonEncoder.Encode(ErrorResponse{Error: "invalid mode"})
		if err != nil {
			logger.Error("error in jsonEncoder.Encode", "err", err)
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
			logger.Error("error in jsonEncoder.Encode", "err", err)
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
				logger.Error("error in jsonEncoder.Encode", "err", err)
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		limit = int(limitI64)
	}

	raw_results := dictmgr.LookupHTML(query, conf, mode, flags, limit)
	// pass resultFlags to LookupHTML
	results := make([]Result, len(raw_results))
	for i, res := range raw_results {
		header, err := headerlib.GetHeader(headerTpl, res)
		if err != nil {
			logger.Error("Error formatting header label", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		results[i] = Result{
			DictName:        res.DictName(),
			Terms:           res.Terms(),
			DefinitionsHTML: res.DefinitionsHTML(),
			EntryIndex:      res.EntryIndex(),
			Score:           res.Score(),
			HeaderHTML:      header,
		}
		// entry.ResourceDir()
	}
	logger.Info("LookupHTML running time", "dt", time.Since(t), "query", query)
	err := jsonEncoder.Encode(results)
	if err != nil {
		logger.Error("error in jsonEncoder.Encode", "err", err)
		err2 := jsonEncoder.Encode(ErrorResponse{Error: err.Error()})
		if err2 != nil {
			logger.Error("error in jsonEncoder.Encode", "err2", err2)
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func api_random(w http.ResponseWriter, _ *http.Request) {
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
		logger.Error("error in jsonEncoder.Encode", "err", err)
		err2 := jsonEncoder.Encode(ErrorResponse{Error: err.Error()})
		if err2 != nil {
			logger.Error("error in jsonEncoder.Encode", "err2", err2)
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
	http.HandleFunc("/"+path_api_query, api_query)
	http.HandleFunc("/"+path_api_random, api_random)
	http.HandleFunc("/", home)
	http.HandleFunc(dictmgr.DictResPathBase, dictRes)

	http.Handle("/web/", http.FileServer(&httpFileSystem{
		fs:     web.FS,
		prefix: "web",
	}))
}

func loadIndexTemplate() error {
	file, err := web.FS.Open("web/index.html")
	if err != nil {
		return err
	}
	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	tpl, err := text_template.New("index").Parse(string(data))
	if err != nil {
		return err
	}
	homeTpl = tpl
	return nil
}

func loadHeaderTemplate() error {
	tpl, err := headerlib.LoadHeaderTemplate(conf)
	if err != nil {
		return err
	}
	headerTpl = tpl
	return nil
}

func loadWebTemplates() error {
	err := loadIndexTemplate()
	if err != nil {
		return err
	}
	err = loadHeaderTemplate()
	if err != nil {
		return err
	}
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

	logger.Info("Starting local server", "port", port)
	addr := "127.0.0.1:" + port
	if conf.WebExpose {
		addr = ":" + port
	}
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		logger.Error("error in ListenAndServe: " + err.Error())
	}
}
