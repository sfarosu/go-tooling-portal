package cmd

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/sfarosu/go-tooling-portal/cmd/helper"
	"github.com/sfarosu/go-tooling-portal/cmd/tmpl"
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
	log.Println(r.URL.String(), r.Method, r.RemoteAddr, r.Proto, r.Header.Get("User-Agent"))
	errExec := tmpl.Tpl.ExecuteTemplate(w, "jsonprettify.html", nil)
	if errExec != nil {
		log.Println("error executing template: ", errExec)
	}
}

func jsonprettifyProcess(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/jsonprettify", http.StatusSeeOther)
	}

	insertedText := r.FormValue("text")

	result := helper.PrettyJSON(insertedText)

	data := struct {
		Text   string
		Result string
	}{
		Text:   insertedText,
		Result: result.String(),
	}

	log.Println(r.URL.String(), r.Method, r.RemoteAddr, r.Proto, r.Header.Get("User-Agent"))
	errExec := tmpl.Tpl.ExecuteTemplate(w, "jsonprettify-process.html", data)
	if errExec != nil {
		log.Println("error executing template: ", errExec)
	}

	jsonPrettify.Inc()
}
