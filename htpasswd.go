package main

import (
    "net/http"
    "os/exec"
    "strings"
    "errors"
    "log"

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
    tpl.ExecuteTemplate(w, "htpasswd.html", nil)
}

func htpasswdProcess(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Redirect(w, r, "/htpasswd", http.StatusSeeOther)
    }

    uname := strings.ToLower(r.FormValue("username"))
    pass := r.FormValue("password")
    alg := r.FormValue("algorithm")

    /* call the generateHtpass function */
    ht, err := generateHtpass(uname, pass, alg)
    if err != nil {
        log.Println(err)
    }

    data := struct {
        Username string
        Password string
        Algorithm string
        Result string
    }{
        Username: uname,
        Password: pass,
        Algorithm: alg,
        Result: ht,
    }
    log.Println(r.URL.String(), r.Method, r.RemoteAddr, r.Proto, r.Header.Get("User-Agent"))
    tpl.ExecuteTemplate(w, "htpasswd-process.html", data)

    htpassGen.Inc()
}

func generateHtpass(uname string, pass string, alg string) (string, error) {
    /* error handling for username and password length */
    if uname == "" {
        return "", errors.New("You must specify a username")
    }
    if pass == "" {
        return "", errors.New("You must specify a password")
    }

    htpasswd, err := exec.Command("openssl", "passwd", alg, pass).Output()
    if err != nil {
        log.Println("Openssl command execution failed, please contact support !")
    }

    result := string(uname)+":"+string(htpasswd)

    return string(result), nil
}
