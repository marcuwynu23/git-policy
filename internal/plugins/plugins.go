package plugins

import (
	"fmt"
	"plugin"

	"github.com/marcuwynu23/git-policy/internal/policy"
)

type Plugin interface {
	Policies() []policy.Policy
}

type Loader struct {
	plugins []Plugin
}

func NewLoader() *Loader {
	return &Loader{}
}

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

func (l *Loader) Policies() []policy.Policy {
	var all []policy.Policy
	for _, plug := range l.plugins {
		all = append(all, plug.Policies()...)
	}
	return all
}

func (l *Loader) Count() int {
	return len(l.plugins)
}
