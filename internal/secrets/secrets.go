package secrets

import (
	"bufio"
	"os"
	"strings"
)

type Pattern struct {
	Name   string
	Prefix string
}

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

type Finder struct {
	patterns []Pattern
}

func NewFinder(patterns []Pattern) *Finder {
	if patterns == nil {
		patterns = DefaultPatterns
	}
	return &Finder{patterns: patterns}
}

type Finding struct {
	Pattern Pattern
	File    string
	Line    int
}

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
