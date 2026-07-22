// Package runner wires Git context into the policy engine and executes all policies.
package runner

import (
	"fmt"
	"os"
	"strings"

	"github.com/marcuwynu23/git-policy/internal/config"
	"github.com/marcuwynu23/git-policy/internal/engine"
	"github.com/marcuwynu23/git-policy/internal/git"
	"github.com/marcuwynu23/git-policy/internal/plugins"
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

	skipNames := readSkipList()
	eng.SetSkipList(skipNames)

	eng.Register(policy.NewBlockFilesPolicy(cfg))
	eng.Register(policy.NewCommitMessagePolicy(cfg))
	eng.Register(policy.NewFileSizePolicy(cfg))
	eng.Register(policy.NewBinaryFilePolicy(cfg))
	eng.Register(policy.NewSecretScanPolicy(cfg))
	eng.Register(policy.NewBranchPolicy(cfg))

	if err := registerPluginPolicies(cfg, eng); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: plugin load failed: %v\n", err)
	}

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

	// Auto-clear skip list on success
	if len(skipNames) > 0 {
		_ = git.UnsetConfig("git-policy.skip")
	}

	fmt.Println("All policies passed.")
	return nil
}

// registerPluginPolicies loads all enabled plugins and registers their policies.
func registerPluginPolicies(cfg *config.Config, eng *engine.Engine) error {
	if len(cfg.Plugins) == 0 {
		return nil
	}
	loader := plugins.NewLoader()
	policies, err := loader.PoliciesFromPlugins(cfg.Plugins)
	if err != nil {
		return err
	}
	for _, p := range policies {
		eng.Register(p)
	}
	return nil
}

// readSkipList reads the skip list from local git config and converts
// CLI rule names to internal policy names.
func readSkipList() []string {
	raw, err := git.GetConfig("git-policy.skip")
	if err != nil || raw == "" {
		return nil
	}
	var internalNames []string
	for _, name := range strings.Split(raw, ",") {
		name = strings.TrimSpace(name)
		if name == "" {
			continue
		}
		if internal, ok := config.PolicyNames[name]; ok {
			internalNames = append(internalNames, internal)
		}
	}
	return internalNames
}
