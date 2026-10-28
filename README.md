# forge

A CLI tool for interacting with Forgejo instances from the terminal.

## Features

### Authentication
- [x] `forge auth login` - Authenticate with a Forgejo server (interactive token prompt)
- [x] `forge auth login github` - Authenticate with GitHub (for mirroring features)
- [x] `forge auth status` - Show current authentication status
- [x] `forge auth logout` - Remove stored credentials

### Repositories
- [x] `forge repo list` - List repositories for a user/org
- [x] `forge repo view [name]` - View repository details
- [x] `forge repo create` - Create a new repository
- [x] `forge repo clone [name]` - Clone a repository
- [x] `forge repo fork [name]` - Fork a repository

### Issues
- [x] `forge issues list` - List issues with filters
- [x] `forge issues view <number>` - View issue details
- [x] `forge issues create` - Create a new issue
- [x] `forge issues close <number>` - Close an issue
- [x] `forge issues reopen <number>` - Reopen an issue
- [ ] `forge issues comment <number>` - Add a comment to an issue
- [ ] `forge issues comments <number>` - List comments on an issue

### Pull Requests
- [ ] `forge pr list` - List pull requests
- [ ] `forge pr view <number>` - View PR details
- [ ] `forge pr create` - Create a new pull request
- [ ] `forge pr merge <number>` - Merge a pull request
- [ ] `forge pr close <number>` - Close a pull request
- [ ] `forge pr reopen <number>` - Reopen a pull request
- [ ] `forge pr diff <number>` - View PR diff

### Labels
- [ ] `forge label list` - List labels in a repository
- [ ] `forge label create` - Create a new label
- [ ] `forge label delete <name>` - Delete a label
- [ ] `forge label edit <name>` - Edit a label

### Milestones
- [ ] `forge milestone list` - List milestones
- [ ] `forge milestone view <title>` - View milestone details
- [ ] `forge milestone create` - Create a new milestone
- [ ] `forge milestone close <title>` - Close a milestone
- [ ] `forge milestone delete <title>` - Delete a milestone

### Releases
- [ ] `forge release list` - List releases
- [ ] `forge release view <tag>` - View release details
- [ ] `forge release create` - Create a new release
- [ ] `forge release delete <tag>` - Delete a release

### Notifications
- [ ] `forge notifications list` - List unread notifications
- [ ] `forge notifications read` - Mark notifications as read

### Deploy Keys
- [ ] `forge key list` - List deploy keys for a repository
- [ ] `forge key add <title> <key>` - Add a deploy key to a repository
- [ ] `forge key remove <id>` - Remove a deploy key from a repository

### User
- [x] `forge user info` - Show current user info
- [x] `forge user view <username>` - View a user's profile

### Mirroring
- [x] `forge mirror setup` - Link a Forgejo repo to a GitHub repo
  - Creates a push mirror on the Forgejo repo to GitHub
  - Generates an SSH keypair for the mirror
  - Adds the public key as a deploy key on GitHub (read/write access)
  - Adds the private key to Forgejo as a login key
  - Uses `POST /repos/{owner}/{repo}/push_mirrors` API
- [x] `forge mirror sync <repo>` - Manually trigger a push mirror sync
- [x] `forge mirror list` - List push mirrors for a repo
- [x] `forge mirror remove <name>` - Remove a push mirror configuration

## Configuration

Program configuration is loaded from (highest to lowest priority):
1. Command-line flags (`--server`, etc.)
2. Environment variables (`FORGE_SERVER`)
3. Config files: `./config.toml`, `~/.config/forge/config.toml`, `/etc/forge/config.toml`

Credentials are not stored in the config file. They are stored separately in:
- `$XDG_DATA_HOME/forge/credentials.toml` (defaults to `~/.local/share/forge/credentials.toml`)

Tokens can also be provided with `--token` / `--github-token` or the `FORGE_TOKEN` / `FORGE_GITHUB_TOKEN` environment variables.

Run `forge init` to generate the config and credentials interactively.

## Installation

```sh
# Build from source
go install github.com/nanoteck137/forge/cmd/forge@latest
```

## License

MIT

---

## Implementation Details

### 1. Config (`config/config.go`)

Config only stores program configuration:

```go
type Config struct {
    Server string `mapstructure:"server"`
}
```

Credentials are stored in `credentials/credentials.go`:

```go
type Store struct {
    Token       string `toml:"token"`
    GitHubToken string `toml:"github_token"`
}
```

Credentials are saved to `$XDG_DATA_HOME/forge/credentials.toml` (or `~/.local/share/forge/credentials.toml`).

