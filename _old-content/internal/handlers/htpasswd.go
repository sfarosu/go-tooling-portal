package handlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/sfarosu/go-tooling-portal/internal/tmpl"

	"github.com/GehirnInc/crypt/apr1_crypt"
	"github.com/GehirnInc/crypt/md5_crypt"
	"github.com/GehirnInc/crypt/sha256_crypt"
	"github.com/GehirnInc/crypt/sha512_crypt"
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

func generateHtpass(username string, password string, algorithm string) (string, error) {
	switch algorithm {
	case "apr1":
		hash, err := apr1crypt(username, password)
		if err != nil {
			return "", err
		}
		return hash, nil
	case "1":
		hash, err := md5crypt(username, password)
		if err != nil {
			return "", err
		}
		return hash, nil
	case "5":
		hash, err := sha256crypt(username, password)
		if err != nil {
			return "", err
		}
		return hash, nil
	case "6":
		hash, err := sha512crypt(username, password)
		if err != nil {
			return "", err
		}
		return hash, nil
	default:
		return "", errors.New("unsupported algorithm; use [apr1], [1], [5] or [6] openssl cryptographic options")
	}
}

func apr1crypt(username string, password string) (string, error) {
	hash, err := apr1_crypt.New().Generate([]byte(password), nil)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s:%s", username, hash), nil
}

func md5crypt(username string, password string) (string, error) {
	hash, err := md5_crypt.New().Generate([]byte(password), nil)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s:%s", username, hash), nil
}

func sha256crypt(username string, password string) (string, error) {
	hash, err := sha256_crypt.New().Generate([]byte(password), nil)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s:%s", username, hash), nil
}

func sha512crypt(username string, password string) (string, error) {
	hash, err := sha512_crypt.New().Generate([]byte(password), nil)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s:%s", username, hash), nil
}
