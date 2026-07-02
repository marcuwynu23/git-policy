package cmd

import (
	"fmt"

	"github.com/marcuwynu23/git-policy/internal/hook"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install global Git hooks",
	Long: `Install global Git hooks to enforce policies across all repositories.

This command sets up git-policy hooks globally using git config core.hooksPath
or by installing hooks into the global hooks directory.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		installer := hook.NewInstaller()
		if err := installer.InstallGlobal(); err != nil {
			return fmt.Errorf("failed to install hooks: %w", err)
		}
		fmt.Println("Git hooks installed successfully.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}
