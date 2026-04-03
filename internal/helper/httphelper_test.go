package helper

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestDisableDirListing(t *testing.T) {
	// Create a temporary directory structure for testing
	tmpDir := t.TempDir()
	// Create a directory and a file inside it
	dirPath := filepath.Join(tmpDir, "dir")
	filePath := filepath.Join(dirPath, "file.txt")
	indexPath := filepath.Join(dirPath, "index.html")
	os.Mkdir(dirPath, 0755)
	os.WriteFile(filePath, []byte("file content"), 0644)
	os.WriteFile(indexPath, []byte("index content"), 0644)

	// Handler that just writes "ok"
	finalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	middleware := DisableDirListing(tmpDir)(finalHandler)

	tests := []struct {
		name       string
		path       string
		setupIndex bool
		wantStatus int
		wantBody   string
	}{
		{
			name:       "directory with index.html returns OK",
			path:       "/dir/",
			setupIndex: true,
			wantStatus: http.StatusOK,
			wantBody:   "ok",
		},
		{
			name:       "directory without index.html returns 404",
			path:       "/dir/",
			setupIndex: false,
			wantStatus: http.StatusNotFound,
			wantBody:   "404 page not found\n",
		},
		{
			name:       "file path passes through",
			path:       "/dir/file.txt",
			setupIndex: true,
			wantStatus: http.StatusOK,
			wantBody:   "ok",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Remove or create index.html as needed
			if tt.setupIndex {
				os.WriteFile(indexPath, []byte("index content"), 0644)
			} else {
				os.Remove(indexPath)
			}

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
