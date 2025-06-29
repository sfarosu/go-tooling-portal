package helper

import (
	"net/http"
	"strings"
)

// DisableDirListing is a middleware that disables directory listing by redirecting requests
func DisableDirListing(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			http.Redirect(w, r, "/index", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}
