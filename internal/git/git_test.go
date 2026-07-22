package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestVersion(t *testing.T) {
	version, err := Version()
	if err != nil {
		t.Fatalf("Version() failed: %v", err)
	}
	if version == "" {
		t.Error("expected non-empty version string")
	}
}

func TestIsRepo(t *testing.T) {
	if !IsRepo() {
		t.Error("expected current directory to be a git repo")
	}
}

func TestIsRepo_NotARepo(t *testing.T) {
	dir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(dir)
	defer func() { _ = os.Chdir(origDir) }()

	if IsRepo() {
		t.Error("expected temp dir not to be a git repo")
	}
}

func TestGetBranchName(t *testing.T) {
	dir := t.TempDir()
	origDir, _ := os.Getwd()
	defer func() { _ = os.Chdir(origDir) }()

	initRepo(t, dir)
	_ = os.Chdir(dir)

	branch, err := GetBranchName()
	if err != nil {
		t.Fatalf("GetBranchName() failed: %v", err)
	}
	if branch != "main" && branch != "master" {
		t.Errorf("expected 'main' or 'master', got %q", branch)
	}
}

func TestGetStagedFiles_Empty(t *testing.T) {
	dir := t.TempDir()
	origDir, _ := os.Getwd()
	defer func() { _ = os.Chdir(origDir) }()

	initRepo(t, dir)
	_ = os.Chdir(dir)

	files, err := GetStagedFiles()
	if err != nil {
		t.Fatalf("GetStagedFiles() failed: %v", err)
	}
	if len(files) != 0 {
		t.Errorf("expected no staged files, got %v", files)
	}
}

func TestGetStagedFiles_WithFiles(t *testing.T) {
	dir := t.TempDir()
	origDir, _ := os.Getwd()
	defer func() { _ = os.Chdir(origDir) }()

	initRepo(t, dir)
	_ = os.Chdir(dir)

	newFile := filepath.Join(dir, "new.txt")
	_ = os.WriteFile(newFile, []byte("hello"), 0644)
	_ = exec.Command("git", "add", "new.txt").Run()

	files, err := GetStagedFiles()
	if err != nil {
		t.Fatalf("GetStagedFiles() failed: %v", err)
	}
	if len(files) == 0 {
		t.Fatal("expected staged files, got none")
	}
	found := false
	for _, f := range files {
		if f == "new.txt" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected 'new.txt' in staged files, got %v", files)
	}
}

func TestGetCommitMsgFile(t *testing.T) {
	dir := t.TempDir()
	origDir, _ := os.Getwd()
	defer func() { _ = os.Chdir(origDir) }()

	initRepo(t, dir)
	_ = os.Chdir(dir)

	path, err := GetCommitMsgFile()
	if err != nil {
		t.Fatalf("GetCommitMsgFile() failed: %v", err)
	}
	if path == "" {
		t.Error("expected non-empty commit message file path")
	}
	if !strings.Contains(path, "COMMIT_EDITMSG") {
		t.Errorf("expected path to contain COMMIT_EDITMSG, got %q", path)
	}
}

func TestGetCommitMessage(t *testing.T) {
	dir := t.TempDir()
	origDir, _ := os.Getwd()
	defer func() { _ = os.Chdir(origDir) }()

	initRepo(t, dir)
	_ = os.Chdir(dir)

	msg, err := GetCommitMessage()
	if err != nil {
		t.Fatalf("GetCommitMessage() failed: %v", err)
	}
	if msg == "" {
		t.Error("expected non-empty commit message")
	}
}

func TestGetConfig_NotFound(t *testing.T) {
	dir := t.TempDir()
	origDir, _ := os.Getwd()
	defer func() { _ = os.Chdir(origDir) }()

	initRepo(t, dir)
	_ = os.Chdir(dir)

	val, err := GetConfig("git-policy.nonexistent")
	if err == nil {
		t.Errorf("expected error for missing key, got value %q", val)
	}
}

func TestSetConfigAndGetConfig(t *testing.T) {
	dir := t.TempDir()
	origDir, _ := os.Getwd()
	defer func() { _ = os.Chdir(origDir) }()

	initRepo(t, dir)
	_ = os.Chdir(dir)

	err := SetConfig("git-policy.skip", "block-files,secret-scan")
	if err != nil {
		t.Fatalf("SetConfig failed: %v", err)
	}

	val, err := GetConfig("git-policy.skip")
	if err != nil {
		t.Fatalf("GetConfig failed: %v", err)
	}
	if val != "block-files,secret-scan" {
		t.Errorf("expected 'block-files,secret-scan', got %q", val)
	}
}

func TestSetConfig_ReplacesExisting(t *testing.T) {
	dir := t.TempDir()
	origDir, _ := os.Getwd()
	defer func() { _ = os.Chdir(origDir) }()

	initRepo(t, dir)
	_ = os.Chdir(dir)

	_ = SetConfig("git-policy.skip", "old-value")
	_ = SetConfig("git-policy.skip", "new-value")

	val, err := GetConfig("git-policy.skip")
	if err != nil {
		t.Fatalf("GetConfig failed: %v", err)
	}
	if val != "new-value" {
		t.Errorf("expected 'new-value', got %q", val)
	}
}

func TestUnsetConfig(t *testing.T) {
	dir := t.TempDir()
	origDir, _ := os.Getwd()
	defer func() { _ = os.Chdir(origDir) }()

	initRepo(t, dir)
	_ = os.Chdir(dir)

	_ = SetConfig("git-policy.skip", "block-files")
	_ = UnsetConfig("git-policy.skip")

	_, err := GetConfig("git-policy.skip")
	if err == nil {
		t.Error("expected error after unset")
	}
}

func TestUnsetConfig_NoErrorIfMissing(t *testing.T) {
	dir := t.TempDir()
	origDir, _ := os.Getwd()
	defer func() { _ = os.Chdir(origDir) }()

	initRepo(t, dir)
	_ = os.Chdir(dir)

	err := UnsetConfig("git-policy.nonexistent")
	if err != nil {
		t.Errorf("expected no error for missing key, got %v", err)
	}
}

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
