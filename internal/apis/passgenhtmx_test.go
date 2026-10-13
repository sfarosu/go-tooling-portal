package apis

import (
	"bytes"
	"encoding/json"
	"html"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
)

func TestPassgenHtmx_Valid(t *testing.T) {
	router := setupTestAPI(RegisterPassgenHtmx)

	payload := struct {
		Length    string `json:"length"`
		Uppercase string `json:"uppercase"`
		Lowercase string `json:"lowercase"`
		Numbers   string `json:"numbers"`
		Symbols   string `json:"symbols"`
	}{
		Length:    "12",
		Uppercase: "true",
		Lowercase: "true",
		Numbers:   "true",
		Symbols:   "false",
	}
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("failed to encode request: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/passgen/htmx", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("HX-Request", "true")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body: %s", rr.Code, rr.Body.String())
	}
	if !strings.Contains(rr.Header().Get("Content-Type"), "text/html") {
		t.Fatalf("expected HTML Content-Type, got %q", rr.Header().Get("Content-Type"))
	}

	match := regexp.MustCompile(`<textarea[^>]*>(.*?)</textarea>`).FindStringSubmatch(rr.Body.String())
	if len(match) != 2 {
		t.Fatalf("expected generated password textarea, got: %s", rr.Body.String())
	}
	password := html.UnescapeString(match[1])
	if len(password) != 12 {
		t.Errorf("expected 12-character password, got %q (%d characters)", password, len(password))
	}
}

func TestPassgenHtmx_InvalidInputReturnsErrorFragment(t *testing.T) {
	router := setupTestAPI(RegisterPassgenHtmx)

	body := `{"length":"12","uppercase":"false","lowercase":"false","numbers":"false","symbols":"false"}`
	req := httptest.NewRequest(http.MethodPost, "/api/passgen/htmx", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("HX-Request", "true")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 with HTMX error fragment, got %d body: %s", rr.Code, rr.Body.String())
	}
	if !strings.Contains(rr.Body.String(), "at least one character class must be selected") {
		t.Fatalf("expected character-class validation error, got: %s", rr.Body.String())
	}
}

func TestPassgenHtmx_MissingHtmxHeader(t *testing.T) {
	router := setupTestAPI(RegisterPassgenHtmx)

	body := `{"length":"12","uppercase":"true","lowercase":"true","numbers":"true","symbols":"false"}`
	req := httptest.NewRequest(http.MethodPost, "/api/passgen/htmx", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d body: %s", rr.Code, rr.Body.String())
	}
}
