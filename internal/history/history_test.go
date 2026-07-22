package history

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/marcuwynu23/git-policy/internal/config"
	"github.com/marcuwynu23/git-policy/internal/policy"
)

func TestLogAndQuery(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer func() { _ = os.Chdir(origDir) }()

	initTestRepo(t, tmpDir)
	_ = os.Chdir(tmpDir)

	cfgPath := filepath.Join(tmpDir, ".config", "git-policy", "config.yaml")
	cfg := config.DefaultConfig()
	_ = os.MkdirAll(filepath.Dir(cfgPath), 0755)
	_ = config.Save(cfg, cfgPath)

	testResults := []policy.Result{
		{
			PolicyName: "BlockFiles",
			Status:     policy.StatusPass,
		},
		{
			PolicyName: "SecretScan",
			Status:     policy.StatusFail,
			Message:    "found secret in file",
		},
	}
	if err := Log(cfg, cfgPath, testResults); err != nil {
		t.Fatalf("Log failed: %v", err)
	}

	records, err := Query(cfg, cfgPath, QueryOptions{})
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}
	if len(records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(records))
	}
	if records[0].Overall != "fail" {
		t.Errorf("expected overall status 'fail', got %q", records[0].Overall)
	}
	if len(records[0].Results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(records[0].Results))
	}
}

func TestQueryWithLimit(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer func() { _ = os.Chdir(origDir) }()

	initTestRepo(t, tmpDir)
	_ = os.Chdir(tmpDir)

	cfgPath := filepath.Join(tmpDir, ".config", "git-policy", "config.yaml")
	cfg := config.DefaultConfig()
	_ = os.MkdirAll(filepath.Dir(cfgPath), 0755)
	_ = config.Save(cfg, cfgPath)

	testResults := []policy.Result{{PolicyName: "TestPolicy", Status: policy.StatusPass}}
	for i := 0; i < 5; i++ {
		_ = Log(cfg, cfgPath, testResults)
		time.Sleep(10 * time.Millisecond)
	}

	records, err := Query(cfg, cfgPath, QueryOptions{Limit: 2})
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}
	if len(records) != 2 {
		t.Fatalf("expected 2 records with limit 2, got %d", len(records))
	}
}

func TestQueryWithStatus(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer func() { _ = os.Chdir(origDir) }()

	initTestRepo(t, tmpDir)
	_ = os.Chdir(tmpDir)

	cfgPath := filepath.Join(tmpDir, ".config", "git-policy", "config.yaml")
	cfg := config.DefaultConfig()
	_ = os.MkdirAll(filepath.Dir(cfgPath), 0755)
	_ = config.Save(cfg, cfgPath)

	_ = Log(cfg, cfgPath, []policy.Result{{PolicyName: "Test", Status: policy.StatusPass}})
	_ = Log(cfg, cfgPath, []policy.Result{{PolicyName: "Test", Status: policy.StatusFail}})
	_ = Log(cfg, cfgPath, []policy.Result{{PolicyName: "Test", Status: policy.StatusPass}})

	records, err := Query(cfg, cfgPath, QueryOptions{Status: "fail"})
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}
	if len(records) != 1 {
		t.Fatalf("expected 1 fail record, got %d", len(records))
	}
}

func TestClear(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer func() { _ = os.Chdir(origDir) }()

	initTestRepo(t, tmpDir)
	_ = os.Chdir(tmpDir)

	cfgPath := filepath.Join(tmpDir, ".config", "git-policy", "config.yaml")
	cfg := config.DefaultConfig()
	_ = os.MkdirAll(filepath.Dir(cfgPath), 0755)
	_ = config.Save(cfg, cfgPath)

	testResults := []policy.Result{{PolicyName: "TestPolicy", Status: policy.StatusPass}}
	if err := Log(cfg, cfgPath, testResults); err != nil {
		t.Fatalf("Log failed: %v", err)
	}

	if err := Clear(cfg, cfgPath, ""); err != nil {
		t.Fatalf("Clear failed: %v", err)
	}

	records, err := Query(cfg, cfgPath, QueryOptions{})
	if err != nil {
		t.Fatalf("Query after clear failed: %v", err)
	}
	if len(records) != 0 {
		t.Fatalf("expected 0 records after clear, got %d", len(records))
	}
}

func initTestRepo(t *testing.T, dir string) {
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
