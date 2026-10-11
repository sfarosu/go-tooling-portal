package apis

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestURLDecoder_Decode(t *testing.T) {
	router := setupTestAPI(RegisterURLDecoder)

	body := `{"url":"https://example.com/foo%21?query=a%2Bb"}`
	req, _ := http.NewRequest("POST", "/api/urldecoder", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d, response body %v", rr.Code, rr.Body.String())
	}

	var output struct {
		Result string `json:"result"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &output); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	for _, expected := range []string{"Scheme: https", "Host: example.com", "Path: /foo!", "Value: a+b"} {
		if !strings.Contains(output.Result, expected) {
			t.Errorf("Expected decoded result to contain %q, got:\n%s", expected, output.Result)
		}
	}
}

func TestURLDecoder_InvalidURL(t *testing.T) {
	router := setupTestAPI(RegisterURLDecoder)

	body := `{"url":"https://example.com/?query=%zz"}`
	req, _ := http.NewRequest("POST", "/api/urldecoder", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d, response body %v", rr.Code, rr.Body.String())
	}
}
