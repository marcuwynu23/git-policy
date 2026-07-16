// Package git provides wrappers around Git CLI commands used by the policy engine.
package git

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Version returns the installed Git version string.
func Version() (string, error) {
	cmd := exec.Command("git", "--version")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("git not found: %w", err)
	}
	return strings.TrimSpace(string(output)), nil
}

// GetStagedFiles returns the list of files staged for commit.
func GetStagedFiles() ([]string, error) {
	cmd := exec.Command("git", "diff", "--cached", "--name-only", "--diff-filter=ACM")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("getting staged files: %w", err)
	}
	files := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(files) == 1 && files[0] == "" {
		return []string{}, nil
	}
	return files, nil
}

// GetBranchName returns the current Git branch name.
func GetBranchName() (string, error) {
	cmd := exec.Command("git", "branch", "--show-current")
	output, err := cmd.Output()
	if err == nil {
		name := strings.TrimSpace(string(output))
		if name != "" {
			return name, nil
		}
	}
	cmd = exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	output, err = cmd.Output()
	if err == nil {
		name := strings.TrimSpace(string(output))
		if name != "" && name != "HEAD" {
			return name, nil
		}
	}
	cmd = exec.Command("git", "symbolic-ref", "--short", "HEAD")
	output, err = cmd.Output()
	if err == nil {
		name := strings.TrimSpace(string(output))
		if name != "" {
			return name, nil
		}
	}
	return "unknown", nil
}

// GetCommitMsgFile returns the path to the COMMIT_EDITMSG file.
func GetCommitMsgFile() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("getting git dir: %w", err)
	}
	gitDir := strings.TrimSpace(string(output))
	return gitDir + "/COMMIT_EDITMSG", nil
}

// GetCommitMessage reads the current commit message from COMMIT_EDITMSG.
func GetCommitMessage() (string, error) {
	msgFile, err := GetCommitMsgFile()
	if err != nil {
		return "", err
	}
	data, err := os.ReadFile(msgFile)
	if err != nil {
		return "", nil
	}
	return string(data), nil
}

// IsRepo checks whether the current directory is inside a Git repository.
func IsRepo() bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	return cmd.Run() == nil
}
