package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

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

func HookDir() (string, error) {
	cfgDir, err := ConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(cfgDir, "hooks"), nil
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func EnsureDir(path string) error {
	return os.MkdirAll(path, 0755)
}
