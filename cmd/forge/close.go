package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/nanoteck137/forge/service"
	"github.com/spf13/cobra"
)

var closeCmd = &cobra.Command{
	Use:   "close <issue-number>",
	Short: "Close an issue",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := resolveConfig(cmd)

		if err := service.ValidateTokens(cfg.Server, cfg.Token, cfg.GitHubToken); err != nil {
			return err
		}

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
			State:  "closed",
		})
		if err != nil {
			return fmt.Errorf("closing issue: %w", err)
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

		fmt.Printf("✓ Closed issue #%d\n\n%s\n\n%s\n", issue.Number, issue.Title, issue.HTMLURL)
		return nil
	},
}

func init() {
	closeCmd.Flags().StringVarP(&repoFlag, "repo", "r", gitInfo.Repo, "Repository name")

	issuesCmd.AddCommand(closeCmd)
}
