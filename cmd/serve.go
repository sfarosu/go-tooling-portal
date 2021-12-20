package cmd

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/sfarosu/go-tooling-portal/cmd/helper"
	"github.com/sfarosu/go-tooling-portal/cmd/version"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Serve() {
	fileServerAssets := http.FileServer(http.Dir("web/assets"))
	http.Handle("/assets/", http.StripPrefix("/assets", helper.DisableDirListing(fileServerAssets)))

	fileServerTmp := http.FileServer(http.Dir("web/tmp"))
	http.Handle("/tmp/", http.StripPrefix("/tmp", helper.DisableDirListing(fileServerTmp)))

	fileServerTemplates := http.FileServer(http.Dir("web/templates"))
	http.Handle("/templates/", http.StripPrefix("/templates", helper.DisableDirListing(fileServerTemplates)))

	http.Handle("/metrics", promhttp.Handler())

	http.HandleFunc("/", index)
	http.HandleFunc("/htpasswd", htpasswd)
	http.HandleFunc("/htpasswd-process", htpasswdProcess)
	http.HandleFunc("/passgen", passgen)
	http.HandleFunc("/passgen-process", passgenProcess)
	http.HandleFunc("/ssh", ssh)
	http.HandleFunc("/ssh-process-keygen", sshProcessKeypair)
	http.HandleFunc("/jsonprettify", jsonprettify)
	http.HandleFunc("/jsonprettify-process", jsonprettifyProcess)
	http.HandleFunc("/formatconvert", formatConvert)
	http.HandleFunc("/formatconvert-process", formatConvertProcess)

	// call AfterFunc 3 seconds after app startup to purge ssh keys from disc
	time.AfterFunc(3*time.Second, helper.KeysCleanup)

	hostname, _ := os.Hostname()
	appPath, _ := os.Getwd()
	log.Println("Tooling-portal " + version.Version + " started on host " + hostname + ":8080")
	log.Println("Application path: " + appPath)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err, "\nAnother process running on that port?")
	}
}