### 2. Forgejo Client (`client.go`)

#### New Types

```go
// PullRequest
type PullRequest struct {
    ID       int64  `json:"id"`
    Number   int    `json:"number"`
    Title    string `json:"title"`
    Body     string `json:"body"`
    State    string `json:"state"`
    User     User   `json:"user"`
    Created  string `json:"created_at"`
    Updated  string `json:"updated_at"`
    Merged   bool   `json:"merged"`
    DiffURL  string `json:"diff_url"`
    HTMLURL  string `json:"html_url"`
}

type CreatePullRequestOption struct {
    Title     string `json:"title"`
    Body      string `json:"body"`
    Head      string `json:"head"`
    Base      string `json:"base"`
    Assignee  string `json:"assignee,omitempty"`
    Labels    []int  `json:"labels,omitempty"`
    Milestone int    `json:"milestone,omitempty"`
}

type MergePullRequestOption struct {
    MergeMethod string `json:"merge_method,omitempty"` // merge, rebase, squash
    HeadCommitID string `json:"head_commit_id,omitempty"`
}

// Comment
type Comment struct {
    ID       int64  `json:"id"`
    Body     string `json:"body"`
    User     User   `json:"user"`
    Created  string `json:"created_at"`
    Updated  string `json:"updated_at"`
}

type CreateCommentOption struct {
    Body string `json:"body"`
}

// Repository (extended)
type Repository struct {
    ID            int64  `json:"id"`
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
    Mirror        bool   `json:"mirror"`
    MirrorInterval string `json:"mirror_interval"`
    Created       string `json:"created_at"`
    Updated       string `json:"updated_at"`
}

type CreateRepoOption struct {
    Name          string `json:"name"`
    Description   string `json:"description,omitempty"`
    Private       bool   `json:"private,omitempty"`
    AutoInit      bool   `json:"auto_init,omitempty"`
    DefaultBranch string `json:"default_branch,omitempty"`
    Gitignores    string `json:"gitignores,omitempty"`
    License       string `json:"license,omitempty"`
    Readme        string `json:"readme,omitempty"`
}

type CreateForkOption struct {
    Organization string `json:"organization,omitempty"`
    Name         string `json:"name,omitempty"`
}

// PushMirror
type PushMirror struct {
    ID             int64  `json:"id"`
    RemoteName     string `json:"remote_name"`
    RemoteAddress  string `json:"remote_address"`
    RemoteUsername  string `json:"remote_username"`
    SyncOnCommit   bool   `json:"sync_on_commit"`
    LastUpdate     string `json:"last_update"`
    Interval       string `json:"interval"`
    BranchFilter   string `json:"branch_filter"`
}

type CreatePushMirrorOption struct {
    RemoteAddress  string `json:"remote_address"`
    RemoteUsername  string `json:"remote_username"`
    RemotePassword string `json:"remote_password"`
    SyncOnCommit   bool   `json:"sync_on_commit"`
    Interval       string `json:"interval,omitempty"`
    BranchFilter   string `json:"branch_filter,omitempty"`
    UseSSH         bool   `json:"use_ssh,omitempty"`
}

// DeployKey
type DeployKey struct {
    ID       int64  `json:"id"`
    Title    string `json:"title"`
    Key      string `json:"key"`
    ReadOnly bool   `json:"read_only"`
    Created  string `json:"created_at"`
}

type CreateKeyOption struct {
    Title    string `json:"title"`
    Key      string `json:"key"`
    ReadOnly bool   `json:"read_only,omitempty"`
}

// Milestone (extended)
type Milestone struct {
    ID           int64  `json:"id"`
    Title        string `json:"title"`
    Description  string `json:"description"`
    State        string `json:"state"` // open, closed
    OpenIssues   int    `json:"open_issues"`
    ClosedIssues int    `json:"closed_issues"`
    Created      string `json:"created_at"`
    Updated      string `json:"updated_at"`
    Closed       string `json:"closed_at"`
}

type CreateMilestoneOption struct {
    Title       string `json:"title"`
    Description string `json:"description,omitempty"`
    DueDate     string `json:"due_date,omitempty"`
}

// Release
type Release struct {
    ID       int64  `json:"id"`
    TagName  string `json:"tag_name"`
    Name     string `json:"name"`
    Body     string `json:"body"`
    Draft    bool   `json:"draft"`
    Prerelease bool `json:"prerelease"`
    Created  string `json:"created_at"`
    Assets   []ReleaseAsset `json:"assets"`
}

type CreateReleaseOption struct {
    TagName    string `json:"tag_name"`
    Name       string `json:"name,omitempty"`
    Body       string `json:"body,omitempty"`
    Draft      bool   `json:"draft,omitempty"`
    Prerelease bool   `json:"prerelease,omitempty"`
    TargetCommitish string `json:"target_commitish,omitempty"`
}

type ReleaseAsset struct {
    ID   int64  `json:"id"`
    Name string `json:"name"`
    Size int64  `json:"size"`
    URL  string `json:"browser_download_url"`
}

// Notification
type NotificationThread struct {
    ID       int64  `json:"id"`
    Status   string `json:"status"`
    Subject  NotificationSubject `json:"subject"`
    Repo     Repository `json:"repository"`
    Created  string `json:"created_at"`
    Updated  string `json:"updated_at"`
}

type NotificationSubject struct {
    Title            string `json:"title"`
    URL              string `json:"url"`
    LatestCommentURL string `json:"latest_comment_url"`
    Type             string `json:"type"`
}

// User (extended)
type User struct {
    ID        int64  `json:"id"`
    Login     string `json:"login"`
    FullName  string `json:"full_name"`
    Email     string `json:"email"`
    AvatarURL string `json:"avatar_url"`
    Bio       string `json:"bio"`
    Website   string `json:"website"`
    Location  string `json:"location"`
    Created   string `json:"created_at"`
}
```

