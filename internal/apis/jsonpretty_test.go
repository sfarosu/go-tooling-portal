package apis

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/sfarosu/go-tooling-portal/internal/logger"
)

func TestJsonPretty_JSON_Success(t *testing.T) {
	router := setupTestAPI(RegisterJsonPretty)

	body := `{"input":"{\"name\":\"Alice\",\"age\":30}", "operation":"pretty", "spaces": 4}`
	req, _ := http.NewRequest("POST", "/api/jsonpretty", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Expected 200 got %d body: %s", rr.Code, rr.Body.String())
	}
}

func TestJsonPretty_JSON_ZeroIndentation(t *testing.T) {
	router := setupTestAPI(RegisterJsonPretty)
	body := `{"input":"{\"name\":\"Alice\"}","operation":"pretty","spaces":0}`
	req := httptest.NewRequest(http.MethodPost, "/api/jsonpretty", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body: %s", rr.Code, rr.Body.String())
	}

	var output struct {
		Result string `json:"result"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &output); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if !strings.Contains(output.Result, "\n\"name\"") {
		t.Errorf("expected JSON formatted with zero indentation, got %q", output.Result)
	}
}

func TestJsonPretty_JSON_Invalid(t *testing.T) {
	// initialize logger to avoid nil-pointer when handler logs errors
	logger.Logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	router := setupTestAPI(RegisterJsonPretty)

	body := `{"input":"{name:Alice,age:30}", "operation":"pretty", "spaces": 2}`
	req, _ := http.NewRequest("POST", "/api/jsonpretty", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("Expected 400 got %d body: %s", rr.Code, rr.Body.String())
	}
}

func TestJsonPretty_JSON_Unpretty(t *testing.T) {
	router := setupTestAPI(RegisterJsonPretty)

	body := `{"input":"{\n  \"name\": \"Alice\",\n  \"age\": 30\n}", "operation":"unpretty", "spaces": 2}`
	req, _ := http.NewRequest("POST", "/api/jsonpretty", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Expected 200 got %d body: %s", rr.Code, rr.Body.String())
	}
}
