package cmd

import (
	"errors"
	"fmt"
	"github.com/sfarosu/go-tooling-portal/cmd/helper"
	"github.com/sfarosu/go-tooling-portal/cmd/tmpl"
	"log"
	"net/http"
	"os/exec"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	sshKeyGen = promauto.NewCounter(prometheus.CounterOpts{
		Name: "ssh_key_generated_total",
		Help: "The total number of generated ssh keys",
	})
)

func ssh(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Redirect(w, r, "/ssh", http.StatusSeeOther)
	}
	log.Println(r.URL.String(), r.Method, r.RemoteAddr, r.Proto, r.Header.Get("User-Agent"))
	errExec := tmpl.Tpl.ExecuteTemplate(w, "ssh.html", nil)
	if errExec != nil {
		log.Println("error executing template: ", errExec)
	}
}

func sshProcessKeypair(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/ssh", http.StatusSeeOther)
	}
	email := r.FormValue("email")
	pass := r.FormValue("password")
	usePass := r.FormValue("usepass")
	bitSize := r.FormValue("bitsize")

	privateKeyContent, publicKeyContent, privateKeyFileName, publicKeyFileName, errGenerateKeyPair := generateKeyPair(email, pass, usePass, bitSize)
	if errGenerateKeyPair != nil {
		log.Println("error generating ssh key pair: ", errGenerateKeyPair)
	}

	data := struct {
		Email              string
		Pass               string
		UsePass            string
		BitSize            string
		PrivateKeyContent  string
		PublicKeyContent   string
		PrivateKeyFileName string
		PublicKeyFileName  string
	}{
		Email:              email,
		Pass:               pass,
		UsePass:            usePass,
		BitSize:            bitSize,
		PrivateKeyContent:  privateKeyContent,
		PublicKeyContent:   publicKeyContent,
		PrivateKeyFileName: privateKeyFileName,
		PublicKeyFileName:  publicKeyFileName,
	}

	log.Println(r.URL.String(), r.Method, r.RemoteAddr, r.Proto, r.Header.Get("User-Agent"))

	errExec := tmpl.Tpl.ExecuteTemplate(w, "ssh-process-keygen.html", data)
	if errExec != nil {
		log.Println("error executing template: ", errExec)
	}

	// call the KeysCleanup() function to delete the keys older than 5 minutes
	time.AfterFunc(5*time.Minute, helper.KeysCleanup)

	sshKeyGen.Inc()
}

func generateKeyPair(email string, pass string, usePass string, bitSize string) (string, string, string, string, error) {
	// error handling for email and password length
	if email == "" {
		return "", "", "", "", errors.New("email field can't be empty")
	}

	if usePass == "yes" {
		if (pass == "") || (len(pass) < 8) || (len(pass) > 64) {
			return "", "", "", "", errors.New("you must specify a password and it must be at least 8 characters long")
		}
	}

	// generate a randomNumber to make sure keys filenames are unique
	randomNumber, errRandom := helper.RandomString(16, false, false, true, false)
	if errRandom != nil {
		log.Println("error generating a random string: ", errRandom)
	}

	// generate the key pair and place it under tmp/ with the format id_rsa-randomNumber / id_rsa-sameRandomNumber.pub
	keygenCmd := exec.Command("ssh-keygen", "-t", "rsa", "-b", bitSize, "-C", email, "-N", pass, "-f", "web/tmp/id_rsa-"+randomNumber, "-q")
	outputKeygenCmd, errCmd := keygenCmd.CombinedOutput()
	if errCmd != nil {
		log.Println(fmt.Sprint(errCmd) + ": " + string(outputKeygenCmd))
	}

	privateKeyContent := helper.ReadFile("web/tmp/id_rsa-" + randomNumber)
	publicKeyContent := helper.ReadFile("web/tmp/id_rsa-" + randomNumber + ".pub")

	privateKeyFileName := "id_rsa-" + randomNumber
	publicKeyFileName := "id_rsa-" + randomNumber + ".pub"

	return string(privateKeyContent), string(publicKeyContent), privateKeyFileName, publicKeyFileName, nil
}
