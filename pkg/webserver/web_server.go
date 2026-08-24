package webserver

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

	common "codeberg.org/ilius/go-dict-commons"
	"github.com/ilius/ayandict/v3/pkg/appinfo"
	"github.com/ilius/ayandict/v3/pkg/config"
	"github.com/ilius/ayandict/v3/pkg/dictmgr"
	"github.com/ilius/ayandict/v3/pkg/headerlib"
	"github.com/ilius/ayandict/v3/pkg/jsonapi"
	"github.com/ilius/ayandict/v3/pkg/logging"
	"github.com/ilius/ayandict/v3/web"
)

const (
	localhost       = "127.0.0.1"
	path_appName    = "app-name"
	path_api_query  = "api/query"
	path_api_random = "api/random"
)

type webServer struct {
	conf      *config.Config
	homeTpl   *text_template.Template
	headerTpl *html_template.Template
}

// using a different logger here, so that it does not show errors in GUI
// because there is a little risk in showing web-user-input values in GUI
var logger = slog.New(logging.NewColoredHandler(
	os.Getenv("NO_COLOR") != "",
	logging.DefaultLevel,
))

const resultFlags = uint32(common.ResultFlag_FixAudio |
	common.ResultFlag_FixFileSrc |
	common.ResultFlag_Web)

func writeMsg(w http.ResponseWriter, msg string) {
	_, err := w.Write([]byte(msg))
	if err != nil {
		logger.Error("error in Write", "err", err)
	}
}

func getAppName(w http.ResponseWriter, _ *http.Request) {
	writeMsg(w, appinfo.APP_NAME)
}

func badRequest(w http.ResponseWriter, msg string) {
	jsonEncoder := json.NewEncoder(w)
	err := jsonEncoder.Encode(jsonapi.ErrorResponse{Error: msg})
	if err != nil {
		logger.Error("error in jsonEncoder.Encode", "err", err)
	}
	w.WriteHeader(http.StatusBadRequest)
}

func (srv *webServer) encodeResults(
	w http.ResponseWriter,
	raw_results []common.SearchResultIface,
) []jsonapi.Result {
	results := make([]jsonapi.Result, len(raw_results))
	for i, res := range raw_results {
		header, err := headerlib.GetHeader(srv.headerTpl, res, 200)
		if err != nil {
			logger.Error("Error formatting header label", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			return nil
		}
		results[i] = jsonapi.Result{
			DictName:        res.DictName(),
			Terms:           res.Terms(),
			DefinitionsHTML: res.DefinitionsHTML(),
			EntryIndex:      res.EntryIndex(),
			Score:           res.Score(),
			HeaderHTML:      header,
		}
		// entry.ResourceDir()
	}
	return results
}

func (srv *webServer) api_query(w http.ResponseWriter, r *http.Request) {
	t := time.Now()

	w.Header().Set("Content-Type", "application/json")

	query := r.FormValue("query")
	if query == "" {
		badRequest(w, "missing query")
		return
	}
	mode, ok := dictmgr.SearchModeByName(r.FormValue("mode"))
	if !ok {
		badRequest(w, "invalid mode")
		return
	}

	flags := resultFlags
	switch r.FormValue("qt") {
	case "":
	case "5", "6":
		flags = (flags |
			common.ResultFlag_FixWordLink |
			common.ResultFlag_ColorMapping |
			common.ResultFlag_QTextBrowser)
	default:
		badRequest(w, "invalid qt version, must be 5 or 6")
		return
	}

	limit := 0
	limitStr := r.FormValue("limit")
	if limitStr != "" {
		limitI64, err := strconv.ParseUint(limitStr, 10, 0)
		if err != nil {
			badRequest(w, "invalid limit")
			return
		}
		limit = int(limitI64)
	}

	raw_results := dictmgr.LookupHTML(query, srv.conf, mode, flags, limit)
	// pass resultFlags to LookupHTML
	results := srv.encodeResults(w, raw_results)
	if results == nil {
		return
	}
	logger.Info("LookupHTML running time", "dt", time.Since(t), "query", query)
	jsonEncoder := json.NewEncoder(w)
	err := jsonEncoder.Encode(results)
	if err != nil {
		logger.Error("error in jsonEncoder.Encode", "err", err)
		err2 := jsonEncoder.Encode(jsonapi.ErrorResponse{Error: err.Error()})
		if err2 != nil {
			logger.Error("error in jsonEncoder.Encode", "err2", err2)
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (srv *webServer) api_random(w http.ResponseWriter, _ *http.Request) {
	jsonEncoder := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	entry := dictmgr.RandomEntry(srv.conf, resultFlags)
	err := jsonEncoder.Encode(jsonapi.Result{
		DictName:        entry.DictName(),
		Terms:           entry.Terms(),
		DefinitionsHTML: entry.DefinitionsHTML(),
		EntryIndex:      entry.EntryIndex(),
		Score:           entry.Score(),
	})
	if err != nil {
		logger.Error("error in jsonEncoder.Encode", "err", err)
		err2 := jsonEncoder.Encode(jsonapi.ErrorResponse{Error: err.Error()})
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

func (srv *webServer) home(w http.ResponseWriter, _ *http.Request) {
	err := srv.homeTpl.Execute(w, homeTemplateParams{
		Config: srv.conf,
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

func (srv *webServer) addWebHandlers(mux *http.ServeMux) {
	mux.HandleFunc("/"+path_api_query, srv.api_query)
	mux.HandleFunc("/"+path_api_random, srv.api_random)
	mux.HandleFunc("/", srv.home)
	mux.HandleFunc(dictmgr.DictResPathBase, dictRes)

	mux.Handle("/web/", http.FileServer(&httpFileSystem{
		fs:     web.FS,
		prefix: "web",
	}))
}

func (srv *webServer) loadIndexTemplate() error {
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
	srv.homeTpl = tpl
	return nil
}

func (srv *webServer) loadHeaderTemplate() error {
	tpl, err := headerlib.LoadHeaderTemplate(srv.conf)
	if err != nil {
		return err
	}
	srv.headerTpl = tpl
	return nil
}

func (srv *webServer) loadWebTemplates() error {
	err := srv.loadIndexTemplate()
	if err != nil {
		return err
	}
	err = srv.loadHeaderTemplate()
	if err != nil {
		return err
	}
	return nil
}

func newServer(conf *config.Config, port string) (*http.Server, error) {
	srv := &webServer{conf: conf}
	mux := http.NewServeMux()
	mux.HandleFunc("/"+path_appName, getAppName)
	// Preserve debug handlers such as net/http/pprof when a caller registers them
	// on the default mux.
	mux.Handle("/debug/pprof/", http.DefaultServeMux)

	if conf.WebEnable {
		err := srv.loadWebTemplates()
		if err != nil {
			return nil, err
		}
		srv.addWebHandlers(mux)
	}

	addr := localhost + ":" + port
	if conf.WebExpose {
		addr = ":" + port
	}
	return &http.Server{
		Addr:    addr,
		Handler: mux,
	}, nil
}

func StartServer(conf *config.Config, port string) {
	server, err := newServer(conf, port)
	if err != nil {
		panic(err)
	}

	logger.Info("Starting local server", "port", port)
	err = server.ListenAndServe()
	if err != nil {
		logger.Error("error in ListenAndServe: " + err.Error())
	}
}
