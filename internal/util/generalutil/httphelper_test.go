package generalutil

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/fstest"
)

func TestDisableDirListing(t *testing.T) {
	root := fstest.MapFS{
		"404.html":       &fstest.MapFile{Data: []byte("custom not found")},
		"dir/file.txt":   &fstest.MapFile{Data: []byte("file content")},
		"dir/index.html": &fstest.MapFile{Data: []byte("index content")},
		"empty/file.txt": &fstest.MapFile{Data: []byte("file content")},
	}

	// Handler that just writes "ok"
	finalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	middleware := DisableDirListing(root)(finalHandler)

	tests := []struct {
		name       string
		path       string
		wantStatus int
		wantBody   string
	}{
		{
			name:       "directory with index.html returns OK",
			path:       "/dir/",
			wantStatus: http.StatusOK,
			wantBody:   "ok",
		},
		{
			name:       "directory without index.html returns 404",
			path:       "/empty/",
			wantStatus: http.StatusNotFound,
			wantBody:   "custom not found",
		},
		{
			name:       "missing file returns custom 404",
			path:       "/missing.txt",
			wantStatus: http.StatusNotFound,
			wantBody:   "custom not found",
		},
		{
			name:       "file path passes through",
			path:       "/dir/file.txt",
			wantStatus: http.StatusOK,
			wantBody:   "ok",
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
			if body := rr.Body.String(); body != tt.wantBody {
				t.Errorf("got body [%v], want [%v]", body, tt.wantBody)
			}
		})
	}
}
