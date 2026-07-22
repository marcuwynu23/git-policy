package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/marcuwynu23/git-policy/internal/config"
	"github.com/marcuwynu23/git-policy/internal/git"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
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
	names = append(names, "custom:<name>")
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
	Short: "Add a custom rule to the config",
	Long: `Add a custom rule directly to the configuration file.

The rule is stored in config.yaml under the customRules section and
runs alongside built-in rules on every commit.

Supported rule types:
  file-block      Block files matching a glob pattern
  file-content    Scan file contents for a string pattern
  branch-name     Block commits to branches matching a pattern
  commit-message  Block commits with messages matching a pattern

Example:
  git-policy rule add no-todo --type file-content --pattern "TODO:" --message "No todos" --fix "Resolve TODO"`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ruleType, _ := cmd.Flags().GetString("type")
		pattern, _ := cmd.Flags().GetString("pattern")
		message, _ := cmd.Flags().GetString("message")
		fix, _ := cmd.Flags().GetString("fix")

		if ruleType == "" {
			return fmt.Errorf("--type is required (file-block, file-content, branch-name, commit-message)")
		}
		if pattern == "" {
			return fmt.Errorf("--pattern is required")
		}
		if message == "" {
			return fmt.Errorf("--message is required")
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

		cfg.AddCustomRule(config.CustomRuleDef{
			Name:    args[0],
			Type:    ruleType,
			Pattern: pattern,
			Message: message,
			Fix:     fix,
		})

		if err := config.Save(cfg, path); err != nil {
			return fmt.Errorf("saving config: %w", err)
		}

		cmd.Printf("Custom rule %q added (%s).\n", args[0], ruleType)
		return nil
	},
}

var policyRemoveCmd = &cobra.Command{
	Use:   "remove [name]",
	Short: "Remove a custom rule from the config",
	Long: `Remove a custom rule from the configuration file.

Example:
  git-policy rule remove no-todo`,
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

		if !cfg.RemoveCustomRule(args[0]) {
			return fmt.Errorf("custom rule %q not found", args[0])
		}

		if err := config.Save(cfg, path); err != nil {
			return fmt.Errorf("saving config: %w", err)
		}

		cmd.Printf("Custom rule %q removed.\n", args[0])
		return nil
	},
}

var policyExportCmd = &cobra.Command{
	Use:   "export [name]",
	Short: "Export a custom rule to a YAML file",
	Long: `Export a custom rule as a standalone YAML file.

The exported file can be shared or re-imported on another machine.

Example:
  git-policy rule export no-todo -o ./my-rule.yaml`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		output, _ := cmd.Flags().GetString("output")

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

		var rule *config.CustomRuleDef
		for _, r := range cfg.CustomRules {
			if r.Name == args[0] {
				rule = &r
				break
			}
		}
		if rule == nil {
			return fmt.Errorf("custom rule %q not found", args[0])
		}

		if output == "" {
			output = filepath.Join(config.RulesDir(path), args[0]+".yaml")
			if err := os.MkdirAll(config.RulesDir(path), 0755); err != nil {
				return fmt.Errorf("creating rules directory: %w", err)
			}
		}

		data, err := yaml.Marshal(rule)
		if err != nil {
			return fmt.Errorf("marshaling rule: %w", err)
		}

		if err := os.WriteFile(output, data, 0644); err != nil {
			return fmt.Errorf("writing rule file: %w", err)
		}

		cmd.Printf("Rule %q exported to %s.\n", args[0], output)
		return nil
	},
}

