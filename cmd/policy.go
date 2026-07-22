package cmd

import (
	"fmt"
	"sort"
	"strings"

	"github.com/marcuwynu23/git-policy/internal/config"
	"github.com/marcuwynu23/git-policy/internal/git"
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
	Short: "Install a plugin from a YAML descriptor file",
	Long: `Install a plugin from a YAML descriptor file.

The descriptor file defines the plugin name and custom rules:

  name: my-custom-rules
  rules:
    - name: no-todo
      type: file-content
      pattern: "TODO:"
      message: "Commits containing TODO are not allowed"
      fix: "Resolve the TODO before committing"

Supported rule types:
  file-block      Block files matching a glob pattern
  file-content    Scan file contents for a string pattern
  branch-name     Block commits to branches matching a pattern
  commit-message  Block commits with messages matching a pattern

Use --disabled to install the plugin with all rules disabled by default.

Example:
  git-policy plugins install ./my-plugin.yaml
  git-policy plugins install --disabled ./my-plugin.yaml`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		disabled, _ := cmd.Flags().GetBool("disabled")

		desc, err := config.LoadPluginDescriptor(args[0])
		if err != nil {
			return err
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

		cfg.AddPlugin(config.PluginEntry{
			Name:    desc.Name,
			Enabled: !disabled,
			Rules:   desc.Rules,
		})

		if err := config.Save(cfg, path); err != nil {
			return fmt.Errorf("saving config: %w", err)
		}

		status := "enabled"
		if disabled {
			status = "disabled"
		}
		cmd.Printf("Plugin %q installed (%d rules, %s).\n", desc.Name, len(desc.Rules), status)
		return nil
	},
}

var pluginsUninstallCmd = &cobra.Command{
	Use:   "uninstall [name]",
	Short: "Uninstall a plugin by name",
	Long: `Remove a plugin and all its rules from the configuration.

Example:
  git-policy plugins uninstall my-custom-rules`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
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

		if !cfg.RemovePlugin(args[0]) {
			return fmt.Errorf("plugin %q not found", args[0])
		}

		if err := config.Save(cfg, path); err != nil {
			return fmt.Errorf("saving config: %w", err)
		}

		cmd.Printf("Plugin %q uninstalled.\n", args[0])
		return nil
	},
}

var policySkipCmd = &cobra.Command{
	Use:   "skip [name...]",
	Short: "Temporarily skip rules for the current commit",
	Long: `Skip one or more rules for the current commit.

Rules are stored in the repository's local git config and are
automatically cleared after a successful commit.

Use --list to see currently skipped rules.
Use --clear to remove all skipped rules.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if cmd.Flags().Changed("list") {
			return listSkippedRules(cmd)
		}
		if cmd.Flags().Changed("clear") {
			return clearSkippedRules(cmd)
		}
		if len(args) == 0 {
			return listSkippedRules(cmd)
		}
		return addSkippedRules(cmd, args)
	},
}

func listSkippedRules(cmd *cobra.Command) error {
	if !git.IsRepo() {
		return fmt.Errorf("not a git repository")
	}
	raw, err := git.GetConfig("git-policy.skip")
	if err != nil || raw == "" {
		cmd.Println("No rules currently skipped.")
		return nil
	}
	names := strings.Split(raw, ",")
	cmd.Println("Skipped rules:")
	for _, name := range names {
		name = strings.TrimSpace(name)
		if name != "" {
			cmd.Printf("  %s\n", name)
		}
	}
	return nil
}

func clearSkippedRules(cmd *cobra.Command) error {
	if !git.IsRepo() {
		return fmt.Errorf("not a git repository")
	}
	if err := git.UnsetConfig("git-policy.skip"); err != nil {
		return fmt.Errorf("clearing skip list: %w", err)
	}
	cmd.Println("All skipped rules cleared.")
	return nil
}

func addSkippedRules(cmd *cobra.Command, cliNames []string) error {
	if !git.IsRepo() {
		return fmt.Errorf("not a git repository")
	}
	for _, name := range cliNames {
		if _, ok := config.PolicyNames[name]; !ok {
			return fmt.Errorf("unknown rule %q\n\nAvailable: %s", name, availablePolicyNames())
		}
	}
	existing := getExistingSkipList()
	for _, name := range cliNames {
		found := false
		for _, e := range existing {
			if e == name {
				found = true
				break
			}
		}
		if !found {
			existing = append(existing, name)
		}
	}
	value := strings.Join(existing, ",")
	if err := git.SetConfig("git-policy.skip", value); err != nil {
		return fmt.Errorf("setting skip list: %w", err)
	}
	cmd.Printf("Skipped rules: %s\n", strings.Join(existing, ", "))
	cmd.Println("Skipped rules will be automatically cleared after a successful commit.")
	return nil
}

func getExistingSkipList() []string {
	raw, err := git.GetConfig("git-policy.skip")
	if err != nil || raw == "" {
		return nil
	}
	var names []string
	for _, n := range strings.Split(raw, ",") {
		n = strings.TrimSpace(n)
		if n != "" {
			names = append(names, n)
		}
	}
	return names
}

var pluginsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List installed plugins",
	RunE: func(cmd *cobra.Command, args []string) error {
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

		if len(cfg.Plugins) == 0 {
			cmd.Println("No plugins installed.")
			return nil
		}

		cmd.Println("Installed plugins:")
		for _, p := range cfg.Plugins {
			status := "enabled"
			if !p.Enabled {
				status = "disabled"
			}
			ruleCount := len(p.Rules)
			cmd.Printf("  %-20s %-6s  %d rule(s)\n", p.Name, status, ruleCount)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(policyCmd)
	policyCmd.AddCommand(policyListCmd)
	policyCmd.AddCommand(policyEnableCmd)
	policyCmd.AddCommand(policyDisableCmd)
	policyCmd.AddCommand(policySkipCmd)
	policyCmd.AddCommand(policyAddCmd)
	policyCmd.AddCommand(policyRemoveCmd)
	policyCmd.AddCommand(policyExportCmd)
	policyCmd.AddCommand(policyImportCmd)

	policySkipCmd.Flags().Bool("list", false, "Show currently skipped rules")
	policySkipCmd.Flags().Bool("clear", false, "Clear all skipped rules")

	rootCmd.AddCommand(pluginsCmd)
	pluginsCmd.AddCommand(pluginsInstallCmd)
	pluginsCmd.AddCommand(pluginsUninstallCmd)
	pluginsCmd.AddCommand(pluginsListCmd)
	pluginsInstallCmd.Flags().Bool("disabled", false, "Install the plugin with rules disabled")
}
