package apis

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBase64Converter_Encode(t *testing.T) {
	router := setupTestAPI(RegisterBase64Converter)

	body := `{"input":"hello","operation":"encode","format":"standard"}`
	req, _ := http.NewRequest("POST", "/api/base64converter", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	t.Logf("Response body: %s", rr.Body.String())

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}

	var out struct {
		Result string `json:"result"`
	}

	if err := json.Unmarshal(rr.Body.Bytes(), &out); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	want := base64.StdEncoding.EncodeToString([]byte("hello"))
	if out.Result != want {
		t.Errorf("Expected result %q, got %q", want, out.Result)
	}
}

func TestBase64Converter_Decode(t *testing.T) {
	router := setupTestAPI(RegisterBase64Converter)

	body := `{"input":"aGVsbG8=","operation":"decode","format":"standard"}`
	req, _ := http.NewRequest("POST", "/api/base64converter", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	t.Logf("Response body: %s", rr.Body.String())

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}

	var out struct {
		Result string `json:"result"`
	}

	if err := json.Unmarshal(rr.Body.Bytes(), &out); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if out.Result != "hello" {
		t.Errorf("Expected result 'hello', got %q", out.Result)
	}
}

func TestBase64Converter_InvalidOperation(t *testing.T) {
	router := setupTestAPI(RegisterBase64Converter)

	body := `{"input":"hello","operation":"invalid","format":"standard"}`
	req, _ := http.NewRequest("POST", "/api/base64converter", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	t.Logf("Response body: %s", rr.Body.String())

	if rr.Code != http.StatusUnprocessableEntity {
		t.Errorf("Expected status 422, got %d", rr.Code)
	}
}
