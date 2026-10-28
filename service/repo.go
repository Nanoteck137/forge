package service

import (
	"github.com/nanoteck137/forge/forgejo"
)

type ListReposInput struct {
	Server   string
	Token    string
	Username string
	Org      string
	Page     int
	Limit    int
}

func ListRepos(input ListReposInput) ([]forgejo.Repository, error) {
	client := forgejo.NewClient(input.Server, input.Token)

	if input.Org != "" {
		return client.ListOrgRepos(input.Org, input.Page, input.Limit)
	} else if input.Username != "" {
		return client.ListUserReposByUser(input.Username, input.Page, input.Limit)
	} else {
		return client.ListUserRepos(forgejo.ListUserReposOptions{
			Page:  input.Page,
			Limit: input.Limit,
		})
	}
}

type GetRepoInput struct {
	Server string
	Token  string
	Owner  string
	Repo   string
}

func GetRepo(input GetRepoInput) (*forgejo.Repository, error) {
	client := forgejo.NewClient(input.Server, input.Token)
	return client.GetRepository(input.Owner, input.Repo)
}

type CreateRepoInput struct {
	Server      string
	Token       string
	Name        string
	Description string
	Private     bool
	AutoInit    bool
	License     string
	Gitignores  string
	Org         string
}

func CreateRepo(input CreateRepoInput) (*forgejo.Repository, error) {
	client := forgejo.NewClient(input.Server, input.Token)

	req := forgejo.CreateRepoRequest{
		Name:        input.Name,
		Description: input.Description,
		Private:     input.Private,
		AutoInit:    input.AutoInit,
		License:     input.License,
		Gitignores:  input.Gitignores,
	}

	if input.Org != "" {
		return client.CreateOrgRepo(input.Org, req)
	}
	return client.CreateRepo(req)
}

type CloneRepoInput struct {
	Server string
	Token  string
	Owner  string
	Repo   string
}

func GetCloneURL(input CloneRepoInput) (string, string, error) {
	client := forgejo.NewClient(input.Server, input.Token)
	repository, err := client.GetRepository(input.Owner, input.Repo)
	if err != nil {
		return "", "", err
	}
	return repository.SSHURL, repository.FullName, nil
}

type ForkRepoInput struct {
	Server string
	Token  string
	Owner  string
	Repo   string
	Org    string
}

func ForkRepo(input ForkRepoInput) (*forgejo.Repository, error) {
	client := forgejo.NewClient(input.Server, input.Token)
	return client.ForkRepo(input.Owner, input.Repo, forgejo.CreateForkRequest{
		Organization: input.Org,
	})
}

type ListAllUserReposInput struct {
	Server string
	Token  string
}

func ListAllUserRepos(input ListAllUserReposInput) ([]forgejo.Repository, error) {
	client := forgejo.NewClient(input.Server, input.Token)

	const pageSize = 50

	var all []forgejo.Repository
	for page := 1; ; page++ {
		repos, err := client.ListUserRepos(forgejo.ListUserReposOptions{
			Page:  page,
			Limit: pageSize,
		})
		if err != nil {
			return nil, err
		}

		all = append(all, repos...)

		if len(repos) < pageSize {
			break
		}
	}

	return all, nil
}
