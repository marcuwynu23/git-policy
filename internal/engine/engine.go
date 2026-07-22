// Package engine orchestrates policy registration and execution.
package engine

import (
	"github.com/marcuwynu23/git-policy/internal/config"
	"github.com/marcuwynu23/git-policy/internal/policy"
)

// Engine manages registered policies and executes them against a given context.
type Engine struct {
	cfg      *config.Config
	policies []policy.Policy
	skipList map[string]bool
}

// New creates a new Engine with the given configuration.
func New(cfg *config.Config) *Engine {
	return &Engine{
		cfg:      cfg,
		skipList: make(map[string]bool),
	}
}

// Register adds a policy to the engine's execution list.
func (e *Engine) Register(p policy.Policy) {
	e.policies = append(e.policies, p)
}

// SetSkipList sets the list of policy names (internal names) to skip during execution.
// These are skipped in addition to any disabled policies.
func (e *Engine) SetSkipList(names []string) {
	e.skipList = make(map[string]bool)
	for _, name := range names {
		e.skipList[name] = true
	}
}

// Execute runs all registered policies using a default context.
func (e *Engine) Execute() []policy.Result {
	return e.ExecuteWith(e.buildContext())
}

// ExecuteWith runs all registered policies against the given context,
// skipping any that are disabled in config or in the skip list.
func (e *Engine) ExecuteWith(ctx policy.Context) []policy.Result {
	var results []policy.Result
	for _, p := range e.policies {
		if e.cfg.Policies.IsDisabled(p.Name()) {
			continue
		}
		if e.skipList[p.Name()] {
			continue
		}
		result := p.Validate(ctx)
		results = append(results, result)
	}
	return results
}

func (e *Engine) buildContext() policy.Context {
	return policy.Context{RepoPath: "."}
}

// PolicyNames returns the names of all registered policies.
func (e *Engine) PolicyNames() []string {
	var names []string
	for _, p := range e.policies {
		names = append(names, p.Name())
	}
	return names
}
