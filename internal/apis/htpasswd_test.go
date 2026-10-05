package apis

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHtpasswd_Valid(t *testing.T) {
	router := setupTestAPI(RegisterHtpasswd)

	body := `{"username":"alice","password":"S3cureP@ssw0rd","algorithm":"apr1"}`
	req, _ := http.NewRequest("POST", "/api/htpasswd", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d, response body %v", rr.Code, rr.Body.String())
	}
	if !strings.Contains(rr.Body.String(), "alice:") {
		t.Errorf("Expected htpasswd entry in response, got: %s", rr.Body.String())
	}
}

func TestHtpasswd_MissingFields(t *testing.T) {
	router := setupTestAPI(RegisterHtpasswd)

	body := `{"username":"","password":"","algorithm":"apr1"}`
	req, _ := http.NewRequest("POST", "/api/htpasswd", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnprocessableEntity {
		t.Errorf("Expected status 422, got %d, response body %v", rr.Code, rr.Body.String())
	}
	if !strings.Contains(rr.Body.String(), "expected length >= 1") {
		t.Errorf("Expected minLength validation error, got: %s", rr.Body.String())
	}
}

func TestHtpasswd_InvalidAlgorithm(t *testing.T) {
	router := setupTestAPI(RegisterHtpasswd)

	body := `{"username":"alice","password":"S3cureP@ssw0rd","algorithm":"invalid"}`
	req, _ := http.NewRequest("POST", "/api/htpasswd", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnprocessableEntity {
		t.Errorf("Expected status 422, got %d, response body %v", rr.Code, rr.Body.String())
	}
	if !strings.Contains(rr.Body.String(), "expected value to be one of") {
		t.Errorf("Expected enum validation error, got: %s", rr.Body.String())
	}
}
