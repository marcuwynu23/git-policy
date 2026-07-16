package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFileExists(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "test.txt")
	_ = os.WriteFile(file, []byte("hello"), 0644)

	if !FileExists(file) {
		t.Error("expected file to exist")
	}
	if FileExists(filepath.Join(dir, "nonexistent.txt")) {
		t.Error("expected file not to exist")
	}
}

func TestEnsureDir(t *testing.T) {
	dir := t.TempDir()
	newDir := filepath.Join(dir, "a", "b", "c")
	if err := EnsureDir(newDir); err != nil {
		t.Fatalf("EnsureDir failed: %v", err)
	}
	if !FileExists(newDir) {
		t.Error("expected directory to exist")
	}
}

func TestConfigDir(t *testing.T) {
	dir, err := ConfigDir()
	if err != nil {
		t.Skipf("ConfigDir failed: %v", err)
	}
	if dir == "" {
		t.Error("expected non-empty config dir")
	}
}

func TestHookDir(t *testing.T) {
	dir, err := HookDir()
	if err != nil {
		t.Skipf("HookDir failed: %v", err)
	}
	if dir == "" {
		t.Error("expected non-empty hook dir")
	}
}
