package httpserver

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/sfarosu/go-tooling-portal/internal/logger"
	"github.com/sfarosu/go-tooling-portal/internal/metrics"
	"github.com/sfarosu/go-tooling-portal/internal/requestid"
)

const maxMCPLogBodyBytes = 4 << 20

type statusRecorder struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

type replayReadCloser struct {
	io.Reader
	io.Closer
}

func (r *statusRecorder) Unwrap() http.ResponseWriter {
	return r.ResponseWriter
}

func (r *statusRecorder) WriteHeader(code int) {
	if r.wroteHeader {
		return // Prevent multiple calls
	}
	r.status = code
	r.wroteHeader = true
	r.ResponseWriter.WriteHeader(code)
}

// loggingMiddleware logs each HTTP request with appropriate log level
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: 200}

		requestID := requestid.FromContext(r.Context())
		component := "api"
		if r.URL.Path == "/mcp" {
			component = "mcp"
		}
		toolName := ""
		if component == "mcp" && r.Method == http.MethodPost {
			toolName = mcpToolName(r)
		}
		requestLogArgs := []any{
			"request-id", requestID,
			"method", r.Method,
			"path", r.URL.Path,
			"protocol", r.Proto,
			"remote-addr", r.RemoteAddr,
			"user-agent", r.UserAgent(),
		}
		if toolName != "" {
			requestLogArgs = append(requestLogArgs, "tool-name", toolName)
		}

		logger.Logger.Info(component+"-incoming-request", requestLogArgs...)

		next.ServeHTTP(rec, r)

		if component == "api" && strings.HasPrefix(r.URL.Path, "/api/") {
			metrics.RecordAPIRequest(r.URL.Path)
		}
		if component == "mcp" && toolName != "" {
			metrics.RecordMCPToolCall(toolName)
		}

		duration := time.Since(start).Milliseconds()
		completionLogArgs := []any{
			"request-id", requestID,
			"method", r.Method,
			"path", r.URL.Path,
			"status", rec.status,
			"duration-ms", duration,
		}
		if toolName != "" {
			completionLogArgs = append(completionLogArgs, "tool-name", toolName)
		}

		message := component + "-request-completed"
		switch {
		case rec.status >= 500:
			logger.Logger.Error(message, completionLogArgs...)
		case rec.status >= 400:
			logger.Logger.Warn(message, completionLogArgs...)
		default:
			logger.Logger.Info(message, completionLogArgs...)
		}
	})
}

func mcpToolName(r *http.Request) string {
	if r.Body == nil {
		return ""
	}

	originalBody := r.Body
	body, err := io.ReadAll(io.LimitReader(originalBody, maxMCPLogBodyBytes+1))
	r.Body = replayReadCloser{
		Reader: io.MultiReader(bytes.NewReader(body), originalBody),
		Closer: originalBody,
	}
	if err != nil || len(body) > maxMCPLogBodyBytes {
		return ""
	}

	var message struct {
		Method string `json:"method"`
		Params struct {
			Name string `json:"name"`
		} `json:"params"`
	}
	if err := json.Unmarshal(body, &message); err != nil || message.Method != "tools/call" {
		return ""
	}
	return message.Params.Name
}
