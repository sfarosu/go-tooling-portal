package apis

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestTimeConverterHtmx_EpochToHuman_Seconds(t *testing.T) {
	router := setupTestAPI(RegisterTimeConverterHtmx)

	body := `{"epochToHumanMode":true,"humanToEpochMode":false,"epochTime":1704067200,"year":0,"month":0,"day":0,"hour":0,"minute":0,"second":0,"timezone":"UTC"}`
	req, _ := http.NewRequest("POST", "/api/timeconverter/htmx", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("HX-Request", "true")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d, response body %v", rr.Code, rr.Body.String())
	}
	if !strings.Contains(rr.Body.String(), "UTC Time") {
		t.Errorf("Expected 'UTC Time' in HTML response, got: %s", rr.Body.String())
	}
	if !strings.Contains(rr.Header().Get("Content-Type"), "text/html") {
		t.Errorf("Expected Content-Type text/html, got: %s", rr.Header().Get("Content-Type"))
	}
}

func TestTimeConverterHtmx_EpochToHuman_Milliseconds(t *testing.T) {
	router := setupTestAPI(RegisterTimeConverterHtmx)

	body := `{"epochToHumanMode":true,"humanToEpochMode":false,"epochTime":1704067200000,"year":0,"month":0,"day":0,"hour":0,"minute":0,"second":0,"timezone":"UTC"}`
	req, _ := http.NewRequest("POST", "/api/timeconverter/htmx", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("HX-Request", "true")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d, response body %v", rr.Code, rr.Body.String())
	}
	if !strings.Contains(rr.Body.String(), "UTC Time") {
		t.Errorf("Expected 'UTC Time' in HTML response, got: %s", rr.Body.String())
	}
	if !strings.Contains(rr.Header().Get("Content-Type"), "text/html") {
		t.Errorf("Expected Content-Type text/html, got: %s", rr.Header().Get("Content-Type"))
	}
}

func TestTimeConverterHtmx_HumanToEpoch_Valid(t *testing.T) {
	router := setupTestAPI(RegisterTimeConverterHtmx)

	body := `{"epochToHumanMode":false,"humanToEpochMode":true,"epochTime":0,"year":2024,"month":1,"day":1,"hour":0,"minute":0,"second":0,"timezone":"UTC"}`
	req, _ := http.NewRequest("POST", "/api/timeconverter/htmx", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("HX-Request", "true")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d, response body %v", rr.Code, rr.Body.String())
	}
	if !strings.Contains(rr.Body.String(), "Epoch (Seconds)") {
		t.Errorf("Expected 'Epoch (Seconds)' in HTML response, got: %s", rr.Body.String())
	}
	if !strings.Contains(rr.Body.String(), "1704067200") {
		t.Errorf("Expected epoch value 1704067200 in HTML, got: %s", rr.Body.String())
	}
	if !strings.Contains(rr.Header().Get("Content-Type"), "text/html") {
		t.Errorf("Expected Content-Type text/html, got: %s", rr.Header().Get("Content-Type"))
	}
}

func TestTimeConverterHtmx_HumanToEpoch_WithTimezone(t *testing.T) {
	router := setupTestAPI(RegisterTimeConverterHtmx)

	// 2024-01-01 00:00:00 in Europe/Paris (UTC+1) = 2023-12-31 23:00:00 UTC
	// Expected epoch for 2023-12-31 23:00:00 UTC = 1704067200 - 3600
	body := `{"epochToHumanMode":false,"humanToEpochMode":true,"epochTime":0,"year":2024,"month":1,"day":1,"hour":0,"minute":0,"second":0,"timezone":"Europe/Paris"}`
	req, _ := http.NewRequest("POST", "/api/timeconverter/htmx", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("HX-Request", "true")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d, response body %v", rr.Code, rr.Body.String())
	}
	if !strings.Contains(rr.Body.String(), "1704063600") {
		t.Errorf("Expected epoch value 1704063600 in HTML, got: %s", rr.Body.String())
	}
}

func TestTimeConverterHtmx_InvalidTimezone(t *testing.T) {
	router := setupTestAPI(RegisterTimeConverterHtmx)

	body := `{"epochToHumanMode":true,"humanToEpochMode":false,"epochTime":1704067200,"year":0,"month":0,"day":0,"hour":0,"minute":0,"second":0,"timezone":"Invalid/Timezone"}`
	req, _ := http.NewRequest("POST", "/api/timeconverter/htmx", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("HX-Request", "true")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200 (error displayed in UI), got %d", rr.Code)
	}
	if !strings.Contains(rr.Body.String(), "ERROR") {
		t.Errorf("Expected ERROR in HTML response, got: %s", rr.Body.String())
	}
	if !strings.Contains(rr.Header().Get("Content-Type"), "text/html") {
		t.Errorf("Expected Content-Type text/html, got: %s", rr.Header().Get("Content-Type"))
	}
}

func TestTimeConverterHtmx_MissingHtmxHeader(t *testing.T) {
	router := setupTestAPI(RegisterTimeConverterHtmx)

	body := `{"epochToHumanMode":true,"humanToEpochMode":false,"epochTime":1704067200,"year":0,"month":0,"day":0,"hour":0,"minute":0,"second":0,"timezone":"UTC"}`
	req, _ := http.NewRequest("POST", "/api/timeconverter/htmx", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	// No HX-Request header

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d, response body %v", rr.Code, rr.Body.String())
	}
	if !strings.Contains(rr.Body.String(), "This endpoint only supports HTMX requests") {
		t.Errorf("Expected error message for missing HX-Request, got: %s", rr.Body.String())
	}
}

func TestTimeConverterHtmx_NeitherModeSpecified(t *testing.T) {
	router := setupTestAPI(RegisterTimeConverterHtmx)

	body := `{"epochToHumanMode":false,"humanToEpochMode":false,"epochTime":0,"year":0,"month":0,"day":0,"hour":0,"minute":0,"second":0,"timezone":"UTC"}`
	req, _ := http.NewRequest("POST", "/api/timeconverter/htmx", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("HX-Request", "true")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200 (error displayed in UI), got %d", rr.Code)
	}
	if !strings.Contains(rr.Body.String(), "ERROR") {
		t.Errorf("Expected ERROR in HTML response, got: %s", rr.Body.String())
	}
}
