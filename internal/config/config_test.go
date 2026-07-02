package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Version != 1 {
		t.Errorf("expected version 1, got %d", cfg.Version)
	}
	if !cfg.Hooks.PreCommit.Enabled {
		t.Error("expected pre-commit to be enabled")
	}
	if len(cfg.Policies.BlockFiles) == 0 {
		t.Error("expected block files")
	}
	if !cfg.Policies.SecretScan {
		t.Error("expected secret scan enabled")
	}
	if !cfg.Policies.ConventionalCommits {
		t.Error("expected conventional commits enabled")
	}
}

func TestLoadNonExistentFile(t *testing.T) {
	cfg, err := Load("/nonexistent/path/config.yaml")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg == nil {
		t.Fatal("expected config, got nil")
	}
}

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")

	cfg := DefaultConfig()
	cfg.Policies.SecretScan = false
	cfg.Policies.ConventionalCommits = false

	if err := Save(cfg, path); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	loaded, err := Load(path)
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}

	if loaded.Policies.SecretScan {
		t.Error("expected secret scan to be false")
	}
	if loaded.Policies.ConventionalCommits {
		t.Error("expected conventional commits to be false")
	}
}

func TestDefaultConfigPath(t *testing.T) {
	path, err := DefaultConfigPath()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if path == "" {
		t.Fatal("expected non-empty path")
	}
}

func TestLoadFromFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	content := []byte(`version: 1
hooks:
  pre-commit:
    enabled: true
policies:
  blockFiles:
    - ".env"
  conventionalCommits: true
`)
	if err := os.WriteFile(path, content, 0644); err != nil {
		t.Fatalf("write failed: %v", err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}

	if len(cfg.Policies.BlockFiles) != 1 || cfg.Policies.BlockFiles[0] != ".env" {
		t.Errorf("unexpected block files: %v", cfg.Policies.BlockFiles)
	}
	if !cfg.Policies.ConventionalCommits {
		t.Error("expected conventional commits enabled")
	}
}
