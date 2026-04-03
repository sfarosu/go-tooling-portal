package apis

import (
	"bytes"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestBase64ConverterHtmx_Encode(t *testing.T) {
	router := setupTestAPI(RegisterBase64ConverterHtmx)

	body := `{"input":"hello","operation":"encode","format":"standard"}`
	req, _ := http.NewRequest("POST", "/api/base64converter/htmx", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("HX-Request", "true")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d, response body %v", rr.Code, rr.Body.String())
	}
	if !strings.Contains(rr.Body.String(), base64.StdEncoding.EncodeToString([]byte("hello"))) {
		t.Errorf("Expected base64 result in HTML, got: %s", rr.Body.String())
	}
	if !strings.Contains(rr.Header().Get("Content-Type"), "text/html") {
		t.Errorf("Expected Content-Type text/html, got: %s", rr.Header().Get("Content-Type"))
	}
}

func TestBase64ConverterHtmx_Decode(t *testing.T) {
	router := setupTestAPI(RegisterBase64ConverterHtmx)

	encoded := base64.StdEncoding.EncodeToString([]byte("hello"))
	body := `{"input":"` + encoded + `","operation":"decode","format":"standard"}`
	req, _ := http.NewRequest("POST", "/api/base64converter/htmx", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("HX-Request", "true")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d, response body %v", rr.Code, rr.Body.String())
	}
	if !strings.Contains(rr.Body.String(), "hello") {
		t.Errorf("Expected decoded result in HTML, got: %s", rr.Body.String())
	}
	if !strings.Contains(rr.Header().Get("Content-Type"), "text/html") {
		t.Errorf("Expected Content-Type text/html, got: %s", rr.Header().Get("Content-Type"))
	}
}

func TestBase64ConverterHtmx_InvalidOperation(t *testing.T) {
	router := setupTestAPI(RegisterBase64ConverterHtmx)

	body := `{"input":"hello","operation":"invalid","format":"standard"}`
	req, _ := http.NewRequest("POST", "/api/base64converter/htmx", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("HX-Request", "true")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnprocessableEntity {
		t.Errorf("Expected status 422, got %d, response body %v", rr.Code, rr.Body.String())
	}
	if !strings.Contains(rr.Body.String(), "expected value to be one of") {
		t.Errorf("Expected enum validation error, got: %s", rr.Body.String())
	}
}

func TestBase64ConverterHtmx_MissingHtmxHeader(t *testing.T) {
	router := setupTestAPI(RegisterBase64ConverterHtmx)

	body := `{"input":"hello","operation":"encode","format":"standard"}`
	req, _ := http.NewRequest("POST", "/api/base64converter/htmx", bytes.NewBufferString(body))
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
