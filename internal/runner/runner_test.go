package runner

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/marcuwynu23/git-policy/internal/config"
	"github.com/marcuwynu23/git-policy/internal/engine"
	"github.com/marcuwynu23/git-policy/internal/git"
	"github.com/marcuwynu23/git-policy/internal/policy"
)

func TestRunner_EngineWiring(t *testing.T) {
	eng := engine.New(config.DefaultConfig())

	passPolicy := &mockPolicy{
		name:   "PassPolicy",
		result: policy.Result{PolicyName: "PassPolicy", Status: policy.StatusPass},
	}
	failPolicy := &mockPolicy{
		name:   "FailPolicy",
		result: policy.Result{PolicyName: "FailPolicy", Status: policy.StatusFail, Message: "fail"},
	}

	eng.Register(passPolicy)
	eng.Register(failPolicy)

	results := eng.Execute()
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if results[0].Status != policy.StatusPass {
		t.Errorf("expected pass, got %s", results[0].Status)
	}
	if results[1].Status != policy.StatusFail {
		t.Errorf("expected fail, got %s", results[1].Status)
	}
}

func TestRunner_DisabledPoliciesAreSkipped(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Policies.SetDisabled("FailPolicy", true)

	eng := engine.New(cfg)

	failPolicy := &mockPolicy{
		name:   "FailPolicy",
		result: policy.Result{PolicyName: "FailPolicy", Status: policy.StatusFail, Message: "should be skipped"},
	}
	passPolicy := &mockPolicy{
		name:   "PassPolicy",
		result: policy.Result{PolicyName: "PassPolicy", Status: policy.StatusPass},
	}

	eng.Register(passPolicy)
	eng.Register(failPolicy)

	results := eng.Execute()
	if len(results) != 1 {
		t.Fatalf("expected 1 result (skipped disabled), got %d", len(results))
	}
	if results[0].PolicyName != "PassPolicy" {
		t.Errorf("expected PassPolicy result, got %s", results[0].PolicyName)
	}
}

func TestReadSkipList_Empty(t *testing.T) {
	dir := t.TempDir()
	origDir, _ := os.Getwd()
	defer func() { _ = os.Chdir(origDir) }()

	initRepo(t, dir)
	_ = os.Chdir(dir)

	names := readSkipList()
	if len(names) != 0 {
		t.Errorf("expected empty skip list, got %v", names)
	}
}

func TestReadSkipList_WithValues(t *testing.T) {
	dir := t.TempDir()
	origDir, _ := os.Getwd()
	defer func() { _ = os.Chdir(origDir) }()

	initRepo(t, dir)
	_ = os.Chdir(dir)

	_ = git.SetConfig("git-policy.skip", "block-files,secret-scan")

	names := readSkipList()
	if len(names) != 2 {
		t.Fatalf("expected 2 names, got %v", names)
	}
	if names[0] != "BlockFiles" {
		t.Errorf("expected BlockFiles, got %s", names[0])
	}
	if names[1] != "SecretScan" {
		t.Errorf("expected SecretScan, got %s", names[1])
	}
}

func TestReadSkipList_UnknownName(t *testing.T) {
	dir := t.TempDir()
	origDir, _ := os.Getwd()
	defer func() { _ = os.Chdir(origDir) }()

	initRepo(t, dir)
	_ = os.Chdir(dir)

	_ = git.SetConfig("git-policy.skip", "unknown-rule")

	names := readSkipList()
	if len(names) != 0 {
		t.Errorf("expected empty skip list for unknown names, got %v", names)
	}
}

type mockPolicy struct {
	name   string
	result policy.Result
}

func (m *mockPolicy) Name() string { return m.name }
func (m *mockPolicy) Validate(ctx policy.Context) policy.Result { return m.result }

func initRepo(t *testing.T, dir string) {
	t.Helper()
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
