package sync

import "fmt"

type Provider interface {
	Fetch(url string) ([]byte, error)
	Name() string
}

type GitProvider struct{}

func (p *GitProvider) Fetch(url string) ([]byte, error) {
	return nil, fmt.Errorf("git sync not yet implemented")
}

func (p *GitProvider) Name() string { return "git" }

type HTTPProvider struct{}

func (p *HTTPProvider) Fetch(url string) ([]byte, error) {
	return nil, fmt.Errorf("http sync not yet implemented")
}

func (p *HTTPProvider) Name() string { return "http" }
