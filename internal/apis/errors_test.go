package apis

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestWrapHTMXValidationErrors(t *testing.T) {
	handler := WrapHTMXValidationErrors(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/problem+json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, _ = w.Write([]byte(`{"detail":"validation failed","errors":[{"message":"expected required property input to be present"}]}`))
	}))

	tests := []struct {
		name       string
		path       string
		hxRequest  string
		wantStatus int
		wantHTML   bool
	}{
		{
			name:       "HTMX endpoint request",
			path:       "/api/tool/htmx",
			hxRequest:  "true",
			wantStatus: http.StatusOK,
			wantHTML:   true,
		},
		{
			name:       "regular API endpoint",
			path:       "/api/tool",
			hxRequest:  "true",
			wantStatus: http.StatusUnprocessableEntity,
		},
		{
			name:       "HTMX endpoint without HTMX header",
			path:       "/api/tool/htmx",
			wantStatus: http.StatusUnprocessableEntity,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, test.path, nil)
			if test.hxRequest != "" {
				req.Header.Set("HX-Request", test.hxRequest)
			}
			recorder := httptest.NewRecorder()
			handler.ServeHTTP(recorder, req)

			if recorder.Code != test.wantStatus {
				t.Fatalf("expected status %d, got %d: %s", test.wantStatus, recorder.Code, recorder.Body.String())
			}
			if test.wantHTML {
				if !strings.Contains(recorder.Header().Get("Content-Type"), "text/html") {
					t.Fatalf("expected HTML content type, got %q", recorder.Header().Get("Content-Type"))
				}
				if !strings.Contains(recorder.Body.String(), "expected required property input to be present") {
					t.Fatalf("expected validation detail in fragment, got %s", recorder.Body.String())
				}
			} else if recorder.Header().Get("Content-Type") != "application/problem+json" {
				t.Fatalf("expected original problem JSON content type, got %q", recorder.Header().Get("Content-Type"))
			}
		})
	}
}

func TestGenerateHTMXErrorEscapesHTML(t *testing.T) {
	result := generateHTMXError(fmt.Errorf("<script>alert(1)</script>"))
	if strings.Contains(result, "<script>") {
		t.Fatalf("error fragment contains unescaped markup: %s", result)
	}
	if !strings.Contains(result, "&lt;script&gt;alert(1)&lt;/script&gt;") {
		t.Errorf("expected escaped error message in fragment, got: %s", result)
	}
}
