package server

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/sfarosu/go-tooling-portal/internal/logger"
)

func Test_statusRecorder_WriteHeader(t *testing.T) {
	tests := []struct {
		name           string
		firstCode      int
		secondCode     int
		expectedStatus int
	}{
		{
			name:           "sets status to 404, ignores second WriteHeader",
			firstCode:      404,
			secondCode:     200,
			expectedStatus: 404,
		},
		{
			name:           "sets status to 201, ignores second WriteHeader",
			firstCode:      201,
			secondCode:     500,
			expectedStatus: 201,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			rec := &statusRecorder{
				ResponseWriter: rr,
			}
			rec.WriteHeader(tt.firstCode)
			rec.WriteHeader(tt.secondCode) // This should be ignored

			if rec.status != tt.expectedStatus {
				t.Errorf("WriteHeader() status = [%v], want [%v]", rec.status, tt.expectedStatus)
			}
			if rr.Code != tt.expectedStatus {
				t.Errorf("underlying ResponseWriter code = [%v], want [%v]", rr.Code, tt.expectedStatus)
			}
		})
	}
}

func Test_loggingMiddlewareIncludesMCPToolName(t *testing.T) {
	previousLogger := logger.Logger
	var logOutput bytes.Buffer
	logger.Logger = slog.New(slog.NewTextHandler(&logOutput, nil))
	t.Cleanup(func() {
		logger.Logger = previousLogger
	})

	body := `{"jsonrpc":"2.0","id":7,"method":"tools/call","params":{"name":"generate-password","arguments":{"length":20}}}`
	handler := loggingMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("read request body: %v", err)
		}
		if string(got) != body {
			t.Errorf("request body = %q, want %q", got, body)
		}
		w.WriteHeader(http.StatusOK)
	}))
	request := httptest.NewRequest(http.MethodPost, "/mcp", strings.NewReader(body))
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)

	logs := logOutput.String()
	for _, event := range []string{"mcp-incoming-request", "mcp-request-completed"} {
		line := ""
		for _, candidate := range strings.Split(logs, "\n") {
			if strings.Contains(candidate, "msg="+event) {
				line = candidate
				break
			}
		}
		if line == "" || !strings.Contains(line, "tool-name=generate-password") {
			t.Errorf("%s does not include MCP tool name: %s", event, logs)
		}
	}
	if strings.Contains(logs, "length=20") {
		t.Fatalf("logs should include only tool name, not request arguments: %s", logs)
	}
}

func Test_loggingMiddleware(t *testing.T) {
	previousLogger := logger.Logger
	logger.Logger = slog.New(slog.NewTextHandler(io.Discard, nil))
	t.Cleanup(func() {
		logger.Logger = previousLogger
	})

	// We'll use a dummy handler that writes a specific status code
	tests := []struct {
		name       string
		statusCode int
	}{
		{"info log (200)", 200},
		{"warn log (404)", 404},
		{"error log (500)", 500},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				w.Write([]byte("ok"))
			})

			req := httptest.NewRequest("GET", "/", nil)
			rr := httptest.NewRecorder()

			loggingMiddleware(handler).ServeHTTP(rr, req)

			if rr.Code != tt.statusCode {
				t.Errorf("loggingMiddleware() status = [%v], want [%v]", rr.Code, tt.statusCode)
			}
			if !strings.Contains(rr.Body.String(), "ok") {
				t.Errorf("loggingMiddleware() body = [%v], want to contain \"ok\"", rr.Body.String())
			}
		})
	}
}
