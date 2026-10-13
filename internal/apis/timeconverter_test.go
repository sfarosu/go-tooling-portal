package apis

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTimeConverter_EpochToHuman_Seconds(t *testing.T) {
	router := setupTestAPI(RegisterTimeConverter)

	body := `{"epochToHumanMode":true,"humanToEpochMode":false,"epochTime":1704067200,"year":0,"month":0,"day":0,"hour":0,"minute":0,"second":0,"timezone":"UTC"}`
	req, _ := http.NewRequest("POST", "/api/timeconverter", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d, response body %v", rr.Code, rr.Body.String())
	}

	var out struct {
		ConvertedEpochToHumanUTC string `json:"convertedEpochToHumanUTC"`
	}

	if err := json.Unmarshal(rr.Body.Bytes(), &out); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if out.ConvertedEpochToHumanUTC == "" {
		t.Errorf("Expected non-empty convertedEpochToHumanUTC")
	}
}

func TestTimeConverter_EpochToHuman_Milliseconds(t *testing.T) {
	router := setupTestAPI(RegisterTimeConverter)

	body := `{"epochToHumanMode":true,"humanToEpochMode":false,"epochTime":1704067200000,"year":0,"month":0,"day":0,"hour":0,"minute":0,"second":0,"timezone":"UTC"}`
	req, _ := http.NewRequest("POST", "/api/timeconverter", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d, response body %v", rr.Code, rr.Body.String())
	}

	var out struct {
		ConvertedEpochToHumanUTC string `json:"convertedEpochToHumanUTC"`
	}

	if err := json.Unmarshal(rr.Body.Bytes(), &out); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if out.ConvertedEpochToHumanUTC == "" {
		t.Errorf("Expected non-empty convertedEpochToHumanUTC")
	}
}

func TestTimeConverter_HumanToEpoch_Valid(t *testing.T) {
	router := setupTestAPI(RegisterTimeConverter)

	body := `{"epochToHumanMode":false,"humanToEpochMode":true,"epochTime":0,"year":2024,"month":1,"day":1,"hour":0,"minute":0,"second":0,"timezone":"UTC"}`
	req, _ := http.NewRequest("POST", "/api/timeconverter", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d, response body %v", rr.Code, rr.Body.String())
	}

	var out struct {
		ConvertedHumanToEpochTimeSeconds int64 `json:"convertedHumanToEpochTimeSeconds"`
	}

	if err := json.Unmarshal(rr.Body.Bytes(), &out); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if out.ConvertedHumanToEpochTimeSeconds != 1704067200 {
		t.Errorf("Expected epoch 1704067200, got %d", out.ConvertedHumanToEpochTimeSeconds)
	}
}

func TestTimeConverter_HumanToEpoch_WithTimezone(t *testing.T) {
	router := setupTestAPI(RegisterTimeConverter)

	// 2024-01-01 00:00:00 in Europe/Paris (UTC+1) = 2023-12-31 23:00:00 UTC
	// Expected epoch for 2023-12-31 23:00:00 UTC = 1704067200 - 3600
	body := `{"epochToHumanMode":false,"humanToEpochMode":true,"epochTime":0,"year":2024,"month":1,"day":1,"hour":0,"minute":0,"second":0,"timezone":"Europe/Paris"}`
	req, _ := http.NewRequest("POST", "/api/timeconverter", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d, response body %v", rr.Code, rr.Body.String())
	}

	var out struct {
		ConvertedHumanToEpochTimeSeconds int64 `json:"convertedHumanToEpochTimeSeconds"`
	}

	if err := json.Unmarshal(rr.Body.Bytes(), &out); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	expectedEpoch := int64(1704063600) // 1704067200 - 3600
	if out.ConvertedHumanToEpochTimeSeconds != expectedEpoch {
		t.Errorf("Expected epoch %d, got %d", expectedEpoch, out.ConvertedHumanToEpochTimeSeconds)
	}
}

func TestTimeConverter_InvalidTimezone(t *testing.T) {
	router := setupTestAPI(RegisterTimeConverter)

	body := `{"epochToHumanMode":true,"humanToEpochMode":false,"epochTime":1704067200,"year":0,"month":0,"day":0,"hour":0,"minute":0,"second":0,"timezone":"Invalid/Timezone"}`
	req, _ := http.NewRequest("POST", "/api/timeconverter", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d, response body %v", rr.Code, rr.Body.String())
	}
}

func TestTimeConverter_NeitherModeSpecified(t *testing.T) {
	router := setupTestAPI(RegisterTimeConverter)

	body := `{"epochToHumanMode":false,"humanToEpochMode":false,"epochTime":0,"year":0,"month":0,"day":0,"hour":0,"minute":0,"second":0,"timezone":"UTC"}`
	req, _ := http.NewRequest("POST", "/api/timeconverter", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d, response body %v", rr.Code, rr.Body.String())
	}
}

func TestTimeConverter_ResponseContainsAllFields(t *testing.T) {
	router := setupTestAPI(RegisterTimeConverter)

	body := `{"epochToHumanMode":true,"humanToEpochMode":false,"epochTime":1704067200,"year":0,"month":0,"day":0,"hour":0,"minute":0,"second":0,"timezone":"UTC"}`
	req, _ := http.NewRequest("POST", "/api/timeconverter", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}

	var out struct {
		Timezone                         string `json:"timezone"`
		EpochToHuman                     bool   `json:"epochToHuman"`
		CurrentUTCEpochTimeSeconds       int64  `json:"currentUTCEpochTimeSeconds"`
		ConvertedEpochToHumanUTC         string `json:"convertedEpochToHumanUTC"`
	}

	if err := json.Unmarshal(rr.Body.Bytes(), &out); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if out.Timezone != "UTC" {
		t.Errorf("Expected timezone UTC, got %s", out.Timezone)
	}
	if !out.EpochToHuman {
		t.Errorf("Expected EpochToHuman to be true")
	}
	if out.CurrentUTCEpochTimeSeconds == 0 {
		t.Errorf("Expected non-zero CurrentUTCEpochTimeSeconds")
	}
	if out.ConvertedEpochToHumanUTC == "" {
		t.Errorf("Expected non-empty ConvertedEpochToHumanUTC")
	}
}
