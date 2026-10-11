package logger

import (
	"log/slog"
	"os"
	"strings"
)

var Logger *slog.Logger
var CurrentLevel string

// Init initializes the global logger with the specified log level
func Init(level string) {
	CurrentLevel = level
	var slogLevel slog.Level
	switch strings.ToLower(level) {
	case "debug":
		slogLevel = slog.LevelDebug
	case "warn", "warning":
		slogLevel = slog.LevelWarn
	case "error":
		slogLevel = slog.LevelError
	default:
		slogLevel = slog.LevelInfo
	}

	Logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:       slogLevel,
		AddSource:   false,
		ReplaceAttr: shortenTimestamp,
	}))
}

func shortenTimestamp(_ []string, attr slog.Attr) slog.Attr {
	if attr.Key == slog.TimeKey {
		return slog.String(slog.TimeKey, attr.Value.Time().Format("2006-01-02T15:04:05"))
	}
	return attr
}
