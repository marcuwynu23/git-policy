// Package runner wires Git context into the policy engine and executes all policies.
package runner

import (
	"fmt"
	"os"

	"github.com/marcuwynu23/git-policy/internal/config"
	"github.com/marcuwynu23/git-policy/internal/engine"
	"github.com/marcuwynu23/git-policy/internal/git"
	"github.com/marcuwynu23/git-policy/internal/policy"
)

// Run loads staged files and branch info from Git, then executes all
// enabled policies against them.
func Run(cfg *config.Config) error {
	if !git.IsRepo() {
		fmt.Fprintln(os.Stderr, "Not a git repository.")
		os.Exit(1)
	}

	stagedFiles, err := git.GetStagedFiles()
	if err != nil {
		return fmt.Errorf("getting staged files: %w", err)
	}

	branchName, err := git.GetBranchName()
	if err != nil {
		return fmt.Errorf("getting branch name: %w", err)
	}

	eng := engine.New(cfg)
	eng.Register(policy.NewBlockFilesPolicy(cfg))
	eng.Register(policy.NewCommitMessagePolicy(cfg))
	eng.Register(policy.NewFileSizePolicy(cfg))
	eng.Register(policy.NewBinaryFilePolicy(cfg))
	eng.Register(policy.NewSecretScanPolicy(cfg))
	eng.Register(policy.NewBranchPolicy(cfg))

	results := eng.ExecuteWith(policy.Context{
		RepoPath:    ".",
		StagedFiles: stagedFiles,
		BranchName:  branchName,
	})

	hasErrors := false
	for _, result := range results {
		if result.Status == policy.StatusFail {
			hasErrors = true
			fmt.Fprintf(os.Stderr, "BLOCKED: %s - %s\n", result.PolicyName, result.Message)
			if result.Fix != "" {
				fmt.Fprintf(os.Stderr, "  Fix: %s\n", result.Fix)
			}
		}
	}

	if hasErrors {
		fmt.Fprintln(os.Stderr, "\nCommit blocked by git-policy.")
		os.Exit(1)
	}

	fmt.Println("All policies passed.")
	return nil
}
