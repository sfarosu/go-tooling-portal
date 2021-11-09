package cmd

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/sfarosu/go-tooling-portal/cmd/helper"
	"github.com/sfarosu/go-tooling-portal/cmd/tmpl"
	"log"
	"net/http"
	"strconv"
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
	log.Println(r.URL.String(), r.Method, r.RemoteAddr, r.Proto, r.Header.Get("User-Agent"))
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

	length := r.FormValue("number")

	passLength, errConversion := strconv.Atoi(length) // convert length string to int
	if errConversion != nil {
		log.Println("error converting string to int: ", errConversion)
	}

	generatedRandomPassword, errRandom := helper.RandomString(passLength, upperBool, lowerBool, numbersBool, symbolsBool)
	if errRandom != nil {
		log.Println("error generating a random string: ", errRandom)
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

	log.Println(r.URL.String(), r.Method, r.RemoteAddr, r.Proto, r.Header.Get("User-Agent"))
	errExec := tmpl.Tpl.ExecuteTemplate(w, "passgen-process.html", data)
	if errExec != nil {
		log.Println("error executing template: ", errExec)
	}

	passGen.Inc()
}
