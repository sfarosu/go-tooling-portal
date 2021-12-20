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
	formatConverter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "formats_converted_total",
		Help: "The total number of formats converted",
	})
)

func formatConvert(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Redirect(w, r, "/formatconvert", http.StatusSeeOther)
	}
	log.Println(r.URL.String(), r.Method, r.RemoteAddr, r.Proto, r.Header.Get("User-Agent"))
	errExec := tmpl.Tpl.ExecuteTemplate(w, "formatconvert.html", nil)
	if errExec != nil {
		log.Println("error executing template: ", errExec)
	}
}

func formatConvertProcess(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/formatconvert", http.StatusSeeOther)
	}

	insertedText := r.FormValue("text")

	result := transformToFormat(insertedText)

	data := struct {
		InsertedText string
		Result       string
	}{
		InsertedText: insertedText,
		Result:       result,
	}

	log.Println(r.URL.String(), r.Method, r.RemoteAddr, r.Proto, r.Header.Get("User-Agent"))
	errExec := tmpl.Tpl.ExecuteTemplate(w, "formatconvert-process.html", data)
	if errExec != nil {
		log.Println("error executing template: ", errExec)
	}

	formatConverter.Inc()
}

func transformToFormat(insertedText string) string {
	jsonData, errJSON := helper.UnmarshalJSON([]byte(insertedText))
	yamlData, errYAML := helper.UnmarshalYAML([]byte(insertedText))

	var transformedData []byte

	// determine if the input was json or yaml by unmarshaling both and see which throws an error
	// remember, yaml.unmarshal on a json does NOT throw error
	switch errJSON == nil && errYAML == nil {
	case true:
		transformedData = helper.MarshalYAML(jsonData)
	case false:
		transformedData = helper.MarshalJSON(yamlData)
	}

	return string(transformedData)
}
