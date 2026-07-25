package service

import (
	"fmt"

	"github.com/nanoteck137/forge/credentials"
	"github.com/nanoteck137/forge/forgejo"
	"github.com/nanoteck137/forge/github"
)

type SetupMirrorInput struct {
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

func SetupMirror(input SetupMirrorInput) (*SetupMirrorResult, error) {
	if input.Token == "" {
		return nil, fmt.Errorf("not authenticated, run 'forge auth login' first")
	}

	if input.GitHubToken == "" {
		var err error
		input.GitHubToken, err = credentials.Read("github_token")
		if err != nil || input.GitHubToken == "" {
			return nil, fmt.Errorf("GitHub token not set, run 'forge auth login github' first")
		}
	}

	forgeClient := forgejo.NewClient(input.Server, input.Token)
	githubClient := github.NewGitHubClient(input.GitHubToken)

	_, err := forgeClient.GetRepository(input.Owner, input.Repo)
	if err != nil {
		return nil, fmt.Errorf("getting Forgejo repo: %w", err)
	}

	githubOwner := input.GitHubOrg
	if githubOwner == "" {
		ghUser, err := githubClient.GetCurrentUser()
		if err != nil {
			return nil, fmt.Errorf("getting GitHub user: %w", err)
		}
		githubOwner = ghUser.Login
	}

	githubRepo, err := githubClient.GetRepo(githubOwner, input.GitHubRepo)
	if err != nil {
		createReq := github.CreateGitHubRepoRequest{
			Name:    input.GitHubRepo,
			Private: true,
		}

		githubRepo, err = githubClient.CreateRepo(createReq)
		if err != nil {
			return nil, fmt.Errorf("creating GitHub repo: %w", err)
		}
	}

	mirror, err := forgeClient.AddPushMirror(input.Owner, input.Repo, forgejo.CreatePushMirrorRequest{
		RemoteAddress: githubRepo.SSHURL,
		SyncOnCommit:  true,
		Interval:      "1h0m0s",
		UseSSH:        true,
	})
	if err != nil {
		return nil, fmt.Errorf("creating push mirror: %w", err)
	}

	mirrorDetails, err := forgeClient.GetPushMirror(input.Owner, input.Repo, mirror.RemoteName)
	if err != nil {
		return nil, fmt.Errorf("getting mirror details: %w", err)
	}

	if mirrorDetails.PublicKey == "" {
		return nil, fmt.Errorf("no public key returned from Forgejo, please check Forgejo version")
	}

	_, err = githubClient.AddDeployKey(githubOwner, input.GitHubRepo, github.CreateDeployKeyRequest{
		Title:    fmt.Sprintf("forge-mirror-%s/%s", input.Owner, input.Repo),
		Key:      mirrorDetails.PublicKey,
		ReadOnly: false,
	})
	if err != nil {
		return nil, fmt.Errorf("adding GitHub deploy key: %w", err)
	}

	return &SetupMirrorResult{
		ForgejoRepo: fmt.Sprintf("%s/%s", input.Owner, input.Repo),
		GitHubRepo:  fmt.Sprintf("%s/%s", githubOwner, input.GitHubRepo),
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