#### New Client Methods

```go
// Repositories
func (c *Client) ListUserRepos(username string, page, limit int) ([]Repository, error)
func (c *Client) ListOrgRepos(org string, page, limit int) ([]Repository, error)
func (c *Client) CreateRepo(req CreateRepoOption) (*Repository, error)
func (c *Client) CreateOrgRepo(org string, req CreateRepoOption) (*Repository, error)
func (c *Client) ForkRepo(owner, repo string, req CreateForkOption) (*Repository, error)
func (c *Client) DeleteRepo(owner, repo string) error

// Pull Requests
func (c *Client) ListPullRequests(owner, repo string, opts ListPullRequestsOptions) ([]PullRequest, error)
func (c *Client) GetPullRequest(owner, repo string, number int) (*PullRequest, error)
func (c *Client) CreatePullRequest(owner, repo string, req CreatePullRequestOption) (*PullRequest, error)
func (c *Client) MergePullRequest(owner, repo string, number int, req MergePullRequestOption) error
func (c *Client) EditPullRequest(owner, repo string, number int, req EditPullRequestOption) (*PullRequest, error)

// Comments
func (c *Client) ListIssueComments(owner, repo string, number int) ([]Comment, error)
func (c *Client) CreateIssueComment(owner, repo string, number int, req CreateCommentOption) (*Comment, error)

// Labels
func (c *Client) CreateLabel(owner, repo string, req CreateLabelOption) (*Label, error)
func (c *Client) DeleteLabel(owner, repo string, id int) error
func (c *Client) EditLabel(owner, repo string, id int, req EditLabelOption) (*Label, error)

// Milestones
func (c *Client) CreateMilestone(owner, repo string, req CreateMilestoneOption) (*Milestone, error)
func (c *Client) GetMilestone(owner, repo string, id int) (*Milestone, error)
func (c *Client) EditMilestone(owner, repo string, id int, req EditMilestoneOption) (*Milestone, error)
func (c *Client) DeleteMilestone(owner, repo string, id int) error

// Releases
func (c *Client) ListReleases(owner, repo string, page, limit int) ([]Release, error)
func (c *Client) GetRelease(owner, repo string, id int) (*Release, error)
func (c *Client) GetReleaseByTag(owner, repo, tag string) (*Release, error)
func (c *Client) CreateRelease(owner, repo string, req CreateReleaseOption) (*Release, error)
func (c *Client) DeleteRelease(owner, repo string, id int) error
func (c *Client) DeleteReleaseByTag(owner, repo, tag string) error

// Deploy Keys
func (c *Client) ListDeployKeys(owner, repo string, page, limit int) ([]DeployKey, error)
func (c *Client) AddDeployKey(owner, repo string, req CreateKeyOption) (*DeployKey, error)
func (c *Client) DeleteDeployKey(owner, repo string, id int) error

// Notifications
func (c *Client) ListNotifications(opts ListNotificationsOptions) ([]NotificationThread, error)
func (c *Client) MarkNotificationsRead(lastReadAt string, all bool) error
func (c *Client) MarkNotificationRead(id int64) error

// User
func (c *Client) GetCurrentUser() (*User, error)
func (c *Client) GetUser(username string) (*User, error)

// Push Mirrors
func (c *Client) ListPushMirrors(owner, repo string, page, limit int) ([]PushMirror, error)
func (c *Client) GetPushMirror(owner, repo, name string) (*PushMirror, error)
func (c *Client) AddPushMirror(owner, repo string, req CreatePushMirrorOption) (*PushMirror, error)
func (c *Client) SyncPushMirrors(owner, repo string) error
func (c *Client) DeletePushMirror(owner, repo, name string) error
```

