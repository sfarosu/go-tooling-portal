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

func TestFormatConverter_JSONtoYAML(t *testing.T) {
	router := setupTestAPI(RegisterFormatConverter)

	body := `{"input":"{\"name\":\"Alice\",\"age\":30}"}`
	req, _ := http.NewRequest(http.MethodPost, "/api/formatconverter", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d, response body %v", rr.Code, rr.Body.String())
	}

	if !strings.Contains(rr.Body.String(), "age: 30") ||
		!strings.Contains(rr.Body.String(), "name: Alice") {
		t.Errorf("Expected YAML output, got: %s", rr.Body.String())
	}
}

func TestFormatConverter_YAMLtoJSON(t *testing.T) {
	router := setupTestAPI(RegisterFormatConverter)

	yaml := "name: Bob\nage: 25\n"
	inputJSON, _ := json.Marshal(struct {
		Input string `json:"input"`
	}{Input: yaml})

	req, _ := http.NewRequest(http.MethodPost, "/api/formatconverter", bytes.NewReader(inputJSON))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d, response body %v", rr.Code, rr.Body.String())
	}

	// Unmarshal outer response
	var out struct {
		Result string `json:"result"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &out); err != nil {
		t.Fatalf("Failed to unmarshal outer response: %v", err)
	}

	// Unmarshal inner JSON contained in Result
	var inner map[string]any
	if err := json.Unmarshal([]byte(out.Result), &inner); err != nil {
		t.Fatalf("Failed to unmarshal inner JSON: %v", err)
	}

	if inner["name"] != "Bob" || int(inner["age"].(float64)) != 25 {
		t.Errorf("Expected inner JSON with name=Bob and age=25, got: %+v", inner)
	}
}

func TestFormatConverter_InvalidJSON(t *testing.T) {
	logger.Logger = slog.New(slog.NewTextHandler(os.Stdout, nil))

	router := setupTestAPI(RegisterFormatConverter)

	// invalid JSON string
	body := `{"input":"{name:Alice,age:30"}`
	req, _ := http.NewRequest("POST", "/api/formatconverter", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// service.FormatConvert should error -> API returns 400
	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d, response body %v", rr.Code, rr.Body.String())
	}
}

func TestFormatConverter_InvalidYAML(t *testing.T) {
	logger.Logger = slog.New(slog.NewTextHandler(os.Stdout, nil))

	router := setupTestAPI(RegisterFormatConverter)

	yaml := "name Bob\nage 30" // invalid YAML
	inputJSON, _ := json.Marshal(struct {
		Input string `json:"input"`
	}{Input: yaml})

	req, _ := http.NewRequest("POST", "/api/formatconverter", bytes.NewReader(inputJSON))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d, response body %v", rr.Code, rr.Body.String())
	}

	if !strings.Contains(rr.Body.String(), "format conversion failed") {
		t.Errorf("Expected error message, got: %s", rr.Body.String())
	}
}
