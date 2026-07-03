package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version = "dev"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Long:  `Print the current version of git-policy.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("git-policy v%s\n", version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
