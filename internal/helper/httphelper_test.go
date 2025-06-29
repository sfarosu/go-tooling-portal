package helper

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDisableDirListing(t *testing.T) {
	// Handler that just writes "ok"
	finalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	middleware := DisableDirListing(finalHandler)

	tests := []struct {
		name         string
		path         string
		wantStatus   int
		wantLocation string
		wantBody     string
	}{
		{
			name:         "directory path gets redirected",
			path:         "/some/dir/",
			wantStatus:   http.StatusSeeOther,
			wantLocation: "/index",
			wantBody:     "",
		},
		{
			name:         "file path passes through",
			path:         "/some/dir/file.txt",
			wantStatus:   http.StatusOK,
			wantLocation: "",
			wantBody:     "ok",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.path, nil)
			rr := httptest.NewRecorder()

			middleware.ServeHTTP(rr, req)

			if rr.Code != tt.wantStatus {
				t.Errorf("got status [%v], want [%v]", rr.Code, tt.wantStatus)
			}
			if loc := rr.Header().Get("Location"); loc != tt.wantLocation {
				t.Errorf("got Location [%v], want [%v]", loc, tt.wantLocation)
			}
			// Only check body for non-redirect responses
			if tt.wantStatus != http.StatusSeeOther {
				if body := rr.Body.String(); body != tt.wantBody {
					t.Errorf("got body [%v], want [%v]", body, tt.wantBody)
				}
			}
		})
	}
}
