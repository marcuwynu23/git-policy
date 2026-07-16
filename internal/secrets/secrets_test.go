package secrets

import (
	"os"
	"path/filepath"
	"testing"
)

func TestScanFile_NoSecrets(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "safe.txt")
	_ = os.WriteFile(file, []byte("hello world\nthis is safe"), 0644)

	finder := NewFinder(nil)
	findings, err := finder.ScanFile(file)
	if err != nil {
		t.Fatalf("scan failed: %v", err)
	}
	if len(findings) != 0 {
		t.Errorf("expected 0 findings, got %d", len(findings))
	}
}

func TestScanFile_WithSecrets(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "secret.txt")
	_ = os.WriteFile(file, []byte("aws_key=AKIA1234567890\n"), 0644)

	finder := NewFinder(nil)
	findings, err := finder.ScanFile(file)
	if err != nil {
		t.Fatalf("scan failed: %v", err)
	}
	if len(findings) == 0 {
		t.Fatal("expected findings, got 0")
	}
	if findings[0].Pattern.Name != "AWS Access Key ID" {
		t.Errorf("expected AWS Access Key ID, got %s", findings[0].Pattern.Name)
	}
}

func TestScanFiles_Multiple(t *testing.T) {
	dir := t.TempDir()
	f1 := filepath.Join(dir, "f1.txt")
	f2 := filepath.Join(dir, "f2.txt")
	os.WriteFile(f1, []byte("ghp_abcdefghijklmnop"), 0644)
	os.WriteFile(f2, []byte("safe content"), 0644)

	finder := NewFinder(nil)
	findings, err := finder.ScanFiles([]string{f1, f2})
	if err != nil {
		t.Fatalf("scan failed: %v", err)
	}
	if len(findings) != 1 {
		t.Errorf("expected 1 finding, got %d", len(findings))
	}
}
