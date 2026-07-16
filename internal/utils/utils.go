// Package utils provides shared utility functions used across the project.
package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// ConfigDir returns the OS-specific path to the git-policy config directory.
func ConfigDir() (string, error) {
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
		return "", fmt.Errorf("cannot determine config directory")
	}
	return filepath.Join(base, "git-policy"), nil
}

// HookDir returns the OS-specific path to the git-policy hooks directory.
func HookDir() (string, error) {
	cfgDir, err := ConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(cfgDir, "hooks"), nil
}

// FileExists checks whether a file exists at the given path.
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// EnsureDir creates a directory and all parents if they don't already exist.
func EnsureDir(path string) error {
	return os.MkdirAll(path, 0755)
}
