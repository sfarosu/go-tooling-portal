package apis

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSSHGen_Ed25519(t *testing.T) {
	router := setupTestAPI(RegisterSSHGen)

	body := `{"algorithm":"ed25519","ecdsaBits":"","rsaBits":"","email":"alice@example.com"}`
	req, _ := http.NewRequest("POST", "/api/sshgen", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d, response body %v", rr.Code, rr.Body.String())
	}
	if !strings.Contains(rr.Body.String(), "ssh-ed25519") {
		t.Errorf("Expected ed25519 public key in response, got: %s", rr.Body.String())
	}
}

func TestSSHGen_ECDSA(t *testing.T) {
	router := setupTestAPI(RegisterSSHGen)

	body := `{"algorithm":"ecdsa","ecdsaBits":"256","rsaBits":"","email":"bob@example.com"}`
	req, _ := http.NewRequest("POST", "/api/sshgen", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d, response body %v", rr.Code, rr.Body.String())
	}
	if !strings.Contains(rr.Body.String(), "ecdsa-sha2-nistp256") {
		t.Errorf("Expected ecdsa public key in response, got: %s", rr.Body.String())
	}
}

func TestSSHGen_RSA(t *testing.T) {
	router := setupTestAPI(RegisterSSHGen)

	body := `{"algorithm":"rsa","ecdsaBits":"","rsaBits":"2048","email":"carol@example.com"}`
	req, _ := http.NewRequest("POST", "/api/sshgen", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d, response body %v", rr.Code, rr.Body.String())
	}
	if !strings.Contains(rr.Body.String(), "ssh-rsa") {
		t.Errorf("Expected rsa public key in response, got: %s", rr.Body.String())
	}
}

func TestSSHGen_MissingBits(t *testing.T) {
	router := setupTestAPI(RegisterSSHGen)

	body := `{"algorithm":"ecdsa","ecdsaBits":"","rsaBits":"","email":"nobody@example.com"}`
	req, _ := http.NewRequest("POST", "/api/sshgen", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d, response body %v", rr.Code, rr.Body.String())
	}
	if !strings.Contains(rr.Body.String(), "ecdsaBits is required for ECDSA algorithm") {
		t.Errorf("Expected error message for missing ecdsaBits, got: %s", rr.Body.String())
	}
}

func TestSSHGen_InvalidAlgorithm(t *testing.T) {
	router := setupTestAPI(RegisterSSHGen)

	body := `{"algorithm":"invalid","ecdsaBits":"","rsaBits":"","email":"nobody@example.com"}`
	req, _ := http.NewRequest("POST", "/api/sshgen", bytes.NewBufferString(body))
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
