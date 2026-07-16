// Package sync provides stubs for remote policy synchronization.
package sync

import "fmt"

// Provider defines the interface for remote config providers.
type Provider interface {
	Fetch(url string) ([]byte, error)
	Name() string
}

// GitProvider syncs policies from a Git repository (not yet implemented).
type GitProvider struct{}

func (p *GitProvider) Fetch(url string) ([]byte, error) {
	return nil, fmt.Errorf("git sync not yet implemented")
}

func (p *GitProvider) Name() string { return "git" }

// HTTPProvider syncs policies from an HTTP endpoint (not yet implemented).
type HTTPProvider struct{}

func (p *HTTPProvider) Fetch(url string) ([]byte, error) {
	return nil, fmt.Errorf("http sync not yet implemented")
}

func (p *HTTPProvider) Name() string { return "http" }
