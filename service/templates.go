package service

import (
	"github.com/nanoteck137/forge/forgejo"
)

type ListTemplatesInput struct {
	Server string
	Token  string
}

func ListLicenseTemplates(input ListTemplatesInput) ([]forgejo.LicenseTemplate, error) {
	client := forgejo.NewClient(input.Server, input.Token)
	return client.ListLicenseTemplates()
}

func ListGitignoreTemplates(input ListTemplatesInput) ([]string, error) {
	client := forgejo.NewClient(input.Server, input.Token)
	return client.ListGitignoreTemplates()
}
