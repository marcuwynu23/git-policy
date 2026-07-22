package cmd

import (
	"fmt"
	"os"

	"github.com/marcuwynu23/git-policy/internal/config"
	"github.com/marcuwynu23/git-policy/internal/hook"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install global Git hooks",
	Long: `Install global Git hooks to enforce policies across all repositories.

This command sets up git-policy hooks globally using git config core.hooksPath
or by installing hooks into the global hooks directory.

It also creates the plugins/ and rules/ directories alongside the config file
for managing custom rules and plugin descriptors.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		installer := hook.NewInstaller()
		if err := installer.InstallGlobal(); err != nil {
			return fmt.Errorf("failed to install hooks: %w", err)
		}

		cfgPath, err := config.DefaultConfigPath()
		if err != nil {
			return fmt.Errorf("determining config path: %w", err)
		}

		pluginsDir := config.PluginsDir(cfgPath)
		if err := os.MkdirAll(pluginsDir, 0755); err != nil {
			return fmt.Errorf("creating plugins directory: %w", err)
		}

		rulesDir := config.RulesDir(cfgPath)
		if err := os.MkdirAll(rulesDir, 0755); err != nil {
			return fmt.Errorf("creating rules directory: %w", err)
		}

		fmt.Printf("Git hooks installed successfully.\n  Config: %s\n  Plugins: %s\n  Rules: %s\n", cfgPath, pluginsDir, rulesDir)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}
