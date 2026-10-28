package main

import (
	"errors"
	"fmt"
	"path/filepath"
	"sort"

	"github.com/charmbracelet/huh"
	"github.com/nanoteck137/forge/forgejo"
	"github.com/nanoteck137/forge/service"
	"github.com/spf13/cobra"
)

var cloneCmd = &cobra.Command{
	Use:   "clone",
	Short: "Interactively clone one of your repositories",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := resolveConfig(cmd)

		if err := service.ValidateTokens(cfg.Server, cfg.Token, cfg.GitHubToken); err != nil {
			return err
		}

		fmt.Println("Fetching repositories...")

		repos, err := service.ListAllUserRepos(service.ListAllUserReposInput{
			Server: cfg.Server,
			Token:  cfg.Token,
		})
		if err != nil {
			return fmt.Errorf("listing repositories: %w", err)
		}

		if len(repos) == 0 {
			return fmt.Errorf("no repositories found")
		}

		sort.Slice(repos, func(i, j int) bool {
			return repos[i].FullName < repos[j].FullName
		})

		byFullName := make(map[string]forgejo.Repository, len(repos))
		options := make([]huh.Option[string], 0, len(repos))
		for _, repo := range repos {
			byFullName[repo.FullName] = repo
			options = append(options, huh.NewOption(repoLabel(repo), repo.FullName))
		}

		var selected string
		err = huh.NewSelect[string]().
			Title("Clone a repository").
			Description("Type to filter, enter to select").
			Options(options...).
			Filtering(true).
			Height(15).
			Value(&selected).
			Run()
		if err != nil {
			if errors.Is(err, huh.ErrUserAborted) {
				return nil
			}
			return err
		}

		repo, ok := byFullName[selected]
		if !ok {
			return fmt.Errorf("repository not found: %s", selected)
		}

		dest, err := defaultCloneDest(cmd, cfg, repo.Name)
		if err != nil {
			return err
		}

		err = huh.NewInput().
			Title("Clone destination").
			Description("Where should " + repo.FullName + " be cloned?").
			Value(&dest).
			Run()
		if err != nil {
			if errors.Is(err, huh.ErrUserAborted) {
				return nil
			}
			return err
		}

		cloneURL := repo.SSHURL
		if cloneURL == "" {
			cloneURL = repo.CloneURL
		}

		return cloneRepository(cloneURL, repo.FullName, dest)
	},
}

func repoLabel(repo forgejo.Repository) string {
	label := repo.FullName
	if repo.Private {
		label += " (private)"
	}
	if repo.Description != "" {
		desc := repo.Description
		if len(desc) > 60 {
			desc = desc[:57] + "..."
		}
		label += " — " + desc
	}
	return label
}

func defaultCloneDest(cmd *cobra.Command, cfg ResolvedConfig, name string) (string, error) {
	directory, _ := cmd.Flags().GetString("directory")
	currentDir, _ := cmd.Flags().GetBool("current-dir")

	if directory != "" && currentDir {
		return "", fmt.Errorf("cannot use --directory together with --current-dir")
	}

	if directory != "" || currentDir {
		base, err := resolveCloneDir(cfg, directory, currentDir)
		if err != nil {
			return "", err
		}
		return filepath.Join(base, name), nil
	}

	if cfg.CloneDirectory != "" {
		base, err := expandPath(cfg.CloneDirectory)
		if err != nil {
			return "", err
		}
		return filepath.Join(base, name), nil
	}

	return name, nil
}

func init() {
	cloneCmd.Flags().StringP("directory", "d", "", "Base directory to clone into (overrides clone.directory)")
	cloneCmd.Flags().Bool("current-dir", false, "Use the current directory instead of clone.directory")

	rootCmd.AddCommand(cloneCmd)
}
