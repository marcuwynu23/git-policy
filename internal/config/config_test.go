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

func TestAddPlugin(t *testing.T) {
	cfg := DefaultConfig()
	entry := PluginEntry{Name: "test", Path: "/test.so", Enabled: true}
	cfg.AddPlugin(entry)
	if len(cfg.Plugins) != 1 {
		t.Fatalf("expected 1 plugin, got %d", len(cfg.Plugins))
	}
	if cfg.Plugins[0].Name != "test" {
		t.Errorf("expected name 'test', got %q", cfg.Plugins[0].Name)
	}
}

func TestAddPlugin_UpdatesExisting(t *testing.T) {
	cfg := DefaultConfig()
	cfg.AddPlugin(PluginEntry{Name: "test", Path: "/old.so", Enabled: true})
	cfg.AddPlugin(PluginEntry{Name: "test", Path: "/new.so", Enabled: false})
	if len(cfg.Plugins) != 1 {
		t.Fatalf("expected 1 plugin, got %d", len(cfg.Plugins))
	}
	if cfg.Plugins[0].Path != "/new.so" {
		t.Errorf("expected path '/new.so', got %q", cfg.Plugins[0].Path)
	}
	if cfg.Plugins[0].Enabled {
		t.Error("expected plugin to be disabled after update")
	}
}

func TestRemovePlugin(t *testing.T) {
	cfg := DefaultConfig()
	cfg.AddPlugin(PluginEntry{Name: "a", Path: "/a.so", Enabled: true})
	cfg.AddPlugin(PluginEntry{Name: "b", Path: "/b.so", Enabled: true})
	cfg.AddPlugin(PluginEntry{Name: "c", Path: "/c.so", Enabled: true})

	if !cfg.RemovePlugin("b") {
		t.Error("expected RemovePlugin to return true")
	}
	if len(cfg.Plugins) != 2 {
		t.Fatalf("expected 2 plugins, got %d", len(cfg.Plugins))
	}
	if cfg.Plugins[0].Name != "a" || cfg.Plugins[1].Name != "c" {
		t.Errorf("unexpected plugins after remove: %v", cfg.Plugins)
	}

	if cfg.RemovePlugin("nonexistent") {
		t.Error("expected RemovePlugin to return false for missing name")
	}
}

func TestLoadPluginDescriptor(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "plugin.yaml")
	_ = os.WriteFile(path, []byte("name: my-plugin\npath: /tmp/my.so\nenabled: true\n"), 0644)

	desc, err := LoadPluginDescriptor(path)
	if err != nil {
		t.Fatalf("LoadPluginDescriptor failed: %v", err)
	}
	if desc.Name != "my-plugin" {
		t.Errorf("expected name 'my-plugin', got %q", desc.Name)
	}
	if desc.Path != "/tmp/my.so" {
		t.Errorf("expected path '/tmp/my.so', got %q", desc.Path)
	}
	if !desc.Enabled {
		t.Error("expected enabled")
	}
}

func TestLoadPluginDescriptor_MissingName(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "plugin.yaml")
	_ = os.WriteFile(path, []byte("path: /tmp/my.so\n"), 0644)

	_, err := LoadPluginDescriptor(path)
	if err == nil {
		t.Fatal("expected error for missing name")
	}
}

func TestLoadPluginDescriptor_MissingPath(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "plugin.yaml")
	_ = os.WriteFile(path, []byte("name: my-plugin\n"), 0644)

	_, err := LoadPluginDescriptor(path)
	if err == nil {
		t.Fatal("expected error for missing path")
	}
}

func TestLoadPluginDescriptor_FileNotFound(t *testing.T) {
	_, err := LoadPluginDescriptor("/nonexistent/plugin.yaml")
	if err == nil {
		t.Fatal("expected error for nonexistent file")
	}
}

func TestLoadPluginsFromYAML(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	content := []byte(`version: 1
hooks:
  pre-commit:
    enabled: true
policies:
  blockFiles: []
  conventionalCommits: true
plugins:
  - name: plugin-a
    path: /tmp/a.so
    enabled: true
  - name: plugin-b
    path: /tmp/b.so
    enabled: false
`)
	_ = os.WriteFile(path, content, 0644)

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}
	if len(cfg.Plugins) != 2 {
		t.Fatalf("expected 2 plugins, got %d", len(cfg.Plugins))
	}
	if cfg.Plugins[0].Name != "plugin-a" {
		t.Errorf("expected 'plugin-a', got %q", cfg.Plugins[0].Name)
	}
	if cfg.Plugins[1].Enabled {
		t.Error("expected plugin-b to be disabled")
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
