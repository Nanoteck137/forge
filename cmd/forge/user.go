package main

import (
	"fmt"

	"github.com/nanoteck137/forge/service"
	"github.com/spf13/cobra"
)

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Manage users",
}

var userInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show current user info",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := resolveConfig(cmd)

		user, err := service.GetCurrentUser(service.GetCurrentUserInput{
			Server: cfg.Server,
			Token:  cfg.Token,
		})
		if err != nil {
			return fmt.Errorf("getting user: %w", err)
		}

		if cfg.JSON {
			fmt.Println(toJSON(user))
			return nil
		}

		fmt.Printf("Username:    %s\n", user.Login)
		fmt.Printf("Full Name:   %s\n", user.FullName)
		fmt.Printf("Email:       %s\n", user.Email)
		fmt.Printf("Bio:         %s\n", user.Bio)
		fmt.Printf("Website:     %s\n", user.Website)
		fmt.Printf("Location:    %s\n", user.Location)
		fmt.Printf("Stars:       %d\n", user.Stars)
		fmt.Printf("Forks:       %d\n", user.Forks)
		fmt.Printf("Created:     %s\n", formatTime(user.Created))
		fmt.Printf("Updated:     %s\n", formatTime(user.Updated))

		return nil
	},
}

var userViewCmd = &cobra.Command{
	Use:   "view <username>",
	Short: "View a user's profile",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := resolveConfig(cmd)

		user, err := service.GetUser(service.GetUserInput{
			Server:   cfg.Server,
			Token:    cfg.Token,
			Username: args[0],
		})
		if err != nil {
			return fmt.Errorf("getting user: %w", err)
		}

		if cfg.JSON {
			fmt.Println(toJSON(user))
			return nil
		}

		fmt.Printf("Username:    %s\n", user.Login)
		fmt.Printf("Full Name:   %s\n", user.FullName)
		fmt.Printf("Email:       %s\n", user.Email)
		fmt.Printf("Bio:         %s\n", user.Bio)
		fmt.Printf("Website:     %s\n", user.Website)
		fmt.Printf("Location:    %s\n", user.Location)
		fmt.Printf("Stars:       %d\n", user.Stars)
		fmt.Printf("Forks:       %d\n", user.Forks)
		fmt.Printf("Created:     %s\n", formatTime(user.Created))
		fmt.Printf("Updated:     %s\n", formatTime(user.Updated))

		return nil
	},
}

func init() {
	userCmd.AddCommand(userInfoCmd)
	userCmd.AddCommand(userViewCmd)

	rootCmd.AddCommand(userCmd)
}
