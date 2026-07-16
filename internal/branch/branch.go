// Package branch provides branch protection checking.
package branch

// Protector checks whether a branch is protected from direct commits.
type Protector struct {
	protectedBranches []string
}

// NewProtector creates a new Protector with the given list of protected branches.
func NewProtector(branches []string) *Protector {
	if branches == nil {
		branches = []string{"main", "master", "production"}
	}
	return &Protector{protectedBranches: branches}
}

// IsProtected returns true if the given branch is in the protected list.
func (p *Protector) IsProtected(branch string) bool {
	for _, protected := range p.protectedBranches {
		if branch == protected {
			return true
		}
	}
	return false
}

// ProtectedBranches returns a copy of the protected branches list.
func (p *Protector) ProtectedBranches() []string {
	result := make([]string, len(p.protectedBranches))
	copy(result, p.protectedBranches)
	return result
}
