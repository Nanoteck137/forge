package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/nanoteck137/forge/service"
	"github.com/spf13/cobra"
)

var reopenCmd = &cobra.Command{
	Use:   "reopen <issue-number>",
	Short: "Reopen a closed issue",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := resolveConfig(cmd)

		number, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid issue number: %s", args[0])
		}

		if cfg.Repo == "" {
			return fmt.Errorf("--repo is required")
		}

		issue, err := service.EditIssue(service.EditIssueInput{
			Server: cfg.Server,
			Token:  cfg.Token,
			Owner:  cfg.Owner,
			Repo:   cfg.Repo,
			Number: number,
			State:  "open",
		})
		if err != nil {
			return fmt.Errorf("reopening issue: %w", err)
		}

		if cfg.JSON {
			out := struct {
				Number int    `json:"number"`
				Title  string `json:"title"`
				State  string `json:"state"`
				URL    string `json:"url"`
			}{
				Number: issue.Number,
				Title:  issue.Title,
				State:  issue.State,
				URL:    issue.HTMLURL,
			}
			return json.NewEncoder(os.Stdout).Encode(out)
		}

		fmt.Printf("✓ Reopened issue #%d\n\n%s\n\n%s\n", issue.Number, issue.Title, issue.HTMLURL)
		return nil
	},
}

func init() {
	reopenCmd.Flags().StringVarP(&repoFlag, "repo", "r", gitInfo.Repo, "Repository name")

	issuesCmd.AddCommand(reopenCmd)
}
