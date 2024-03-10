package handlers

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"runtime"
	"time"

	"github.com/sfarosu/go-tooling-portal/internal/helper"
	"github.com/sfarosu/go-tooling-portal/internal/tmpl"
	"github.com/sfarosu/go-tooling-portal/internal/version"
	"go.uber.org/automaxprocs/maxprocs"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Serve() {
	// Flags
	addr := flag.String("addr", ":8080", "Network address and port to start on")
	flag.Parse()

	// GOMAXPROCS - respect k8s cpu quota
	_, errMax := maxprocs.Set()
	if errMax != nil {
		log.Printf("Error setting maxprocs: %v", errMax)
	}

	tmpl.Tpl = template.Must(template.ParseGlob("web/templates/*html"))

	// Http handlers
	http.Handle("/assets/", http.StripPrefix("/assets", helper.DisableDirListing(http.FileServer(http.Dir("web/assets")))))
	http.Handle("/tmp/", http.StripPrefix("/tmp", helper.DisableDirListing(http.FileServer(http.Dir("web/tmp")))))
	http.Handle("/templates/", http.StripPrefix("/templates", helper.DisableDirListing(http.FileServer(http.Dir("web/templates")))))

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
	http.HandleFunc("/timeconvert", timeconvert)
	http.HandleFunc("/timeconvert-process", timeconvertProcess)
	http.HandleFunc("/base64convert", base64convert)
	http.HandleFunc("/base64convert-process", base64convertProcess)
	http.HandleFunc("/urldecode", urlDecode)
	http.HandleFunc("/urldecode-process", urlDecodeProcess)

	// Call AfterFunc 3 seconds after app startup to purge ssh keys from disc
	time.AfterFunc(3*time.Second, helper.KeysCleanup)

	// Startup logging
	log.Printf("Program listening on '%v', path '%v', GOMAXPROCS '%v'", *addr, helper.CurrentWorkingDirectory(), runtime.GOMAXPROCS(0))
	log.Printf("Version '%v', BuildDate '%v', GitShortHash '%v', GoVersion '%v'", version.Version, version.BuildDate, version.GitShortHash, runtime.Version())

	// Start http server
	srv := &http.Server{
		Addr:         *addr,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatalf("Error starting the http server on '%v': %v", *addr, err)
	}
}
