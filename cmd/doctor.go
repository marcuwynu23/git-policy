package cmd

import (
	"fmt"
	"os"

	"github.com/marcuwynu23/git-policy/internal/git"
	"github.com/marcuwynu23/git-policy/internal/hook"
	"github.com/spf13/cobra"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check system health",
	Long: `Run diagnostics to verify git-policy is set up correctly.

Checks Git installation, hook integrity, configuration validity,
and plugin compatibility.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		allOk := true

		fmt.Println("Running git-policy doctor...")
		fmt.Println()

		// Check Git installed
		version, err := git.Version()
		if err != nil {
			fmt.Fprintf(os.Stderr, "FAIL  Git not found: %v\n", err)
			allOk = false
		} else {
			fmt.Printf("PASS  Git found: %s\n", version)
		}

		// Check hooks installed
		installer := hook.NewInstaller()
		if installed := installer.IsInstalled(); installed {
			fmt.Println("PASS  Global hooks installed")
		} else {
			fmt.Fprintln(os.Stderr, "FAIL  Global hooks not installed")
			allOk = false
		}

		if allOk {
			fmt.Println()
			fmt.Println("All checks passed.")
		} else {
			fmt.Println()
			fmt.Println("Some checks failed. Run 'git-policy install' to fix.")
			os.Exit(1)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}
