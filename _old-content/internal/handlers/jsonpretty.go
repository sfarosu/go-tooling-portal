package handlers

import (
	"log"
	"net/http"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/sfarosu/go-tooling-portal/internal/helper"
	"github.com/sfarosu/go-tooling-portal/internal/tmpl"
)

var (
	jsonPrettify = promauto.NewCounter(prometheus.CounterOpts{
		Name: "json_prettified_total",
		Help: "The total number of prettified jsons",
	})
)

func jsonprettify(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Redirect(w, r, "/jsonprettify", http.StatusSeeOther)
	}
	log.Println(r.Method, r.URL.String(), r.Proto, r.RemoteAddr, r.Header.Get("User-Agent"))
	errExec := tmpl.Tpl.ExecuteTemplate(w, "jsonprettify.html", nil)
	if errExec != nil {
		log.Println("error executing template: ", errExec)
	}
}

func jsonprettifyProcess(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/jsonprettify", http.StatusSeeOther)
	}

	result := helper.PrettyJSON(strings.TrimSpace(r.FormValue("text")))

	data := struct {
		Text   string
		Result string
	}{
		Text:   r.FormValue("text"),
		Result: result.String(),
	}

	log.Println(r.Method, r.URL.String(), r.Proto, r.RemoteAddr, r.Header.Get("User-Agent"))
	errExec := tmpl.Tpl.ExecuteTemplate(w, "jsonprettify-process.html", data)
	if errExec != nil {
		log.Println("error executing template: ", errExec)
	}

	jsonPrettify.Inc()
}