var policyImportCmd = &cobra.Command{
	Use:   "import [file]",
	Short: "Import a custom rule from a YAML file",
	Long: `Import a custom rule from a YAML file into the configuration.

The file should contain a single rule definition:

  name: no-todo
  type: file-content
  pattern: "TODO:"
  message: "No todos"
  fix: "Resolve TODO"

Example:
  git-policy rule import ./my-rule.yaml`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := os.ReadFile(args[0])
		if err != nil {
			return fmt.Errorf("reading rule file: %w", err)
		}

		var rule config.CustomRuleDef
		if err := yaml.Unmarshal(data, &rule); err != nil {
			return fmt.Errorf("parsing rule file: %w", err)
		}
		if rule.Name == "" {
			return fmt.Errorf("rule file %q: name is required", args[0])
		}
		if rule.Type == "" {
			return fmt.Errorf("rule file %q: type is required", args[0])
		}
		if rule.Pattern == "" {
			return fmt.Errorf("rule file %q: pattern is required", args[0])
		}
		if rule.Message == "" {
			return fmt.Errorf("rule file %q: message is required", args[0])
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

		cfg.AddCustomRule(rule)

		if err := config.Save(cfg, path); err != nil {
			return fmt.Errorf("saving config: %w", err)
		}

		cmd.Printf("Custom rule %q imported (%s).\n", rule.Name, rule.Type)
		return nil
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

		pluginsDir := config.PluginsDir(path)
		if err := os.MkdirAll(pluginsDir, 0755); err != nil {
			return fmt.Errorf("creating plugins directory: %w", err)
		}

		dstName := desc.Name + ".yaml"
		dstPath := filepath.Join(pluginsDir, dstName)
		srcData, err := os.ReadFile(args[0])
		if err != nil {
			return fmt.Errorf("reading plugin descriptor: %w", err)
		}
		if err := os.WriteFile(dstPath, srcData, 0644); err != nil {
			return fmt.Errorf("copying plugin to %s: %w", dstPath, err)
		}

		cfg, err := config.Load(path)
		if err != nil {
			return fmt.Errorf("loading config: %w", err)
		}

		cfg.AddPlugin(config.PluginEntry{
			Name:    desc.Name,
			Path:    dstPath,
			Enabled: !disabled,
		})

		if err := config.Save(cfg, path); err != nil {
			return fmt.Errorf("saving config: %w", err)
		}

		status := "enabled"
		if disabled {
			status = "disabled"
		}
		cmd.Printf("Plugin %q installed (%d rules, %s).\n  File: %s\n", desc.Name, len(desc.Rules), status, dstPath)
		return nil
	},
}

var pluginsUninstallCmd = &cobra.Command{
	Use:   "uninstall [name]",
	Short: "Uninstall a plugin by name",
	Long: `Remove a plugin and its descriptor file from the configuration.

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

		entry := cfg.RemovePlugin(args[0])
		if !entry {
			return fmt.Errorf("plugin %q not found", args[0])
		}

		if err := config.Save(cfg, path); err != nil {
			return fmt.Errorf("saving config: %w", err)
		}

		// Remove the descriptor file from the plugins directory
		pluginFile := filepath.Join(config.PluginsDir(path), args[0]+".yaml")
		if rmErr := os.Remove(pluginFile); rmErr != nil && !os.IsNotExist(rmErr) {
			cmd.Printf("Warning: could not remove plugin file %s: %v\n", pluginFile, rmErr)
		} else {
			cmd.Printf("  Removed: %s\n", pluginFile)
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
		_, isBuiltin := config.PolicyNames[name]
		isCustom := strings.HasPrefix(name, "custom:")
		if !isBuiltin && !isCustom {
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
			cmd.Printf("  %-20s %-6s  %s\n", p.Name, status, p.Path)
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

	policyAddCmd.Flags().String("type", "", "Rule type (file-block, file-content, branch-name, commit-message)")
	policyAddCmd.Flags().String("pattern", "", "Pattern to match (glob or text)")
	policyAddCmd.Flags().String("message", "", "Error message when rule blocks")
	policyAddCmd.Flags().String("fix", "", "Suggested fix (optional)")

	policyExportCmd.Flags().StringP("output", "o", "", "Output file path (default: <config-dir>/rules/<name>.yaml)")

	rootCmd.AddCommand(pluginsCmd)
	pluginsCmd.AddCommand(pluginsInstallCmd)
	pluginsCmd.AddCommand(pluginsUninstallCmd)
	pluginsCmd.AddCommand(pluginsListCmd)
	pluginsInstallCmd.Flags().Bool("disabled", false, "Install the plugin with rules disabled")
}
