// Package files provides file blocking, binary detection, and size checking.
package files

import (
	"os"
	"path/filepath"
	"strings"
)

// Checker evaluates files against blocklists and size limits.
type Checker struct {
	blockedPatterns []string
	blockedExts     []string
	maxSize         int64
}

// NewChecker creates a new Checker with the given patterns, extensions, and max size.
func NewChecker(blockedPatterns, blockedExts []string, maxSize int64) *Checker {
	if blockedPatterns == nil {
		blockedPatterns = []string{".env", "*.pem", "*.key", "*.p12", "*.crt"}
	}
	if blockedExts == nil {
		blockedExts = []string{".exe", ".dll", ".so", ".iso", ".zip"}
	}
	if maxSize <= 0 {
		maxSize = 10 * 1024 * 1024
	}
	return &Checker{
		blockedPatterns: blockedPatterns,
		blockedExts:     blockedExts,
		maxSize:         maxSize,
	}
}

// IsBlockedFile checks if a filename matches any blocked file pattern.
func (c *Checker) IsBlockedFile(name string) (bool, string) {
	base := filepath.Base(name)
	for _, pattern := range c.blockedPatterns {
		if matched, _ := filepath.Match(pattern, base); matched {
			return true, pattern
		}
	}
	return false, ""
}

// IsBinaryFile checks if a filename has a blocked binary extension.
func (c *Checker) IsBinaryFile(name string) (bool, string) {
	ext := strings.ToLower(filepath.Ext(name))
	for _, blocked := range c.blockedExts {
		if ext == blocked {
			return true, ext
		}
	}
	return false, ""
}

// IsOverMaxSize checks if a file exceeds the configured maximum size.
func (c *Checker) IsOverMaxSize(path string) (bool, int64, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, 0, err
	}
	return info.Size() > c.maxSize, info.Size(), nil
}
