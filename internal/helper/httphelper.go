package helper

import (
	"io/fs"
	"net/http"
	"path"
	"strconv"
	"strings"
)

// DisableDirListing returns middleware that prevents directory listing for the given static root.
// If a requested path is a directory without an index.html, or if a file does not exist,
// it serves the root's custom 404.html page with a 404 status.
func DisableDirListing(root fs.FS) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestPath := strings.TrimPrefix(path.Clean("/"+r.URL.Path), "/")
			if requestPath == "" {
				requestPath = "."
			}
			fi, err := fs.Stat(root, requestPath)
			if err != nil {
				serveNotFound(w, r, root)
				return
			}
			if fi.IsDir() {
				if _, err := fs.Stat(root, path.Join(requestPath, "index.html")); err != nil {
					serveNotFound(w, r, root)
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}

func serveNotFound(w http.ResponseWriter, r *http.Request, root fs.FS) {
	content, err := fs.ReadFile(root, "404.html")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Content-Length", strconv.Itoa(len(content)))
	w.WriteHeader(http.StatusNotFound)
	if r.Method != http.MethodHead {
		_, _ = w.Write(content)
	}
}
