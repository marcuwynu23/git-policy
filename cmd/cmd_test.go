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
		"uninstall",
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
    path: /tmp/my-plugin.yaml
    enabled: true
  - name: disabled-plugin
    path: /tmp/disabled.yaml
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
	_ = os.WriteFile(pluginDesc, []byte("name: test-plugin\nrules:\n  - name: r1\n    type: file-block\n    pattern: \"*.zip\"\n    message: no zips\n"), 0644)

	output, err := executeCommand(rootCmd, "--config", configPath, "plugins", "install", pluginDesc)
	if err != nil {
		t.Fatalf("plugins install failed: %v", err)
	}
	if !strings.Contains(output, "test-plugin") {
		t.Errorf("expected output to contain 'test-plugin', got: %s", output)
	}

	// Verify it was saved to config and file was copied to plugins dir
	cfg, err := config.Load(configPath)
	if err != nil {
		t.Fatalf("loading config after install: %v", err)
	}
	found := false
	for _, p := range cfg.Plugins {
		if p.Name == "test-plugin" {
			found = true
			expectedPath := filepath.Join(config.PluginsDir(configPath), "test-plugin.yaml")
			if p.Path != expectedPath {
				t.Errorf("expected path %q, got %q", expectedPath, p.Path)
			}
			if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
				t.Errorf("plugin file was not copied to %s", expectedPath)
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

func TestPluginsInstall_DisabledFlag(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.yaml")
	_ = os.WriteFile(configPath, []byte("version: 1\nhooks:\n  pre-commit:\n    enabled: true\npolicies:\n  blockFiles: []\n  conventionalCommits: true\n"), 0644)

	pluginDesc := filepath.Join(dir, "plugin.yaml")
	_ = os.WriteFile(pluginDesc, []byte("name: test-plugin\nrules:\n  - name: r1\n    type: file-block\n    pattern: \"*.zip\"\n    message: no zips\n"), 0644)

	output, err := executeCommand(rootCmd, "--config", configPath, "plugins", "install", "--disabled", pluginDesc)
	if err != nil {
		t.Fatalf("plugins install --disabled failed: %v", err)
	}
	if !strings.Contains(output, "disabled") {
		t.Errorf("expected output to contain 'disabled', got: %s", output)
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		t.Fatalf("loading config after install: %v", err)
	}
	var found *config.PluginEntry
	for _, p := range cfg.Plugins {
		if p.Name == "test-plugin" {
			found = &p
			break
		}
	}
	if found == nil {
		t.Fatal("expected test-plugin to be in config")
	}
	if found.Enabled {
		t.Error("expected plugin to be disabled when --disabled flag used")
	}
	expectedPath := filepath.Join(config.PluginsDir(configPath), "test-plugin.yaml")
	if found.Path != expectedPath {
		t.Errorf("expected path %q, got %q", expectedPath, found.Path)
	}
	if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
		t.Errorf("plugin file was not copied to %s", expectedPath)
	}
}

func TestPluginsUninstall_Success(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.yaml")
	pluginsDir := filepath.Join(dir, "plugins")
	_ = os.MkdirAll(pluginsDir, 0755)

	myPluginFile := filepath.Join(pluginsDir, "my-plugin.yaml")
	otherPluginFile := filepath.Join(pluginsDir, "other-plugin.yaml")
	_ = os.WriteFile(myPluginFile, []byte("name: my-plugin\nrules:\n  - name: r1\n    type: file-block\n    pattern: \"*.zip\"\n    message: no\n"), 0644)
	_ = os.WriteFile(otherPluginFile, []byte("name: other-plugin\nrules: []\n"), 0644)

	yamlContent := "version: 1\nhooks:\n  pre-commit:\n    enabled: true\npolicies:\n  blockFiles: []\n  conventionalCommits: true\nplugins:\n  - name: my-plugin\n    path: " + myPluginFile + "\n    enabled: true\n  - name: other-plugin\n    path: " + otherPluginFile + "\n    enabled: true\n"
	_ = os.WriteFile(configPath, []byte(yamlContent), 0644)

	output, err := executeCommand(rootCmd, "--config", configPath, "plugins", "uninstall", "my-plugin")
	if err != nil {
		t.Fatalf("plugins uninstall failed: %v", err)
	}
	if !strings.Contains(output, "my-plugin") {
		t.Errorf("expected output to contain 'my-plugin', got: %s", output)
	}
	if !strings.Contains(output, "uninstalled") {
		t.Errorf("expected output to contain 'uninstalled', got: %s", output)
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		t.Fatalf("loading config after uninstall: %v", err)
	}
	for _, p := range cfg.Plugins {
		if p.Name == "my-plugin" {
			t.Error("expected my-plugin to be removed from config")
		}
	}
	if len(cfg.Plugins) != 1 {
		t.Errorf("expected 1 plugin remaining, got %d", len(cfg.Plugins))
	}

	if _, err := os.Stat(myPluginFile); !os.IsNotExist(err) {
		t.Errorf("expected plugin file %s to be deleted", myPluginFile)
	}
}

func TestPluginsUninstall_NotFound(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.yaml")
	_ = os.WriteFile(configPath, []byte("version: 1\nhooks:\n  pre-commit:\n    enabled: true\npolicies:\n  blockFiles: []\n  conventionalCommits: true\n"), 0644)

	_, err := executeCommand(rootCmd, "--config", configPath, "plugins", "uninstall", "nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent plugin")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("expected 'not found', got: %v", err)
	}
}

