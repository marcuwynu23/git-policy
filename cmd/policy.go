package cmd

import (
	"fmt"
	"sort"
	"strings"

	"github.com/marcuwynu23/git-policy/internal/config"
	"github.com/spf13/cobra"
)

var policyCmd = &cobra.Command{
	Use:     "rule",
	Aliases: []string{"rules", "policy", "policies"},
	Short:   "Manage policies",
	Long:    `Enable, disable, and list policies.`,
}

var policyListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all rules and their status",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load(cfgFile)
		if err != nil {
			return fmt.Errorf("loading config: %w", err)
		}

		var names []string
		for cli := range config.PolicyNames {
			names = append(names, cli)
		}
		sort.Strings(names)

		fmt.Println("Rules:")
		for _, name := range names {
			internalName := config.PolicyNames[name]
			status := "enabled"
			if cfg.Policies.IsDisabled(internalName) {
				status = "disabled"
			}
			fmt.Printf("  %-20s %s\n", name, status)
		}
		return nil
	},
}

var policyEnableCmd = &cobra.Command{
	Use:   "enable [name]",
	Short: "Enable a rule",
	Args:  cobra.ExactArgs(1),
	Long: `Enable a rule by name.

Available rules: ` + availablePolicyNames(),
	RunE: func(cmd *cobra.Command, args []string) error {
		return setPolicyEnabled(args[0], false)
	},
}

var policyDisableCmd = &cobra.Command{
	Use:   "disable [name]",
	Short: "Disable a rule",
	Args:  cobra.ExactArgs(1),
	Long: `Disable a rule by name.

Available rules: ` + availablePolicyNames(),
	RunE: func(cmd *cobra.Command, args []string) error {
		return setPolicyEnabled(args[0], true)
	},
}

func availablePolicyNames() string {
	var names []string
	for cli := range config.PolicyNames {
		names = append(names, cli)
	}
	sort.Strings(names)
	return strings.Join(names, ", ")
}

func setPolicyEnabled(cliName string, disabled bool) error {
	internalName, ok := config.PolicyNames[cliName]
	if !ok {
		return fmt.Errorf("unknown rule %q\n\nAvailable: %s", cliName, availablePolicyNames())
	}

	path := cfgFile
	if path == "" {
		var err error
		path, err = config.DefaultConfigPath()
		if err != nil {
			return fmt.Errorf("determining config path: %w", err)
		}
	}

	cfg, err := config.Load(path)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	cfg.Policies.SetDisabled(internalName, disabled)

	if err := config.Save(cfg, path); err != nil {
		return fmt.Errorf("saving config: %w", err)
	}

	action := "enabled"
	if disabled {
		action = "disabled"
	}
	fmt.Printf("Rule %q %s.\n", cliName, action)
	return nil
}

var policyAddCmd = &cobra.Command{
	Use:   "add [name]",
	Short: "Add a custom rule (not yet implemented)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return fmt.Errorf("rule add not yet implemented")
	},
}

var policyRemoveCmd = &cobra.Command{
	Use:   "remove [name]",
	Short: "Remove a rule (not yet implemented)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return fmt.Errorf("rule remove not yet implemented")
	},
}

var policyExportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export rules to a file (not yet implemented)",
	RunE: func(cmd *cobra.Command, args []string) error {
		return fmt.Errorf("rule export not yet implemented")
	},
}

var policyImportCmd = &cobra.Command{
	Use:   "import [file]",
	Short: "Import rules from a file (not yet implemented)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return fmt.Errorf("rule import not yet implemented")
	},
}

var pluginsCmd = &cobra.Command{
	Use:   "plugins",
	Short: "Manage plugins",
	Long:  `Install, list, and manage git-policy plugins.`,
}

var pluginsInstallCmd = &cobra.Command{
	Use:   "install [path]",
	Short: "Install a plugin (not yet implemented)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return fmt.Errorf("plugins install not yet implemented")
	},
}

var pluginsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List installed plugins",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("No plugins installed.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(policyCmd)
	policyCmd.AddCommand(policyListCmd)
	policyCmd.AddCommand(policyEnableCmd)
	policyCmd.AddCommand(policyDisableCmd)
	policyCmd.AddCommand(policyAddCmd)
	policyCmd.AddCommand(policyRemoveCmd)
	policyCmd.AddCommand(policyExportCmd)
	policyCmd.AddCommand(policyImportCmd)

	rootCmd.AddCommand(pluginsCmd)
	pluginsCmd.AddCommand(pluginsInstallCmd)
	pluginsCmd.AddCommand(pluginsListCmd)
}
