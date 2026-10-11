package logger

import (
	"context"
	"log/slog"
	"testing"
	"time"
)

func TestInit(t *testing.T) {
	tests := []struct {
		name         string
		inputLevel   string
		expectedSlog slog.Level
		expectedStr  string
	}{
		{"debug", "debug", slog.LevelDebug, "debug"},
		{"warn", "warn", slog.LevelWarn, "warn"},
		{"warning", "warning", slog.LevelWarn, "warning"},
		{"error", "error", slog.LevelError, "error"},
		{"info", "info", slog.LevelInfo, "info"},
		{"default", "somethingelse", slog.LevelInfo, "somethingelse"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Init(tt.inputLevel)
			if Logger == nil {
				t.Error("Logger should not be nil after Init")
			}
			if CurrentLevel != tt.expectedStr {
				t.Errorf("CurrentLevel = %q, want %q", CurrentLevel, tt.expectedStr)
			}
			// We can't directly check slogLevel, but we can check Logger.Handler().Enabled for a given level
			if !Logger.Handler().Enabled(context.Background(), tt.expectedSlog) {
				t.Errorf("Logger should be enabled for level %v", tt.expectedSlog)
			}
		})
	}
}

func TestShortenTimestamp(t *testing.T) {
	loc := time.FixedZone("UTC+3", 3*60*60)
	attr := slog.Time(slog.TimeKey, time.Date(2026, 10, 10, 16, 42, 11, 395_000_000, loc))

	got := shortenTimestamp(nil, attr).Value.String()
	want := "2026-10-10T16:42:11"
	if got != want {
		t.Errorf("shortenTimestamp() = %q, want %q", got, want)
	}
}
