package cmd

import (
	"log"
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/sfarosu/go-tooling-portal/cmd/helper"
	"github.com/sfarosu/go-tooling-portal/cmd/tmpl"
)

var (
	passGen = promauto.NewCounter(prometheus.CounterOpts{
		Name: "passwords_generated_total",
		Help: "The total number of generated passwords",
	})
)

func passgen(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Redirect(w, r, "/passgen", http.StatusSeeOther)
	}
	log.Println(r.Method, r.URL.String(), r.Proto, r.RemoteAddr, r.Header.Get("User-Agent"))
	errExec := tmpl.Tpl.ExecuteTemplate(w, "passgen.html", nil)
	if errExec != nil {
		log.Println("error executing template: ", errExec)
	}
}

func passgenProcess(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/passgen", http.StatusSeeOther)
	}

	// Get the checkbox status from the user and convert it to bool before sending it to the randomString() function
	upper := r.FormValue("uppercase")
	upperBool := false
	if len(upper) != 0 {
		upperBool = true
	}

	lower := r.FormValue("lowercase")
	lowerBool := false
	if len(lower) != 0 {
		lowerBool = true
	}

	numbers := r.FormValue("numbers")
	numbersBool := false
	if len(numbers) != 0 {
		numbersBool = true
	}

	symbols := r.FormValue("symbols")
	symbolsBool := false
	if len(symbols) != 0 {
		symbolsBool = true
	}

	length := r.FormValue("length")

	passLength, err := strconv.Atoi(length)
	if err != nil {
		log.Println("error converting string to int: ", err)
	}

	generatedRandomPassword, err := helper.RandomString(passLength, upperBool, lowerBool, numbersBool, symbolsBool)
	if err != nil {
		log.Println("error generating a random string: ", err)
	}

	data := struct {
		Uppercase string
		Lowercase string
		Numbers   string
		Symbols   string
		Length    string
		Result    string
	}{
		Uppercase: upper,
		Lowercase: lower,
		Numbers:   numbers,
		Symbols:   symbols,
		Length:    length,
		Result:    generatedRandomPassword,
	}

	log.Println(r.Method, r.URL.String(), r.Proto, r.RemoteAddr, r.Header.Get("User-Agent"))
	errExec := tmpl.Tpl.ExecuteTemplate(w, "passgen-process.html", data)
	if errExec != nil {
		log.Println("error executing template: ", errExec)
	}

	passGen.Inc()
}