### 3. GitHub Client (`github.go` - NEW)

New file for GitHub API interactions:

```go
type GitHubClient struct {
    token string
    client *http.Client
}

type GitHubDeployKey struct {
    ID   int64  `json:"id"`
    Key  string `json:"key"`
    Title string `json:"title"`
    ReadOnly bool `json:"read_only"`
}

func NewGitHubClient(token string) *GitHubClient
func (c *GitHubClient) AddDeployKey(owner, repo string, req CreateDeployKeyRequest) (*GitHubDeployKey, error)
func (c *GitHubClient) ListDeployKeys(owner, repo string) ([]GitHubDeployKey, error)
func (c *GitHubClient) DeleteDeployKey(owner, repo string, id int) error
```

### 4. SSH Key Generation (`ssh.go` - NEW)

New file for SSH key operations:

```go
func GenerateSSHKeyPair(comment string) (privKey []byte, pubKey []byte, err error)
func SavePrivateKey(path string, key []byte) error
func GetPublicKey(path string) (string, error)
```

### 5. Commands (`cmd/forge/`)

#### `cmd/forge/auth.go` (NEW)

```
forge auth login
  --server <url>     # Forgejo server URL
  --token <token>    # API token (optional, prompt if not provided)

forge auth login --github
  --token <token>    # GitHub token (optional, prompt if not provided)

forge auth status
  # Shows: server URL, Forgejo user, GitHub user (if configured)

forge auth logout
  --github           # Remove only GitHub credentials
```

#### `cmd/forge/repo.go` (NEW)

```
forge repo list
  -u / --user <username>    # List user's repos (default: current user)
  -o / --org <org>          # List org's repos
  --limit <n>               # Max results (default: 30)
  --page <n>                # Page number
  --json                    # JSON output

forge repo view [<owner>/<repo>]
  # Auto-detects from git remote if not specified
  --json                    # JSON output

forge repo create <name>
  -d / --description <desc>
  -p / --private
  --auto-init               # Initialize with README
  --license <spdx-id>       # License template
  --gitignore <template>    # Gitignore template

forge repo clone <owner>/<repo> [<directory>]
  # Runs: git clone <clone_url>

forge repo fork <owner>/<repo>
  --org <org>               # Fork to organization
```

#### `cmd/forge/pr.go` (NEW)

```
forge pr list
  -r / --repo <repo>
  -s / --state <state>      # open, closed, all (default: open)
  --limit <n>
  --page <n>
  --json

forge pr view <number>
  -r / --repo <repo>
  --json

forge pr create
  -r / --repo <repo>
  -t / --title <title>
  -b / --body <body>
  --body-file <path>
  --editor
  --head <branch>           # Source branch (default: current)
  --base <branch>           # Target branch (default: default branch)
  -l / --label <label>      # Repeatable
  -m / --milestone <title>
  -a / --assignee <user>    # Repeatable

forge pr merge <number>
  -r / --repo <repo>
  --method <method>         # merge, rebase, squash (default: merge)

forge pr close <number>
  -r / --repo <repo>

forge pr reopen <number>
  -r / --repo <repo>

forge pr diff <number>
  -r / --repo <repo>
  # Opens in pager or outputs to stdout
```

#### `cmd/forge/issues.go` (MODIFY)

Add to existing file:

```
forge issues comment <number>
  -r / --repo <repo>
  -b / --body <body>
  --body-file <path>
  --editor

forge issues comments <number>
  -r / --repo <repo>
  --json
```

#### `cmd/forge/label.go` (NEW)

```
forge label list
  -r / --repo <repo>
  --json

forge label create
  -r / --repo <repo>
  -n / --name <name>
  -c / --color <hex>        # e.g., "ff0000" or "#ff0000"
  -d / --description <desc>

forge label delete <name>
  -r / --repo <repo>

forge label edit <name>
  -r / --repo <repo>
  -n / --name <new-name>
  -c / --color <hex>
  -d / --description <desc>
```

#### `cmd/forge/milestone.go` (NEW)

