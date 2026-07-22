package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/marcuwynu23/git-policy/internal/config"
	"github.com/marcuwynu23/git-policy/internal/history"
)

func init() {
	historyCmd := &cobra.Command{
		Use:   "history",
		Short: "Show git-policy run history",
		Long:  `Show history of git-policy runs`,
		Run: func(cmd *cobra.Command, args []string) {
			limit, _ := cmd.Flags().GetInt("limit")
			repo, _ := cmd.Flags().GetString("repo")
			status, _ := cmd.Flags().GetString("status")
			clear, _ := cmd.Flags().GetBool("clear")

			cfgPath, err := config.DefaultConfigPath()
			if err != nil {
				cfgPath = ""
			}
			cfg, _ := config.Load(cfgPath)
			if clear {
				_ = history.Clear(cfg, cfgPath, repo)
				fmt.Println("History cleared.")
				return
			}
			records, err := history.Query(cfg, cfgPath, history.QueryOptions{
				Limit:    limit,
				RepoPath: repo,
				Status:   status,
			})
			if err != nil {
				fmt.Printf("Error querying history: %v\n", err)
				return
			}
			for _, r := range records {
				fmt.Printf("%s - %s @ %s (%s) - %s\n", r.Timestamp, r.Repo, r.Branch, r.Commit, r.Overall)
				for _, res := range r.Results {
					fmt.Printf("  - %s: %s", res.Rule, res.Status)
					if res.Message != "" {
						fmt.Printf(" - %s", res.Message)
					}
					fmt.Println()
				}
			}
		},
	}
	historyCmd.Flags().IntP("limit", "l", 20, "Limit number of history entries")
	historyCmd.Flags().StringP("repo", "r", "", "Filter by repository path (default current repo)")
	historyCmd.Flags().StringP("status", "s", "", "Filter by overall status (pass/fail)")
	historyCmd.Flags().BoolP("clear", "c", false, "Clear history for repository")
	rootCmd.AddCommand(historyCmd)
}
