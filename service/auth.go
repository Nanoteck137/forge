package service

import (
	"github.com/nanoteck137/forge/forgejo"
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
