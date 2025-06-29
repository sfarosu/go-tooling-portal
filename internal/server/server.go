package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	"github.com/sfarosu/go-tooling-portal/internal/apis"
	"github.com/sfarosu/go-tooling-portal/internal/logger"
	"github.com/sfarosu/go-tooling-portal/internal/version"
	"go.uber.org/automaxprocs/maxprocs"
)

// server timeouts and configuration values.
const (
	serverReadTimeout     = 5 * time.Second
	serverWriteTimeout    = 10 * time.Second
	serverIdleTimeout     = 120 * time.Second
	serverShutdownTimeout = 10 * time.Second
)

// Start initializes and runs the HTTP server, sets up routing, logging, and handles graceful shutdown
func Start(addr string) {
	// GOMAXPROCS - respect K8S cpu quota
	_, errMax := maxprocs.Set()
	if errMax != nil {
		logger.Logger.Error("error setting maxprocs", "error", errMax)
	}

	router := setupRouter()

	srv := setupServer(addr, requestIDMiddleware(loggingMiddleware(router)))

	err := startupLogging(addr)
	if err != nil {
		logger.Logger.Error("error during startup logging", "error", err)
		os.Exit(1)
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

	<-stop
	logger.Logger.Info("shutting down server")
	waitForShutdown(srv)
}

// setupRouter configures the HTTP router, registers the API endpoints, and static file serving
func setupRouter() *http.ServeMux {
	router := http.NewServeMux()

	// Register API endpoints
	humaConfig := huma.DefaultConfig("Go Tooling API", "1.0.0")
	humaAPI := humago.New(router, humaConfig)
	apis.RegisterVersion(humaAPI)
	apis.RegisterHtpasswd(humaAPI)

	// Serve static files
	fileServer := http.FileServer(http.Dir("web"))
	router.Handle("/", fileServer)

	return router
}

// setupServer creates and configures the HTTP server.
func setupServer(addr string, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  serverReadTimeout,
		WriteTimeout: serverWriteTimeout,
		IdleTimeout:  serverIdleTimeout,
	}
}

// startupLogging logs startup information about the server and binary file
func startupLogging(addr string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error establishing current working directory: %v", err)
	}

	logger.Logger.Info(
		"server started",
		"address", addr,
		"gomaxprocs", runtime.GOMAXPROCS(0),
		"verbosity", logger.CurrentLevel,
	)
	logger.Logger.Info(
		"binary info",
		"binary_path", cwd,
		"version", version.Version,
		"build_date", version.BuildDate,
		"git_short_hash", version.GitShortHash,
		"go_version", runtime.Version(),
	)

	return nil
}

// waitForShutdown gracefully shuts down the HTTP server on interrupt
func waitForShutdown(srv *http.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), serverShutdownTimeout)
	defer cancel()
	err := srv.Shutdown(ctx)
	if err != nil {
		logger.Logger.Error("server forced to shutdown due to context timeout", "error", err)
		os.Exit(1)
	}
	logger.Logger.Info("server exited gracefully")
}
