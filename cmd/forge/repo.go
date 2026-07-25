package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/nanoteck137/forge/service"
	"github.com/spf13/cobra"
)

var repoCmd = &cobra.Command{
	Use:   "repo",
	Short: "Manage repositories",
}

var repoListCmd = &cobra.Command{
	Use:   "list",
	Short: "List repositories",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := resolveConfig(cmd)

		username, _ := cmd.Flags().GetString("user")
		org, _ := cmd.Flags().GetString("org")
		limit, _ := cmd.Flags().GetInt("limit")
		page, _ := cmd.Flags().GetInt("page")

		repos, err := service.ListRepos(service.ListReposInput{
			Server:   cfg.Server,
			Token:    cfg.Token,
			Username: username,
			Org:      org,
			Page:     page,
			Limit:    limit,
		})
		if err != nil {
			return fmt.Errorf("listing repos: %w", err)
		}

		if cfg.JSON {
			fmt.Println(toJSON(repos))
			return nil
		}

		if len(repos) == 0 {
			fmt.Println("No repositories found.")
			return nil
		}

		printTable("Name", "Description", "Stars", "Forks", "Updated")
		for _, r := range repos {
			desc := r.Description
			if len(desc) > 50 {
				desc = desc[:47] + "..."
			}
			printRow(r.Name, desc, strconv.Itoa(r.Stars), strconv.Itoa(r.Forks), formatTime(r.Updated))
		}

		return nil
	},
}

var repoViewCmd = &cobra.Command{
	Use:   "view [<owner>/<repo>]",
	Short: "View repository details",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := resolveConfig(cmd)

		var owner, repo string

		if len(args) > 0 {
			parts := strings.SplitN(args[0], "/", 2)
			if len(parts) != 2 {
				return fmt.Errorf("invalid repository format, expected owner/repo")
			}
			owner = parts[0]
			repo = parts[1]
		} else {
			owner = cfg.Owner
			repo = cfg.Repo
		}

		if owner == "" || repo == "" {
			return fmt.Errorf("repository required: provide owner/repo argument or run inside a git repository")
		}

		repository, err := service.GetRepo(service.GetRepoInput{
			Server: cfg.Server,
			Token:  cfg.Token,
			Owner:  owner,
			Repo:   repo,
		})
		if err != nil {
			return fmt.Errorf("getting repo: %w", err)
		}

		if cfg.JSON {
			fmt.Println(toJSON(repository))
			return nil
		}

		fmt.Printf("Name:        %s\n", repository.FullName)
		fmt.Printf("Description: %s\n", repository.Description)
		fmt.Printf("URL:         %s\n", repository.HTMLURL)
		fmt.Printf("Clone URL:   %s\n", repository.CloneURL)
		fmt.Printf("SSH URL:     %s\n", repository.SSHURL)
		fmt.Printf("Default:     %s\n", repository.DefaultBranch)
		fmt.Printf("Private:     %v\n", repository.Private)
		fmt.Printf("Fork:        %v\n", repository.Fork)
		fmt.Printf("Stars:       %d\n", repository.Stars)
		fmt.Printf("Forks:       %d\n", repository.Forks)
		fmt.Printf("Open Issues: %d\n", repository.OpenIssues)
		fmt.Printf("Created:     %s\n", formatTime(repository.Created))
		fmt.Printf("Updated:     %s\n", formatTime(repository.Updated))

		return nil
	},
}

var repoCreateCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "Create a new repository",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := resolveConfig(cmd)

		name := args[0]
		desc, _ := cmd.Flags().GetString("description")
		private, _ := cmd.Flags().GetBool("private")
		autoInit, _ := cmd.Flags().GetBool("auto-init")
		license, _ := cmd.Flags().GetString("license")
		gitignore, _ := cmd.Flags().GetString("gitignore")
		org, _ := cmd.Flags().GetString("org")

		repo, err := service.CreateRepo(service.CreateRepoInput{
			Server:      cfg.Server,
			Token:       cfg.Token,
			Name:        name,
			Description: desc,
			Private:     private,
			AutoInit:    autoInit,
			License:     license,
			Gitignores:  gitignore,
			Org:         org,
		})
		if err != nil {
			return fmt.Errorf("creating repo: %w", err)
		}

		if cfg.JSON {
			fmt.Println(toJSON(repo))
			return nil
		}

		fmt.Printf("Created %s (%s)\n", repo.FullName, repo.HTMLURL)
		return nil
	},
}

var repoCloneCmd = &cobra.Command{
	Use:   "clone <owner>/<repo> [<directory>]",
	Short: "Clone a repository",
	Args:  cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := resolveConfig(cmd)

		parts := strings.SplitN(args[0], "/", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid repository format, expected owner/repo")
		}
		owner := parts[0]
		repo := parts[1]

		cloneURL, fullName, err := service.GetCloneURL(service.CloneRepoInput{
			Server: cfg.Server,
			Token:  cfg.Token,
			Owner:  owner,
			Repo:   repo,
		})
		if err != nil {
			return fmt.Errorf("getting repo: %w", err)
		}

		dir := ""
		if len(args) > 1 {
			dir = args[1]
		}

		gitArgs := []string{"clone", cloneURL}
		if dir != "" {
			gitArgs = append(gitArgs, dir)
		}

		fmt.Printf("Cloning %s...\n", fullName)

		gitCmd := exec.Command("git", gitArgs...)
		gitCmd.Stdout = os.Stdout
		gitCmd.Stderr = os.Stderr

		if err := gitCmd.Run(); err != nil {
			return fmt.Errorf("cloning repo: %w", err)
		}

		fmt.Println("Done.")
		return nil
	},
}

var repoForkCmd = &cobra.Command{
	Use:   "fork <owner>/<repo>",
	Short: "Fork a repository",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := resolveConfig(cmd)

		parts := strings.SplitN(args[0], "/", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid repository format, expected owner/repo")
		}
		owner := parts[0]
		repo := parts[1]

		org, _ := cmd.Flags().GetString("org")

		fork, err := service.ForkRepo(service.ForkRepoInput{
			Server: cfg.Server,
			Token:  cfg.Token,
			Owner:  owner,
			Repo:   repo,
			Org:    org,
		})
		if err != nil {
			return fmt.Errorf("forking repo: %w", err)
		}

		if cfg.JSON {
			fmt.Println(toJSON(fork))
			return nil
		}

		fmt.Printf("Forked %s/%s to %s (%s)\n", owner, repo, fork.FullName, fork.HTMLURL)
		return nil
	},
}

func init() {
	repoListCmd.Flags().StringP("user", "u", "", "List user's repos")
	repoListCmd.Flags().StringP("org", "o", "", "List org's repos")
	repoListCmd.Flags().Int("limit", 30, "Max results")
	repoListCmd.Flags().Int("page", 1, "Page number")

	repoCreateCmd.Flags().StringP("description", "d", "", "Repository description")
	repoCreateCmd.Flags().BoolP("private", "p", false, "Make repository private")
	repoCreateCmd.Flags().Bool("auto-init", false, "Initialize with README")
	repoCreateCmd.Flags().String("license", "", "License template")
	repoCreateCmd.Flags().String("gitignore", "", "Gitignore template")
	repoCreateCmd.Flags().String("org", "", "Create in organization")

	repoForkCmd.Flags().String("org", "", "Fork to organization")

	repoCmd.AddCommand(repoListCmd)
	repoCmd.AddCommand(repoViewCmd)
	repoCmd.AddCommand(repoCreateCmd)
	repoCmd.AddCommand(repoCloneCmd)
	repoCmd.AddCommand(repoForkCmd)

	rootCmd.AddCommand(repoCmd)
}
