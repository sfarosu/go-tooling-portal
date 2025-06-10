package server

import (
	"net/http"
	"time"

	"github.com/sfarosu/go-tooling-portal/internal/logger"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

// loggingMiddleware logs each HTTP request with appropriate log level.
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: 200}

		logger.Logger.Info(
			"incoming request",
			"method", r.Method,
			"path", r.URL.Path,
			"protocol", r.Proto,
			"remote_addr", r.RemoteAddr,
			"user_agent", r.UserAgent(),
		)

		next.ServeHTTP(rec, r)

		duration := time.Since(start).Milliseconds()

		switch {
		case rec.status >= 500:
			logger.Logger.Error(
				"request completed",
				"method", r.Method,
				"path", r.URL.Path,
				"status", rec.status,
				"duration_ms", duration,
			)
		case rec.status >= 400:
			logger.Logger.Warn(
				"request completed",
				"method", r.Method,
				"path", r.URL.Path,
				"status", rec.status,
				"duration_ms", duration,
			)
		default:
			logger.Logger.Info(
				"request completed",
				"method", r.Method,
				"path", r.URL.Path,
				"status", rec.status,
				"duration_ms", duration,
			)
		}
	})
}
