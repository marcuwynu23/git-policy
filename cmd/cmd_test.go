package cmd

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/marcuwynu23/git-policy/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func executeCommand(root *cobra.Command, args ...string) (output string, err error) {
	// Reset all flags on all commands to avoid stale state between tests
	resetFlags(root)

	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)
	err = root.Execute()
	return buf.String(), err
}

func resetFlags(cmd *cobra.Command) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		f.Changed = false
		_ = f.Value.Set(f.DefValue)
	})
	for _, sub := range cmd.Commands() {
		resetFlags(sub)
	}
}

func TestCommandsAdded(t *testing.T) {
	expected := []string{
		"install",
		"uninstall",
		"run",
		"doctor",
		"validate",
		"version",
		"sync",
		"rule",
		"plugins",
	}
	for _, name := range expected {
		cmd, _, err := rootCmd.Find([]string{name})
		if err != nil {
			t.Errorf("expected command %q to be registered, got error: %v", name, err)
		}
		if cmd == nil {
			t.Errorf("expected command %q to be non-nil", name)
		}
	}
}

func TestRuleSubcommandsAdded(t *testing.T) {
	expected := []string{
		"list",
		"enable",
		"disable",
		"skip",
		"add",
		"remove",
		"export",
		"import",
	}
	for _, name := range expected {
		cmd, _, err := rootCmd.Find([]string{"rule", name})
		if err != nil {
			t.Errorf("expected rule subcommand %q to be registered, got error: %v", name, err)
		}
		if cmd == nil {
			t.Errorf("expected rule subcommand %q to be non-nil", name)
		}
	}
}

func TestPluginsSubcommandsAdded(t *testing.T) {
	expected := []string{
		"install",
		"list",
	}
	for _, name := range expected {
		cmd, _, err := rootCmd.Find([]string{"plugins", name})
		if err != nil {
			t.Errorf("expected plugins subcommand %q to be registered, got error: %v", name, err)
		}
		if cmd == nil {
			t.Errorf("expected plugins subcommand %q to be non-nil", name)
		}
	}
}

func TestVersionCommand(t *testing.T) {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	rootCmd.SetArgs([]string{"version"})
	_ = rootCmd.Execute()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	output := buf.String()

	if !strings.Contains(output, "git-policy vdev") {
		t.Errorf("expected output containing 'git-policy vdev', got %q", output)
	}
}

func TestRuleSkip_NotInRepo(t *testing.T) {
	dir := t.TempDir()
	origDir, _ := os.Getwd()
	defer func() { _ = os.Chdir(origDir) }()
	_ = os.Chdir(dir)

	_, err := executeCommand(rootCmd, "rule", "skip", "block-files")
	if err == nil {
		t.Error("expected error when not in a git repo")
	}
	if err != nil && !strings.Contains(err.Error(), "not a git repository") {
		t.Errorf("expected 'not a git repository', got %v", err)
	}
}

func TestRuleSkip_SkipAndList(t *testing.T) {
	dir := t.TempDir()
	origDir, _ := os.Getwd()
	defer func() { _ = os.Chdir(origDir) }()

	initCmdRepo(t, dir)
	_ = os.Chdir(dir)

	// Skip two rules
	_, err := executeCommand(rootCmd, "rule", "skip", "block-files", "secret-scan")
	if err != nil {
		t.Fatalf("skip failed: %v", err)
	}

	// List skipped rules
	output, err := executeCommand(rootCmd, "rule", "skip", "--list")
	if err != nil {
		t.Fatalf("skip --list failed: %v", err)
	}
	if !strings.Contains(output, "block-files") || !strings.Contains(output, "secret-scan") {
		t.Errorf("expected output to contain skipped rules, got: %s", output)
	}
}

func TestRuleSkip_Clear(t *testing.T) {
	dir := t.TempDir()
	origDir, _ := os.Getwd()
	defer func() { _ = os.Chdir(origDir) }()

	initCmdRepo(t, dir)
	_ = os.Chdir(dir)

	_, _ = executeCommand(rootCmd, "rule", "skip", "block-files")
	_, err := executeCommand(rootCmd, "rule", "skip", "--clear")
	if err != nil {
		t.Fatalf("skip --clear failed: %v", err)
	}

	output, err := executeCommand(rootCmd, "rule", "skip", "--list")
	if err != nil {
		t.Fatalf("skip --list failed: %v", err)
	}
	if !strings.Contains(output, "No rules currently skipped") {
		t.Errorf("expected 'No rules currently skipped', got: %s", output)
	}
}

