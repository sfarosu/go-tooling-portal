package main

import (
	"github.com/sfarosu/go-tooling-portal/cmd"
	"github.com/sfarosu/go-tooling-portal/cmd/tmpl"
	"html/template"
	"log"
	"os"
)

func init() {
	tmpl.Tpl = template.Must(template.ParseGlob("web/templates/*html"))
	log.SetOutput(os.Stdout) // Change the device for logging to stdout
}

func main() {
	cmd.Serve()
}
