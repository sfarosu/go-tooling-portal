package cmd

import (
	"log"
	"net/http"
	"os/exec"
	"strings"

	"github.com/sfarosu/go-tooling-portal/cmd/tmpl"

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
	log.Println(r.Method, r.URL.String(), r.Proto, r.RemoteAddr, r.Header.Get("User-Agent"))
	errExec := tmpl.Tpl.ExecuteTemplate(w, "htpasswd.html", nil)
	if errExec != nil {
		log.Println("error executing template: ", errExec)
	}
}

func htpasswdProcess(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/htpasswd", http.StatusSeeOther)
	}

	ht, err := generateHtpass(strings.ToLower(strings.TrimSpace(r.FormValue("username"))), strings.TrimSpace(r.FormValue("password")), r.FormValue("algorithm"))
	if err != nil {
		log.Println("error generating htpassword", err)
	}

	data := struct {
		Username  string
		Password  string
		Algorithm string
		Result    string
	}{
		Username:  strings.ToLower(r.FormValue("username")),
		Password:  r.FormValue("password"),
		Algorithm: r.FormValue("algorithm"),
		Result:    ht,
	}

	log.Println(r.Method, r.URL.String(), r.Proto, r.RemoteAddr, r.Header.Get("User-Agent"))

	errExec := tmpl.Tpl.ExecuteTemplate(w, "htpasswd-process.html", data)
	if errExec != nil {
		log.Println("error executing template: ", errExec)
	}

	htpassGen.Inc()
}

func generateHtpass(uname string, pass string, alg string) (string, error) {
	htpasswd, errCmd := exec.Command("openssl", "passwd", alg, pass).Output()
	if errCmd != nil {
		log.Println("Openssl command execution failed !")
	}

	result := string(uname) + ":" + string(htpasswd)

	return result, nil
}