func TestRuleSkip_InvalidName(t *testing.T) {
	dir := t.TempDir()
	origDir, _ := os.Getwd()
	defer func() { _ = os.Chdir(origDir) }()

	initCmdRepo(t, dir)
	_ = os.Chdir(dir)

	_, err := executeCommand(rootCmd, "rule", "skip", "nonexistent-rule")
	if err == nil {
		t.Error("expected error for invalid rule name")
	}
	if err != nil && !strings.Contains(err.Error(), "unknown rule") {
		t.Errorf("expected 'unknown rule', got %v", err)
	}
}

func TestRuleSkip_NoArgs(t *testing.T) {
	dir := t.TempDir()
	origDir, _ := os.Getwd()
	defer func() { _ = os.Chdir(origDir) }()

	initCmdRepo(t, dir)
	_ = os.Chdir(dir)

	// No args with no skips set
	output, err := executeCommand(rootCmd, "rule", "skip")
	if err != nil {
		t.Fatalf("skip with no args failed: %v", err)
	}
	if !strings.Contains(output, "No rules currently skipped") {
		t.Errorf("expected 'No rules currently skipped', got: %s", output)
	}

	// Add a skip, then run with no args
	_, _ = executeCommand(rootCmd, "rule", "skip", "block-files")
	output, err = executeCommand(rootCmd, "rule", "skip")
	if err != nil {
		t.Fatalf("skip with no args failed: %v", err)
	}
	if !strings.Contains(output, "block-files") {
		t.Errorf("expected output to contain 'block-files', got: %s", output)
	}
}

func TestRuleSkip_Duplicate(t *testing.T) {
	dir := t.TempDir()
	origDir, _ := os.Getwd()
	defer func() { _ = os.Chdir(origDir) }()

	initCmdRepo(t, dir)
	_ = os.Chdir(dir)

	_, _ = executeCommand(rootCmd, "rule", "skip", "block-files")
	_, err := executeCommand(rootCmd, "rule", "skip", "block-files")
	if err != nil {
		t.Fatalf("skip duplicate failed: %v", err)
	}

	output, _ := executeCommand(rootCmd, "rule", "skip", "--list")
	count := strings.Count(output, "block-files")
	if count != 1 {
		t.Errorf("expected 'block-files' to appear once, got %d", count)
	}
}

func initCmdRepo(t *testing.T, dir string) {
	t.Helper()
	_ = os.Chdir(dir)
	cmd := exec.Command("git", "init")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git init failed: %s: %v", string(out), err)
	}
	_ = exec.Command("git", "-C", dir, "config", "user.email", "test@test.com").Run()
	_ = exec.Command("git", "-C", dir, "config", "user.name", "Test").Run()
	_ = exec.Command("git", "-C", dir, "config", "commit.gpgSign", "false").Run()
	readme := filepath.Join(dir, "README.md")
	_ = os.WriteFile(readme, []byte("# test"), 0644)
	_ = exec.Command("git", "-C", dir, "add", "README.md").Run()
	_ = exec.Command("git", "-C", dir, "commit", "-m", "feat: initial commit").Run()
}

func TestPluginsInstall_InvalidFile(t *testing.T) {
	_, err := executeCommand(rootCmd, "plugins", "install", "nonexistent.yaml")
	if err == nil {
		t.Error("expected error for nonexistent file")
	}
	if err != nil && !strings.Contains(err.Error(), "reading plugin descriptor") {
		t.Errorf("expected 'reading plugin descriptor', got %v", err)
	}
}

