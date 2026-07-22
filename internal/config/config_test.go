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
	entry := PluginEntry{
		Name:    "test",
		Enabled: true,
		Rules:   []CustomRuleDef{{Name: "r1", Type: "file-block", Pattern: "*.zip", Message: "no zips"}},
	}
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
	cfg.AddPlugin(PluginEntry{
		Name:    "test",
		Enabled: true,
		Rules:   []CustomRuleDef{{Name: "r1", Type: "file-block", Pattern: "*.zip", Message: "no zips"}},
	})
	cfg.AddPlugin(PluginEntry{
		Name:    "test",
		Enabled: false,
		Rules:   []CustomRuleDef{{Name: "r2", Type: "file-content", Pattern: "TODO:", Message: "no todos"}},
	})
	if len(cfg.Plugins) != 1 {
		t.Fatalf("expected 1 plugin, got %d", len(cfg.Plugins))
	}
	if len(cfg.Plugins[0].Rules) != 1 || cfg.Plugins[0].Rules[0].Name != "r2" {
		t.Errorf("expected updated rules, got %+v", cfg.Plugins[0].Rules)
	}
	if cfg.Plugins[0].Enabled {
		t.Error("expected plugin to be disabled after update")
	}
}

func TestRemovePlugin(t *testing.T) {
	cfg := DefaultConfig()
	rule := []CustomRuleDef{{Name: "r", Type: "file-block", Pattern: "*.zip", Message: "no"}}
	cfg.AddPlugin(PluginEntry{Name: "a", Enabled: true, Rules: rule})
	cfg.AddPlugin(PluginEntry{Name: "b", Enabled: true, Rules: rule})
	cfg.AddPlugin(PluginEntry{Name: "c", Enabled: true, Rules: rule})

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
	_ = os.WriteFile(path, []byte("name: my-plugin\nrules:\n  - name: r1\n    type: file-block\n    pattern: \"*.zip\"\n    message: no zips\n"), 0644)

	desc, err := LoadPluginDescriptor(path)
	if err != nil {
		t.Fatalf("LoadPluginDescriptor failed: %v", err)
	}
	if desc.Name != "my-plugin" {
		t.Errorf("expected name 'my-plugin', got %q", desc.Name)
	}
	if len(desc.Rules) != 1 || desc.Rules[0].Name != "r1" {
		t.Errorf("unexpected rules: %+v", desc.Rules)
	}
}

func TestLoadPluginDescriptor_MissingName(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "plugin.yaml")
	_ = os.WriteFile(path, []byte("rules:\n  - name: r1\n    type: file-block\n    pattern: \"*.zip\"\n    message: no\n"), 0644)

	_, err := LoadPluginDescriptor(path)
	if err == nil {
		t.Fatal("expected error for missing name")
	}
}

func TestLoadPluginDescriptor_MissingRules(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "plugin.yaml")
	_ = os.WriteFile(path, []byte("name: my-plugin\n"), 0644)

	_, err := LoadPluginDescriptor(path)
	if err == nil {
		t.Fatal("expected error for missing rules")
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
    enabled: true
    rules:
      - name: r1
        type: file-block
        pattern: "*.zip"
        message: no zips
  - name: plugin-b
    enabled: false
    rules:
      - name: r2
        type: file-content
        pattern: "TODO:"
        message: no todos
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
