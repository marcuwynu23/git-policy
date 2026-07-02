package cmd

import (
	"fmt"

	"github.com/marcuwynu23/git-policy/internal/config"
	"github.com/marcuwynu23/git-policy/internal/runner"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run policies against the current repository",
	Long: `Execute all enabled policies against the current Git repository.

This is invoked by Git hooks but can also be run manually to check
your repository against configured policies.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load(cfgFile)
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}
		return runner.Run(cfg)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
