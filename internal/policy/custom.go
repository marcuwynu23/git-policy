package policy

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/marcuwynu23/git-policy/internal/config"
)

// CustomPolicy validates a single YAML-defined custom rule.
type CustomPolicy struct {
	def config.CustomRuleDef
}

// NewCustomPolicy creates a new CustomPolicy from a rule definition.
func NewCustomPolicy(def config.CustomRuleDef) *CustomPolicy {
	return &CustomPolicy{def: def}
}

// Name returns the policy name (plugin rule name).
func (p *CustomPolicy) Name() string {
	return "Custom:" + p.def.Name
}

// Validate runs the rule against the given context.
func (p *CustomPolicy) Validate(ctx Context) Result {
	switch p.def.Type {
	case "file-block":
		return p.validateFileBlock(ctx)
	case "file-content":
		return p.validateFileContent(ctx)
	case "branch-name":
		return p.validateBranchName(ctx)
	case "commit-message":
		return p.validateCommitMessage(ctx)
	default:
		return Result{PolicyName: p.Name(), Status: StatusPass,
			Message: fmt.Sprintf("unknown rule type %q", p.def.Type)}
	}
}

func (p *CustomPolicy) validateFileBlock(ctx Context) Result {
	for _, file := range ctx.StagedFiles {
		if matched, _ := filepath.Match(p.def.Pattern, filepath.Base(file)); matched {
			return Result{
				PolicyName: p.Name(),
				Status:     StatusFail,
				Message:    p.def.Message,
				Fix:        p.def.Fix,
			}
		}
	}
	return Result{PolicyName: p.Name(), Status: StatusPass}
}

func (p *CustomPolicy) validateFileContent(ctx Context) Result {
	for _, file := range ctx.StagedFiles {
		f, err := os.Open(file)
		if err != nil {
			continue
		}
		defer f.Close()
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			if strings.Contains(scanner.Text(), p.def.Pattern) {
				f.Close()
				return Result{
					PolicyName: p.Name(),
					Status:     StatusFail,
					Message:    p.def.Message,
					Fix:        p.def.Fix,
				}
			}
		}
		f.Close()
	}
	return Result{PolicyName: p.Name(), Status: StatusPass}
}

func (p *CustomPolicy) validateBranchName(ctx Context) Result {
	if matched, _ := filepath.Match(p.def.Pattern, ctx.BranchName); matched {
		return Result{
			PolicyName: p.Name(),
			Status:     StatusFail,
			Message:    p.def.Message,
			Fix:        p.def.Fix,
		}
	}
	return Result{PolicyName: p.Name(), Status: StatusPass}
}

func (p *CustomPolicy) validateCommitMessage(ctx Context) Result {
	if matched, _ := filepath.Match(p.def.Pattern, ctx.CommitMsg); matched {
		return Result{
			PolicyName: p.Name(),
			Status:     StatusFail,
			Message:    p.def.Message,
			Fix:        p.def.Fix,
		}
	}
	return Result{PolicyName: p.Name(), Status: StatusPass}
}
