package httpserver

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sfarosu/go-tooling-portal/internal/requestid"
)

func Test_requestIDMiddleware(t *testing.T) {
	var called bool
	var contextID string

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		contextID = requestid.FromContext(r.Context())
	})

	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	requestIDMiddleware(handler).ServeHTTP(rr, req)

	if !called {
		t.Fatal("handler was not called")
	}
	headerID := rr.Header().Get("X-Request-ID")
	if headerID == "" {
		t.Fatal("X-Request-ID header not set")
	}
	if contextID != headerID {
		t.Errorf("request IDs do not match: context=%q response=%q", contextID, headerID)
	}
}
