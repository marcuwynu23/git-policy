package files

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIsBlockedFile(t *testing.T) {
	checker := NewChecker(nil, nil, 0)
	tests := []struct {
		name     string
		expected bool
	}{
		{".env", true},
		{"file.pem", true},
		{"file.key", true},
		{"main.go", false},
		{"README.md", false},
	}
	for _, tt := range tests {
		blocked, pattern := checker.IsBlockedFile(tt.name)
		if blocked != tt.expected {
			t.Errorf("IsBlockedFile(%q) = %v, want %v (pattern: %s)", tt.name, blocked, tt.expected, pattern)
		}
	}
}

func TestIsBinaryFile(t *testing.T) {
	checker := NewChecker(nil, nil, 0)
	tests := []struct {
		name     string
		expected bool
	}{
		{"malware.exe", true},
		{"library.dll", true},
		{"image.iso", true},
		{"source.go", false},
		{"text.txt", false},
	}
	for _, tt := range tests {
		blocked, ext := checker.IsBinaryFile(tt.name)
		if blocked != tt.expected {
			t.Errorf("IsBinaryFile(%q) = %v, want %v (ext: %s)", tt.name, blocked, tt.expected, ext)
		}
	}
}

func TestIsOverMaxSize(t *testing.T) {
	checker := NewChecker(nil, nil, 100)
	dir := t.TempDir()
	smallFile := filepath.Join(dir, "small.txt")
	_ = os.WriteFile(smallFile, []byte("hello"), 0644)

	over, size, err := checker.IsOverMaxSize(smallFile)
	if err != nil {
		t.Fatalf("check failed: %v", err)
	}
	if over {
		t.Errorf("expected under max size, size=%d", size)
	}
}
