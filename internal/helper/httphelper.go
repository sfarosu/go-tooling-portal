package helper

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// DisableDirListing returns middleware that prevents directory listing for the given static root.
// If a requested path is a directory without an index.html, or if a file does not exist,
// it serves a global custom 404.html page from the web root with a 404 status.
// Otherwise, it allows the request to proceed to the next handler.
func DisableDirListing(root string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			localPath := filepath.Join(root, filepath.FromSlash(strings.TrimPrefix(r.URL.Path, "/")))
			fi, err := os.Stat(localPath)
			if err != nil {
				// Not found: serve global 404.html
				w.WriteHeader(http.StatusNotFound)
				http.ServeFile(w, r, filepath.Join("web", "404.html"))
				return
			}
			if fi.IsDir() {
				indexPath := filepath.Join(localPath, "index.html")
				if _, err := os.Stat(indexPath); os.IsNotExist(err) {
					// Directory without index.html: serve global 404.html
					w.WriteHeader(http.StatusNotFound)
					http.ServeFile(w, r, filepath.Join("web", "404.html"))
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}
