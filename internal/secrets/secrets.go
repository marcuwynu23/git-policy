// Package secrets provides secret pattern scanning for staged files.
package secrets

import (
	"bufio"
	"os"
	"strings"
)

// Pattern defines a secret pattern to detect in files.
type Pattern struct {
	Name   string
	Prefix string
}

// DefaultPatterns is the built-in list of secret patterns to scan for.
var DefaultPatterns = []Pattern{
	{Name: "AWS Access Key ID", Prefix: "AKIA"},
	{Name: "AWS Secret Access Key", Prefix: "SecretAccessKey"},
	{Name: "GitHub PAT", Prefix: "ghp_"},
	{Name: "GitLab PAT", Prefix: "glpat-"},
	{Name: "Google API Key", Prefix: "AIza"},
	{Name: "Stripe Live Key", Prefix: "sk_live_"},
	{Name: "Stripe Test Key", Prefix: "sk_test_"},
	{Name: "OpenAI Key", Prefix: "sk-"},
	{Name: "Slack Bot Token", Prefix: "xoxb-"},
	{Name: "Slack App Token", Prefix: "xapp-"},
	{Name: "JWT", Prefix: "eyJ"},
}

// Finder scans files for secret patterns.
type Finder struct {
	patterns []Pattern
}

// NewFinder creates a new Finder with the given patterns (defaults if nil).
func NewFinder(patterns []Pattern) *Finder {
	if patterns == nil {
		patterns = DefaultPatterns
	}
	return &Finder{patterns: patterns}
}

// Finding represents a single secret match in a file.
type Finding struct {
	Pattern Pattern
	File    string
	Line    int
}

// ScanFile scans a single file for secret patterns and returns all findings.
func (f *Finder) ScanFile(path string) ([]Finding, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var findings []Finding
	scanner := bufio.NewScanner(file)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		for _, pattern := range f.patterns {
			if strings.Contains(line, pattern.Prefix) {
				findings = append(findings, Finding{
					Pattern: pattern,
					File:    path,
					Line:    lineNum,
				})
			}
		}
	}
	return findings, scanner.Err()
}

// ScanFiles scans multiple files and aggregates all secret findings.
func (f *Finder) ScanFiles(files []string) ([]Finding, error) {
	var all []Finding
	for _, file := range files {
		findings, err := f.ScanFile(file)
		if err != nil {
			continue
		}
		all = append(all, findings...)
	}
	return all, nil
}
