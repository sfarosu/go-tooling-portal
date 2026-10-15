package server

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/sfarosu/go-tooling-portal/internal/logger"
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
			router, err := setupRouter("http://127.0.0.1")
			if err != nil {
				t.Fatalf("setupRouter() returned an error: %v", err)
			}
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

func Test_setupRouterServesEmbeddedAssets(t *testing.T) {
	router, err := setupRouter("http://127.0.0.1")
	if err != nil {
		t.Fatalf("setupRouter() returned an error: %v", err)
	}
	tests := []struct {
		path         string
		contentType  string
		bodyContains string
		status       int
	}{
		{path: "/", contentType: "text/html", bodyContains: "Go Tooling Portal", status: http.StatusOK},
		{path: "/assets/styles.css", contentType: "text/css", bodyContains: "body {", status: http.StatusOK},
		{path: "/assets/", contentType: "text/html", bodyContains: "404", status: http.StatusNotFound},
	}

	for _, test := range tests {
		t.Run(test.path, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, test.path, nil)
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			if rr.Code != test.status {
				t.Fatalf("GET %s returned status %d, want %d: %s", test.path, rr.Code, test.status, rr.Body.String())
			}
			if contentType := rr.Header().Get("Content-Type"); !strings.Contains(contentType, test.contentType) {
				t.Errorf("GET %s Content-Type = %q, want it to contain %q", test.path, contentType, test.contentType)
			}
			if !strings.Contains(rr.Body.String(), test.bodyContains) {
				t.Errorf("GET %s body does not contain %q", test.path, test.bodyContains)
			}
		})
	}
}

func Test_requestIDCorrelatesHandlerAndCompletionLogs(t *testing.T) {
	previousLogger := logger.Logger
	var logOutput bytes.Buffer
	logger.Logger = slog.New(slog.NewTextHandler(&logOutput, nil))
	t.Cleanup(func() {
		logger.Logger = previousLogger
	})

	router, err := setupRouter("http://127.0.0.1")
	if err != nil {
		t.Fatalf("setupRouter() returned an error: %v", err)
	}
	handler := requestIDMiddleware(loggingMiddleware(router))
	req := httptest.NewRequest(http.MethodPost, "/api/formatconverter", strings.NewReader(`{"input":"not: [valid"}`))
	req.Header.Set("Content-Type", "application/json")
	mcpRequestID := "e6512b20-9ec9-4c42-b93b-d7fcf32440e7"
	req.Header.Set("MCP-Request-ID", mcpRequestID)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	requestID := rr.Header().Get("X-Request-ID")
	if requestID == "" {
		t.Fatal("response is missing X-Request-ID")
	}
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("request returned status %d, want %d: %s", rr.Code, http.StatusBadRequest, rr.Body.String())
	}
	if got := rr.Header().Get("MCP-Request-ID"); got != mcpRequestID {
		t.Fatalf("response MCP-Request-ID = %q, want %q", got, mcpRequestID)
	}

	logs := logOutput.String()
	if !strings.Contains(logs, `msg="failed to perform format conversion"`) {
		t.Fatalf("handler error log missing: %s", logs)
	}
	if strings.Contains(logs, `"input"`) || strings.Contains(logs, "not: [valid") {
		t.Fatalf("API logs include the submitted payload: %s", logs)
	}
	if !strings.Contains(logs, `msg=api-request-completed`) {
		t.Fatalf("request completion log missing: %s", logs)
	}
	if !strings.Contains(logs, `msg=api-incoming-request`) {
		t.Fatalf("API incoming request log missing: %s", logs)
	}
	if requestID != mcpRequestID {
		t.Fatalf("API request ID = %q, want propagated MCP request ID %q", requestID, mcpRequestID)
	}
	if count := strings.Count(logs, `request-id=`+requestID); count < 3 {
		t.Fatalf("request ID %q appears in %d log entries, want handler, incoming, and completion logs: %s", requestID, count, logs)
	}
}

