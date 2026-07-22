// Package plugins manages YAML-defined custom rules.
package plugins

import (
	"fmt"
	"path/filepath"

	"github.com/marcuwynu23/git-policy/internal/config"
	"github.com/marcuwynu23/git-policy/internal/policy"
)

// Loader creates custom policy instances from YAML rule definitions.
type Loader struct{}

// NewLoader creates a new Loader.
func NewLoader() *Loader {
	return &Loader{}
}

// PoliciesFromRules creates Policy instances from a list of rule definitions.
func (l *Loader) PoliciesFromRules(rules []config.CustomRuleDef) []policy.Policy {
	var policies []policy.Policy
	for _, rule := range rules {
		policies = append(policies, policy.NewCustomPolicy(rule))
	}
	return policies
}

// PoliciesFromPlugins creates Policy instances from all enabled plugin entries.
// configDir is used to resolve relative plugin paths. Empty means CWD.
func (l *Loader) PoliciesFromPlugins(entries []config.PluginEntry, configDir string) ([]policy.Policy, error) {
	var all []policy.Policy
	for _, entry := range entries {
		if !entry.Enabled {
			continue
		}
		pluginPath := entry.Path
		if !filepath.IsAbs(pluginPath) && configDir != "" {
			pluginPath = filepath.Join(configDir, pluginPath)
		}
		desc, err := config.LoadPluginDescriptor(pluginPath)
		if err != nil {
			return nil, fmt.Errorf("plugin %q: %w", entry.Name, err)
		}
		all = append(all, l.PoliciesFromRules(desc.Rules)...)
	}
	return all, nil
}
