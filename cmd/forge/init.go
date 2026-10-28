package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/huh"
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
			overwrite := false
			if isStdinPiped() {
				fmt.Print("Config already exists. Overwrite? [y/N]: ")
				var answer string
				fmt.Scanln(&answer)
				overwrite = strings.ToLower(answer) == "y"
			} else {
				if err := huh.NewConfirm().
					Title("Config already exists").
					Description(configPath).
					Affirmative("Overwrite").
					Negative("Abort").
					Value(&overwrite).
					Run(); err != nil {
					if errors.Is(err, huh.ErrUserAborted) {
						fmt.Println("Aborting.")
						return nil
					}
					return err
				}
			}

			if !overwrite {
				fmt.Println("Aborting.")
				return nil
			}
		}

		if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
			return fmt.Errorf("creating config directory: %w", err)
		}

		existingCreds := loadCredentials()
		existingCfg := loadConfig()

		server := existingCfg.Server
		if server == "" {
			server = "https://forgejo.nanoteck137.net"
		}
		cloneDirectory := existingCfg.Clone.Directory
		license := existingCfg.Create.License
		gitignore := existingCfg.Create.Gitignore
		private := existingCfg.Create.Private
		autoInit := existingCfg.Create.AutoInit
		token := existingCreds.Token
		githubToken := existingCreds.GitHubToken

		if isStdinPiped() {
			scanner := bufio.NewScanner(os.Stdin)

			server = prompt(scanner, "Server", server)
			token = promptSecret(scanner, "Token", token)
			githubToken = promptSecret(scanner, "GitHub Token (optional, for mirroring)", githubToken)
			cloneDirectory = prompt(scanner, "Clone directory (optional)", cloneDirectory)
			license = prompt(scanner, "Default license for new repos (optional)", license)
			gitignore = prompt(scanner, "Default gitignore template for new repos (optional)", gitignore)
			private = promptBool(scanner, "Create private repos by default", private)
			autoInit = promptBool(scanner, "Initialize new repos with a README by default", autoInit)
		} else {
			settingsForm := huh.NewForm(
				huh.NewGroup(
					huh.NewInput().
						Title("Forgejo server").
						Value(&server),
				),
				huh.NewGroup(
					huh.NewInput().
						Title("Forgejo token").
						EchoMode(huh.EchoModePassword).
						Value(&token),
				),
				huh.NewGroup(
					huh.NewInput().
						Title("GitHub token (optional, for mirroring)").
						EchoMode(huh.EchoModePassword).
						Value(&githubToken),
				),
				huh.NewGroup(
					huh.NewInput().
						Title("Clone directory (optional)").
						Value(&cloneDirectory),
				),
				huh.NewGroup(
					huh.NewConfirm().
						Title("Create private repos by default?").
						Value(&private),
				),
				huh.NewGroup(
					huh.NewConfirm().
						Title("Initialize new repos with a README?").
						Value(&autoInit),
				),
			)

			if err := settingsForm.Run(); err != nil {
				if errors.Is(err, huh.ErrUserAborted) {
					fmt.Println("Aborting.")
					return nil
				}
				return err
			}

			templatesForm := huh.NewForm(
				huh.NewGroup(
					licenseSelectField(server, token, "Default license", &license),
				),
				huh.NewGroup(
					gitignoreSelectField(server, token, "Default gitignore template", &gitignore),
				),
			)

			if err := templatesForm.Run(); err != nil {
				if errors.Is(err, huh.ErrUserAborted) {
					fmt.Println("Aborting.")
					return nil
				}
				return err
			}
		}

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
