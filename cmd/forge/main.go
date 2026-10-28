package main

import (
	"os"

	"github.com/kr/pretty"
	"github.com/nanoteck137/forge"
	"github.com/nanoteck137/forge/config"
	"github.com/nanoteck137/forge/credentials"
	"github.com/spf13/cobra"
)

var (
	serverFlag      string
	ownerFlag       string
	tokenFlag       string
	githubTokenFlag string
	jsonFlag        bool
	verboseFlag     bool
	configFlag      string

	gitInfo = forge.DetectFromGit()
)

type ResolvedConfig struct {
	Server      string
	Owner       string
	Token       string
	GitHubToken string
	Repo        string
	JSON        bool
	Verbose     bool
}

var rootCmd = &cobra.Command{
	Use:     forge.AppName,
	Version: forge.Version,
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
}

func init() {
	rootCmd.SetVersionTemplate(forge.VersionTemplate(forge.AppName))

	rootCmd.PersistentFlags().StringVar(&configFlag, "config", "", "Config file path")
	rootCmd.PersistentFlags().StringVar(&serverFlag, "server", "", "Forgejo server URL")
	rootCmd.PersistentFlags().StringVar(&ownerFlag, "owner", "", "Repository owner")
	rootCmd.PersistentFlags().StringVar(&tokenFlag, "token", "", "Forgejo authentication token")
	rootCmd.PersistentFlags().StringVar(&githubTokenFlag, "github-token", "", "GitHub authentication token")
	rootCmd.PersistentFlags().BoolVar(&jsonFlag, "json", false, "Output in JSON format")
	rootCmd.PersistentFlags().BoolVarP(&verboseFlag, "verbose", "v", false, "Verbose output")
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		pretty.Println(err)
		os.Exit(1)
	}
}

func loadConfig() *config.Config {
	cfg, err := config.Load(configFlag)
	if err != nil {
		return &config.Config{
			Server: "https://forgejo.nanoteck137.net",
		}
	}
	return cfg
}

func loadCredentials() *credentials.Store {
	store, err := credentials.Load()
	if err != nil {
		return &credentials.Store{}
	}
	return store
}

func resolveConfig(cmd *cobra.Command) ResolvedConfig {
	cfg := loadConfig()
	creds := loadCredentials()

	server := serverFlag
	if !cmd.Flags().Changed("server") {
		if cfg.Server != "" {
			server = cfg.Server
		} else if gitInfo.Server != "" {
			server = gitInfo.Server
		} else {
			server = "https://forgejo.nanoteck137.net"
		}
	}

	owner := ownerFlag
	if !cmd.Flags().Changed("owner") && gitInfo.Owner != "" {
		owner = gitInfo.Owner
	}

	token := creds.Token
	if envToken := os.Getenv("FORGE_TOKEN"); envToken != "" {
		token = envToken
	}
	if cmd.Flags().Changed("token") {
		token = tokenFlag
	}

	githubToken := creds.GitHubToken
	if envToken := os.Getenv("FORGE_GITHUB_TOKEN"); envToken != "" {
		githubToken = envToken
	}
	if cmd.Flags().Changed("github-token") {
		githubToken = githubTokenFlag
	}

	repo := ""
	if cmd.Flags().Lookup("repo") != nil {
		repo = cmd.Flags().Lookup("repo").Value.String()
		if !cmd.Flags().Changed("repo") && gitInfo.Repo != "" {
			repo = gitInfo.Repo
		}
	}

	return ResolvedConfig{
		Server:      server,
		Owner:       owner,
		Token:       token,
		GitHubToken: githubToken,
		Repo:        repo,
		JSON:        jsonFlag,
		Verbose:     verboseFlag,
	}
}
