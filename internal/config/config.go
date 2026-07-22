// Package config handles loading, saving, and managing git-policy configuration.
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"gopkg.in/yaml.v3"
)

// CustomRuleDef defines a single custom rule within a plugin.
type CustomRuleDef struct {
	Name    string `yaml:"name"`
	Type    string `yaml:"type"`
	Pattern string `yaml:"pattern"`
	Message string `yaml:"message"`
	Fix     string `yaml:"fix,omitempty"`
}

// PluginEntry represents a single plugin in the configuration.
// Path is the location of the plugin descriptor YAML file.
type PluginEntry struct {
	Name    string `yaml:"name"`
	Path    string `yaml:"path"`
	Enabled bool   `yaml:"enabled"`
}

// Config represents the complete git-policy configuration.
type Config struct {
	Version  int             `yaml:"version"`
	Hooks    HooksConfig     `yaml:"hooks"`
	Policies PoliciesConfig  `yaml:"policies"`
	Plugins  []PluginEntry   `yaml:"plugins,omitempty"`
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

// AddPlugin appends a plugin entry or updates an existing one by name.
func (c *Config) AddPlugin(entry PluginEntry) {
	for i, p := range c.Plugins {
		if p.Name == entry.Name {
			c.Plugins[i] = entry
			return
		}
	}
	c.Plugins = append(c.Plugins, entry)
}

// RemovePlugin removes a plugin entry by name.
func (c *Config) RemovePlugin(name string) bool {
	for i, p := range c.Plugins {
		if p.Name == name {
			c.Plugins = append(c.Plugins[:i], c.Plugins[i+1:]...)
			return true
		}
	}
	return false
}

// PluginDescriptor is the YAML structure for a standalone plugin file.
type PluginDescriptor struct {
	Name  string          `yaml:"name"`
	Rules []CustomRuleDef `yaml:"rules"`
}

// LoadPluginDescriptor reads a plugin descriptor from a YAML file.
func LoadPluginDescriptor(path string) (*PluginDescriptor, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading plugin descriptor %q: %w", path, err)
	}
	var desc PluginDescriptor
	if err := yaml.Unmarshal(data, &desc); err != nil {
		return nil, fmt.Errorf("parsing plugin descriptor %q: %w", path, err)
	}
	if desc.Name == "" {
		return nil, fmt.Errorf("plugin descriptor %q: name is required", path)
	}
	if len(desc.Rules) == 0 {
		return nil, fmt.Errorf("plugin descriptor %q: at least one rule is required", path)
	}
	for i, rule := range desc.Rules {
		if rule.Name == "" {
			return nil, fmt.Errorf("plugin descriptor %q: rule %d: name is required", path, i)
		}
		if rule.Type == "" {
			return nil, fmt.Errorf("plugin descriptor %q: rule %q: type is required", path, rule.Name)
		}
		if rule.Pattern == "" {
			return nil, fmt.Errorf("plugin descriptor %q: rule %q: pattern is required", path, rule.Name)
		}
		if rule.Message == "" {
			return nil, fmt.Errorf("plugin descriptor %q: rule %q: message is required", path, rule.Name)
		}
	}
	return &desc, nil
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
