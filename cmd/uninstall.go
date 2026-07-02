package cmd

import (
	"fmt"

	"github.com/marcuwynu23/git-policy/internal/hook"
	"github.com/spf13/cobra"
)

var uninstallAll bool

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall global Git hooks",
	Long: `Remove git-policy global hooks and restore previous Git hook configuration.

By default this removes hook files and unsets the global core.hooksPath.
Use --all to also delete the config file and config directory.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		installer := hook.NewInstaller()
		var err error
		if uninstallAll {
			err = installer.UninstallAll()
			if err != nil {
				return fmt.Errorf("failed to fully uninstall: %w", err)
			}
			fmt.Println("git-policy fully uninstalled (hooks + config).")
		} else {
			err = installer.UninstallGlobal()
			if err != nil {
				return fmt.Errorf("failed to uninstall hooks: %w", err)
			}
			fmt.Println("Git hooks uninstalled successfully.")
			fmt.Println("Config kept at ~/.config/git-policy/config.yaml")
			fmt.Println("Use 'git-policy uninstall --all' to also remove config.")
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(uninstallCmd)
	uninstallCmd.Flags().BoolVarP(&uninstallAll, "all", "a", false, "Remove hooks, config file, and config directory")
}
