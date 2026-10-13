package server

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/sfarosu/go-tooling-portal/internal/requestid"
)

// requestIDMiddleware generates a unique request ID for each incoming HTTP request
func requestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := uuid.NewString()
		ctx := requestid.WithContext(r.Context(), reqID)
		r = r.WithContext(ctx)
		w.Header().Set("X-Request-ID", reqID)
		next.ServeHTTP(w, r)
	})
}
