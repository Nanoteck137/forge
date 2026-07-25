package forge

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type GitHubClient struct {
	token  string
	client *http.Client
}

func NewGitHubClient(token string) *GitHubClient {
	return &GitHubClient{
		token: token,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *GitHubClient) newRequest(method, path string, body io.Reader) (*http.Request, error) {
	url := fmt.Sprintf("https://api.github.com%s", path)
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	if c.token != "" {
		req.Header.Set("Authorization", "token "+c.token)
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	return req, nil
}

func (c *GitHubClient) do(req *http.Request, v any) error {
	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		msg := strings.TrimSpace(string(bodyBytes))
		return fmt.Errorf("GitHub API error (%d): %s", resp.StatusCode, msg)
	}

	if v != nil {
		return json.NewDecoder(resp.Body).Decode(v)
	}
	return nil
}

type GitHubRepository struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	HTMLURL  string `json:"html_url"`
	CloneURL string `json:"clone_url"`
	SSHURL   string `json:"ssh_url"`
	Private  bool   `json:"private"`
}

type CreateGitHubRepoRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Private     bool   `json:"private,omitempty"`
	AutoInit    bool   `json:"auto_init,omitempty"`
}

func (c *GitHubClient) CreateRepo(req CreateGitHubRepoRequest) (*GitHubRepository, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshaling request: %w", err)
	}

	httpReq, err := c.newRequest("POST", "/user/repos", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	var repo GitHubRepository
	if err := c.do(httpReq, &repo); err != nil {
		return nil, err
	}
	return &repo, nil
}

func (c *GitHubClient) GetRepo(owner, name string) (*GitHubRepository, error) {
	httpReq, err := c.newRequest("GET", fmt.Sprintf("/repos/%s/%s", owner, name), nil)
	if err != nil {
		return nil, err
	}

	var repo GitHubRepository
	if err := c.do(httpReq, &repo); err != nil {
		return nil, err
	}
	return &repo, nil
}

type GitHubDeployKey struct {
	ID       int64  `json:"id"`
	Key      string `json:"key"`
	Title    string `json:"title"`
	ReadOnly bool   `json:"read_only"`
}

type CreateDeployKeyRequest struct {
	Title    string `json:"title"`
	Key      string `json:"key"`
	ReadOnly bool   `json:"read_only"`
}

func (c *GitHubClient) AddDeployKey(owner, repo string, req CreateDeployKeyRequest) (*GitHubDeployKey, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshaling request: %w", err)
	}

	httpReq, err := c.newRequest("POST", fmt.Sprintf("/repos/%s/%s/keys", owner, repo), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	var key GitHubDeployKey
	if err := c.do(httpReq, &key); err != nil {
		return nil, err
	}
	return &key, nil
}

func (c *GitHubClient) ListDeployKeys(owner, repo string) ([]GitHubDeployKey, error) {
	httpReq, err := c.newRequest("GET", fmt.Sprintf("/repos/%s/%s/keys", owner, repo), nil)
	if err != nil {
		return nil, err
	}

	var keys []GitHubDeployKey
	if err := c.do(httpReq, &keys); err != nil {
		return nil, err
	}
	return keys, nil
}

func (c *GitHubClient) DeleteDeployKey(owner, repo string, id int64) error {
	httpReq, err := c.newRequest("DELETE", fmt.Sprintf("/repos/%s/%s/keys/%d", owner, repo, id), nil)
	if err != nil {
		return err
	}

	return c.do(httpReq, nil)
}

type GitHubUser struct {
	Login string `json:"login"`
}

func (c *GitHubClient) GetCurrentUser() (*GitHubUser, error) {
	httpReq, err := c.newRequest("GET", "/user", nil)
	if err != nil {
		return nil, err
	}

	var user GitHubUser
	if err := c.do(httpReq, &user); err != nil {
		return nil, err
	}
	return &user, nil
}
