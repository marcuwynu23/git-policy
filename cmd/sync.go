package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync policies from a remote source",
	Long: `Synchronize policies from a remote provider.

Supports Git repositories, HTTP endpoints, S3, and GCS as
policy sources.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Sync command is not yet implemented.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
}
