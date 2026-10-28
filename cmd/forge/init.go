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

		credentialsPath, err := credentials.Path()
		if err != nil {
			return fmt.Errorf("getting credentials path: %w", err)
		}

		_, configErr := os.Stat(configPath)
		configExists := configErr == nil

		_, credsErr := os.Stat(credentialsPath)
		credsExists := credsErr == nil

		reinitConfig := !configExists
		overwriteCreds := !credsExists

		piped := isStdinPiped()
		var scanner *bufio.Scanner
		if piped {
			scanner = bufio.NewScanner(os.Stdin)
		}

		if configExists || credsExists {
			if piped {
				if configExists {
					reinitConfig = promptBool(scanner, "Config already exists, re-initialize it", false)
				}
				if credsExists {
					overwriteCreds = promptBool(scanner, "Credentials already exist, overwrite them", false)
				}
			} else {
				var fields []huh.Field
				if configExists {
					fields = append(fields, huh.NewConfirm().
						Title("Config already exists").
						Description(configPath+"\nRe-initialize config?").
						Affirmative("Re-initialize").
						Negative("Keep").
						Value(&reinitConfig))
				}
				if credsExists {
					fields = append(fields, huh.NewConfirm().
						Title("Credentials already exist").
						Description(credentialsPath+"\nOverwrite credentials?").
						Affirmative("Overwrite").
						Negative("Keep").
						Value(&overwriteCreds))
				}

				if err := huh.NewForm(huh.NewGroup(fields...)).Run(); err != nil {
					if errors.Is(err, huh.ErrUserAborted) {
						fmt.Println("Aborting.")
						return nil
					}
					return err
				}
			}
		}

		if !reinitConfig && !overwriteCreds {
			fmt.Println("Nothing to do.")
			return nil
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

		if piped {
			if reinitConfig {
				server = prompt(scanner, "Server", server)
				cloneDirectory = prompt(scanner, "Clone directory (optional)", cloneDirectory)
				private = promptBool(scanner, "Create private repos by default", private)
				autoInit = promptBool(scanner, "Initialize new repos with a README by default", autoInit)
			}

			if overwriteCreds {
				token = promptSecret(scanner, "Token", token)
				githubToken = promptSecret(scanner, "GitHub Token (optional, for mirroring)", githubToken)
			}

			if reinitConfig {
				license = prompt(scanner, "Default license for new repos (optional)", license)
				gitignore = prompt(scanner, "Default gitignore template for new repos (optional)", gitignore)
			}
		} else {
			if reinitConfig {
				settingsForm := huh.NewForm(
					huh.NewGroup(
						huh.NewInput().
							Title("Forgejo server").
							Value(&server),
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
			}

			if overwriteCreds {
				credentialsForm := huh.NewForm(
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
				)

				if err := credentialsForm.Run(); err != nil {
					if errors.Is(err, huh.ErrUserAborted) {
						fmt.Println("Aborting.")
						return nil
					}
					return err
				}
			}

			if reinitConfig {
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
		}

		if reinitConfig {
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

			fmt.Printf("Config written to %s\n", configPath)
		}

		if overwriteCreds {
			creds := credentials.Store{
				Token:       token,
				GitHubToken: githubToken,
			}

			if err := credentials.Save(&creds); err != nil {
				return fmt.Errorf("saving credentials: %w", err)
			}

			fmt.Printf("Credentials written to %s\n", credentialsPath)
		}

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