```
forge milestone list
  -r / --repo <repo>
  -s / --state <state>      # open, closed, all (default: open)
  --json

forge milestone view <title>
  -r / --repo <repo>
  --json

forge milestone create
  -r / --repo <repo>
  -t / --title <title>
  -d / --description <desc>
  --due-date <YYYY-MM-DD>

forge milestone close <title>
  -r / --repo <repo>

forge milestone delete <title>
  -r / --repo <repo>
```

#### `cmd/forge/release.go` (NEW)

```
forge release list
  -r / --repo <repo>
  --limit <n>
  --page <n>
  --json

forge release view <tag>
  -r / --repo <repo>
  --json

forge release create
  -r / --repo <repo>
  -t / --tag <tag>
  -n / --name <name>
  -b / --body <body>
  --body-file <path>
  --editor
  --draft
  --prerelease
  --target <commit-ish>

forge release delete <tag>
  -r / --repo <repo>
  -y / --yes               # Skip confirmation
```

#### `cmd/forge/notifications.go` (NEW)

```
forge notifications list
  --all                    # Show read notifications too
  --limit <n>
  --json

forge notifications read
  --all                    # Mark all as read
  --id <id>                # Mark specific notification as read
```

#### `cmd/forge/key.go` (NEW)

```
forge key list
  -r / --repo <repo>
  --json

forge key add <title> <public-key-path>
  -r / --repo <repo>
  --read-only              # Read-only access (default: read-write)

forge key remove <id>
  -r / --repo <repo>
  -y / --yes               # Skip confirmation
```

#### `cmd/forge/user.go` (NEW)

```
forge user info
  --json                   # Shows current authenticated user

forge user view <username>
  --json
```

#### `cmd/forge/mirror.go` (NEW)

```
forge mirror setup
  -r / --repo <repo>           # Forgejo repo (owner/repo)
  --remote-url <url>           # Remote repo URL (https or ssh)
  --remote-user <user>         # Remote username (for HTTPS auth)
  --interval <duration>        # Sync interval (default: 6h)
  --branch <branch>            # Branch filter (default: all branches)
  --sync-on-commit             # Sync on every commit (default: false)
  --use-ssh                    # Use SSH for remote (default: false)

forge mirror sync <owner>/<repo>
  # Triggers immediate push mirror sync

forge mirror list -r / --repo <repo>
  --json                       # List push mirrors for a repo

forge mirror remove <owner>/<repo> <remote-name>
  -y / --yes                   # Skip confirmation
```

### 6. Mirror Setup Flow

When `forge mirror setup` runs:

1. Validate inputs (Forgejo repo exists, remote repo accessible)
2. **Option A: HTTPS with token**
   - Prompt for GitHub username and token
   - Use `POST /repos/{owner}/{repo}/push_mirrors` with:
     - `remote_address`: `https://<user>:<token>@github.com/<user>/<repo>.git`
     - `remote_username`: GitHub username
     - `remote_password`: GitHub token
     - `interval`: user-specified interval
     - `sync_on_commit`: user-specified
     - `branch_filter`: user-specified
3. **Option B: SSH with deploy key**
   - Generate SSH keypair in `~/.config/forge/keys/<owner>_<repo>/`
   - Add public key as deploy key on GitHub via GitHub API (read/write)
   - Add private key to Forgejo user via `POST /user/keys`
   - Use `POST /repos/{owner}/{repo}/push_mirrors` with:
     - `remote_address`: `git@github.com:<user>/<repo>.git`
     - `use_ssh`: true
     - `interval`: user-specified interval
     - `sync_on_commit`: user-specified
     - `branch_filter`: user-specified
4. Print success message with sync status

### 7. Shared Flags

Global flags (already exist):
- `--config <path>` - Config file path
- `--server <url>` - Forgejo server URL
- `--owner <owner>` - Repository owner
- `--token <token>` - Forgejo API token
- `--json` - JSON output
- `--verbose / -v` - Verbose output

Per-repo commands use:
- `-r / --repo <repo>` - Repository name (required if not in git repo)
- Auto-detect from git remote when inside a repository

### 8. Output Formatting

Table output (default):
```
#  Title                    State   Author   Created
1  Fix login bug            open    alice    2 hours ago
2  Add dark mode            closed  bob      3 days ago
```

JSON output (`--json`):
```json
[
  {
    "number": 1,
    "title": "Fix login bug",
    "state": "open",
    "user": {"login": "alice"},
    "created_at": "2026-07-24T10:00:00Z"
  }
]
```
