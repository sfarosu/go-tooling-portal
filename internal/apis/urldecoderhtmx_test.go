package apis

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestURLDecoderHtmx_Decode(t *testing.T) {
	router := setupTestAPI(RegisterURLDecoderHtmx)

	body := `{"url":"https://example.com/foo+bar%21?service=ab%2Bc&product=de%24f&service=second"}`
	req, _ := http.NewRequest("POST", "/api/urldecoder/htmx", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("HX-Request", "true")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d, response body %v", rr.Code, rr.Body.String())
	}
	for _, expected := range []string{"Scheme: https", "Path: /foo&#43;bar!", "Value: de$f", "Value: ab&#43;c", "Value: second"} {
		if !strings.Contains(rr.Body.String(), expected) {
			t.Errorf("Expected decoded URL response to contain %q, got: %s", expected, rr.Body.String())
		}
	}
	if !strings.Contains(rr.Body.String(), `class="input-group w-100"`) {
		t.Errorf("Expected decoded URL result area to fill its available width, got: %s", rr.Body.String())
	}
	if !strings.Contains(rr.Body.String(), `id="urldecoder-result-input" readonly rows="20"`) {
		t.Errorf("Expected decoded URL result area to have 20 rows, got: %s", rr.Body.String())
	}
	if !strings.Contains(rr.Header().Get("Content-Type"), "text/html") {
		t.Errorf("Expected Content-Type text/html, got: %s", rr.Header().Get("Content-Type"))
	}
}

func TestURLDecoderHtmx_InvalidURL(t *testing.T) {
	router := setupTestAPI(RegisterURLDecoderHtmx)

	body := `{"url":"https://example.com/?query=%zz"}`
	req, _ := http.NewRequest("POST", "/api/urldecoder/htmx", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("HX-Request", "true")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200 with HTMX error fragment, got %d", rr.Code)
	}
	if !strings.Contains(rr.Body.String(), "ERROR:") {
		t.Errorf("Expected error fragment in response, got: %s", rr.Body.String())
	}
}

func TestURLDecoderHtmx_MissingHtmxHeader(t *testing.T) {
	router := setupTestAPI(RegisterURLDecoderHtmx)

	body := `{"url":"https://example.com"}`
	req, _ := http.NewRequest("POST", "/api/urldecoder/htmx", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d, response body %v", rr.Code, rr.Body.String())
	}
	if !strings.Contains(rr.Body.String(), "only supports HTMX requests") {
		t.Errorf("Expected missing HTMX header error, got: %s", rr.Body.String())
	}
}