func TestRuleAdd_RequiresFlags(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.yaml")
	_ = os.WriteFile(configPath, []byte("version: 1\nhooks:\n  pre-commit:\n    enabled: true\npolicies:\n  blockFiles: []\n  conventionalCommits: true\n"), 0644)

	_, err := executeCommand(rootCmd, "--config", configPath, "rule", "add", "myrule")
	if err == nil {
		t.Fatal("expected error for missing --type flag")
	}
}

func TestRuleAdd_Success(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.yaml")
	_ = os.WriteFile(configPath, []byte("version: 1\nhooks:\n  pre-commit:\n    enabled: true\npolicies:\n  blockFiles: []\n  conventionalCommits: true\n"), 0644)

	output, err := executeCommand(rootCmd, "--config", configPath,
		"rule", "add", "no-zips", "--type", "file-block", "--pattern", "*.zip", "--message", "No zips", "--fix", "Remove zips")
	if err != nil {
		t.Fatalf("rule add failed: %v", err)
	}
	if !strings.Contains(output, "no-zips") {
		t.Errorf("expected output to contain 'no-zips', got: %s", output)
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		t.Fatalf("loading config: %v", err)
	}
	found := false
	for _, r := range cfg.CustomRules {
		if r.Name == "no-zips" {
			found = true
			if r.Type != "file-block" {
				t.Errorf("expected type 'file-block', got %q", r.Type)
			}
			if r.Pattern != "*.zip" {
				t.Errorf("expected pattern '*.zip', got %q", r.Pattern)
			}
			if r.Message != "No zips" {
				t.Errorf("expected message 'No zips', got %q", r.Message)
			}
			if r.Fix != "Remove zips" {
				t.Errorf("expected fix 'Remove zips', got %q", r.Fix)
			}
			break
		}
	}
	if !found {
		t.Error("expected no-zips to be in custom rules")
	}
}

func TestRuleRemove_Success(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.yaml")
	yamlContent := "version: 1\nhooks:\n  pre-commit:\n    enabled: true\npolicies:\n  blockFiles: []\n  conventionalCommits: true\ncustomRules:\n  - name: my-rule\n    type: file-block\n    pattern: \"*.zip\"\n    message: no\n"
	_ = os.WriteFile(configPath, []byte(yamlContent), 0644)

	output, err := executeCommand(rootCmd, "--config", configPath, "rule", "remove", "my-rule")
	if err != nil {
		t.Fatalf("rule remove failed: %v", err)
	}
	if !strings.Contains(output, "my-rule") || !strings.Contains(output, "removed") {
		t.Errorf("unexpected output: %s", output)
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		t.Fatalf("loading config: %v", err)
	}
	for _, r := range cfg.CustomRules {
		if r.Name == "my-rule" {
			t.Error("expected my-rule to be removed")
		}
	}
}

func TestRuleRemove_NotFound(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.yaml")
	_ = os.WriteFile(configPath, []byte("version: 1\nhooks:\n  pre-commit:\n    enabled: true\npolicies:\n  blockFiles: []\n  conventionalCommits: true\n"), 0644)

	_, err := executeCommand(rootCmd, "--config", configPath, "rule", "remove", "nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent rule")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("expected 'not found', got: %v", err)
	}
}

func TestRuleExport_Success(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.yaml")
	yamlContent := "version: 1\nhooks:\n  pre-commit:\n    enabled: true\npolicies:\n  blockFiles: []\n  conventionalCommits: true\ncustomRules:\n  - name: my-rule\n    type: file-block\n    pattern: \"*.zip\"\n    message: no\n"
	_ = os.WriteFile(configPath, []byte(yamlContent), 0644)

	outFile := filepath.Join(dir, "exported.yaml")
	output, err := executeCommand(rootCmd, "--config", configPath, "rule", "export", "my-rule", "--output", outFile)
	if err != nil {
		t.Fatalf("rule export failed: %v", err)
	}
	if !strings.Contains(output, "my-rule") {
		t.Errorf("expected output to contain 'my-rule', got: %s", output)
	}

	data, err := os.ReadFile(outFile)
	if err != nil {
		t.Fatalf("reading exported file: %v", err)
	}
	if !strings.Contains(string(data), "my-rule") {
		t.Errorf("exported file missing rule name: %s", string(data))
	}
}

func TestRuleImport_Success(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.yaml")
	_ = os.WriteFile(configPath, []byte("version: 1\nhooks:\n  pre-commit:\n    enabled: true\npolicies:\n  blockFiles: []\n  conventionalCommits: true\n"), 0644)

	ruleFile := filepath.Join(dir, "rule.yaml")
	_ = os.WriteFile(ruleFile, []byte("name: imported-rule\ntype: file-content\npattern: \"TODO:\"\nmessage: No todos\nfix: Remove TODO\n"), 0644)

	output, err := executeCommand(rootCmd, "--config", configPath, "rule", "import", ruleFile)
	if err != nil {
		t.Fatalf("rule import failed: %v", err)
	}
	if !strings.Contains(output, "imported-rule") {
		t.Errorf("expected output to contain 'imported-rule', got: %s", output)
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		t.Fatalf("loading config: %v", err)
	}
	found := false
	for _, r := range cfg.CustomRules {
		if r.Name == "imported-rule" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected imported-rule to be in config")
	}
}
