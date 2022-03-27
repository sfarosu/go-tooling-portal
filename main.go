package main

import (
	"html/template"
	"log"
	"os"

	"github.com/sfarosu/go-tooling-portal/cmd"
	"github.com/sfarosu/go-tooling-portal/cmd/tmpl"
)

func init() {
	tmpl.Tpl = template.Must(template.ParseGlob("web/templates/*html"))
	log.SetOutput(os.Stdout)
}

func main() {
	cmd.Serve()
}
