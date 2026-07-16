// Package config handles loading, saving, and managing git-policy configuration.
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"gopkg.in/yaml.v3"
)

// Config represents the complete git-policy configuration.
type Config struct {
	Version  int           `yaml:"version"`
	Hooks    HooksConfig   `yaml:"hooks"`
	Policies PoliciesConfig `yaml:"policies"`
}

// HooksConfig controls which Git hooks are enabled.
type HooksConfig struct {
	PreCommit  HookConfig `yaml:"pre-commit"`
	CommitMsg  HookConfig `yaml:"commit-msg"`
	PrePush    HookConfig `yaml:"pre-push"`
	PostMerge  HookConfig `yaml:"post-merge"`
}

// HookConfig controls a single Git hook's enabled state.
type HookConfig struct {
	Enabled bool `yaml:"enabled"`
}

// PoliciesConfig holds all configurable policy settings.
type PoliciesConfig struct {
	BlockFiles          []string `yaml:"blockFiles"`
	MaxFileSize         string   `yaml:"maxFileSize"`
	SecretScan          bool     `yaml:"secretScan"`
	ProtectedBranches   []string `yaml:"protectedBranches"`
	ConventionalCommits bool     `yaml:"conventionalCommits"`
	BlockBinaries       []string `yaml:"blockBinaries"`
	RequiredFiles       []string `yaml:"requiredFiles"`
	DisabledPolicies    []string `yaml:"disabledPolicies"`
}

// PolicyNames maps CLI-friendly policy names to their internal names.
var PolicyNames = map[string]string{
	"block-files":        "BlockFiles",
	"commit-message":     "CommitMessage",
	"file-size":          "FileSize",
	"binary-file":        "BinaryFile",
	"secret-scan":        "SecretScan",
	"branch-protection":  "BranchProtection",
}

// PolicyCLIName returns the CLI-friendly name for a given internal policy name.
func PolicyCLIName(internalName string) string {
	for cli, internal := range PolicyNames {
		if internal == internalName {
			return cli
		}
	}
	return internalName
}

// IsDisabled checks if a named policy is in the disabled list.
func (p *PoliciesConfig) IsDisabled(name string) bool {
	for _, d := range p.DisabledPolicies {
		if d == name {
			return true
		}
	}
	return false
}

// SetDisabled adds or removes a policy from the disabled list.
func (p *PoliciesConfig) SetDisabled(name string, disabled bool) {
	if disabled {
		for _, d := range p.DisabledPolicies {
			if d == name {
				return
			}
		}
		p.DisabledPolicies = append(p.DisabledPolicies, name)
	} else {
		var updated []string
		for _, d := range p.DisabledPolicies {
			if d != name {
				updated = append(updated, d)
			}
		}
		p.DisabledPolicies = updated
	}
}

// DefaultConfig returns a configuration with sensible default values.
func DefaultConfig() *Config {
	return &Config{
		Version: 1,
		Hooks: HooksConfig{
			PreCommit:  HookConfig{Enabled: true},
			CommitMsg:  HookConfig{Enabled: true},
			PrePush:    HookConfig{Enabled: true},
			PostMerge:  HookConfig{Enabled: false},
		},
		Policies: PoliciesConfig{
			BlockFiles:          []string{".env", "*.pem", "*.key", "*.p12", "*.crt"},
			MaxFileSize:         "10MB",
			SecretScan:          true,
			ProtectedBranches:   []string{"main", "master", "production"},
			ConventionalCommits: true,
			BlockBinaries:       []string{".exe", ".dll", ".so", ".iso", ".zip"},
			RequiredFiles:       []string{"README.md", "LICENSE"},
		},
	}
}

// DefaultConfigPath returns the default OS-specific path for the config file.
func DefaultConfigPath() (string, error) {
	var dir string
	switch runtime.GOOS {
	case "windows":
		dir = os.Getenv("APPDATA")
		if dir == "" {
			dir = filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Roaming")
		}
	default:
		dir = filepath.Join(os.Getenv("HOME"), ".config")
	}
	if dir == "" {
		return "", fmt.Errorf("cannot determine config directory")
	}
	return filepath.Join(dir, "git-policy", "config.yaml"), nil
}

// Load reads and parses a YAML config file, falling back to defaults if needed.
func Load(path string) (*Config, error) {
	if path == "" {
		var err error
		path, err = DefaultConfigPath()
		if err != nil {
			return DefaultConfig(), nil
		}
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return DefaultConfig(), nil
		}
		return nil, fmt.Errorf("reading config: %w", err)
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}
	return &cfg, nil
}

// Save writes a config to the specified file path as YAML.
func Save(cfg *Config, path string) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("marshaling config: %w", err)
	}
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("creating config directory: %w", err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("writing config: %w", err)
	}
	return nil
}
