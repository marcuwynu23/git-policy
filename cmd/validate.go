package cmd

import (
	"fmt"
	"os"

	"github.com/marcuwynu23/git-policy/internal/config"
	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate the configuration",
	Long:  `Check the git-policy configuration file for errors and correctness.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		_, err := config.Load(cfgFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "FAIL  Invalid configuration: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("PASS  Configuration is valid.")
		if cfgFile == "" {
			fmt.Println("       Using default configuration.")
		} else {
			fmt.Printf("       File: %s\n", cfgFile)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)
}
