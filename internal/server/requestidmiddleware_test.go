package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_requestIDMiddleware(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "sets X-Request-ID header and context"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var called bool
			var gotReqID string

			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				called = true
				gotReqID = getRequestID(r.Context())
			})

			req := httptest.NewRequest("GET", "/", nil)
			rr := httptest.NewRecorder()

			requestIDMiddleware(handler).ServeHTTP(rr, req)

			if !called {
				t.Error("handler was not called")
			}
			headerID := rr.Header().Get("X-Request-ID")
			if headerID == "" {
				t.Error("X-Request-ID header not set")
			}
			if gotReqID == "" {
				t.Error("request ID not set in context")
			}
			if gotReqID != headerID {
				t.Errorf("context request ID [%v] does not match header [%v]", gotReqID, headerID)
			}
		})
	}
}

func Test_getRequestID(t *testing.T) {
	tests := []struct {
		name      string
		ctx       context.Context
		wantValue string
	}{
		{
			name:      "no request ID in context",
			ctx:       context.Background(),
			wantValue: "",
		},
		{
			name:      "request ID present in context",
			ctx:       context.WithValue(context.Background(), requestIDKey, "test-id-123"),
			wantValue: "test-id-123",
		},
		{
			name:      "wrong type in context",
			ctx:       context.WithValue(context.Background(), requestIDKey, 12345),
			wantValue: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getRequestID(tt.ctx)
			if got != tt.wantValue {
				t.Errorf("getRequestID() = [%v], want [%v]", got, tt.wantValue)
			}
		})
	}
}
