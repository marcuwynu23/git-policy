package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "git-policy",
	Short: "Global Git policy management for developers and teams",
	Long: `Git Policy is a cross-platform CLI that provides global Git policy management.

Install once and protect every repository on your machine with policies
for blocking secrets, enforcing commit conventions, protecting branches,
and more.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default ~/.config/git-policy/config.yaml)")
}

func initConfig() {
}
