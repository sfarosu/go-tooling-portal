package apis

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestJsonPrettyHtmx_Valid(t *testing.T) {
	router := setupTestAPI(RegisterJsonPrettyHtmx)

	// Use lowercase json field names (form submissions use these)
	body := `{"input":"{\"name\":\"Alice\",\"age\":30}","spaces":"2"}`
	req, _ := http.NewRequest("POST", "/api/jsonpretty/htmx", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("HX-Request", "true")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d body: %s", rr.Code, rr.Body.String())
	}

	if !strings.Contains(rr.Header().Get("Content-Type"), "text/html") {
		t.Fatalf("expected Content-Type text/html, got: %s", rr.Header().Get("Content-Type"))
	}
}

func TestJsonPrettyHtmx_MissingFields(t *testing.T) {
	router := setupTestAPI(RegisterJsonPrettyHtmx)

	// Missing/empty input should trigger validation (minLength)
	body := `{"input":"","spaces":"2"}`
	req, _ := http.NewRequest("POST", "/api/jsonpretty/htmx", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("HX-Request", "true")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422 got %d body: %s", rr.Code, rr.Body.String())
	}

	if !strings.Contains(rr.Body.String(), "expected length >= 1") {
		t.Fatalf("expected minLength validation error, got: %s", rr.Body.String())
	}
}

func TestJsonPrettyHtmx_InvalidInput(t *testing.T) {
	router := setupTestAPI(RegisterJsonPrettyHtmx)

	// Invalid JSON should be handled and returned inside the HTMX HTML fragment
	body := `{"input":"not a json","spaces":"2"}`
	req, _ := http.NewRequest("POST", "/api/jsonpretty/htmx", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("HX-Request", "true")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// HTMX endpoint returns 200 even on errors (error shown inside the fragment)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d body: %s", rr.Code, rr.Body.String())
	}

	if !strings.Contains(rr.Body.String(), "ERROR:") {
		t.Fatalf("expected an ERROR: message in the HTML fragment, got: %s", rr.Body.String())
	}
}

func TestJsonPrettyHtmx_MissingHtmxHeader(t *testing.T) {
	router := setupTestAPI(RegisterJsonPrettyHtmx)

	body := `{"input":"{\"name\":\"Alice\",\"age\":30}","spaces":"2"}`
	req, _ := http.NewRequest("POST", "/api/jsonpretty/htmx", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	// No HX-Request header

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d body: %s", rr.Code, rr.Body.String())
	}

	if !strings.Contains(rr.Body.String(), "This endpoint only supports HTMX requests") {
		t.Fatalf("expected missing htmx header error, got: %s", rr.Body.String())
	}
}
