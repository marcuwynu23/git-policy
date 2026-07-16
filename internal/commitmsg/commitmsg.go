// Package commitmsg validates commit messages against the conventional commits format.
package commitmsg

import "strings"

// ValidPrefixes lists all accepted conventional commit prefixes.
var ValidPrefixes = []string{
	"feat", "fix", "refactor", "docs", "test",
	"build", "ci", "style", "perf", "chore", "revert",
}

// Validator checks commit messages for conventional commit compliance.
type Validator struct {
	enforceConventional bool
}

// NewValidator creates a new Validator with the given enforcement setting.
func NewValidator(enforceConventional bool) *Validator {
	return &Validator{enforceConventional: enforceConventional}
}

// ValidationResult holds the outcome of a commit message validation.
type ValidationResult struct {
	Valid   bool
	Message string
}

// Validate checks whether a commit message follows the conventional commits format.
func (v *Validator) Validate(msg string) ValidationResult {
	if !v.enforceConventional {
		return ValidationResult{Valid: true}
	}

	firstLine := strings.SplitN(strings.TrimSpace(msg), "\n", 2)[0]
	if firstLine == "" {
		return ValidationResult{
			Valid:   false,
			Message: "commit message is empty",
		}
	}

	for _, prefix := range ValidPrefixes {
		if strings.HasPrefix(firstLine, prefix+":") ||
			strings.HasPrefix(firstLine, strings.ToUpper(prefix[:1])+prefix[1:]+":") {
			return ValidationResult{Valid: true}
		}
	}

	return ValidationResult{
		Valid:   false,
		Message: "does not follow conventional commits format",
	}
}
