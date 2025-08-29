package apis

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"strings"

	"github.com/sfarosu/go-tooling-portal/internal/logger"
	"github.com/sfarosu/go-tooling-portal/internal/version"
)

func TestVersion_Valid(t *testing.T) {
	version.Version = "1.2.3"

	router := setupTestAPI(RegisterVersion)

	req, _ := http.NewRequest("GET", "/api/version", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	t.Logf("Response body: %s", rr.Body.String())

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}
	var out struct {
		Version string `json:"version"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &out); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	if out.Version != "1.2.3" {
		t.Errorf("Expected version '1.2.3', got %q", out.Version)
	}
}

func TestVersion_Empty(t *testing.T) {
	logger.Logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	version.Version = ""

	router := setupTestAPI(RegisterVersion)

	req, _ := http.NewRequest("GET", "/api/version", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	t.Logf("Response body: %s", rr.Body.String())

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", rr.Code)
	}
	if !json.Valid(rr.Body.Bytes()) {
		t.Errorf("Expected JSON error response, got: %s", rr.Body.String())
	}
	if !strings.Contains(rr.Body.String(), "Unable to retrieve api version information") {
		t.Errorf("Expected error message for empty version, got: %s", rr.Body.String())
	}
}
