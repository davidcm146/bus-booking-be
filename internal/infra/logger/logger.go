package logger

import (
	"log/slog"
	"os"
)

// New creates a structured logger. Uses JSON output for production
// and human-readable text output for development.
func New(env string) *slog.Logger {
	var handler slog.Handler

	switch env {
	case "production", "staging":
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	default:
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
	}

	return slog.New(handler)
}
