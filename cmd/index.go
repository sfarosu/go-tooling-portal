package cmd

import (
	"log"
	"net/http"

	"github.com/sfarosu/go-tooling-portal/cmd/tmpl"
	"github.com/sfarosu/go-tooling-portal/cmd/version"
)

func index(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Redirect(w, r, "/index", http.StatusSeeOther)
	}

	data := struct {
		Version string
	}{
		Version: version.Version,
	}

	log.Println(r.Method, r.URL.String(), r.Proto, r.RemoteAddr, r.Header.Get("User-Agent"))
	errExec := tmpl.Tpl.ExecuteTemplate(w, "index.html", data)
	if errExec != nil {
		log.Println("error executing template: ", errExec)
	}
}
