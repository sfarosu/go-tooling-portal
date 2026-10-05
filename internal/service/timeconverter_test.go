package service

import (
	"testing"
	"time"
)

func TestTimeConvert_EpochToHuman_Seconds(t *testing.T) {
	// 2024-01-01 00:00:00 UTC = 1704067200 seconds
	req := &TimeConvertRequest{
		EpochToHumanMode: true,
		EpochTime:        1704067200,
		Timezone:         "UTC",
	}

	result, err := TimeConvert(req)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !result.EpochToHuman {
		t.Errorf("Expected EpochToHuman to be true")
	}
	if result.HumanToEpoch {
		t.Errorf("Expected HumanToEpoch to be false")
	}
	if result.Timezone != "UTC" {
		t.Errorf("Expected timezone UTC, got %s", result.Timezone)
	}

	// Check that converted times are not zero
	if result.ConvertedEpochToHumanUTC.IsZero() {
		t.Errorf("Expected ConvertedEpochToHumanUTC to be non-zero")
	}
	if result.ConvertedEpochToHumanLOCAL.IsZero() {
		t.Errorf("Expected ConvertedEpochToHumanLOCAL to be non-zero")
	}

	// Verify the converted time is correct
	expectedYear := 2024
	if result.ConvertedEpochToHumanUTC.Year() != expectedYear {
		t.Errorf("Expected year %d, got %d", expectedYear, result.ConvertedEpochToHumanUTC.Year())
	}
	if result.ConvertedEpochToHumanUTC.Month() != time.January {
		t.Errorf("Expected month January, got %v", result.ConvertedEpochToHumanUTC.Month())
	}
	if result.ConvertedEpochToHumanUTC.Day() != 1 {
		t.Errorf("Expected day 1, got %d", result.ConvertedEpochToHumanUTC.Day())
	}
}

func TestTimeConvert_EpochToHuman_Milliseconds(t *testing.T) {
	// Same timestamp in milliseconds
	req := &TimeConvertRequest{
		EpochToHumanMode: true,
		EpochTime:        1704067200000,
		Timezone:         "UTC",
	}

	result, err := TimeConvert(req)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !result.EpochToHuman {
		t.Errorf("Expected EpochToHuman to be true")
	}

	// Verify the converted time is correct
	if result.ConvertedEpochToHumanUTC.Year() != 2024 {
		t.Errorf("Expected year 2024, got %d", result.ConvertedEpochToHumanUTC.Year())
	}
}

func TestTimeConvert_EpochToHuman_Microseconds(t *testing.T) {
	// Same timestamp in microseconds
	req := &TimeConvertRequest{
		EpochToHumanMode: true,
		EpochTime:        1704067200000000,
		Timezone:         "UTC",
	}

	result, err := TimeConvert(req)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !result.EpochToHuman {
		t.Errorf("Expected EpochToHuman to be true")
	}

	// Verify the converted time is correct
	if result.ConvertedEpochToHumanUTC.Year() != 2024 {
		t.Errorf("Expected year 2024, got %d", result.ConvertedEpochToHumanUTC.Year())
	}
}

func TestTimeConvert_EpochToHuman_WithTimezone(t *testing.T) {
	// Test with different timezone (Europe/Paris)
	req := &TimeConvertRequest{
		EpochToHumanMode: true,
		EpochTime:        1704067200,
		Timezone:         "Europe/Paris",
	}

	result, err := TimeConvert(req)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if result.Timezone != "Europe/Paris" {
		t.Errorf("Expected timezone Europe/Paris, got %s", result.Timezone)
	}

	// UTC time and local time should be different due to timezone
	if result.ConvertedEpochToHumanUTC.Hour() == result.ConvertedEpochToHumanLOCAL.Hour() {
		t.Errorf("Expected different hours for UTC vs local timezone")
	}
}

func TestTimeConvert_EpochToHuman_InvalidTimezone(t *testing.T) {
	req := &TimeConvertRequest{
		EpochToHumanMode: true,
		EpochTime:        1704067200,
		Timezone:         "Invalid/Timezone",
	}

	_, err := TimeConvert(req)
	if err == nil {
		t.Fatalf("Expected error for invalid timezone, got nil")
	}
}

func TestTimeConvert_EpochToHuman_InvalidPrecision(t *testing.T) {
	// 5-digit number should fail precision detection
	req := &TimeConvertRequest{
		EpochToHumanMode: true,
		EpochTime:        12345,
		Timezone:         "UTC",
	}

	_, err := TimeConvert(req)
	if err == nil {
		t.Fatalf("Expected error for invalid epoch precision, got nil")
	}
}

func TestTimeConvert_HumanToEpoch_Valid(t *testing.T) {
	// 2024-01-01 00:00:00 UTC
	req := &TimeConvertRequest{
		HumanToEpochMode: true,
		Year:             2024,
		Month:            1,
		Day:              1,
		Hour:             0,
		Minute:           0,
		Second:           0,
		Timezone:         "UTC",
	}

	result, err := TimeConvert(req)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !result.HumanToEpoch {
		t.Errorf("Expected HumanToEpoch to be true")
	}
	if result.EpochToHuman {
		t.Errorf("Expected EpochToHuman to be false")
	}

	// Check epoch time is correct (2024-01-01 00:00:00 UTC = 1704067200)
	expectedEpochSeconds := int64(1704067200)
	if result.ConvertedHumanToEpochTimeSeconds != expectedEpochSeconds {
		t.Errorf("Expected epoch seconds %d, got %d", expectedEpochSeconds, result.ConvertedHumanToEpochTimeSeconds)
	}
}

