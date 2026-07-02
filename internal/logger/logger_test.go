package logger

import (
	"bytes"
	"log/slog"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	var buf bytes.Buffer
	logger := New(LevelDebug, &buf)
	logger.Debug("test debug")
	if !strings.Contains(buf.String(), "test debug") {
		t.Error("expected debug message in output")
	}
}

func TestNew_LevelFiltering(t *testing.T) {
	var buf bytes.Buffer
	logger := New(LevelError, &buf)
	logger.Info("should not appear")
	if buf.Len() > 0 {
		t.Error("expected no output for info when level is error")
	}
}

func TestNew_NilWriter(t *testing.T) {
	logger := New(LevelInfo, nil)
	if logger == nil {
		t.Fatal("expected logger, got nil")
	}
	logger.Info("test")
}

func TestLevels(t *testing.T) {
	var buf bytes.Buffer
	logger := New(LevelWarn, &buf)
	logger.Warn("warn message", slog.String("key", "value"))
	if !strings.Contains(buf.String(), "warn message") {
		t.Error("expected warn message in output")
	}
}
