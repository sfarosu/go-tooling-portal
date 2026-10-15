package httpserver

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/sfarosu/go-tooling-portal/internal/requestid"
)

// requestIDMiddleware generates a unique request ID for each incoming HTTP request
func requestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := ""
		mcpRequestID := r.Header.Get("MCP-Request-ID")
		if r.URL.Path == "/mcp" {
			requestID = uuid.NewString()
			mcpRequestID = requestID
		} else if _, err := uuid.Parse(mcpRequestID); err == nil {
			requestID = mcpRequestID
		}
		if requestID == "" {
			requestID = uuid.NewString()
		}

		ctx := requestid.WithContext(r.Context(), requestID)
		if r.URL.Path == "/mcp" || mcpRequestID == requestID {
			ctx = requestid.WithMCPContext(ctx, requestID)
			w.Header().Set("MCP-Request-ID", requestID)
		}

		r = r.WithContext(ctx)
		w.Header().Set("X-Request-ID", requestID)
		next.ServeHTTP(w, r)
	})
}
