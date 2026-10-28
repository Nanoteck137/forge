package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/nanoteck137/forge/config"
	"github.com/nanoteck137/forge/credentials"
	"github.com/spf13/cobra"
)

var initOutput string

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		configPath := initOutput
		if configPath == "" {
			configDir, err := config.DefaultConfigDir()
			if err != nil {
				return fmt.Errorf("getting config directory: %w", err)
			}
			configPath = filepath.Join(configDir, "config.toml")
		}

		if _, err := os.Stat(configPath); err == nil {
			fmt.Print("Config already exists. Overwrite? [y/N]: ")
			var answer string
			fmt.Scanln(&answer)
			if strings.ToLower(answer) != "y" {
				fmt.Println("Aborting.")
				return nil
			}
		}

		if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
			return fmt.Errorf("creating config directory: %w", err)
		}

		existingCreds := loadCredentials()
		existingCfg := loadConfig()

		scanner := bufio.NewScanner(os.Stdin)

		server := prompt(scanner, "Server", "https://forgejo.nanoteck137.net")
		cloneDirectory := prompt(scanner, "Clone directory (optional)", existingCfg.Clone.Directory)
		license := prompt(scanner, "Default license for new repos (optional)", existingCfg.Create.License)
		gitignore := prompt(scanner, "Default gitignore template for new repos (optional)", existingCfg.Create.Gitignore)
		private := promptBool(scanner, "Create private repos by default", existingCfg.Create.Private)
		autoInit := promptBool(scanner, "Initialize new repos with a README by default", existingCfg.Create.AutoInit)
		token := promptSecret(scanner, "Token", existingCreds.Token)
		githubToken := promptSecret(scanner, "GitHub Token (optional, for mirroring)", existingCreds.GitHubToken)

		cfg := config.Config{
			Server: server,
			Clone: config.ConfigClone{
				Directory: cloneDirectory,
			},
			Create: config.ConfigCreate{
				License:   license,
				Gitignore: gitignore,
				Private:   private,
				AutoInit:  autoInit,
			},
		}

		if err := config.Save(configPath, &cfg); err != nil {
			return fmt.Errorf("saving config: %w", err)
		}

		creds := credentials.Store{
			Token:       token,
			GitHubToken: githubToken,
		}

		if err := credentials.Save(&creds); err != nil {
			return fmt.Errorf("saving credentials: %w", err)
		}

		credentialsPath, err := credentials.Path()
		if err != nil {
			return fmt.Errorf("getting credentials path: %w", err)
		}

		fmt.Printf("Config written to %s\n", configPath)
		fmt.Printf("Credentials written to %s\n", credentialsPath)
		return nil
	},
}

func prompt(scanner *bufio.Scanner, label, defaultVal string) string {
	if defaultVal != "" {
		fmt.Printf("%s [%s]: ", label, defaultVal)
	} else {
		fmt.Printf("%s: ", label)
	}
	scanner.Scan()
	input := strings.TrimSpace(scanner.Text())
	if input == "" {
		return defaultVal
	}
	return input
}

func promptSecret(scanner *bufio.Scanner, label, current string) string {
	if current != "" {
		fmt.Printf("%s [leave blank to keep existing]: ", label)
	} else {
		fmt.Printf("%s: ", label)
	}
	scanner.Scan()
	input := strings.TrimSpace(scanner.Text())
	if input == "" {
		return current
	}
	return input
}

func promptBool(scanner *bufio.Scanner, label string, defaultVal bool) bool {
	if defaultVal {
		fmt.Printf("%s [Y/n]: ", label)
	} else {
		fmt.Printf("%s [y/N]: ", label)
	}

	scanner.Scan()
	input := strings.ToLower(strings.TrimSpace(scanner.Text()))
	if input == "" {
		return defaultVal
	}

	return input == "y" || input == "yes"
}

func init() {
	initCmd.Flags().StringVarP(&initOutput, "output", "o", "", "Output path for the config file")

	rootCmd.AddCommand(initCmd)
}
