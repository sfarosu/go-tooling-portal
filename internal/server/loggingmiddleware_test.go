package server

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/sfarosu/go-tooling-portal/internal/logger"
)

func Test_statusRecorder_WriteHeader(t *testing.T) {
	tests := []struct {
		name           string
		initialStatus  int
		writeCode      int
		expectedStatus int
	}{
		{
			name:           "sets status to 404",
			initialStatus:  200,
			writeCode:      404,
			expectedStatus: 404,
		},
		{
			name:           "sets status to 201",
			initialStatus:  200,
			writeCode:      201,
			expectedStatus: 201,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			rec := &statusRecorder{
				ResponseWriter: rr,
				status:         tt.initialStatus,
			}
			rec.WriteHeader(tt.writeCode)
			if rec.status != tt.expectedStatus {
				t.Errorf("WriteHeader() status = [%v], want [%v]", rec.status, tt.expectedStatus)
			}
			if rr.Code != tt.expectedStatus {
				t.Errorf("WriteHeader() underlying ResponseWriter code = [%v], want [%v]", rr.Code, tt.expectedStatus)
			}
		})
	}
}

func Test_loggingMiddleware(t *testing.T) {
	// Initialize a test logger to avoid nil pointer panic
	logger.Logger = slog.New(slog.NewTextHandler(os.Stdout, nil))

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
