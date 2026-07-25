package main

import (
	"fmt"

	"github.com/nanoteck137/forge/service"
	"github.com/spf13/cobra"
)

var fixCmd = &cobra.Command{
	Use:   "fix",
	Short: "Fix issues with literal \\n in their body",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := resolveConfig(cmd)

		if cfg.Repo == "" {
			return fmt.Errorf("--repo is required")
		}

		result, err := service.FixIssues(service.FixIssuesInput{
			Server: cfg.Server,
			Token:  cfg.Token,
			Owner:  cfg.Owner,
			Repo:   cfg.Repo,
		})
		if err != nil {
			return err
		}

		fmt.Printf("Done. Fixed %d issues.\n", result.FixedCount)
		return nil
	},
}

func init() {
	fixCmd.Flags().StringVarP(&repoFlag, "repo", "r", gitInfo.Repo, "Repository name")

	issuesCmd.AddCommand(fixCmd)
}
