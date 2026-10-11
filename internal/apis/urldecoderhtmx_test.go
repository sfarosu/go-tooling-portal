package apis

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestURLDecoderHtmx_Decode(t *testing.T) {
	router := setupTestAPI(RegisterURLDecoderHtmx)

	body := `{"url":"https://example.com/foo%21?query=a%2Bb"}`
	req, _ := http.NewRequest("POST", "/api/urldecoder/htmx", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("HX-Request", "true")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d, response body %v", rr.Code, rr.Body.String())
	}
	if !strings.Contains(rr.Body.String(), "Scheme: https") || !strings.Contains(rr.Body.String(), "Value: a&#43;b") {
		t.Errorf("Expected decoded URL components in HTML response, got: %s", rr.Body.String())
	}
	if !strings.Contains(rr.Header().Get("Content-Type"), "text/html") {
		t.Errorf("Expected Content-Type text/html, got: %s", rr.Header().Get("Content-Type"))
	}
}

func TestURLDecoderHtmx_InvalidURL(t *testing.T) {
	router := setupTestAPI(RegisterURLDecoderHtmx)

	body := `{"url":"https://example.com/?query=%zz"}`
	req, _ := http.NewRequest("POST", "/api/urldecoder/htmx", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("HX-Request", "true")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200 with HTMX error fragment, got %d", rr.Code)
	}
	if !strings.Contains(rr.Body.String(), "ERROR:") {
		t.Errorf("Expected error fragment in response, got: %s", rr.Body.String())
	}
}

func TestURLDecoderHtmx_MissingHtmxHeader(t *testing.T) {
	router := setupTestAPI(RegisterURLDecoderHtmx)

	body := `{"url":"https://example.com"}`
	req, _ := http.NewRequest("POST", "/api/urldecoder/htmx", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d, response body %v", rr.Code, rr.Body.String())
	}
	if !strings.Contains(rr.Body.String(), "only supports HTMX requests") {
		t.Errorf("Expected missing HTMX header error, got: %s", rr.Body.String())
	}
}
