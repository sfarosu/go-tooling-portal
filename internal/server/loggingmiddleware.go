package server

import (
	"net/http"
	"time"

	"github.com/sfarosu/go-tooling-portal/internal/logger"
)

type statusRecorder struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
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

		reqID := getRequestID(r.Context())

		logger.Logger.Info(
			"incoming-request",
			"request-id", reqID,
			"method", r.Method,
			"path", r.URL.Path,
			"protocol", r.Proto,
			"remote-addr", r.RemoteAddr,
			"user-agent", r.UserAgent(),
		)

		next.ServeHTTP(rec, r)

		duration := time.Since(start).Milliseconds()

		switch {
		case rec.status >= 500:
			logger.Logger.Error(
				"request-completed",
				"request-id", reqID,
				"method", r.Method,
				"path", r.URL.Path,
				"status", rec.status,
				"duration-ms", duration,
			)
		case rec.status >= 400:
			logger.Logger.Warn(
				"request-completed",
				"request-id", reqID,
				"method", r.Method,
				"path", r.URL.Path,
				"status", rec.status,
				"duration-ms", duration,
			)
		default:
			logger.Logger.Info(
				"request-completed",
				"request-id", reqID,
				"method", r.Method,
				"path", r.URL.Path,
				"status", rec.status,
				"duration-ms", duration,
			)
		}
	})
}
