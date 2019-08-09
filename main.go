package main

import (
	"html/template"
	"net/http"
	"log"
	"os"
	"strings"
	"time"
)

// BuildVersion for the app
const BuildVersion string = "version 1.0"

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*html"))
	log.SetOutput(os.Stdout) //Change the device for logging to stdout
}

func index(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Redirect(w, r, "/index", http.StatusSeeOther)
	}
	
	data := struct {
		AppVersion string
	}{
		AppVersion: BuildVersion,
	}
	
	log.Println(r.URL.String(), r.Method, r.RemoteAddr, r.Proto, r.Header.Get("User-Agent"))
	tpl.ExecuteTemplate(w, "index.html", data)
}

func main() {
	fileServerAssets := http.FileServer(http.Dir("./assets"))
	http.Handle("/assets/", http.StripPrefix("/assets", disableDirListing(fileServerAssets)))

	fileServerTmp := http.FileServer(http.Dir("./tmp"))
	http.Handle("/tmp/", http.StripPrefix("/tmp", disableDirListing(fileServerTmp)))

	http.HandleFunc("/", index)
	http.HandleFunc("/htpasswd", htpasswd)
	http.HandleFunc("/htpasswd-process", htpasswdProcess)
	http.HandleFunc("/passgen", passgen)
	http.HandleFunc("/passgen-process", passgenProcess)
	http.HandleFunc("/ssh", ssh)
	http.HandleFunc("/ssh-process-keygen", sshProcessKeypair)

	/* call the keysCleanup() function to purge any keys from disc */
	clean := keysCleanup
	time.AfterFunc(3 * time.Second, clean) /* call AfterFunc 3 seconds after app startup */

	hostname, _ := os.Hostname()
	appPath, _ := os.Getwd()
	log.Println("Tooling-portal " + BuildVersion + " started on host " + hostname +":8080")
	log.Println("Application path: " + appPath)

	

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err, "\nAnother process running on that port?")
	}
}

func disableDirListing(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if strings.HasSuffix(r.URL.Path, "/") {
            http.Redirect(w, r, "/index", http.StatusSeeOther)
            return
        }
        next.ServeHTTP(w, r)
    })
}