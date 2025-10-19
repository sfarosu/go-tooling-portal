package apis

import (
    "bytes"
    "io"
    "net/http"
    "net/http/httptest"
    "testing"
    "encoding/json"
)

func TestPassgenHtmx_Valid(t *testing.T) {
    router := setupTestAPI(RegisterPassgenHtmx)

    // simulate form-like string values (e.g. "on" for checked checkboxes)
    payload := map[string]interface{}{"length": "12", "uppercase": "on", "lowercase": "on", "numbers": "on", "symbols": "off"}
    b, _ := json.Marshal(payload)

    req := httptest.NewRequest(http.MethodPost, "/api/passgen/htmx", bytes.NewReader(b))
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("HX-Request", "true")
    rr := httptest.NewRecorder()
    router.ServeHTTP(rr, req)

    if rr.Code != http.StatusOK {
        t.Fatalf("expected 200 got %d body: %s", rr.Code, rr.Body.String())
    }

    if ct := rr.Header().Get("Content-Type"); ct == "" {
        t.Fatalf("expected Content-Type header set for HTML response")
    }

    body, _ := io.ReadAll(rr.Body)
    if len(body) == 0 {
        t.Fatalf("expected non-empty body")
    }
}

func TestPassgenHtmx_MissingHtmxHeader(t *testing.T) {
    router := setupTestAPI(RegisterPassgenHtmx)

    // simulate form-like string values
    payload := map[string]interface{}{"length": "12", "uppercase": "on", "lowercase": "on", "numbers": "on", "symbols": "off"}
    b, _ := json.Marshal(payload)

    req := httptest.NewRequest(http.MethodPost, "/api/passgen/htmx", bytes.NewReader(b))
    req.Header.Set("Content-Type", "application/json")
    rr := httptest.NewRecorder()
    router.ServeHTTP(rr, req)

    if rr.Code != http.StatusBadRequest {
        t.Fatalf("expected 400 got %d body: %s", rr.Code, rr.Body.String())
    }
}

// NOTE: use encoding/json directly in tests above
