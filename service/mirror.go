package service

import (
	"fmt"

	"github.com/nanoteck137/forge/forgejo"
	"github.com/nanoteck137/forge/github"
)

type SetupMirrorParams struct {
	Server      string
	Token       string
	Owner       string
	Repo        string
	GitHubToken string
	GitHubRepo  string
	GitHubOrg   string
}

type SetupMirrorResult struct {
	ForgejoRepo string
	GitHubRepo  string
	GitHubURL   string
}

func SetupMirror(params SetupMirrorParams) (*SetupMirrorResult, error) {
	if params.Token == "" {
		return nil, fmt.Errorf("not authenticated, run 'forge init' to set up your config")
	}

	if params.GitHubToken == "" {
		return nil, fmt.Errorf("GitHub token not set, run 'forge init' to set up your config")
	}

	forgeClient := forgejo.NewClient(params.Server, params.Token)
	githubClient := github.NewGitHubClient(params.GitHubToken)

	_, err := forgeClient.GetRepository(params.Owner, params.Repo)
	if err != nil {
		return nil, fmt.Errorf("getting Forgejo repo: %w", err)
	}

	githubOwner := params.GitHubOrg
	if githubOwner == "" {
		ghUser, err := githubClient.GetCurrentUser()
		if err != nil {
			return nil, fmt.Errorf("getting GitHub user: %w", err)
		}
		githubOwner = ghUser.Login
	}

	githubRepo, err := githubClient.GetRepo(githubOwner, params.GitHubRepo)
	if err != nil {
		createReq := github.CreateGitHubRepoRequest{
			Name:    params.GitHubRepo,
			Private: true,
		}

		githubRepo, err = githubClient.CreateRepo(createReq)
		if err != nil {
			return nil, fmt.Errorf("creating GitHub repo: %w", err)
		}
	}

	mirror, err := forgeClient.AddPushMirror(params.Owner, params.Repo, forgejo.CreatePushMirrorRequest{
		RemoteAddress: githubRepo.SSHURL,
		SyncOnCommit:  true,
		Interval:      "1h0m0s",
		UseSSH:        true,
	})
	if err != nil {
		return nil, fmt.Errorf("creating push mirror: %w", err)
	}

	mirrorDetails, err := forgeClient.GetPushMirror(params.Owner, params.Repo, mirror.RemoteName)
	if err != nil {
		return nil, fmt.Errorf("getting mirror details: %w", err)
	}

	if mirrorDetails.PublicKey == "" {
		return nil, fmt.Errorf("no public key returned from Forgejo, please check Forgejo version")
	}

	_, err = githubClient.AddDeployKey(githubOwner, params.GitHubRepo, github.CreateDeployKeyRequest{
		Title:    fmt.Sprintf("forge-mirror-%s/%s", params.Owner, params.Repo),
		Key:      mirrorDetails.PublicKey,
		ReadOnly: false,
	})
	if err != nil {
		return nil, fmt.Errorf("adding GitHub deploy key: %w", err)
	}

	return &SetupMirrorResult{
		ForgejoRepo: fmt.Sprintf("%s/%s", params.Owner, params.Repo),
		GitHubRepo:  fmt.Sprintf("%s/%s", githubOwner, params.GitHubRepo),
		GitHubURL:   githubRepo.HTMLURL,
	}, nil
}

type SyncMirrorInput struct {
	Server string
	Token  string
	Owner  string
	Repo   string
}

func SyncMirror(input SyncMirrorInput) error {
	client := forgejo.NewClient(input.Server, input.Token)
	return client.SyncPushMirrors(input.Owner, input.Repo)
}

type ListMirrorsInput struct {
	Server string
	Token  string
	Owner  string
	Repo   string
	Page   int
	Limit  int
}

func ListMirrors(input ListMirrorsInput) ([]forgejo.PushMirror, error) {
	client := forgejo.NewClient(input.Server, input.Token)
	return client.ListPushMirrors(input.Owner, input.Repo, input.Page, input.Limit)
}

type RemoveMirrorInput struct {
	Server     string
	Token      string
	Owner      string
	Repo       string
	RemoteName string
}

func RemoveMirror(input RemoveMirrorInput) error {
	client := forgejo.NewClient(input.Server, input.Token)
	return client.DeletePushMirror(input.Owner, input.Repo, input.RemoteName)
}
