package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/nanoteck137/forge/service"
	"github.com/spf13/cobra"
)

var mirrorCmd = &cobra.Command{
	Use:   "mirror",
	Short: "Manage push mirrors",
}

var mirrorSetupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Set up a push mirror from Forgejo to GitHub",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := resolveConfig(cmd)

		if err := service.ValidateTokens(cfg.Server, cfg.Token, cfg.GitHubToken); err != nil {
			return err
		}

		scanner := bufio.NewScanner(os.Stdin)

		owner := cfg.Owner
		repo := cfg.Repo

		if owner == "" || repo == "" {
			fmt.Print("Forgejo repository (owner/repo): ")
			scanner.Scan()
			input := strings.TrimSpace(scanner.Text())
			parts := strings.SplitN(input, "/", 2)
			if len(parts) != 2 {
				return fmt.Errorf("invalid format, expected owner/repo")
			}
			owner = parts[0]
			repo = parts[1]
		} else {
			fmt.Printf("Forgejo repository: %s/%s\n", owner, repo)
		}

		fmt.Print("GitHub repository name (leave blank to use same name): ")
		scanner.Scan()
		githubRepoName := strings.TrimSpace(scanner.Text())
		if githubRepoName == "" {
			githubRepoName = repo
		}

		fmt.Print("GitHub organization (leave blank for personal): ")
		scanner.Scan()
		githubOrg := strings.TrimSpace(scanner.Text())

		result, err := service.SetupMirror(service.SetupMirrorParams{
			Server:      cfg.Server,
			Token:       cfg.Token,
			Owner:       owner,
			Repo:        repo,
			GitHubToken: cfg.GitHubToken,
			GitHubRepo:  githubRepoName,
			GitHubOrg:   githubOrg,
		})
		if err != nil {
			return err
		}

		fmt.Println("\nMirror setup complete!")
		fmt.Printf("  Forgejo: %s\n", result.ForgejoRepo)
		fmt.Printf("  GitHub:  %s\n", result.GitHubURL)
		fmt.Println("\nTo trigger an initial sync, run:")
		fmt.Printf("  forge mirror sync %s/%s\n", owner, repo)

		return nil
	},
}

var mirrorSyncCmd = &cobra.Command{
	Use:   "sync <owner>/<repo>",
	Short: "Manually trigger a push mirror sync",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := resolveConfig(cmd)

		if err := service.ValidateTokens(cfg.Server, cfg.Token, cfg.GitHubToken); err != nil {
			return err
		}

		parts := strings.SplitN(args[0], "/", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid format, expected owner/repo")
		}
		owner := parts[0]
		repo := parts[1]

		fmt.Printf("Syncing push mirrors for %s/%s...\n", owner, repo)

		if err := service.SyncMirror(service.SyncMirrorInput{
			Server: cfg.Server,
			Token:  cfg.Token,
			Owner:  owner,
			Repo:   repo,
		}); err != nil {
			return fmt.Errorf("syncing mirrors: %w", err)
		}

		fmt.Println("Sync triggered successfully.")
		return nil
	},
}

var mirrorListCmd = &cobra.Command{
	Use:   "list",
	Short: "List push mirrors for a repository",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := resolveConfig(cmd)

		if err := service.ValidateTokens(cfg.Server, cfg.Token, cfg.GitHubToken); err != nil {
			return err
		}

		owner := cfg.Owner
		repo := cfg.Repo

		if owner == "" || repo == "" {
			return fmt.Errorf("repository required: provide --owner and --repo flags or run inside a git repository")
		}

		mirrors, err := service.ListMirrors(service.ListMirrorsInput{
			Server: cfg.Server,
			Token:  cfg.Token,
			Owner:  owner,
			Repo:   repo,
			Page:   1,
			Limit:  50,
		})
		if err != nil {
			return fmt.Errorf("listing mirrors: %w", err)
		}

		if cfg.JSON {
			fmt.Println(toJSON(mirrors))
			return nil
		}

		if len(mirrors) == 0 {
			fmt.Println("No push mirrors configured.")
			return nil
		}

		printTable("Remote", "Address", "Sync On Commit", "Last Update")
		for _, m := range mirrors {
			sync := "no"
			if m.SyncOnCommit {
				sync = "yes"
			}
			printRow(m.RemoteName, m.RemoteAddress, sync, formatTime(m.LastUpdate))
		}

		return nil
	},
}

var mirrorRemoveCmd = &cobra.Command{
	Use:   "remove <owner>/<repo> <remote-name>",
	Short: "Remove a push mirror",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := resolveConfig(cmd)

		if err := service.ValidateTokens(cfg.Server, cfg.Token, cfg.GitHubToken); err != nil {
			return err
		}

		parts := strings.SplitN(args[0], "/", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid format, expected owner/repo")
		}
		owner := parts[0]
		repo := parts[1]
		remoteName := args[1]

		yes, _ := cmd.Flags().GetBool("yes")
		if !yes {
			fmt.Printf("Remove push mirror '%s' from %s/%s? [y/N]: ", remoteName, owner, repo)
			var answer string
			fmt.Scanln(&answer)
			if strings.ToLower(answer) != "y" {
				fmt.Println("Aborted.")
				return nil
			}
		}

		if err := service.RemoveMirror(service.RemoveMirrorInput{
			Server:     cfg.Server,
			Token:      cfg.Token,
			Owner:      owner,
			Repo:       repo,
			RemoteName: remoteName,
		}); err != nil {
			return fmt.Errorf("removing mirror: %w", err)
		}

		fmt.Println("Push mirror removed.")
		return nil
	},
}

func init() {
	mirrorListCmd.Flags().Int("limit", 50, "Max results")
	mirrorListCmd.Flags().Int("page", 1, "Page number")

	mirrorRemoveCmd.Flags().BoolP("yes", "y", false, "Skip confirmation")

	mirrorCmd.AddCommand(mirrorSetupCmd)
	mirrorCmd.AddCommand(mirrorSyncCmd)
	mirrorCmd.AddCommand(mirrorListCmd)
	mirrorCmd.AddCommand(mirrorRemoveCmd)

	rootCmd.AddCommand(mirrorCmd)
}
