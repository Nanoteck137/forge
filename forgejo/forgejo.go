package forgejo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type User struct {
	Id        int    `json:"id"`
	Login     string `json:"login"`
	FullName  string `json:"full_name"`
	Email     string `json:"email"`
	Bio       string `json:"bio"`
	Website   string `json:"website"`
	Location  string `json:"location"`
	AvatarURL string `json:"avatar_url"`
	Stars     int    `json:"stars_count"`
	Forks     int    `json:"forks_count"`
	Created   string `json:"created_at"`
	Updated   string `json:"updated_at"`
}

type Repository struct {
	Id            int64  `json:"id"`
	Name          string `json:"name"`
	FullName      string `json:"full_name"`
	Description   string `json:"description"`
	Private       bool   `json:"private"`
	Fork          bool   `json:"fork"`
	URL           string `json:"url"`
	HTMLURL       string `json:"html_url"`
	CloneURL      string `json:"clone_url"`
	SSHURL        string `json:"ssh_url"`
	DefaultBranch string `json:"default_branch"`
	Stars         int    `json:"stargazers_count"`
	Forks         int    `json:"forks_count"`
	OpenIssues    int    `json:"open_issues_count"`
	Watchers      int    `json:"watchers_count"`
	Mirror        bool   `json:"mirror"`
	Created       string `json:"created_at"`
	Updated       string `json:"updated_at"`
}

type Issue struct {
	Number    int        `json:"number"`
	Title     string     `json:"title"`
	State     string     `json:"state"`
	Body      string     `json:"body"`
	HTMLURL   string     `json:"html_url"`
	Labels    []Label    `json:"labels"`
	Milestone *Milestone `json:"milestone"`
	Assignees []User     `json:"assignees"`
}

type Label struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Milestone struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
}

type CreateIssueRequest struct {
	Title     string   `json:"title"`
	Body      string   `json:"body,omitempty"`
	Assignees []string `json:"assignees,omitempty"`
	Milestone int      `json:"milestone,omitempty"`
	Labels    []int64  `json:"labels,omitempty"`
}

type Client struct {
	server string
	token  string
	http   *http.Client
}

