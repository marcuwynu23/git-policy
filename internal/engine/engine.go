package engine

import (
	"github.com/marcuwynu23/git-policy/internal/config"
	"github.com/marcuwynu23/git-policy/internal/policy"
)

type Engine struct {
	cfg      *config.Config
	policies []policy.Policy
}

func New(cfg *config.Config) *Engine {
	return &Engine{
		cfg: cfg,
	}
}

func (e *Engine) Register(p policy.Policy) {
	e.policies = append(e.policies, p)
}

func (e *Engine) Execute() []policy.Result {
	return e.ExecuteWith(e.buildContext())
}

func (e *Engine) ExecuteWith(ctx policy.Context) []policy.Result {
	var results []policy.Result
	for _, p := range e.policies {
		if e.cfg.Policies.IsDisabled(p.Name()) {
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

func (e *Engine) PolicyNames() []string {
	var names []string
	for _, p := range e.policies {
		names = append(names, p.Name())
	}
	return names
}
