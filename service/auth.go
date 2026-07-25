package service

import (
	"fmt"

	"github.com/nanoteck137/forge/forgejo"
	"github.com/nanoteck137/forge/github"
)

type AuthStatusResult struct {
	Server       string
	TokenValid   bool
	Username     string
	GitHubToken  bool
	ErrorMessage string
}

func GetAuthStatus(server, token, githubToken string) (*AuthStatusResult, error) {
	result := &AuthStatusResult{
		Server:      server,
		GitHubToken: githubToken != "",
	}

	if token != "" {
		client := forgejo.NewClient(server, token)
		user, err := client.GetCurrentUser()
		if err != nil {
			result.ErrorMessage = err.Error()
		} else {
			result.TokenValid = true
			result.Username = user.Login
		}
	}

	return result, nil
}

func ValidateTokens(server, token, githubToken string) error {
	if token != "" {
		client := forgejo.NewClient(server, token)
		user, err := client.GetCurrentUser()
		if err != nil {
			return fmt.Errorf("forgejo token invalid: %w", err)
		}
		if user.FullName != "" {
			fmt.Printf("Forgejo: authenticated as %s (%s)\n", user.Login, user.FullName)
		} else {
			fmt.Printf("Forgejo: authenticated as %s\n", user.Login)
		}
	} else {
		fmt.Println("Forgejo: not configured")
	}

	if githubToken != "" {
		client := github.NewGitHubClient(githubToken)
		user, err := client.GetCurrentUser()
		if err != nil {
			return fmt.Errorf("github token invalid: %w", err)
		}
		fmt.Printf("GitHub:  authenticated as %s\n", user.Login)
	} else {
		fmt.Println("GitHub:  not configured")
	}

	return nil
}
