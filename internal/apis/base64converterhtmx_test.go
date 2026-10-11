package apis

import (
	"bytes"
	"encoding/base64"
	"html"
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

func TestBase64ConverterHtmx_EscapesHTMLInResult(t *testing.T) {
	router := setupTestAPI(RegisterBase64ConverterHtmx)
	input := base64.StdEncoding.EncodeToString([]byte("<script>alert(1)</script>"))
	body := `{"input":"` + input + `","operation":"decode","format":"standard"}`
	req := httptest.NewRequest(http.MethodPost, "/api/base64converter/htmx", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("HX-Request", "true")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body: %s", rr.Code, rr.Body.String())
	}
	if strings.Contains(rr.Body.String(), "<script>") {
		t.Fatalf("result contains unescaped markup: %s", rr.Body.String())
	}
	if !strings.Contains(rr.Body.String(), "&lt;script&gt;") {
		t.Errorf("expected HTML-escaped output: %s", rr.Body.String())
	}
	if !strings.Contains(html.UnescapeString(rr.Body.String()), "<script>alert(1)</script>") {
		t.Errorf("expected escaped result to display its original text: %s", rr.Body.String())
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

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200 with HTMX validation fragment, got %d, response body %v", rr.Code, rr.Body.String())
	}
	if !strings.Contains(rr.Body.String(), "ERROR:") || !strings.Contains(rr.Body.String(), "expected value to be one of") {
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
