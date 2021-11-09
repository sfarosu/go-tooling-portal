package cmd

import (
	"github.com/sfarosu/go-tooling-portal/cmd/tmpl"
	"github.com/sfarosu/go-tooling-portal/cmd/version"
	"log"
	"net/http"
)

func index(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Redirect(w, r, "/index", http.StatusSeeOther)
	}

	data := struct {
		AppVersion string
	}{
		AppVersion: version.Version,
	}

	log.Println(r.URL.String(), r.Method, r.RemoteAddr, r.Proto, r.Header.Get("User-Agent"))
	errExec := tmpl.Tpl.ExecuteTemplate(w, "index.html", data)
	if errExec != nil {
		log.Println("error executing template: ", errExec)
	}
}
