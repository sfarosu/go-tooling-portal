package apis

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHtpasswdHtmx_Valid(t *testing.T) {
	router := setupTestAPI(RegisterHtpasswdHtmx)

	body := `{"username":"alice","password":"S3cureP@ssw0rd","algorithm":"apr1"}`
	req, _ := http.NewRequest("POST", "/api/htpasswd/htmx", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("HX-Request", "true")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	t.Logf("Response body: %s", rr.Body.String())

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}
	if !strings.Contains(rr.Body.String(), "alice:") {
		t.Errorf("Expected htpasswd entry in HTML, got: %s", rr.Body.String())
	}
	if !strings.Contains(rr.Header().Get("Content-Type"), "text/html") {
		t.Errorf("Expected Content-Type text/html, got: %s", rr.Header().Get("Content-Type"))
	}
}

func TestHtpasswdHtmx_MissingFields(t *testing.T) {
	router := setupTestAPI(RegisterHtpasswdHtmx)

	body := `{"username":"","password":"","algorithm":"apr1"}`
	req, _ := http.NewRequest("POST", "/api/htpasswd/htmx", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("HX-Request", "true")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	t.Logf("Response body: %s", rr.Body.String())

	if rr.Code != http.StatusUnprocessableEntity {
		t.Errorf("Expected status 422, got %d", rr.Code)
	}
	if !strings.Contains(rr.Body.String(), "expected length >= 1") {
		t.Errorf("Expected minLength validation error, got: %s", rr.Body.String())
	}
}

func TestHtpasswdHtmx_InvalidAlgorithm(t *testing.T) {
	router := setupTestAPI(RegisterHtpasswdHtmx)

	body := `{"username":"alice","password":"S3cureP@ssw0rd","algorithm":"invalid"}`
	req, _ := http.NewRequest("POST", "/api/htpasswd/htmx", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("HX-Request", "true")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	t.Logf("Response body: %s", rr.Body.String())

	if rr.Code != http.StatusUnprocessableEntity {
		t.Errorf("Expected status 422, got %d", rr.Code)
	}
	if !strings.Contains(rr.Body.String(), "expected value to be one of") {
		t.Errorf("Expected enum validation error, got: %s", rr.Body.String())
	}
}

func TestHtpasswdHtmx_MissingHtmxHeader(t *testing.T) {
	router := setupTestAPI(RegisterHtpasswdHtmx)

	body := `{"username":"alice","password":"S3cureP@ssw0rd","algorithm":"apr1"}`
	req, _ := http.NewRequest("POST", "/api/htpasswd/htmx", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	// No HX-Request header

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	t.Logf("Response body: %s", rr.Body.String())

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", rr.Code)
	}
	if !strings.Contains(rr.Body.String(), "This endpoint only supports HTMX requests") {
		t.Errorf("Expected error message for missing HX-Request, got: %s", rr.Body.String())
	}
}
