package helper

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// DisableDirListing returns a middleware that disables directory listing for the given root directory.
func DisableDirListing(root string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Build the local path
			localPath := filepath.Join(root, filepath.FromSlash(strings.TrimPrefix(r.URL.Path, "/")))
			fi, err := os.Stat(localPath)
			if err == nil && fi.IsDir() {
				indexPath := filepath.Join(localPath, "index.html")
				if _, err := os.Stat(indexPath); os.IsNotExist(err) {
					http.NotFound(w, r)
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}
