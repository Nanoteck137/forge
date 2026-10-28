package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
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

		if err := service.ValidateTokens(cfg.Server, cfg.Token, cfg.GitHubToken); err != nil {
			return err
		}

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
	Use:   "view [[<owner>/]<repo>]",
	Short: "View repository details",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := resolveConfig(cmd)

		if err := service.ValidateTokens(cfg.Server, cfg.Token, cfg.GitHubToken); err != nil {
			return err
		}

		var owner, repo string

		if len(args) > 0 {
			var err error
			owner, repo, err = splitRepoArg(cfg, args[0])
			if err != nil {
				return err
			}
		} else {
			owner = cfg.Owner
			repo = cfg.Repo
		}

		if owner == "" || repo == "" {
			return fmt.Errorf("repository required: provide a repository argument or run inside a git repository")
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

		if err := service.ValidateTokens(cfg.Server, cfg.Token, cfg.GitHubToken); err != nil {
			return err
		}

		name := args[0]
		desc, _ := cmd.Flags().GetString("description")
		private, _ := cmd.Flags().GetBool("private")
		autoInit, _ := cmd.Flags().GetBool("auto-init")
		license, _ := cmd.Flags().GetString("license")
		gitignore, _ := cmd.Flags().GetString("gitignore")
		org, _ := cmd.Flags().GetString("org")
		mirror, _ := cmd.Flags().GetBool("mirror")

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

		scanner := bufio.NewScanner(os.Stdin)

		if mirror {
			if cfg.GitHubToken == "" {
				return fmt.Errorf("GitHub token not set, run 'forge init' to set up your config")
			}

			githubRepoName := prompt(scanner, "GitHub repository name", repo.Name)

			fmt.Print("GitHub organization (leave blank for personal): ")
			scanner.Scan()
			githubOrg := strings.TrimSpace(scanner.Text())

			fmt.Println("Setting up mirror...")

			mirrorResult, err := service.SetupMirror(service.SetupMirrorParams{
				Server:      cfg.Server,
				Token:       cfg.Token,
				Owner:       strings.SplitN(repo.FullName, "/", 2)[0],
				Repo:        repo.Name,
				GitHubToken: cfg.GitHubToken,
				GitHubRepo:  githubRepoName,
				GitHubOrg:   githubOrg,
			})
			if err != nil {
				return fmt.Errorf("setting up mirror: %w", err)
			}

			fmt.Println("Mirror setup complete!")
			fmt.Printf("  Forgejo: %s\n", mirrorResult.ForgejoRepo)
			fmt.Printf("  GitHub:  %s\n", mirrorResult.GitHubURL)
		}

		fmt.Print("Clone repository? [y/N]: ")
		scanner.Scan()
		answer := strings.ToLower(strings.TrimSpace(scanner.Text()))
		if answer == "y" || answer == "yes" {
			defaultDir := repo.Name
			if cfg.CloneDirectory != "" {
				base, err := expandPath(cfg.CloneDirectory)
				if err != nil {
					return err
				}
				defaultDir = filepath.Join(base, repo.Name)
			}

			dir := prompt(scanner, "Clone directory", defaultDir)

			cloneURL := repo.SSHURL
			if cloneURL == "" {
				cloneURL = repo.CloneURL
			}

			if err := cloneRepository(cloneURL, repo.FullName, dir); err != nil {
				return err
			}
		}

		return nil
	},
}

func expandPath(p string) (string, error) {
	if p != "~" && !strings.HasPrefix(p, "~/") {
		return p, nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("getting home directory: %w", err)
	}

	if p == "~" {
		return homeDir, nil
	}

	return filepath.Join(homeDir, strings.TrimPrefix(p, "~/")), nil
}

func resolveCloneDir(cfg ResolvedConfig, directory string, currentDir bool) (string, error) {
	base := directory
	if currentDir {
		base = "."
	}
	if base == "" {
		base = cfg.CloneDirectory
	}
	if base == "" {
		base = "."
	}

	base, err := expandPath(base)
	if err != nil {
		return "", err
	}

	return filepath.Abs(base)
}

