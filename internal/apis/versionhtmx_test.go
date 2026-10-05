package apis

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"strings"

	"github.com/sfarosu/go-tooling-portal/internal/logger"
	"github.com/sfarosu/go-tooling-portal/internal/version"
)

func TestVersionHtmx_Valid(t *testing.T) {
	logger.Logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	version.Version = "2.0.0"

	router := setupTestAPI(RegisterVersionHtmx)

	req, _ := http.NewRequest("GET", "/api/version/htmx", nil)
	req.Header.Set("HX-Request", "true")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d, response body %v", rr.Code, rr.Body.String())
	}
	if !strings.Contains(rr.Body.String(), `<span class="badge">v2.0.0</span>`) {
		t.Errorf("Expected HTML badge with version, got: %s", rr.Body.String())
	}
}

func TestVersionHtmx_Empty(t *testing.T) {
	logger.Logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	version.Version = ""

	router := setupTestAPI(RegisterVersionHtmx)

	req, _ := http.NewRequest("GET", "/api/version/htmx", nil)
	req.Header.Set("HX-Request", "true")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// For HTMX requests, errors are returned with 200 but message in body
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d, response body %v", rr.Code, rr.Body.String())
	}
	if !strings.Contains(rr.Body.String(), "internal server error, unable to retrieve api version information") {
		t.Errorf("Expected error message for empty version, got: %s", rr.Body.String())
	}
}
