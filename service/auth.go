package service

import (
	"fmt"

	"github.com/nanoteck137/forge/credentials"
	"github.com/nanoteck137/forge/forgejo"
)

type AuthStatusResult struct {
	Server       string
	TokenValid   bool
	Username     string
	GitHubToken  bool
	ErrorMessage string
}

func GetAuthStatus(server string) (*AuthStatusResult, error) {
	result := &AuthStatusResult{
		Server: server,
	}

	token, _ := credentials.Read("forgejo_token")
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

	githubToken, _ := credentials.Read("github_token")
	result.GitHubToken = githubToken != ""

	return result, nil
}

func Login(server, token string) error {
	if err := credentials.Write("forgejo_token", token); err != nil {
		return fmt.Errorf("saving token: %w", err)
	}
	return nil
}

func LoginGitHub(token string) error {
	if err := credentials.Write("github_token", token); err != nil {
		return fmt.Errorf("saving token: %w", err)
	}
	return nil
}

func Logout(githubOnly bool) error {
	if githubOnly {
		if err := credentials.Delete("github_token"); err != nil {
			return fmt.Errorf("deleting credential: %w", err)
		}
	} else {
		if err := credentials.Delete("forgejo_token"); err != nil {
			return fmt.Errorf("deleting credential: %w", err)
		}
		if err := credentials.Delete("github_token"); err != nil {
			return fmt.Errorf("deleting credential: %w", err)
		}
	}
	return nil
}
