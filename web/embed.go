package web

import "embed"

// Files contains the static files served by the application.
//
//go:embed *.html *.ico assets
var Files embed.FS
