package cmd

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func executeCommand(root *cobra.Command, args ...string) (output string, err error) {
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)
	err = root.Execute()
	return buf.String(), err
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
	rootCmd.Execute()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	if !strings.Contains(output, "git-policy vdev") {
		t.Errorf("expected output containing 'git-policy vdev', got %q", output)
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
		{[]string{"plugins", "install", "plugin.so"}, "plugins install not yet implemented"},
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
