package main

import (
	"log"
	"os"

	"github.com/sfarosu/go-tooling-portal/internal/handlers"
)

func init() {
	log.SetOutput(os.Stdout)
}

func main() {
	handlers.Serve()
}
