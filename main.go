package main

import (
	"log"
	"os"

	"flag"

	"github.com/sfarosu/go-tooling-portal/internal/httpserver"
	"github.com/sfarosu/go-tooling-portal/internal/logger"
)

func init() {
	log.SetOutput(os.Stdout)
}

func main() {
	addr := flag.String("addr", ":8100", "Network address and port to start on")
	logLevel := flag.String("verbosity", "info", "Verbosity: debug, info, warn, error")
	flag.Parse()

	logger.Init(*logLevel)

	httpserver.Start(*addr)
}
