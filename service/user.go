package service

import (
	"github.com/nanoteck137/forge/forgejo"
)

type GetCurrentUserInput struct {
	Server string
	Token  string
}

func GetCurrentUser(input GetCurrentUserInput) (*forgejo.User, error) {
	client := forgejo.NewClient(input.Server, input.Token)
	return client.GetCurrentUser()
}

type GetUserInput struct {
	Server   string
	Token    string
	Username string
}

func GetUser(input GetUserInput) (*forgejo.User, error) {
	client := forgejo.NewClient(input.Server, input.Token)
	return client.GetUser(input.Username)
}
