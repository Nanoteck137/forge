package service

import (
	"fmt"
	"strings"

	"github.com/nanoteck137/forge/forgejo"
)

type ListIssuesInput struct {
	Server         string
	Token          string
	Owner          string
	Repo           string
	State          string
	Labels         []string
	ExcludeLabels  []string
	Assignee       string
	Milestone      string
	Limit          int
	Page           int
}

type ListIssuesResult struct {
	Issues []forgejo.Issue
}

func ListIssues(input ListIssuesInput) (*ListIssuesResult, error) {
	client := forgejo.NewClient(input.Server, input.Token)

	var labelIDs []int64
	if len(input.Labels) > 0 {
		var err error
		labelIDs, err = client.ResolveLabels(input.Owner, input.Repo, input.Labels)
		if err != nil {
			return nil, fmt.Errorf("resolving labels: %w", err)
		}
	}

	milestoneID, err := client.ResolveMilestone(input.Owner, input.Repo, input.Milestone)
	if err != nil {
		return nil, fmt.Errorf("resolving milestone: %w", err)
	}

	opts := forgejo.ListIssuesOptions{
		State:     input.State,
		Labels:    labelIDs,
		Assignee:  input.Assignee,
		Milestone: milestoneID,
		Limit:     input.Limit,
		Page:      input.Page,
	}

	issuesList, err := client.ListIssues(input.Owner, input.Repo, opts)
	if err != nil {
		return nil, err
	}

	if len(input.ExcludeLabels) > 0 {
		exclude := make(map[string]bool, len(input.ExcludeLabels))
		for _, name := range input.ExcludeLabels {
			exclude[name] = true
		}

		filtered := make([]forgejo.Issue, 0, len(issuesList))
		for _, issue := range issuesList {
			keep := true
			for _, l := range issue.Labels {
				if exclude[l.Name] {
					keep = false
					break
				}
			}
			if keep {
				filtered = append(filtered, issue)
			}
		}
		issuesList = filtered
	}

	return &ListIssuesResult{Issues: issuesList}, nil
}

type CreateIssueInput struct {
	Server     string
	Token      string
	Owner      string
	Repo       string
	Title      string
	Body       string
	Labels     []string
	Milestone  string
	Assignees  []string
}

func CreateIssue(input CreateIssueInput) (*forgejo.Issue, error) {
	client := forgejo.NewClient(input.Server, input.Token)

	labelIDs, err := client.ResolveLabels(input.Owner, input.Repo, input.Labels)
	if err != nil {
		return nil, fmt.Errorf("resolving labels: %w", err)
	}

	milestoneID, err := client.ResolveMilestone(input.Owner, input.Repo, input.Milestone)
	if err != nil {
		return nil, fmt.Errorf("resolving milestone: %w", err)
	}

	body := input.Body
	if body != "" {
		body += "\n\n---\n"
	}
	body += "_Created by forge_"

	req := forgejo.CreateIssueRequest{
		Title:     input.Title,
		Body:      body,
		Assignees: input.Assignees,
		Milestone: milestoneID,
		Labels:    labelIDs,
	}

	return client.CreateIssue(input.Owner, input.Repo, req)
}

type GetIssueInput struct {
	Server string
	Token  string
	Owner  string
	Repo   string
	Number int
}

func GetIssue(input GetIssueInput) (*forgejo.Issue, error) {
	client := forgejo.NewClient(input.Server, input.Token)
	return client.GetIssue(input.Owner, input.Repo, input.Number)
}

type EditIssueInput struct {
	Server string
	Token  string
	Owner  string
	Repo   string
	Number int
	State  string
	Body   string
}

func EditIssue(input EditIssueInput) (*forgejo.Issue, error) {
	client := forgejo.NewClient(input.Server, input.Token)

	req := forgejo.EditIssueRequest{
		State: input.State,
		Body:  input.Body,
	}

	return client.EditIssue(input.Owner, input.Repo, input.Number, req)
}

type FixIssuesInput struct {
	Server string
	Token  string
	Owner  string
	Repo   string
}

type FixIssuesResult struct {
	FixedCount int
}

func FixIssues(input FixIssuesInput) (*FixIssuesResult, error) {
	client := forgejo.NewClient(input.Server, input.Token)

	result := &FixIssuesResult{}
	page := 1
	for {
		opts := forgejo.ListIssuesOptions{
			State: "all",
			Limit: 100,
			Page:  page,
		}

		issuesList, err := client.ListIssues(input.Owner, input.Repo, opts)
		if err != nil {
			return nil, fmt.Errorf("listing issues: %w", err)
		}

		if len(issuesList) == 0 {
			break
		}

		for _, issue := range issuesList {
			if !strings.Contains(issue.Body, "\\n") {
				continue
			}

			newBody := strings.ReplaceAll(issue.Body, "\\n", "\n")

			req := forgejo.EditIssueRequest{Body: newBody}
			_, err := client.EditIssue(input.Owner, input.Repo, issue.Number, req)
			if err != nil {
				continue
			}

			result.FixedCount++
		}

		page++
	}

	return result, nil
}
