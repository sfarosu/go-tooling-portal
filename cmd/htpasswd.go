package cmd

import (
	"errors"
	"github.com/sfarosu/go-tooling-portal/cmd/tmpl"
	"log"
	"net/http"
	"os/exec"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	htpassGen = promauto.NewCounter(prometheus.CounterOpts{
		Name: "ht_passwords_generated_total",
		Help: "The total number of generated htpasswords",
	})
)

func htpasswd(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Redirect(w, r, "/htpasswd", http.StatusSeeOther)
	}
	log.Println(r.URL.String(), r.Method, r.RemoteAddr, r.Proto, r.Header.Get("User-Agent"))
	errExec := tmpl.Tpl.ExecuteTemplate(w, "htpasswd.html", nil)
	if errExec != nil {
		log.Println("error executing template: ", errExec)
	}
}

func htpasswdProcess(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/htpasswd", http.StatusSeeOther)
	}

	uname := strings.ToLower(r.FormValue("username"))
	pass := r.FormValue("password")
	alg := r.FormValue("algorithm")

	// call the generateHtpass function
	ht, errGenerateHtpass := generateHtpass(uname, pass, alg)
	if errGenerateHtpass != nil {
		log.Println("error generating htpassword", errGenerateHtpass)
	}

	data := struct {
		Username  string
		Password  string
		Algorithm string
		Result    string
	}{
		Username:  uname,
		Password:  pass,
		Algorithm: alg,
		Result:    ht,
	}

	log.Println(r.URL.String(), r.Method, r.RemoteAddr, r.Proto, r.Header.Get("User-Agent"))

	errExec := tmpl.Tpl.ExecuteTemplate(w, "htpasswd-process.html", data)
	if errExec != nil {
		log.Println("error executing template: ", errExec)
	}

	htpassGen.Inc()
}

func generateHtpass(uname string, pass string, alg string) (string, error) {
	// error handling for username and password length
	if uname == "" {
		return "", errors.New("you must specify a username")
	}
	if pass == "" {
		return "", errors.New("you must specify a password")
	}

	htpasswd, errCmd := exec.Command("openssl", "passwd", alg, pass).Output()
	if errCmd != nil {
		log.Println("Openssl command execution failed !")
	}

	result := string(uname) + ":" + string(htpasswd)

	return result, nil
}