func Test_mcpRequestLogsIdentifyMCPComponent(t *testing.T) {
	previousLogger := logger.Logger
	var logOutput bytes.Buffer
	logger.Logger = slog.New(slog.NewTextHandler(&logOutput, nil))
	t.Cleanup(func() {
		logger.Logger = previousLogger
	})

	router, err := setupRouter("http://127.0.0.1")
	if err != nil {
		t.Fatalf("setupRouter() returned an error: %v", err)
	}
	handler := requestIDMiddleware(loggingMiddleware(router))
	req := httptest.NewRequest(http.MethodGet, "/mcp", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("GET /mcp returned status %d, want %d", rr.Code, http.StatusMethodNotAllowed)
	}
	logs := logOutput.String()
	if !strings.Contains(logs, `msg=mcp-incoming-request`) || !strings.Contains(logs, `msg=mcp-request-completed`) {
		t.Fatalf("MCP request logs do not identify the component: %s", logs)
	}
	requestID := rr.Header().Get("X-Request-ID")
	if requestID == "" || rr.Header().Get("MCP-Request-ID") != requestID {
		t.Fatalf("MCP response identifiers do not share a request ID: X-Request-ID=%q MCP-Request-ID=%q", requestID, rr.Header().Get("MCP-Request-ID"))
	}
	if count := strings.Count(logs, `request-id=`+requestID); count != 2 {
		t.Fatalf("MCP request ID %q appears in %d log entries, want exactly two transport entries: %s", requestID, count, logs)
	}
	if strings.Contains(logs, `mcp-request-id=`) || strings.Contains(logs, `mcp-http-request-id=`) {
		t.Fatalf("MCP logs contain redundant request ID fields: %s", logs)
	}
}

func Test_directAPIRequestUsesOneRequestIDAcrossLogs(t *testing.T) {
	previousLogger := logger.Logger
	var logOutput bytes.Buffer
	logger.Logger = slog.New(slog.NewTextHandler(&logOutput, nil))
	t.Cleanup(func() {
		logger.Logger = previousLogger
	})

	router, err := setupRouter("http://127.0.0.1")
	if err != nil {
		t.Fatalf("setupRouter() returned an error: %v", err)
	}
	handler := requestIDMiddleware(loggingMiddleware(router))
	req := httptest.NewRequest(http.MethodGet, "/api/version", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("GET /api/version returned status %d: %s", rr.Code, rr.Body.String())
	}
	requestID := rr.Header().Get("X-Request-ID")
	if requestID == "" {
		t.Fatal("API response is missing X-Request-ID")
	}
	logs := logOutput.String()
	if count := strings.Count(logs, `request-id=`+requestID); count != 2 {
		t.Fatalf("API request ID %q appears in %d log entries, want exactly two: %s", requestID, count, logs)
	}
	if !strings.Contains(logs, `msg=api-incoming-request`) || !strings.Contains(logs, `msg=api-request-completed`) {
		t.Fatalf("API request logs are missing lifecycle events: %s", logs)
	}
	if strings.Contains(logs, `mcp-request-id=`) {
		t.Fatalf("direct API request unexpectedly has an MCP request ID: %s", logs)
	}
}

func Test_setupRouterOpenAPIDocumentation(t *testing.T) {
	router, err := setupRouter("http://127.0.0.1")
	if err != nil {
		t.Fatalf("setupRouter() returned an error: %v", err)
	}
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

func Test_mcpAPIBaseURLUsesReachableLoopback(t *testing.T) {
	tests := []struct {
		name string
		addr net.Addr
		want string
	}{
		{
			name: "IPv4 wildcard",
			addr: &net.TCPAddr{IP: net.IPv4zero, Port: 8100},
			want: "http://127.0.0.1:8100",
		},
		{
			name: "IPv6 wildcard",
			addr: &net.TCPAddr{IP: net.IPv6zero, Port: 8100},
			want: "http://[::1]:8100",
		},
		{
			name: "loopback",
			addr: &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 8100},
			want: "http://127.0.0.1:8100",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := mcpAPIBaseURL(test.addr)
			if err != nil {
				t.Fatalf("mcpAPIBaseURL() error = %v", err)
			}
			if got != test.want {
				t.Fatalf("mcpAPIBaseURL() = %q, want %q", got, test.want)
			}
		})
	}
}