func TestTimeConvert_HumanToEpoch_WithTimezone(t *testing.T) {
	// Same human time but in different timezone (Europe/Paris which is UTC+1 in January)
	req := &TimeConvertRequest{
		HumanToEpochMode: true,
		Year:             2024,
		Month:            1,
		Day:              1,
		Hour:             0,
		Minute:           0,
		Second:           0,
		Timezone:         "Europe/Paris",
	}

	result, err := TimeConvert(req)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if result.Timezone != "Europe/Paris" {
		t.Errorf("Expected timezone Europe/Paris, got %s", result.Timezone)
	}

	// 2024-01-01 00:00:00 CET (UTC+1) should be earlier in epoch than UTC
	// 2024-01-01 00:00:00 CET = 2023-12-31 23:00:00 UTC = 1704067200 - 3600
	expectedEpochSeconds := int64(1704067200) - 3600
	if result.ConvertedHumanToEpochTimeSeconds != expectedEpochSeconds {
		t.Errorf("Expected epoch seconds %d, got %d", expectedEpochSeconds, result.ConvertedHumanToEpochTimeSeconds)
	}
}

func TestTimeConvert_HumanToEpoch_MidYear(t *testing.T) {
	req := &TimeConvertRequest{
		HumanToEpochMode: true,
		Year:             2024,
		Month:            6,
		Day:              15,
		Hour:             12,
		Minute:           30,
		Second:           45,
		Timezone:         "UTC",
	}

	result, err := TimeConvert(req)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Verify the time components are reasonable
	if result.ConvertedHumanToEpochTimeSeconds == 0 {
		t.Errorf("Expected non-zero epoch seconds")
	}
}

func TestTimeConvert_HumanToEpoch_InvalidTimezone(t *testing.T) {
	req := &TimeConvertRequest{
		HumanToEpochMode: true,
		Year:             2024,
		Month:            1,
		Day:              1,
		Hour:             0,
		Minute:           0,
		Second:           0,
		Timezone:         "Invalid/Timezone",
	}

	_, err := TimeConvert(req)
	if err == nil {
		t.Fatalf("Expected error for invalid timezone, got nil")
	}
}

func TestTimeConvert_NeitherModeSpecified(t *testing.T) {
	req := &TimeConvertRequest{
		EpochToHumanMode: false,
		HumanToEpochMode: false,
		Timezone:         "UTC",
	}

	_, err := TimeConvert(req)
	if err == nil {
		t.Fatalf("Expected error when neither mode is specified, got nil")
	}
}

func TestTimeConvert_BothModesSpecified(t *testing.T) {
	// If both are true, EpochToHuman takes precedence based on the if-else logic
	req := &TimeConvertRequest{
		EpochToHumanMode: true,
		EpochTime:        1704067200,
		HumanToEpochMode: true,
		Year:             2024,
		Month:            1,
		Day:              1,
		Hour:             0,
		Minute:           0,
		Second:           0,
		Timezone:         "UTC",
	}

	result, err := TimeConvert(req)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// EpochToHuman should take precedence
	if !result.EpochToHuman {
		t.Errorf("Expected EpochToHuman to be true when both modes specified")
	}
}

func TestTimeConvert_CurrentTimePopulated(t *testing.T) {
	req := &TimeConvertRequest{
		EpochToHumanMode: true,
		EpochTime:        1704067200,
		Timezone:         "UTC",
	}

	result, err := TimeConvert(req)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Current time fields should be populated
	if result.CurrentUTCEpochTimeSeconds == 0 {
		t.Errorf("Expected CurrentUTCEpochTimeSeconds to be non-zero")
	}
	if result.CurrentUTCEpochTimeMilli == 0 {
		t.Errorf("Expected CurrentUTCEpochTimeMilli to be non-zero")
	}
	if result.CurrentUTCEpochTimeMicro == 0 {
		t.Errorf("Expected CurrentUTCEpochTimeMicro to be non-zero")
	}
	if result.CurrentUTCHumanTime == "" {
		t.Errorf("Expected CurrentUTCHumanTime to be non-empty")
	}
	if result.CurrentLOCALHumanTime == "" {
		t.Errorf("Expected CurrentLOCALHumanTime to be non-empty")
	}
}

func TestTimeConvert_EpochPrecisionConsistency(t *testing.T) {
	// Test that millis and micros are consistent with seconds
	epochSeconds := int64(1704067200)
	expectedMillis := epochSeconds * 1000

	req := &TimeConvertRequest{
		EpochToHumanMode: true,
		EpochTime:        epochSeconds,
		Timezone:         "UTC",
	}

	result, err := TimeConvert(req)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// The human timestamps should be the same regardless of input precision
	req2 := &TimeConvertRequest{
		EpochToHumanMode: true,
		EpochTime:        expectedMillis,
		Timezone:         "UTC",
	}

	result2, err := TimeConvert(req2)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Both conversions should produce the same UTC time
	if !result.ConvertedEpochToHumanUTC.Equal(result2.ConvertedEpochToHumanUTC) {
		t.Errorf("Expected same UTC time from seconds and millis conversion")
	}
}
