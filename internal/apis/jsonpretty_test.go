package apis

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestJsonPretty_JSON_Success(t *testing.T) {
	router := setupTestAPI(RegisterJsonPretty)

	body := `{"Input":"{\"name\":\"Alice\",\"age\":30}", "spaces": 4}`
	req, _ := http.NewRequest("POST", "/api/jsonpretty", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Expected 200 got %d body: %s", rr.Code, rr.Body.String())
	}
}

func TestJsonPretty_JSON_Invalid(t *testing.T) {
	router := setupTestAPI(RegisterJsonPretty)

	body := `{"Input":"{name:Alice,age:30}", "spaces": 2}`
	req, _ := http.NewRequest("POST", "/api/jsonpretty", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("Expected 400 got %d body: %s", rr.Code, rr.Body.String())
	}
}
