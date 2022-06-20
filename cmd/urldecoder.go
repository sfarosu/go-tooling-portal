package cmd

import (
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/sfarosu/go-tooling-portal/cmd/helper"
	"github.com/sfarosu/go-tooling-portal/cmd/tmpl"
)

type URLDecodedData struct {
	Scheme    string `yaml:"Scheme"`
	Host      string `yaml:"Host"`
	Path      string `yaml:"Path"`
	QueryArgs []struct {
		Key   string `yaml:"Key"`
		Value string `yaml:"Value"`
	} `yaml:"QueryArgs"`
}

var (
	urlDecoder = promauto.NewCounter(prometheus.CounterOpts{
		Name: "urls_decoded_total",
		Help: "The total number of urls decoded",
	})
)

func urlDecode(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Redirect(w, r, "/urldecode", http.StatusSeeOther)
	}
	log.Println(r.Method, r.URL.String(), r.Proto, r.RemoteAddr, r.Header.Get("User-Agent"))
	errExec := tmpl.Tpl.ExecuteTemplate(w, "urldecode.html", nil)
	if errExec != nil {
		log.Println("error executing template: ", errExec)
	}
}

func urlDecodeProcess(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/urldecode-process", http.StatusSeeOther)
	}

	result := decode(strings.TrimSpace(r.FormValue("text")))

	data := struct {
		InsertedText string
		Result       string
	}{
		InsertedText: r.FormValue("text"),
		Result:       string(result),
	}

	log.Println(r.Method, r.URL.String(), r.Proto, r.RemoteAddr, r.Header.Get("User-Agent"))
	errExec := tmpl.Tpl.ExecuteTemplate(w, "urldecode-process.html", data)
	if errExec != nil {
		log.Println("error executing template: ", errExec)
	}

	urlDecoder.Inc()
}

func decode(insertedData string) []byte {
	parsedURL, err := url.Parse(insertedData)
	if err != nil {
		return []byte("Error decoding data: " + err.Error())
	}

	populateDecodedData := URLDecodedData{}
	populateDecodedData.Scheme = parsedURL.Scheme
	populateDecodedData.Host = parsedURL.Host
	populateDecodedData.Path = parsedURL.Path
	for k, v := range parsedURL.Query() {
		populateDecodedData.QueryArgs = append(populateDecodedData.QueryArgs, struct {
			Key   string `yaml:"Key"`
			Value string `yaml:"Value"`
		}{k, v[0]})
	}

	return helper.MarshalYAML(populateDecodedData)
}
