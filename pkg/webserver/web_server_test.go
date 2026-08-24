package webserver

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ilius/ayandict/v3/pkg/config"
)

func handlerPattern(t *testing.T, handler http.Handler, target string) string {
	t.Helper()

	mux, ok := handler.(*http.ServeMux)
	if !ok {
		t.Fatalf("expected *http.ServeMux, got %T", handler)
	}
	request := httptest.NewRequest(http.MethodGet, target, nil)
	_, pattern := mux.Handler(request)
	return pattern
}

func TestNewServerRoutes(t *testing.T) {
	conf := config.Default()
	server, err := newServer(conf, "8357")
	if err != nil {
		t.Fatal(err)
	}
	if pattern := handlerPattern(t, server.Handler, "/app-name"); pattern != "/app-name" {
		t.Fatalf("unexpected app-name pattern: %q", pattern)
	}
	for _, target := range []string{"/", "/api/query", "/api/random", "/web/style.css"} {
		if pattern := handlerPattern(t, server.Handler, target); pattern != "" {
			t.Errorf("unexpected pattern %q for disabled web route %q", pattern, target)
		}
	}

	conf.WebEnable = true
	server, err = newServer(conf, "8357")
	if err != nil {
		t.Fatal(err)
	}
	for target, expected := range map[string]string{
		"/":              "/",
		"/api/query":     "/api/query",
		"/api/random":    "/api/random",
		"/web/style.css": "/web/",
	} {
		if pattern := handlerPattern(t, server.Handler, target); pattern != expected {
			t.Errorf("pattern for %q: got %q, want %q", target, pattern, expected)
		}
	}
}

func TestNewServerUsesPassedConfig(t *testing.T) {
	conf := config.Default()
	conf.WebEnable = true
	conf.WebSearchOnTypeMinLength = 73
	conf.WebShowPoweredBy = false

	server, err := newServer(conf, "8357")
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	server.Handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("unexpected status: got %d, want %d", recorder.Code, http.StatusOK)
	}
	body := recorder.Body.String()
	if !strings.Contains(body, "if len(query) < 73:") {
		t.Error("response does not contain the passed search-on-type minimum length")
	}
	if strings.Contains(body, "Powered by") {
		t.Error("response contains Powered by footer disabled in the passed config")
	}
}

func TestNewServerAddress(t *testing.T) {
	for _, test := range []struct {
		name   string
		expose bool
		want   string
	}{
		{name: "local", expose: false, want: "127.0.0.1:8357"},
		{name: "exposed", expose: true, want: ":8357"},
	} {
		t.Run(test.name, func(t *testing.T) {
			conf := config.Default()
			conf.WebExpose = test.expose
			server, err := newServer(conf, "8357")
			if err != nil {
				t.Fatal(err)
			}
			if server.Addr != test.want {
				t.Fatalf("got address %q, want %q", server.Addr, test.want)
			}
		})
	}
}
