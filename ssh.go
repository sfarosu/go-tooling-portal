package main

import (
	"net/http"
	"fmt"
	"errors"
	"io/ioutil"
	"os/exec"
	"log"
	"time"
)

func ssh(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Redirect(w, r, "/ssh", http.StatusSeeOther)
	}
	log.Println(r.URL.String(), r.Method, r.RemoteAddr, r.Proto, r.Header.Get("User-Agent"))
	tpl.ExecuteTemplate(w, "ssh.html", nil)
}

func sshProcessKeypair(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/ssh", http.StatusSeeOther)
	}
	email := r.FormValue("email")
	pass := r.FormValue("password")
	usePass := r.FormValue("usepass")
	
	privateKeyContent, publicKeyContent, privateKeyFileName, publicKeyFileName, err := generateKeyPair(email, pass, usePass)
	if err != nil {
		log.Println(err)
	}

	data := struct {
		Email string
		Pass string
		UsePass string
		PrivateKeyContent string
		PublicKeyContent string
		PrivateKeyFileName string
		PublicKeyFileName string
	}{
		Email: email,
		Pass: pass,
		UsePass: usePass,
		PrivateKeyContent: privateKeyContent,
		PublicKeyContent: publicKeyContent,
		PrivateKeyFileName: privateKeyFileName,
		PublicKeyFileName: publicKeyFileName,
	}

	log.Println(r.URL.String(), r.Method, r.RemoteAddr, r.Proto, r.Header.Get("User-Agent"))
	tpl.ExecuteTemplate(w, "ssh-process-keygen.html", data)

	/* call the keysCleanup() function to delete the keys older than 5 minutes */
	clean := keysCleanup
	time.AfterFunc(5 * time.Minute, clean)
}

/* Helper functions */
func generateKeyPair(email string, pass string, usePass string) (string, string, string, string, error) {
	/* error handling for email and password length */
	if email == "" {
		return "", "", "", "", errors.New("Email field can't be empty")
	}
	
	if usePass == "yes" {
		if (pass == "") || (len(pass) < 8) || (len(pass) > 64) {
			return "", "", "", "", errors.New("You must specify a password and it must be at least 8 characters long")
		}
	}
	
	/* generate a randomnumber to make sure keys are unique */
	randomNumber, err := randomString(16, false, false, true, false)
	if err != nil {
		log.Println(err)
	}

	/* generate the key pair and place it under tmp/ with the format id_rsa-randomnumber / id_rsa-samerandomnumber.pub */
	keygenCmd := exec.Command("ssh-keygen", "-t", "rsa", "-b", "4096", "-C", email, "-N", pass, "-f", "tmp/id_rsa-"+randomNumber, "-q")
	outputKeygenCmd, err := keygenCmd.CombinedOutput()
	if err != nil {
		log.Println(fmt.Sprint(err) + ": " + string(outputKeygenCmd))
	}

	privateKeyContent, err := ioutil.ReadFile("tmp/id_rsa-"+randomNumber)
	if err != nil {
		log.Println(err)
		
	}
	publicKeyContent, err := ioutil.ReadFile("tmp/id_rsa-"+randomNumber+".pub")
	if err != nil {
		log.Println(err)
	}

	privateKeyFileName := "id_rsa-"+randomNumber
	publicKeyFileName := "id_rsa-"+randomNumber+".pub"

	return string(privateKeyContent), string(publicKeyContent), string(privateKeyFileName), string(publicKeyFileName), nil	
}

func keysCleanup() {
	/* cleanup the tmp folder; key age deletion is defined in the time.Afterfunc function */
	cleanupCmd := exec.Command("find", "tmp/", "-type", "f", "-name", "id_rsa*", "-mmin", "+0", "-exec", "rm", "{}", ";")
	outputCleanupCmd, err := cleanupCmd.CombinedOutput()
	if err != nil {
		log.Println(fmt.Sprint(err) + ": " + string(outputCleanupCmd))
	}
	log.Println("SSH keys stored locally have been purged !")
}