func TestPluginsList_Empty(t *testing.T) {
	// Need a temp config to avoid modifying real config
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.yaml")
	_ = os.WriteFile(configPath, []byte("version: 1\nhooks:\n  pre-commit:\n    enabled: true\n  commit-msg:\n    enabled: true\n  pre-push:\n    enabled: true\n  post-merge:\n    enabled: false\npolicies:\n  blockFiles: []\n  maxFileSize: 10MB\n  secretScan: true\n  protectedBranches: []\n  conventionalCommits: true\n  blockBinaries: []\n"), 0644)

	output, err := executeCommand(rootCmd, "--config", configPath, "plugins", "list")
	if err != nil {
		t.Fatalf("plugins list failed: %v", err)
	}
	if !strings.Contains(output, "No plugins installed") {
		t.Errorf("expected 'No plugins installed', got: %s", output)
	}
}

func TestPluginsList_WithPlugins(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.yaml")
	yamlContent := `version: 1
hooks:
  pre-commit:
    enabled: true
  commit-msg:
    enabled: true
  pre-push:
    enabled: true
  post-merge:
    enabled: false
policies:
  blockFiles: []
  maxFileSize: 10MB
  secretScan: true
  protectedBranches: []
  conventionalCommits: true
  blockBinaries: []
plugins:
  - name: my-plugin
    path: /tmp/my-plugin.so
    enabled: true
  - name: disabled-plugin
    path: /tmp/disabled.so
    enabled: false
`
	_ = os.WriteFile(configPath, []byte(yamlContent), 0644)

	output, err := executeCommand(rootCmd, "--config", configPath, "plugins", "list")
	if err != nil {
		t.Fatalf("plugins list failed: %v", err)
	}
	if !strings.Contains(output, "my-plugin") {
		t.Errorf("expected output to contain 'my-plugin', got: %s", output)
	}
	if !strings.Contains(output, "disabled") {
		t.Errorf("expected output to contain 'disabled', got: %s", output)
	}
}

func TestPluginsInstall_FromDescriptor(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.yaml")
	_ = os.WriteFile(configPath, []byte("version: 1\nhooks:\n  pre-commit:\n    enabled: true\n  commit-msg:\n    enabled: true\n  pre-push:\n    enabled: true\n  post-merge:\n    enabled: false\npolicies:\n  blockFiles: []\n  maxFileSize: 10MB\n  secretScan: true\n  protectedBranches: []\n  conventionalCommits: true\n  blockBinaries: []\n"), 0644)

	pluginDesc := filepath.Join(dir, "plugin.yaml")
	_ = os.WriteFile(pluginDesc, []byte("name: test-plugin\npath: /tmp/test.so\nenabled: true\n"), 0644)

	output, err := executeCommand(rootCmd, "--config", configPath, "plugins", "install", pluginDesc)
	if err != nil {
		t.Fatalf("plugins install failed: %v", err)
	}
	if !strings.Contains(output, "test-plugin") {
		t.Errorf("expected output to contain 'test-plugin', got: %s", output)
	}

	// Verify it was saved to config
	cfg, err := config.Load(configPath)
	if err != nil {
		t.Fatalf("loading config after install: %v", err)
	}
	found := false
	for _, p := range cfg.Plugins {
		if p.Name == "test-plugin" {
			found = true
			if p.Path != "/tmp/test.so" {
				t.Errorf("expected path '/tmp/test.so', got %q", p.Path)
			}
			if !p.Enabled {
				t.Error("expected plugin to be enabled")
			}
			break
		}
	}
	if !found {
		t.Error("expected test-plugin to be in config")
	}
}

func TestUnimplementedCommandsReturnError(t *testing.T) {
	tests := []struct {
		args []string
		msg  string
	}{
		{[]string{"rule", "add", "myrule"}, "rule add not yet implemented"},
		{[]string{"rule", "remove", "myrule"}, "rule remove not yet implemented"},
		{[]string{"rule", "export"}, "rule export not yet implemented"},
		{[]string{"rule", "import", "file.yaml"}, "rule import not yet implemented"},
	}
	for _, tt := range tests {
		_, err := executeCommand(rootCmd, tt.args...)
		if err == nil {
			t.Errorf("expected error for %v, got nil", tt.args)
		}
		if !strings.Contains(err.Error(), tt.msg) {
			t.Errorf("expected error %q to contain %q", err.Error(), tt.msg)
		}
	}
}
