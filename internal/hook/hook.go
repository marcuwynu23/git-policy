package hook

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/marcuwynu23/git-policy/internal/config"
)

type Installer struct {
	cfg *config.Config
}

func NewInstaller() *Installer {
	return &Installer{}
}

func (i *Installer) globalHookDir() (string, error) {
	var base string
	switch runtime.GOOS {
	case "windows":
		base = os.Getenv("APPDATA")
		if base == "" {
			base = filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Roaming")
		}
	default:
		base = filepath.Join(os.Getenv("HOME"), ".config")
	}
	if base == "" {
		return "", fmt.Errorf("cannot determine home directory")
	}
	return filepath.Join(base, "git-policy", "hooks"), nil
}

func (i *Installer) InstallGlobal() error {
	hookDir, err := i.globalHookDir()
	if err != nil {
		return fmt.Errorf("determining hook directory: %w", err)
	}
	if err := os.MkdirAll(hookDir, 0755); err != nil {
		return fmt.Errorf("creating hook directory: %w", err)
	}

	hooks := []string{"pre-commit", "pre-push", "commit-msg", "post-merge"}
	for _, hook := range hooks {
		hookPath := filepath.Join(hookDir, hook)
		content := i.hookScript(hook)
		if err := os.WriteFile(hookPath, []byte(content), 0755); err != nil {
			return fmt.Errorf("writing hook %s: %w", hook, err)
		}
	}
	return i.setGlobalHooksPath(hookDir)
}

func (i *Installer) setGlobalHooksPath(hookDir string) error {
	cmd := exec.Command("git", "config", "--global", "core.hooksPath", hookDir)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("setting git hooksPath: %s: %w", string(output), err)
	}
	return nil
}

func (i *Installer) globalConfigDir() (string, error) {
	var base string
	switch runtime.GOOS {
	case "windows":
		base = os.Getenv("APPDATA")
		if base == "" {
			base = filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Roaming")
		}
	default:
		base = filepath.Join(os.Getenv("HOME"), ".config")
	}
	if base == "" {
		return "", fmt.Errorf("cannot determine home directory")
	}
	return filepath.Join(base, "git-policy"), nil
}

func (i *Installer) UninstallGlobal() error {
	hookDir, err := i.globalHookDir()
	if err != nil {
		return fmt.Errorf("determining hook directory: %w", err)
	}

	hooks := []string{"pre-commit", "pre-push", "commit-msg", "post-merge"}
	for _, hook := range hooks {
		hookPath := filepath.Join(hookDir, hook)
		if err := os.Remove(hookPath); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("removing hook %s: %w", hook, err)
		}
	}

	if err := os.Remove(hookDir); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("removing hook directory: %w", err)
	}

	cfgDir, err := i.globalConfigDir()
	if err == nil {
		entries, _ := os.ReadDir(cfgDir)
		if len(entries) == 0 {
			os.Remove(cfgDir)
		}
	}

	cmd := exec.Command("git", "config", "--global", "--unset", "core.hooksPath")
	cmd.Run()
	return nil
}

func (i *Installer) UninstallAll() error {
	if err := i.UninstallGlobal(); err != nil {
		return err
	}

	cfgDir, err := i.globalConfigDir()
	if err != nil {
		return fmt.Errorf("determining config directory: %w", err)
	}

	configPath := filepath.Join(cfgDir, "config.yaml")
	if err := os.Remove(configPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("removing config: %w", err)
	}

	os.Remove(cfgDir)
	return nil
}

func (i *Installer) IsInstalled() bool {
	cmd := exec.Command("git", "config", "--global", "core.hooksPath")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	hookDir := strings.TrimSpace(string(output))
	if hookDir == "" {
		return false
	}
	info, err := os.Stat(filepath.Join(hookDir, "pre-commit"))
	if err != nil {
		return false
	}
	return !info.IsDir()
}

func (i *Installer) hookScript(hookName string) string {
	return fmt.Sprintf(`#!/bin/sh
# git-policy hook: %s
GIT_POLICY=$(command -v git-policy 2>/dev/null || command -v git-policy.exe 2>/dev/null)
if [ -z "$GIT_POLICY" ]; then
	echo "error: git-policy not found in PATH" >&2
	echo "Run 'make dev' from the git-policy project directory to set up." >&2
	echo "Or add git-policy to your PATH manually." >&2
	exit 1
fi
exec "$GIT_POLICY" run
`, hookName)
}
