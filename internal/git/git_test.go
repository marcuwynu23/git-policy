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
	os.Chdir(dir)
	defer os.Chdir(origDir)

	if IsRepo() {
		t.Error("expected temp dir not to be a git repo")
	}
}

func TestGetBranchName(t *testing.T) {
	dir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)

	initRepo(t, dir)
	os.Chdir(dir)

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
	defer os.Chdir(origDir)

	initRepo(t, dir)
	os.Chdir(dir)

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
	defer os.Chdir(origDir)

	initRepo(t, dir)
	os.Chdir(dir)

	newFile := filepath.Join(dir, "new.txt")
	os.WriteFile(newFile, []byte("hello"), 0644)
	exec.Command("git", "add", "new.txt").Run()

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
	defer os.Chdir(origDir)

	initRepo(t, dir)
	os.Chdir(dir)

	path, err := GetCommitMsgFile()
	if err != nil {
		t.Fatalf("GetCommitMsgFile() failed: %v", err)
	}
	if path == "" {
		t.Error("expected non-empty commit message file path")
	}
	absPath, _ := filepath.Abs(path)
	absDir, _ := filepath.Abs(dir)
	if !strings.HasPrefix(absPath, absDir) {
		t.Errorf("expected path inside repo, got %q (abs: %q)", path, absPath)
	}
}

func TestGetCommitMessage(t *testing.T) {
	dir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)

	initRepo(t, dir)
	os.Chdir(dir)

	msg, err := GetCommitMessage()
	if err != nil {
		t.Fatalf("GetCommitMessage() failed: %v", err)
	}
	if msg == "" {
		t.Error("expected non-empty commit message")
	}
}

func initRepo(t *testing.T, dir string) {
	t.Helper()
	cmd := exec.Command("git", "init")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git init failed: %s: %v", string(out), err)
	}

	exec.Command("git", "-C", dir, "config", "user.email", "test@test.com").Run()
	exec.Command("git", "-C", dir, "config", "user.name", "Test").Run()
	exec.Command("git", "-C", dir, "config", "commit.gpgSign", "false").Run()

	readme := filepath.Join(dir, "README.md")
	os.WriteFile(readme, []byte("# test"), 0644)
	exec.Command("git", "-C", dir, "add", "README.md").Run()
	exec.Command("git", "-C", dir, "commit", "-m", "feat: initial commit").Run()
}
