package commitmsg

import "strings"

var ValidPrefixes = []string{
	"feat", "fix", "refactor", "docs", "test",
	"build", "ci", "style", "perf", "chore", "revert",
}

type Validator struct {
	enforceConventional bool
}

func NewValidator(enforceConventional bool) *Validator {
	return &Validator{enforceConventional: enforceConventional}
}

type ValidationResult struct {
	Valid   bool
	Message string
}

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
