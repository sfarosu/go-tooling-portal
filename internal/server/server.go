package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/sfarosu/go-tooling-portal/internal/helper"
	"github.com/sfarosu/go-tooling-portal/internal/logger"
	"github.com/sfarosu/go-tooling-portal/internal/version"
	"go.uber.org/automaxprocs/maxprocs"
)

func Start(addr string) {
	// GOMAXPROCS - respect k8s cpu quota
	_, errMax := maxprocs.Set()
	if errMax != nil {
		logger.Logger.Error("error setting maxprocs", "error", errMax)
	}

	// Set up a new ServeMux
	mux := http.NewServeMux()

	// Serve static files from the "web" directory
	fileServer := http.FileServer(http.Dir("web"))
	mux.Handle("/", fileServer)

	// Configure the http server
	srv := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Channel to listen for interrupt or terminate signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// Run server in a goroutine
	go func() {
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			logger.Logger.Error("error starting the server", "address", addr, "error", err)
			os.Exit(1)
		}
	}()

	// Startup logging
	logger.Logger.Info(
		"server started",
		"address", addr,
		"gomaxprocs", runtime.GOMAXPROCS(0),
		"verbosity", logger.CurrentLevel,
	)

	logger.Logger.Info(
		"binary info",
		"binary_path", helper.CurrentWorkingDirectory(),
		"version", version.Version,
		"build_date", version.BuildDate,
		"git_short_hash", version.GitShortHash,
		"go_version", runtime.Version(),
	)

	// Block until a signal is received
	<-stop
	logger.Logger.Info("shutting down server")

	// Create a context with timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := srv.Shutdown(ctx)
	if err != nil {
		logger.Logger.Error("server forced to shutdown due to context timeout", "error", err)
		os.Exit(1)
	}

	logger.Logger.Info("server exited gracefully")
}
