package branch

import "testing"

func TestIsProtected(t *testing.T) {
	p := NewProtector(nil)
	tests := []struct {
		branch   string
		protected bool
	}{
		{"main", true},
		{"master", true},
		{"production", true},
		{"feature/my-feature", false},
		{"develop", false},
		{"main-fix", false},
	}
	for _, tt := range tests {
		got := p.IsProtected(tt.branch)
		if got != tt.protected {
			t.Errorf("IsProtected(%q) = %v, want %v", tt.branch, got, tt.protected)
		}
	}
}

func TestIsProtected_Custom(t *testing.T) {
	p := NewProtector([]string{"develop", "staging"})
	if !p.IsProtected("develop") {
		t.Error("expected develop to be protected")
	}
	if p.IsProtected("main") {
		t.Error("expected main not to be protected")
	}
}

func TestProtectedBranches(t *testing.T) {
	p := NewProtector([]string{"main", "develop"})
	branches := p.ProtectedBranches()
	if len(branches) != 2 {
		t.Fatalf("expected 2 branches, got %d", len(branches))
	}
	// Ensure returned slice is a copy
	branches[0] = "changed"
	if p.ProtectedBranches()[0] != "main" {
		t.Error("ProtectedBranches should return a copy")
	}
}
