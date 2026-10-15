package httpserver

import (
	"context"
	"fmt"
	"net"
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
	"github.com/sfarosu/go-tooling-portal/internal/mcpserver"
	"github.com/sfarosu/go-tooling-portal/internal/metrics"
	"github.com/sfarosu/go-tooling-portal/internal/util/generalutil"
	"github.com/sfarosu/go-tooling-portal/internal/version"
	webassets "github.com/sfarosu/go-tooling-portal/web"
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

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Logger.Error("error starting the server", "address", addr, "error", err)
		os.Exit(1)
	}

	apiBaseURL, err := mcpAPIBaseURL(listener.Addr())
	if err != nil {
		_ = listener.Close()
		logger.Logger.Error("error determining MCP API address", "address", listener.Addr(), "error", err)
		os.Exit(1)
	}

	router, err := setupRouter(apiBaseURL)
	if err != nil {
		_ = listener.Close()
		logger.Logger.Error("error setting up router", "error", err)
		os.Exit(1)
	}

	srv := setupServer(addr, requestIDMiddleware(loggingMiddleware(router)))

	err = startupLogging(addr)
	if err != nil {
		_ = listener.Close()
		logger.Logger.Error("error during startup logging", "error", err)
		os.Exit(1)
	}

	// Channel to listen for interrupt or terminate signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// Run server in a goroutine
	go func() {
		err := srv.Serve(listener)
		if err != nil && err != http.ErrServerClosed {
			logger.Logger.Error("error starting the server", "address", addr, "error", err)
			os.Exit(1)
		}
	}()

	<-stop
	logger.Logger.Info("shutting-down-server")
	waitForShutdown(srv)
}

// setupRouter configures the HTTP router, registers the API endpoints, and static file serving
func setupRouter(apiBaseURL string) (http.Handler, error) {
	router := http.NewServeMux()

	// Register API endpoints
	humaConfig := huma.DefaultConfig("Go Tooling API", version.Version)
	humaConfig.Info.Description = "API for the Go Tooling Portal's generators, converters, formatters, and time utilities. JSON endpoints are intended for API clients; HTMX endpoints return HTML fragments and require the HX-Request: true header."
	humaConfig.Tags = []*huma.Tag{
		{Name: "Base64 Converter", Description: "Encode and decode standard or URL-safe Base64."},
		{Name: "Format Converter", Description: "Convert data between JSON and YAML."},
		{Name: "Htpasswd Generator", Description: "Generate Apache htpasswd entries."},
		{Name: "Password Generator", Description: "Generate random passwords with selected character sets."},
		{Name: "JSON Prettifier", Description: "Format or minify JSON."},
		{Name: "SSH Key Generator", Description: "Generate SSH key pairs."},
		{Name: "Time Converter", Description: "Convert between epoch timestamps and local date/time values."},
		{Name: "URL Decoder", Description: "Parse URLs and decode their components and query parameters."},
		{Name: "API Version", Description: "Retrieve the running application version."},
	}
	humaAPI := humago.New(router, humaConfig)
	apis.RegisterBase64Converter(humaAPI)
	apis.RegisterBase64ConverterHtmx(humaAPI)
	apis.RegisterFormatConverter(humaAPI)
	apis.RegisterFormatConverterHtmx(humaAPI)
	apis.RegisterHtpasswd(humaAPI)
	apis.RegisterHtpasswdHtmx(humaAPI)
	apis.RegisterPassgen(humaAPI)
	apis.RegisterPassgenHtmx(humaAPI)
	apis.RegisterJsonPretty(humaAPI)
	apis.RegisterJsonPrettyHtmx(humaAPI)
	apis.RegisterSSHGen(humaAPI)
	apis.RegisterSSHGenHtmx(humaAPI)
	apis.RegisterTimeConverter(humaAPI)
	apis.RegisterTimeConverterHtmx(humaAPI)
	apis.RegisterURLDecoder(humaAPI)
	apis.RegisterURLDecoderHtmx(humaAPI)
	apis.RegisterVersion(humaAPI)
	apis.RegisterVersionHtmx(humaAPI)

	apiEndpoints := make([]string, 0, len(humaAPI.OpenAPI().Paths))
	for path := range humaAPI.OpenAPI().Paths {
		apiEndpoints = append(apiEndpoints, path)
	}
	metrics.InitializeAPIEndpoints(apiEndpoints)

	mcpHandler, err := mcpserver.NewHandler(apiBaseURL)
	if err != nil {
		return nil, fmt.Errorf("create MCP handler: %w", err)
	}
	router.Handle("/mcp", mcpHandler)
	router.Handle("/metrics", metrics.Handler())

	// Serve embedded static files and disable directory listing.
	router.Handle("/", generalutil.DisableDirListing(webassets.Files)(http.FileServer(http.FS(webassets.Files))))

	endpoints := []string{"/", "/docs", "/metrics", "/mcp"}
	if logger.Logger != nil {
		logger.Logger.Info("http-endpoints-registered", "endpoints", endpoints)
	}

	return apis.WrapHTMXValidationErrors(router), nil
}

func mcpAPIBaseURL(addr net.Addr) (string, error) {
	tcpAddr, ok := addr.(*net.TCPAddr)
	if !ok {
		return "", fmt.Errorf("expected TCP listener address, got %T", addr)
	}

	ip := tcpAddr.IP
	if ip == nil || ip.IsUnspecified() {
		if ip != nil && ip.To4() == nil {
			ip = net.IPv6loopback
		} else {
			ip = net.IPv4(127, 0, 0, 1)
		}
	}
	return "http://" + net.JoinHostPort(ip.String(), fmt.Sprint(tcpAddr.Port)), nil
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
		"server-started",
		"address", addr,
		"gomaxprocs", runtime.GOMAXPROCS(0),
		"verbosity", logger.CurrentLevel,
	)
	logger.Logger.Info(
		"binary-info",
		"binary-path", cwd,
		"version", version.Version,
		"build-date", version.BuildDate,
		"git-short-hash", version.GitShortHash,
		"go-version", runtime.Version(),
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
	logger.Logger.Info("server-exited-gracefully")
}
