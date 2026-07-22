package policy

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/marcuwynu23/git-policy/internal/config"
)

func TestCustomPolicy_FileBlock_Pass(t *testing.T) {
	p := NewCustomPolicy(config.CustomRuleDef{
		Name: "no-zips", Type: "file-block", Pattern: "*.zip", Message: "no zips",
	})
	result := p.Validate(Context{StagedFiles: []string{"main.go"}})
	if result.Status != StatusPass {
		t.Errorf("expected pass, got %s", result.Status)
	}
}

func TestCustomPolicy_FileBlock_Fail(t *testing.T) {
	p := NewCustomPolicy(config.CustomRuleDef{
		Name: "no-zips", Type: "file-block", Pattern: "*.zip", Message: "no zips", Fix: "remove zips",
	})
	result := p.Validate(Context{StagedFiles: []string{"archive.zip"}})
	if result.Status != StatusFail {
		t.Errorf("expected fail, got %s", result.Status)
	}
	if result.Message != "no zips" {
		t.Errorf("expected message 'no zips', got %q", result.Message)
	}
	if result.Fix != "remove zips" {
		t.Errorf("expected fix 'remove zips', got %q", result.Fix)
	}
}

func TestCustomPolicy_FileContent_Pass(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "safe.go")
	_ = os.WriteFile(f, []byte("package main\n"), 0644)

	p := NewCustomPolicy(config.CustomRuleDef{
		Name: "no-todo", Type: "file-content", Pattern: "TODO:", Message: "no todos",
	})
	result := p.Validate(Context{StagedFiles: []string{f}})
	if result.Status != StatusPass {
		t.Errorf("expected pass, got %s", result.Status)
	}
}

func TestCustomPolicy_FileContent_Fail(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "todo.go")
	_ = os.WriteFile(f, []byte("// TODO: implement this\n"), 0644)

	p := NewCustomPolicy(config.CustomRuleDef{
		Name: "no-todo", Type: "file-content", Pattern: "TODO:", Message: "no todos",
	})
	result := p.Validate(Context{StagedFiles: []string{f}})
	if result.Status != StatusFail {
		t.Errorf("expected fail, got %s", result.Status)
	}
}

func TestCustomPolicy_BranchName_Pass(t *testing.T) {
	p := NewCustomPolicy(config.CustomRuleDef{
		Name: "no-draft", Type: "branch-name", Pattern: "draft-*", Message: "no drafts",
	})
	result := p.Validate(Context{BranchName: "feature/foo"})
	if result.Status != StatusPass {
		t.Errorf("expected pass, got %s", result.Status)
	}
}

func TestCustomPolicy_BranchName_Fail(t *testing.T) {
	p := NewCustomPolicy(config.CustomRuleDef{
		Name: "no-draft", Type: "branch-name", Pattern: "draft-*", Message: "no drafts",
	})
	result := p.Validate(Context{BranchName: "draft-experiment"})
	if result.Status != StatusFail {
		t.Errorf("expected fail, got %s", result.Status)
	}
}

func TestCustomPolicy_CommitMessage_Pass(t *testing.T) {
	p := NewCustomPolicy(config.CustomRuleDef{
		Name: "no-wip", Type: "commit-message", Pattern: "wip:*", Message: "no wip",
	})
	result := p.Validate(Context{CommitMsg: "feat: add feature"})
	if result.Status != StatusPass {
		t.Errorf("expected pass, got %s", result.Status)
	}
}

func TestCustomPolicy_CommitMessage_Fail(t *testing.T) {
	p := NewCustomPolicy(config.CustomRuleDef{
		Name: "no-wip", Type: "commit-message", Pattern: "WIP:*", Message: "no wip",
	})
	result := p.Validate(Context{CommitMsg: "WIP: stuff"})
	if result.Status != StatusFail {
		t.Errorf("expected fail, got %s", result.Status)
	}
}

func TestCustomPolicy_UnknownType(t *testing.T) {
	p := NewCustomPolicy(config.CustomRuleDef{
		Name: "weird", Type: "invalid-type", Pattern: "x", Message: "unknown",
	})
	result := p.Validate(Context{})
	if result.Status != StatusPass {
		t.Errorf("expected pass for unknown type, got %s", result.Status)
	}
}

func TestCustomPolicy_Name(t *testing.T) {
	p := NewCustomPolicy(config.CustomRuleDef{Name: "my-rule"})
	if p.Name() != "Custom:my-rule" {
		t.Errorf("expected 'Custom:my-rule', got %q", p.Name())
	}
}