func splitRepoArg(cfg ResolvedConfig, arg string) (string, string, error) {
	if parts := strings.SplitN(arg, "/", 2); len(parts) == 2 {
		return parts[0], parts[1], nil
	}

	user, err := service.GetCurrentUser(service.GetCurrentUserInput{
		Server: cfg.Server,
		Token:  cfg.Token,
	})
	if err != nil {
		return "", "", fmt.Errorf("getting current user: %w", err)
	}

	return user.Login, arg, nil
}

func cloneRepository(cloneURL, fullName, dir string) error {
	if dir != "" {
		expanded, err := expandPath(dir)
		if err != nil {
			return err
		}
		dir = expanded
	}

	gitArgs := []string{"clone", cloneURL}
	if dir != "" {
		if err := os.MkdirAll(filepath.Dir(dir), 0755); err != nil {
			return fmt.Errorf("creating clone directory: %w", err)
		}
		gitArgs = append(gitArgs, dir)
	}

	if dir != "" {
		fmt.Printf("Cloning %s into %s...\n", fullName, dir)
	} else {
		fmt.Printf("Cloning %s...\n", fullName)
	}

	gitCmd := exec.Command("git", gitArgs...)
	gitCmd.Stdout = os.Stdout
	gitCmd.Stderr = os.Stderr

	if err := gitCmd.Run(); err != nil {
		return fmt.Errorf("cloning repo: %w", err)
	}

	fmt.Println("Done.")
	return nil
}

var repoCloneCmd = &cobra.Command{
	Use:   "clone [<owner>/]<repo> [<directory>]",
	Short: "Clone a repository",
	Args:  cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := resolveConfig(cmd)

		if err := service.ValidateTokens(cfg.Server, cfg.Token, cfg.GitHubToken); err != nil {
			return err
		}

		var owner, repo string
		var err error
		owner, repo, err = splitRepoArg(cfg, args[0])
		if err != nil {
			return err
		}

		cloneURL, fullName, err := service.GetCloneURL(service.CloneRepoInput{
			Server: cfg.Server,
			Token:  cfg.Token,
			Owner:  owner,
			Repo:   repo,
		})
		if err != nil {
			return fmt.Errorf("getting repo: %w", err)
		}

		directory, _ := cmd.Flags().GetString("directory")
		currentDir, _ := cmd.Flags().GetBool("current-dir")

		if directory != "" && currentDir {
			return fmt.Errorf("cannot use --directory together with --current-dir")
		}

		var dest string
		if len(args) > 1 {
			if directory != "" || currentDir {
				return fmt.Errorf("cannot use a directory argument together with --directory or --current-dir")
			}

			dest, err = expandPath(args[1])
			if err != nil {
				return err
			}
		} else {
			base, err := resolveCloneDir(cfg, directory, currentDir)
			if err != nil {
				return fmt.Errorf("resolving clone directory: %w", err)
			}
			dest = filepath.Join(base, path.Base(fullName))
		}

		return cloneRepository(cloneURL, fullName, dest)
	},
}

var repoForkCmd = &cobra.Command{
	Use:   "fork [<owner>/]<repo>",
	Short: "Fork a repository",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := resolveConfig(cmd)

		if err := service.ValidateTokens(cfg.Server, cfg.Token, cfg.GitHubToken); err != nil {
			return err
		}

		owner, repo, err := splitRepoArg(cfg, args[0])
		if err != nil {
			return err
		}

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
	repoCreateCmd.Flags().Bool("mirror", false, "Set up GitHub push mirror after creation")

	repoCloneCmd.Flags().StringP("directory", "d", "", "Base directory to clone into (overrides clone.directory)")
	repoCloneCmd.Flags().Bool("current-dir", false, "Use the current directory instead of clone.directory")

	repoForkCmd.Flags().String("org", "", "Fork to organization")

	repoCmd.AddCommand(repoListCmd)
	repoCmd.AddCommand(repoViewCmd)
	repoCmd.AddCommand(repoCreateCmd)
	repoCmd.AddCommand(repoCloneCmd)
	repoCmd.AddCommand(repoForkCmd)

	rootCmd.AddCommand(repoCmd)
}
