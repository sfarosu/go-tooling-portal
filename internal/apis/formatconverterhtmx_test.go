package apis

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/sfarosu/go-tooling-portal/internal/logger"
)

func TestFormatConverterHtmx_JSONtoYAML(t *testing.T) {
	router := setupTestAPI(RegisterFormatConverterHtmx)

	bodyStruct := struct {
		Input string `json:"input"`
	}{Input: `{"name":"Alice","age":30}`}

	bodyJSON, _ := json.Marshal(bodyStruct)

	req, _ := http.NewRequest("POST", "/api/formatconverter/htmx", bytes.NewReader(bodyJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("HX-Request", "true")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d, response body %v", rr.Code, rr.Body.String())
	}

	if !strings.Contains(rr.Body.String(), "age: 30") ||
		!strings.Contains(rr.Body.String(), "name: Alice") {
		t.Errorf("Expected YAML output in HTML, got: %s", rr.Body.String())
	}

	if !strings.Contains(rr.Header().Get("Content-Type"), "text/html") {
		t.Errorf("Expected Content-Type text/html, got: %s", rr.Header().Get("Content-Type"))
	}
}

func TestFormatConverterHtmx_YAMLtoJSON(t *testing.T) {
	router := setupTestAPI(RegisterFormatConverterHtmx)

	yamlInput := "name: Bob\nage: 25\n"
	bodyStruct := struct {
		Input string `json:"input"`
	}{Input: yamlInput}

	bodyJSON, _ := json.Marshal(bodyStruct)

	req, _ := http.NewRequest("POST", "/api/formatconverter/htmx", bytes.NewReader(bodyJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("HX-Request", "true")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d, response body %v", rr.Code, rr.Body.String())
	}

	// Extract JSON from the <textarea> in the HTML
	re := regexp.MustCompile(`<textarea[^>]*>([\s\S]*?)</textarea>`)
	matches := re.FindStringSubmatch(rr.Body.String())
	if len(matches) < 2 {
		t.Fatalf("Could not find <textarea> content in HTML")
	}
	jsonStr := matches[1]

	var out map[string]any
	if err := json.Unmarshal([]byte(jsonStr), &out); err != nil {
		t.Fatalf("Failed to unmarshal JSON from textarea: %v", err)
	}

	if out["name"] != "Bob" || out["age"] != float64(25) { // JSON numbers are float64
		t.Errorf("Expected JSON with name=Bob, age=25, got: %+v", out)
	}

	if !strings.Contains(rr.Header().Get("Content-Type"), "text/html") {
		t.Errorf("Expected Content-Type text/html, got: %s", rr.Header().Get("Content-Type"))
	}
}

func TestFormatConverterHtmx_InvalidInput(t *testing.T) {
	logger.Logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	router := setupTestAPI(RegisterFormatConverterHtmx)

	bodyStruct := struct {
		Input string `json:"input"`
	}{Input: "invalid data !!!"}

	bodyJSON, _ := json.Marshal(bodyStruct)
	req, _ := http.NewRequest("POST", "/api/formatconverter/htmx", bytes.NewReader(bodyJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("HX-Request", "true")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// HTMX endpoint returns 200 even on errors
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d, response body %v", rr.Code, rr.Body.String())
	}

	// Extract error text from <textarea>
	re := regexp.MustCompile(`<textarea[^>]*>([\s\S]*?)</textarea>`)
	matches := re.FindStringSubmatch(rr.Body.String())
	if len(matches) < 2 {
		t.Fatalf("Could not find <textarea> content in HTML")
	}
	errorText := matches[1]

	if !strings.Contains(errorText, "input is neither valid JSON nor YAML") {
		t.Errorf("Expected error message in HTML, got: %s", errorText)
	}

	if !strings.Contains(rr.Header().Get("Content-Type"), "text/html") {
		t.Errorf("Expected Content-Type text/html, got: %s", rr.Header().Get("Content-Type"))
	}
}

func TestFormatConverterHtmx_MissingHtmxHeader(t *testing.T) {
	router := setupTestAPI(RegisterFormatConverterHtmx)

	bodyStruct := struct {
		Input string `json:"input"`
	}{Input: `{"name":"Alice","age":30}`}

	bodyJSON, _ := json.Marshal(bodyStruct)

	req, _ := http.NewRequest("POST", "/api/formatconverter/htmx", bytes.NewReader(bodyJSON))
	req.Header.Set("Content-Type", "application/json")
	// No HX-Request header

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d, response body %v", rr.Code, rr.Body.String())
	}

	if !strings.Contains(rr.Body.String(), "This endpoint only supports HTMX requests") {
		t.Errorf("Expected error message for missing HX-Request, got: %s", rr.Body.String())
	}
}