func NewClient(server, token string) *Client {
	return &Client{
		server: strings.TrimRight(server, "/"),
		token:  token,
		http: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) newRequest(method, path string, body io.Reader) (*http.Request, error) {
	url := fmt.Sprintf("%s/api/v1%s", c.server, path)
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	if c.token != "" {
		req.Header.Set("Authorization", "token "+c.token)
	}
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

func (c *Client) do(req *http.Request, v any) error {
	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		msg := strings.TrimSpace(string(bodyBytes))
		return fmt.Errorf("API error (%d): %s", resp.StatusCode, msg)
	}

	if v != nil {
		return json.NewDecoder(resp.Body).Decode(v)
	}
	return nil
}

type ListIssuesOptions struct {
	State     string
	Labels    []int64
	Assignee  string
	Milestone int
	Limit     int
	Page      int
}

func (c *Client) GetCurrentUser() (*User, error) {
	httpReq, err := c.newRequest("GET", "/user", nil)
	if err != nil {
		return nil, err
	}

	var user User
	if err := c.do(httpReq, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (c *Client) GetUser(username string) (*User, error) {
	httpReq, err := c.newRequest("GET", fmt.Sprintf("/users/%s", username), nil)
	if err != nil {
		return nil, err
	}

	var user User
	if err := c.do(httpReq, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (c *Client) ListIssues(owner, repo string, opts ListIssuesOptions) ([]Issue, error) {
	path := fmt.Sprintf("/repos/%s/%s/issues", owner, repo)

	params := url.Values{}
	if opts.State != "" {
		params.Set("state", opts.State)
	}
	if len(opts.Labels) > 0 {
		strs := make([]string, len(opts.Labels))
		for i, id := range opts.Labels {
			strs[i] = strconv.FormatInt(id, 10)
		}
		params.Set("labels", strings.Join(strs, ","))
	}
	if opts.Assignee != "" {
		params.Set("assignee", opts.Assignee)
	}
	if opts.Milestone > 0 {
		params.Set("milestone", strconv.Itoa(opts.Milestone))
	}
	if opts.Limit > 0 {
		params.Set("limit", strconv.Itoa(opts.Limit))
	}
	if opts.Page > 0 {
		params.Set("page", strconv.Itoa(opts.Page))
	}

	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	httpReq, err := c.newRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var issues []Issue
	if err := c.do(httpReq, &issues); err != nil {
		return nil, err
	}
	return issues, nil
}

func (c *Client) CreateIssue(owner, repo string, req CreateIssueRequest) (*Issue, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshaling request: %w", err)
	}

	httpReq, err := c.newRequest("POST", fmt.Sprintf("/repos/%s/%s/issues", owner, repo), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	var issue Issue
	if err := c.do(httpReq, &issue); err != nil {
		return nil, err
	}

	return &issue, nil
}

func (c *Client) GetIssue(owner, repo string, number int) (*Issue, error) {
	httpReq, err := c.newRequest("GET", fmt.Sprintf("/repos/%s/%s/issues/%d", owner, repo, number), nil)
	if err != nil {
		return nil, err
	}

	var issue Issue
	if err := c.do(httpReq, &issue); err != nil {
		return nil, err
	}
	return &issue, nil
}

type EditIssueRequest struct {
	State string `json:"state,omitempty"`
	Body  string `json:"body,omitempty"`
}

func (c *Client) EditIssue(owner, repo string, number int, req EditIssueRequest) (*Issue, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshaling request: %w", err)
	}

	httpReq, err := c.newRequest("PATCH", fmt.Sprintf("/repos/%s/%s/issues/%d", owner, repo, number), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	var issue Issue
	if err := c.do(httpReq, &issue); err != nil {
		return nil, err
	}
	return &issue, nil
}

func (c *Client) ListLabels(owner, repo string) ([]Label, error) {
	httpReq, err := c.newRequest("GET", fmt.Sprintf("/repos/%s/%s/labels", owner, repo), nil)
	if err != nil {
		return nil, err
	}

	var labels []Label
	if err := c.do(httpReq, &labels); err != nil {
		return nil, err
	}
	return labels, nil
}

func (c *Client) ListMilestones(owner, repo string) ([]Milestone, error) {
	httpReq, err := c.newRequest("GET", fmt.Sprintf("/repos/%s/%s/milestones", owner, repo), nil)
	if err != nil {
		return nil, err
	}

	var milestones []Milestone
	if err := c.do(httpReq, &milestones); err != nil {
		return nil, err
	}
	return milestones, nil
}

func (c *Client) ResolveLabels(owner, repo string, names []string) ([]int64, error) {
	if len(names) == 0 {
		return nil, nil
	}

	labels, err := c.ListLabels(owner, repo)
	if err != nil {
		return nil, fmt.Errorf("fetching labels: %w", err)
	}

	nameToID := make(map[string]int64, len(labels))
	for _, l := range labels {
		nameToID[l.Name] = int64(l.Id)
	}

	var ids []int64
	for _, name := range names {
		id, ok := nameToID[name]
		if !ok {
			return nil, fmt.Errorf("unknown label: %q", name)
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func (c *Client) ResolveMilestone(owner, repo, name string) (int, error) {
	if name == "" {
		return 0, nil
	}

	milestones, err := c.ListMilestones(owner, repo)
	if err != nil {
		return 0, fmt.Errorf("fetching milestones: %w", err)
	}

	for _, m := range milestones {
		if m.Title == name {
			return m.Id, nil
		}
	}
	return 0, fmt.Errorf("milestone not found: %q", name)
}

// Repository methods

func (c *Client) GetRepository(owner, repo string) (*Repository, error) {
	httpReq, err := c.newRequest("GET", fmt.Sprintf("/repos/%s/%s", owner, repo), nil)
	if err != nil {
		return nil, err
	}

	var repository Repository
	if err := c.do(httpReq, &repository); err != nil {
		return nil, err
	}
	return &repository, nil
}

type ListUserReposOptions struct {
	Page  int
	Limit int
	Order string
}

func (c *Client) ListUserRepos(opts ListUserReposOptions) ([]Repository, error) {
	params := url.Values{}
	if opts.Page > 0 {
		params.Set("page", strconv.Itoa(opts.Page))
	}
	if opts.Limit > 0 {
		params.Set("limit", strconv.Itoa(opts.Limit))
	}
	if opts.Order != "" {
		params.Set("order_by", opts.Order)
	}

	path := "/user/repos"
	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	httpReq, err := c.newRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var repos []Repository
	if err := c.do(httpReq, &repos); err != nil {
		return nil, err
	}
	return repos, nil
}

func (c *Client) ListOrgRepos(org string, page, limit int) ([]Repository, error) {
	params := url.Values{}
	if page > 0 {
		params.Set("page", strconv.Itoa(page))
	}
	if limit > 0 {
		params.Set("limit", strconv.Itoa(limit))
	}

	path := fmt.Sprintf("/orgs/%s/repos", org)
	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	httpReq, err := c.newRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var repos []Repository
	if err := c.do(httpReq, &repos); err != nil {
		return nil, err
	}
	return repos, nil
}

func (c *Client) ListUserReposByUser(username string, page, limit int) ([]Repository, error) {
	params := url.Values{}
	if page > 0 {
		params.Set("page", strconv.Itoa(page))
	}
	if limit > 0 {
		params.Set("limit", strconv.Itoa(limit))
	}

	path := fmt.Sprintf("/users/%s/repos", username)
	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	httpReq, err := c.newRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var repos []Repository
	if err := c.do(httpReq, &repos); err != nil {
		return nil, err
	}
	return repos, nil
}

type CreateRepoRequest struct {
	Name          string `json:"name"`
	Description   string `json:"description,omitempty"`
	Private       bool   `json:"private,omitempty"`
	AutoInit      bool   `json:"auto_init,omitempty"`
	DefaultBranch string `json:"default_branch,omitempty"`
	Gitignores    string `json:"gitignores,omitempty"`
	License       string `json:"license,omitempty"`
	Readme        string `json:"readme,omitempty"`
}

func (c *Client) CreateRepo(req CreateRepoRequest) (*Repository, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshaling request: %w", err)
	}

	httpReq, err := c.newRequest("POST", "/user/repos", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	var repo Repository
	if err := c.do(httpReq, &repo); err != nil {
		return nil, err
	}
	return &repo, nil
}

func (c *Client) CreateOrgRepo(org string, req CreateRepoRequest) (*Repository, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshaling request: %w", err)
	}

	httpReq, err := c.newRequest("POST", fmt.Sprintf("/orgs/%s/repos", org), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	var repo Repository
	if err := c.do(httpReq, &repo); err != nil {
		return nil, err
	}
	return &repo, nil
}

type CreateForkRequest struct {
	Organization string `json:"organization,omitempty"`
}

func (c *Client) ForkRepo(owner, repo string, req CreateForkRequest) (*Repository, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshaling request: %w", err)
	}

	httpReq, err := c.newRequest("POST", fmt.Sprintf("/repos/%s/%s/forks", owner, repo), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	var fork Repository
	if err := c.do(httpReq, &fork); err != nil {
		return nil, err
	}
	return &fork, nil
}

func (c *Client) DeleteRepo(owner, repo string) error {
	httpReq, err := c.newRequest("DELETE", fmt.Sprintf("/repos/%s/%s", owner, repo), nil)
	if err != nil {
		return err
	}

	return c.do(httpReq, nil)
}

// Template methods

type LicenseTemplate struct {
	Key  string `json:"key"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

func (c *Client) ListLicenseTemplates() ([]LicenseTemplate, error) {
	httpReq, err := c.newRequest("GET", "/licenses", nil)
	if err != nil {
		return nil, err
	}

	var templates []LicenseTemplate
	if err := c.do(httpReq, &templates); err != nil {
		return nil, err
	}
	return templates, nil
}

func (c *Client) ListGitignoreTemplates() ([]string, error) {
	httpReq, err := c.newRequest("GET", "/gitignore/templates", nil)
	if err != nil {
		return nil, err
	}

	var templates []string
	if err := c.do(httpReq, &templates); err != nil {
		return nil, err
	}
	return templates, nil
}

// Push Mirror methods

type PushMirror struct {
	ID            int64  `json:"id"`
	RemoteName    string `json:"remote_name"`
	RemoteAddress string `json:"remote_address"`
	RemoteUser    string `json:"remote_username"`
	SyncOnCommit  bool   `json:"sync_on_commit"`
	LastUpdate    string `json:"last_update"`
	Interval      string `json:"interval"`
	BranchFilter  string `json:"branch_filter"`
	PublicKey     string `json:"public_key"`
}

type CreatePushMirrorRequest struct {
	RemoteAddress  string `json:"remote_address"`
	RemoteUsername  string `json:"remote_username"`
	RemotePassword string `json:"remote_password"`
	SyncOnCommit   bool   `json:"sync_on_commit"`
	Interval       string `json:"interval,omitempty"`
	BranchFilter   string `json:"branch_filter,omitempty"`
	UseSSH         bool   `json:"use_ssh,omitempty"`
}

func (c *Client) ListPushMirrors(owner, repo string, page, limit int) ([]PushMirror, error) {
	params := url.Values{}
	if page > 0 {
		params.Set("page", strconv.Itoa(page))
	}
	if limit > 0 {
		params.Set("limit", strconv.Itoa(limit))
	}

	path := fmt.Sprintf("/repos/%s/%s/push_mirrors", owner, repo)
	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	httpReq, err := c.newRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var mirrors []PushMirror
	if err := c.do(httpReq, &mirrors); err != nil {
		return nil, err
	}
	return mirrors, nil
}

func (c *Client) AddPushMirror(owner, repo string, req CreatePushMirrorRequest) (*PushMirror, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshaling request: %w", err)
	}

	httpReq, err := c.newRequest("POST", fmt.Sprintf("/repos/%s/%s/push_mirrors", owner, repo), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	var mirror PushMirror
	if err := c.do(httpReq, &mirror); err != nil {
		return nil, err
	}
	return &mirror, nil
}

func (c *Client) SyncPushMirrors(owner, repo string) error {
	httpReq, err := c.newRequest("POST", fmt.Sprintf("/repos/%s/%s/push_mirrors-sync", owner, repo), nil)
	if err != nil {
		return err
	}

	return c.do(httpReq, nil)
}

func (c *Client) DeletePushMirror(owner, repo, name string) error {
	httpReq, err := c.newRequest("DELETE", fmt.Sprintf("/repos/%s/%s/push_mirrors/%s", owner, repo, name), nil)
	if err != nil {
		return err
	}

	return c.do(httpReq, nil)
}

func (c *Client) GetPushMirror(owner, repo, name string) (*PushMirror, error) {
	httpReq, err := c.newRequest("GET", fmt.Sprintf("/repos/%s/%s/push_mirrors/%s", owner, repo, name), nil)
	if err != nil {
		return nil, err
	}

	var mirror PushMirror
	if err := c.do(httpReq, &mirror); err != nil {
		return nil, err
	}
	return &mirror, nil
}

// User SSH Key methods

type PublicKey struct {
	ID      int    `json:"id"`
	Key     string `json:"key"`
	Title   string `json:"title"`
	Created string `json:"created_at"`
}

func (c *Client) ListUserKeys(page, limit int) ([]PublicKey, error) {
	params := url.Values{}
	if page > 0 {
		params.Set("page", strconv.Itoa(page))
	}
	if limit > 0 {
		params.Set("limit", strconv.Itoa(limit))
	}

	path := "/user/keys"
	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	httpReq, err := c.newRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var keys []PublicKey
	if err := c.do(httpReq, &keys); err != nil {
		return nil, err
	}
	return keys, nil
}

type CreateKeyRequest struct {
	Title string `json:"title"`
	Key   string `json:"key"`
}

func (c *Client) AddUserKey(req CreateKeyRequest) (*PublicKey, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshaling request: %w", err)
	}

	httpReq, err := c.newRequest("POST", "/user/keys", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	var key PublicKey
	if err := c.do(httpReq, &key); err != nil {
		return nil, err
	}
	return &key, nil
}

func (c *Client) DeleteUserKey(id int) error {
	httpReq, err := c.newRequest("DELETE", fmt.Sprintf("/user/keys/%d", id), nil)
	if err != nil {
		return err
	}

	return c.do(httpReq, nil)
}
