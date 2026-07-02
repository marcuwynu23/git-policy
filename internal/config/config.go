package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Version  int      `yaml:"version"`
	Hooks    HooksConfig    `yaml:"hooks"`
	Policies PoliciesConfig `yaml:"policies"`
}

type HooksConfig struct {
	PreCommit  HookConfig `yaml:"pre-commit"`
	CommitMsg  HookConfig `yaml:"commit-msg"`
	PrePush    HookConfig `yaml:"pre-push"`
	PostMerge  HookConfig `yaml:"post-merge"`
}

type HookConfig struct {
	Enabled bool `yaml:"enabled"`
}

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

var PolicyNames = map[string]string{
	"block-files":        "BlockFiles",
	"commit-message":     "CommitMessage",
	"file-size":          "FileSize",
	"binary-file":        "BinaryFile",
	"secret-scan":        "SecretScan",
	"branch-protection":  "BranchProtection",
}

func PolicyCLIName(internalName string) string {
	for cli, internal := range PolicyNames {
		if internal == internalName {
			return cli
		}
	}
	return internalName
}

func (p *PoliciesConfig) IsDisabled(name string) bool {
	for _, d := range p.DisabledPolicies {
		if d == name {
			return true
		}
	}
	return false
}

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
