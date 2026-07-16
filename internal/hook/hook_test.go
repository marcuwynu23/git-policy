package hook

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestNewInstaller(t *testing.T) {
	inst := NewInstaller()
	if inst == nil {
		t.Fatal("expected non-nil installer")
	}
}

func TestInstallGlobalAndUninstall(t *testing.T) {
	origAppData := os.Getenv("APPDATA")
	origUserProfile := os.Getenv("USERPROFILE")
	origHome := os.Getenv("HOME")
	defer func() {
		os.Setenv("APPDATA", origAppData)
		os.Setenv("USERPROFILE", origUserProfile)
		os.Setenv("HOME", origHome)
	}()

	testDir := t.TempDir()
	if runtime.GOOS == "windows" {
		os.Setenv("APPDATA", testDir)
	} else {
		os.Setenv("HOME", testDir)
	}

	inst := NewInstaller()

	err := inst.InstallGlobal()
	if err != nil {
		t.Fatalf("InstallGlobal() failed: %v", err)
	}

	hookDir, _ := inst.globalHookDir()
	hooks := []string{"pre-commit", "pre-push", "commit-msg", "post-merge"}
	for _, hook := range hooks {
		path := filepath.Join(hookDir, hook)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("expected hook file %s to exist", path)
		}
	}

	err = inst.UninstallGlobal()
	if err != nil {
		t.Fatalf("UninstallGlobal() failed: %v", err)
	}

	for _, hook := range hooks {
		path := filepath.Join(hookDir, hook)
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			t.Errorf("expected hook file %s to be removed", path)
		}
	}
}

func TestGlobalHookDir(t *testing.T) {
	origAppData := os.Getenv("APPDATA")
	origUserProfile := os.Getenv("USERPROFILE")
	origHome := os.Getenv("HOME")
	defer func() {
		os.Setenv("APPDATA", origAppData)
		os.Setenv("USERPROFILE", origUserProfile)
		os.Setenv("HOME", origHome)
	}()

	testDir := t.TempDir()
	if runtime.GOOS == "windows" {
		os.Setenv("APPDATA", testDir)
	} else {
		os.Setenv("HOME", testDir)
	}

	inst := NewInstaller()
	dir, err := inst.globalHookDir()
	if err != nil {
		t.Fatalf("globalHookDir() failed: %v", err)
	}
	if dir == "" {
		t.Error("expected non-empty hook directory")
	}
	if !strings.HasSuffix(filepath.ToSlash(dir), "git-policy/hooks") {
		t.Errorf("expected path to end with 'git-policy/hooks', got %q", dir)
	}
}

func TestGlobalConfigDir(t *testing.T) {
	origAppData := os.Getenv("APPDATA")
	origUserProfile := os.Getenv("USERPROFILE")
	origHome := os.Getenv("HOME")
	defer func() {
		os.Setenv("APPDATA", origAppData)
		os.Setenv("USERPROFILE", origUserProfile)
		os.Setenv("HOME", origHome)
	}()

	testDir := t.TempDir()
	if runtime.GOOS == "windows" {
		os.Setenv("APPDATA", testDir)
	} else {
		os.Setenv("HOME", testDir)
	}

	inst := NewInstaller()
	dir, err := inst.globalConfigDir()
	if err != nil {
		t.Fatalf("globalConfigDir() failed: %v", err)
	}
	if dir == "" {
		t.Error("expected non-empty config directory")
	}
	if !strings.HasSuffix(dir, "git-policy") {
		t.Errorf("expected path to end with 'git-policy', got %q", dir)
	}
}

func TestHookScript(t *testing.T) {
	inst := NewInstaller()
	script := inst.hookScript("pre-commit")
	if !strings.Contains(script, "pre-commit") {
		t.Error("expected script to contain hook name")
	}
	if !strings.Contains(script, "git-policy") {
		t.Error("expected script to reference git-policy")
	}
	if !strings.HasPrefix(script, "#!/bin/sh") {
		t.Error("expected script to start with shebang")
	}
}

func TestUninstallAll(t *testing.T) {
	origAppData := os.Getenv("APPDATA")
	origUserProfile := os.Getenv("USERPROFILE")
	origHome := os.Getenv("HOME")
	defer func() {
		os.Setenv("APPDATA", origAppData)
		os.Setenv("USERPROFILE", origUserProfile)
		os.Setenv("HOME", origHome)
	}()

	testDir := t.TempDir()
	if runtime.GOOS == "windows" {
		os.Setenv("APPDATA", testDir)
	} else {
		os.Setenv("HOME", testDir)
	}

	inst := NewInstaller()
	inst.InstallGlobal()

	err := inst.UninstallAll()
	if err != nil {
		t.Fatalf("UninstallAll() failed: %v", err)
	}

	hookDir, _ := inst.globalHookDir()
	if _, err := os.Stat(hookDir); !os.IsNotExist(err) {
		t.Errorf("expected hook directory to be removed")
	}
}

func TestIsInstalled(t *testing.T) {
	inst := NewInstaller()
	result := inst.IsInstalled()
	_ = result
}
