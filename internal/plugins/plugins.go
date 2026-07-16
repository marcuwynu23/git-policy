// Package plugins provides a stub for loading external policy plugins.
package plugins

import (
	"fmt"
	"plugin"

	"github.com/marcuwynu23/git-policy/internal/policy"
)

// Plugin defines the interface that external plugins must implement.
type Plugin interface {
	Policies() []policy.Policy
}

// Loader manages loading and querying external Go plugins.
type Loader struct {
	plugins []Plugin
}

// NewLoader creates a new Loader.
func NewLoader() *Loader {
	return &Loader{}
}

// Load opens and validates a Go plugin file.
func (l *Loader) Load(path string) error {
	p, err := plugin.Open(path)
	if err != nil {
		return fmt.Errorf("opening plugin %s: %w", path, err)
	}
	sym, err := p.Lookup("Plugin")
	if err != nil {
		return fmt.Errorf("plugin %s does not export Plugin symbol: %w", path, err)
	}
	plug, ok := sym.(Plugin)
	if !ok {
		return fmt.Errorf("plugin %s does not implement Plugin interface", path)
	}
	l.plugins = append(l.plugins, plug)
	return nil
}

// Policies returns all policies from all loaded plugins.
func (l *Loader) Policies() []policy.Policy {
	var all []policy.Policy
	for _, plug := range l.plugins {
		all = append(all, plug.Policies()...)
	}
	return all
}

// Count returns the number of loaded plugins.
func (l *Loader) Count() int {
	return len(l.plugins)
}
