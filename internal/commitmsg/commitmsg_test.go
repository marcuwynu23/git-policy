package commitmsg

import "testing"

func TestValidate_ValidMessages(t *testing.T) {
	v := NewValidator(true)
	tests := []string{
		"feat: add new feature",
		"fix: resolve bug",
		"docs: update readme",
		"refactor: clean up code",
		"test: add tests",
		"build: update deps",
		"ci: fix pipeline",
		"style: format code",
		"perf: optimize loop",
		"chore: bump version",
		"revert: previous change",
	}
	for _, msg := range tests {
		result := v.Validate(msg)
		if !result.Valid {
			t.Errorf("expected valid for %q, got invalid: %s", msg, result.Message)
		}
	}
}

func TestValidate_InvalidMessages(t *testing.T) {
	v := NewValidator(true)
	tests := []string{
		"added stuff",
		"fix bugs",
		"WIP",
	}
	for _, msg := range tests {
		result := v.Validate(msg)
		if result.Valid {
			t.Errorf("expected invalid for %q, got valid", msg)
		}
	}
}

func TestValidate_Disabled(t *testing.T) {
	v := NewValidator(false)
	result := v.Validate("random message")
	if !result.Valid {
		t.Errorf("expected valid when disabled, got invalid")
	}
}

func TestValidate_Empty(t *testing.T) {
	v := NewValidator(true)
	result := v.Validate("")
	if result.Valid {
		t.Error("expected invalid for empty message")
	}
}
