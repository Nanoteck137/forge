package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/nanoteck137/forge/service"
	"github.com/spf13/cobra"
)

var (
	listState        string
	listLabels       []string
	listExcludeLabels []string
	listAssignee     string
	listMilestone    string
	listLimit        int
	listPage         int
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List issues",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := resolveConfig(cmd)

		if cfg.Repo == "" {
			return fmt.Errorf("--repo is required")
		}

		result, err := service.ListIssues(service.ListIssuesInput{
			Server:        cfg.Server,
			Token:         cfg.Token,
			Owner:         cfg.Owner,
			Repo:          cfg.Repo,
			State:         listState,
			Labels:        listLabels,
			ExcludeLabels: listExcludeLabels,
			Assignee:      listAssignee,
			Milestone:     listMilestone,
			Limit:         listLimit,
			Page:          listPage,
		})
		if err != nil {
			return err
		}

		issuesList := result.Issues

		if cfg.JSON {
			type issueOut struct {
				Number int    `json:"number"`
				Title  string `json:"title"`
				State  string `json:"state"`
				URL    string `json:"url"`
			}
			out := make([]issueOut, len(issuesList))
			for i, issue := range issuesList {
				out[i] = issueOut{
					Number: issue.Number,
					Title:  issue.Title,
					State:  issue.State,
					URL:    issue.HTMLURL,
				}
			}
			return json.NewEncoder(os.Stdout).Encode(out)
		}

		if len(issuesList) == 0 {
			fmt.Println("No issues found")
			return nil
		}

		for _, issue := range issuesList {
			labels := ""
			for i, l := range issue.Labels {
				if i > 0 {
					labels += ", "
				}
				labels += l.Name
			}
			if labels != "" {
				labels = "  [" + labels + "]"
			}
			fmt.Printf("#%-5d %-50s %s%s\n", issue.Number, issue.Title, issue.State, labels)
		}

		return nil
	},
}

func init() {
	listCmd.Flags().StringVarP(&repoFlag, "repo", "r", gitInfo.Repo, "Repository name")
	listCmd.Flags().StringVarP(&listState, "state", "s", "", "Filter by state (open, closed, all)")
	listCmd.Flags().StringSliceVarP(&listLabels, "label", "l", nil, "Filter by labels (by name)")
	listCmd.Flags().StringSliceVar(&listExcludeLabels, "exclude-label", nil, "Exclude issues with these labels")
	listCmd.Flags().StringVar(&listAssignee, "assignee", "", "Filter by assignee")
	listCmd.Flags().StringVarP(&listMilestone, "milestone", "m", "", "Filter by milestone (by name)")
	listCmd.Flags().IntVar(&listLimit, "limit", 0, "Maximum number of issues")
	listCmd.Flags().IntVar(&listPage, "page", 0, "Page number")

	issuesCmd.AddCommand(listCmd)
}
