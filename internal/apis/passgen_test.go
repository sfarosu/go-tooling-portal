package apis

import (
    "bytes"
    "encoding/json"
    "io"
    "net/http"
    "net/http/httptest"
    "testing"

)

func TestPassgen_JSON_Success(t *testing.T) {
    router := setupTestAPI(RegisterPassgen)

    payload := map[string]interface{}{
        "length":    12,
        "uppercase": true,
        "lowercase": true,
        "numbers":   true,
        "symbols":   false,
    }
    b, _ := json.Marshal(payload)

    req := httptest.NewRequest(http.MethodPost, "/api/passgen", bytes.NewReader(b))
    req.Header.Set("Content-Type", "application/json")
    rr := httptest.NewRecorder()
    router.ServeHTTP(rr, req)

    if rr.Code != http.StatusOK {
        t.Fatalf("expected 200 got %d body: %s", rr.Code, rr.Body.String())
    }

    var resp struct{
        Result string `json:"result"`
    }
    body, _ := io.ReadAll(rr.Body)
    if err := json.Unmarshal(body, &resp); err != nil {
        t.Fatalf("invalid json response: %v", err)
    }

    if resp.Result == "" {
        t.Fatalf("expected non-empty result")
    }
}

func TestPassgen_JSON_InvalidInput(t *testing.T) {
    router := setupTestAPI(RegisterPassgen)

    payload := map[string]interface{}{"length": 0, "uppercase": true, "lowercase": true, "numbers": true, "symbols": false}
    b, _ := json.Marshal(payload)

    req := httptest.NewRequest(http.MethodPost, "/api/passgen", bytes.NewReader(b))
    req.Header.Set("Content-Type", "application/json")
    rr := httptest.NewRecorder()
    router.ServeHTTP(rr, req)

    if rr.Code != http.StatusBadRequest {
        t.Fatalf("expected 400 got %d body: %s", rr.Code, rr.Body.String())
    }
}
