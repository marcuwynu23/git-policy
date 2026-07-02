package logger

import (
	"io"
	"log/slog"
	"os"
	"strings"
)

type Level string

const (
	LevelDebug Level = "debug"
	LevelInfo  Level = "info"
	LevelWarn  Level = "warn"
	LevelError Level = "error"
)

func New(level Level, w io.Writer) *slog.Logger {
	if w == nil {
		w = os.Stderr
	}
	var slogLevel slog.Level
	switch strings.ToLower(string(level)) {
	case "debug":
		slogLevel = slog.LevelDebug
	case "info":
		slogLevel = slog.LevelInfo
	case "warn":
		slogLevel = slog.LevelWarn
	case "error":
		slogLevel = slog.LevelError
	default:
		slogLevel = slog.LevelInfo
	}
	handler := slog.NewTextHandler(w, &slog.HandlerOptions{
		Level: slogLevel,
	})
	return slog.New(handler)
}
