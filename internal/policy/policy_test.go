package policy

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/marcuwynu23/git-policy/internal/config"
)

func TestBlockFilesPolicy_Pass(t *testing.T) {
	cfg := config.DefaultConfig()
	p := NewBlockFilesPolicy(cfg)
	result := p.Validate(Context{
		StagedFiles: []string{"main.go", "go.mod"},
	})
	if result.Status != StatusPass {
		t.Errorf("expected pass, got %s: %s", result.Status, result.Message)
	}
}

func TestBlockFilesPolicy_Fail(t *testing.T) {
	cfg := config.DefaultConfig()
	p := NewBlockFilesPolicy(cfg)
	result := p.Validate(Context{
		StagedFiles: []string{".env"},
	})
	if result.Status != StatusFail {
		t.Errorf("expected fail, got %s", result.Status)
	}
}

func TestBlockFilesPolicy_GlobPattern(t *testing.T) {
	cfg := config.DefaultConfig()
	p := NewBlockFilesPolicy(cfg)
	result := p.Validate(Context{
		StagedFiles: []string{"secret.key", "cert.pem"},
	})
	if result.Status != StatusFail {
		t.Errorf("expected fail for glob pattern, got %s", result.Status)
	}
}

func TestCommitMessagePolicy_Pass(t *testing.T) {
	cfg := config.DefaultConfig()
	p := NewCommitMessagePolicy(cfg)
	tests := []string{
		"feat: add new feature",
		"fix: resolve bug",
		"docs: update readme",
		"refactor: clean up code",
		"test: add tests",
	}
	for _, msg := range tests {
		result := p.Validate(Context{CommitMsg: msg})
		if result.Status != StatusPass {
			t.Errorf("expected pass for %q, got %s", msg, result.Status)
		}
	}
}

func TestCommitMessagePolicy_Fail(t *testing.T) {
	cfg := config.DefaultConfig()
	p := NewCommitMessagePolicy(cfg)
	result := p.Validate(Context{CommitMsg: "added stuff"})
	if result.Status != StatusFail {
		t.Errorf("expected fail, got %s", result.Status)
	}
}

func TestCommitMessagePolicy_Disabled(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Policies.ConventionalCommits = false
	p := NewCommitMessagePolicy(cfg)
	result := p.Validate(Context{CommitMsg: "random message"})
	if result.Status != StatusPass {
		t.Errorf("expected pass when disabled, got %s", result.Status)
	}
}

func TestFileSizePolicy_Pass(t *testing.T) {
	cfg := config.DefaultConfig()
	p := NewFileSizePolicy(cfg)
	dir := t.TempDir()
	smallFile := filepath.Join(dir, "small.txt")
	_ = os.WriteFile(smallFile, []byte("hello"), 0644)
	result := p.Validate(Context{StagedFiles: []string{smallFile}})
	if result.Status != StatusPass {
		t.Errorf("expected pass, got %s: %s", result.Status, result.Message)
	}
}

func TestBinaryFilePolicy_Pass(t *testing.T) {
	cfg := config.DefaultConfig()
	p := NewBinaryFilePolicy(cfg)
	result := p.Validate(Context{StagedFiles: []string{"main.go", "file.txt"}})
	if result.Status != StatusPass {
		t.Errorf("expected pass, got %s", result.Status)
	}
}

func TestBinaryFilePolicy_Fail(t *testing.T) {
	cfg := config.DefaultConfig()
	p := NewBinaryFilePolicy(cfg)
	result := p.Validate(Context{StagedFiles: []string{"malware.exe"}})
	if result.Status != StatusFail {
		t.Errorf("expected fail, got %s", result.Status)
	}
}

func TestSecretScanPolicy_Pass(t *testing.T) {
	cfg := config.DefaultConfig()
	p := NewSecretScanPolicy(cfg)
	dir := t.TempDir()
	safeFile := filepath.Join(dir, "safe.txt")
	_ = os.WriteFile(safeFile, []byte("hello world"), 0644)
	result := p.Validate(Context{StagedFiles: []string{safeFile}})
	if result.Status != StatusPass {
		t.Errorf("expected pass, got %s: %s", result.Status, result.Message)
	}
}

func TestSecretScanPolicy_Fail(t *testing.T) {
	cfg := config.DefaultConfig()
	p := NewSecretScanPolicy(cfg)
	dir := t.TempDir()
	secretFile := filepath.Join(dir, "secret.txt")
	os.WriteFile(secretFile, []byte("AKIA1234567890"), 0644)
	result := p.Validate(Context{StagedFiles: []string{secretFile}})
	if result.Status != StatusFail {
		t.Errorf("expected fail, got %s", result.Status)
	}
}

func TestBranchPolicy_Pass(t *testing.T) {
	cfg := config.DefaultConfig()
	p := NewBranchPolicy(cfg)
	result := p.Validate(Context{BranchName: "feature/my-feature"})
	if result.Status != StatusPass {
		t.Errorf("expected pass, got %s", result.Status)
	}
}

func TestBranchPolicy_Fail(t *testing.T) {
	cfg := config.DefaultConfig()
	p := NewBranchPolicy(cfg)
	result := p.Validate(Context{BranchName: "main"})
	if result.Status != StatusFail {
		t.Errorf("expected fail, got %s", result.Status)
	}
}

func TestParseMaxSize(t *testing.T) {
	tests := []struct {
		input string
		want  int64
	}{
		{"10MB", 10 * 1024 * 1024},
		{"1GB", 1024 * 1024 * 1024},
		{"100KB", 100 * 1024},
		{"500B", 500},
		{"", 10 * 1024 * 1024},
		{"invalid", 10 * 1024 * 1024},
	}
	for _, tt := range tests {
		got, err := parseMaxSize(tt.input)
		if err != nil {
			t.Errorf("parseMaxSize(%q) error: %v", tt.input, err)
		}
		if got != tt.want {
			t.Errorf("parseMaxSize(%q) = %d, want %d", tt.input, got, tt.want)
		}
	}
}
