package main

import (
	"fmt"

	"github.com/nanoteck137/forge/service"
	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Manage authentication",
}

var authStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show current authentication status",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := loadConfig()

		result, err := service.GetAuthStatus(cfg.Server, cfg.Token, cfg.GitHubToken)
		if err != nil {
			return err
		}

		fmt.Printf("Server: %s\n", result.Server)

		if result.TokenValid {
			fmt.Printf("Token: valid (user: %s)\n", result.Username)
		} else if result.ErrorMessage != "" {
			fmt.Printf("Token: invalid (error: %s)\n", result.ErrorMessage)
		} else {
			fmt.Println("Token: not set")
		}

		if result.GitHubToken {
			fmt.Println("GitHub Token: set")
		} else {
			fmt.Println("GitHub Token: not set")
		}

		return nil
	},
}

func init() {
	authCmd.AddCommand(authStatusCmd)

	rootCmd.AddCommand(authCmd)
}
