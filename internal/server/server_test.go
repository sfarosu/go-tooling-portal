package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/sfarosu/go-tooling-portal/internal/version"
)

func Test_setupRouter(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "router is not nil and serves static files"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := setupRouter()
			if router == nil {
				t.Fatal("setupRouter() returned nil")
			}
			// Optionally, check that static file handler is registered
			req := httptest.NewRequest("GET", "/", nil)
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)
			// Accept 200, 403, or 404 depending on your static file setup
			if rr.Code != http.StatusOK && rr.Code != http.StatusNotFound && rr.Code != http.StatusForbidden {
				t.Errorf("setupRouter() root handler returned status [%v], want 200, 403, or 404", rr.Code)
			}
		})
	}
}

func Test_setupRouterOpenAPIDocumentation(t *testing.T) {
	router := setupRouter()
	req := httptest.NewRequest(http.MethodGet, "/openapi.json", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("GET /openapi.json returned status %d: %s", rr.Code, rr.Body.String())
	}

	var document struct {
		Info struct {
			Description string `json:"description"`
			Version     string `json:"version"`
		} `json:"info"`
		Tags []struct {
			Name        string `json:"name"`
			Description string `json:"description"`
		} `json:"tags"`
		Paths map[string]map[string]struct {
			Summary     string   `json:"summary"`
			Description string   `json:"description"`
			Tags        []string `json:"tags"`
			Parameters  []struct {
				Name        string `json:"name"`
				Description string `json:"description"`
			} `json:"parameters"`
			Responses map[string]struct {
				Content map[string]json.RawMessage `json:"content"`
			} `json:"responses"`
		} `json:"paths"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &document); err != nil {
		t.Fatalf("decode OpenAPI document: %v", err)
	}

	if document.Info.Description == "" {
		t.Fatal("OpenAPI document is missing its API description")
	}
	if document.Info.Version != version.Version {
		t.Fatalf("OpenAPI version = %q, want application version %q", document.Info.Version, version.Version)
	}
	hasTagDescription := false
	for _, tag := range document.Tags {
		if tag.Name == "URL Decoder" && tag.Description != "" {
			hasTagDescription = true
		}
	}
	if !hasTagDescription {
		t.Fatal("OpenAPI document is missing the URL Decoder tag description")
	}

	operation, ok := document.Paths["/api/urldecoder/htmx"]["post"]
	if !ok {
		t.Fatal("OpenAPI document is missing POST /api/urldecoder/htmx")
	}
	if !strings.Contains(operation.Summary, "(HTMX)") || !strings.Contains(operation.Description, "HX-Request: true") {
		t.Fatalf("HTMX endpoint metadata is unclear: summary=%q description=%q", operation.Summary, operation.Description)
	}
	if len(operation.Tags) != 1 || operation.Tags[0] != "URL Decoder" {
		t.Fatalf("HTMX endpoint has unexpected tags: %v", operation.Tags)
	}
	hasHeaderDescription := false
	for _, parameter := range operation.Parameters {
		if parameter.Name == "HX-Request" && parameter.Description != "" {
			hasHeaderDescription = true
		}
	}
	if !hasHeaderDescription {
		t.Fatal("HTMX endpoint is missing the HX-Request parameter description")
	}
	if _, ok := operation.Responses["200"].Content["text/html"]; !ok {
		t.Fatal("HTMX endpoint does not document its HTML response")
	}
}

func Test_setupServer(t *testing.T) {
	tests := []struct {
		name string
		addr string
	}{
		{name: "server config", addr: ":12345"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := http.NewServeMux()
			srv := setupServer(tt.addr, handler)
			if srv == nil {
				t.Fatal("setupServer() returned nil")
			}
			if srv.Addr != tt.addr {
				t.Errorf("setupServer() Addr = [%v], want [%v]", srv.Addr, tt.addr)
			}
			if srv.Handler != handler {
				t.Error("setupServer() Handler does not match input handler")
			}
			if srv.ReadTimeout != serverReadTimeout {
				t.Errorf("setupServer() ReadTimeout = [%v], want [%v]", srv.ReadTimeout, serverReadTimeout)
			}
			if srv.WriteTimeout != serverWriteTimeout {
				t.Errorf("setupServer() WriteTimeout = [%v], want [%v]", srv.WriteTimeout, serverWriteTimeout)
			}
			if srv.IdleTimeout != serverIdleTimeout {
				t.Errorf("setupServer() IdleTimeout = [%v], want [%v]", srv.IdleTimeout, serverIdleTimeout)
			}
		})
	}
}
