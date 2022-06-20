package cmd

import (
	b64 "encoding/base64"
	"log"
	"net/http"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/sfarosu/go-tooling-portal/cmd/tmpl"
)

var (
	base64Converter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "base64_converted_total",
		Help: "The total number of base64 converted",
	})
)

func base64convert(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Redirect(w, r, "/base64convert", http.StatusSeeOther)
	}
	log.Println(r.Method, r.URL.String(), r.Proto, r.RemoteAddr, r.Header.Get("User-Agent"))
	errExec := tmpl.Tpl.ExecuteTemplate(w, "base64convert.html", nil)
	if errExec != nil {
		log.Println("error executing template: ", errExec)
	}
}

func base64convertProcess(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/base64convert-process", http.StatusSeeOther)
	}

	result := base64EncDec(strings.TrimSpace(r.FormValue("text")), r.FormValue("operation"), r.FormValue("format"))

	data := struct {
		InsertedText string
		Result       string
	}{
		InsertedText: r.FormValue("text"),
		Result:       result,
	}

	log.Println(r.Method, r.URL.String(), r.Proto, r.RemoteAddr, r.Header.Get("User-Agent"))
	errExec := tmpl.Tpl.ExecuteTemplate(w, "base64convert-process.html", data)
	if errExec != nil {
		log.Println("error executing template: ", errExec)
	}

	base64Converter.Inc()
}

func base64EncDec(insertedData, operation, format string) string {
	var processedData string

	switch operation {
	case "decode":
		if format == "standard" {
			r, err := b64.StdEncoding.DecodeString(insertedData)
			if err != nil {
				processedData = "Error decoding in standard format: " + err.Error()
			} else {
				processedData = string(r)
			}
		} else if format == "url-compatible" {
			r, err := b64.URLEncoding.DecodeString(insertedData)
			if err != nil {
				processedData = "Error decoding in url-compatible format: " + err.Error()
			} else {
				processedData = string(r)
			}
		}
	case "encode":
		if format == "standard" {
			processedData = b64.StdEncoding.EncodeToString([]byte(insertedData))
		} else if format == "url-compatible" {
			processedData = b64.URLEncoding.EncodeToString([]byte(insertedData))
		}
	}

	return processedData
}
