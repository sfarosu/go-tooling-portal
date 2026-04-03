package apis

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
)

func TestJsonPrettyHtmx_Valid(t *testing.T) {
	router := setupTestAPI(RegisterJsonPrettyHtmx)

	// Use lowercase json field names (form submissions use these)
	body := `{"input":"{\"name\":\"Alice\",\"age\":30}","operation":"pretty","spaces":"2"}`
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
	body := `{"input":"","operation":"pretty","spaces":"2"}`
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
	body := `{"input":"not a json","operation":"pretty","spaces":"2"}`
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

	body := `{"input":"{\"name\":\"Alice\",\"age\":30}","operation":"pretty","spaces":"2"}`
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

func TestJsonPrettyHtmx_Unpretty(t *testing.T) {
	router := setupTestAPI(RegisterJsonPrettyHtmx)

	// Input is prettified JSON; operation unpretty should minify it
	// Use escaped newlines so the JSON payload is valid
	prettyInputEscaped := "{\\n  \\\"name\\\": \\\"Alice\\\",\\n  \\\"age\\\": 30\\n}"
	body := `{"input":"` + prettyInputEscaped + `","operation":"unpretty","spaces":"2"}`
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

	// Extract JSON from the <textarea> in the HTML
	re := regexp.MustCompile(`<textarea[^>]*>([\s\S]*?)</textarea>`)
	matches := re.FindStringSubmatch(rr.Body.String())
	if len(matches) < 2 {
		t.Fatalf("Could not find <textarea> content in HTML")
	}
	jsonStr := matches[1]

	if !strings.Contains(jsonStr, `"name":"Alice"`) || !strings.Contains(jsonStr, `"age":30`) {
		t.Fatalf("expected compact JSON components in result textarea, got: %s", jsonStr)
	}
	if strings.Contains(jsonStr, "\n") {
		t.Fatalf("expected no newlines in compact JSON result, got: %s", jsonStr)
	}
}
