package branch

type Protector struct {
	protectedBranches []string
}

func NewProtector(branches []string) *Protector {
	if branches == nil {
		branches = []string{"main", "master", "production"}
	}
	return &Protector{protectedBranches: branches}
}

func (p *Protector) IsProtected(branch string) bool {
	for _, protected := range p.protectedBranches {
		if branch == protected {
			return true
		}
	}
	return false
}

func (p *Protector) ProtectedBranches() []string {
	result := make([]string, len(p.protectedBranches))
	copy(result, p.protectedBranches)
	return result
}
