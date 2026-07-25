# Forgejo API Reference

Base URL: `/api/v1`
Version: 14.0.3+gitea-1.22.0

## ActivityPub (11)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/activitypub/actor` | Returns the instance's Actor |
| POST | `/activitypub/actor/inbox` | Send to the inbox |
| POST | `/activitypub/actor/outbox` | Display the outbox (always empty) |
| GET | `/activitypub/repository-id/{repository-id}` | Returns the Repository actor for a repo |
| POST | `/activitypub/repository-id/{repository-id}/inbox` | Send to the inbox |
| POST | `/activitypub/repository-id/{repository-id}/outbox` | Display the outbox |
| GET | `/activitypub/user-id/{user-id}` | Returns the Person actor for a user |
| GET | `/activitypub/user-id/{user-id}/activities/{activity-id}` | Get a specific activity object of the user |
| GET | `/activitypub/user-id/{user-id}/activities/{activity-id}/activity` | Get a specific activity of the user |
| POST | `/activitypub/user-id/{user-id}/inbox` | Send to the inbox |
| GET | `/activitypub/user-id/{user-id}/outbox` | List the user's recorded activity |

## Admin (42)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/admin/cron` | List cron tasks |
| POST | `/admin/cron/{task}` | Run cron task |
| GET | `/admin/emails` | List all users' email addresses |
| GET | `/admin/emails/search` | Search users' email addresses |
| GET | `/admin/hooks` | List global (system) webhooks |
| POST | `/admin/hooks` | Create a hook |
| GET | `/admin/hooks/{id}` | Get a hook |
| DELETE | `/admin/hooks/{id}` | Delete a hook |
| PATCH | `/admin/hooks/{id}` | Update a hook |
| GET | `/admin/orgs` | List all organizations |
| GET | `/admin/quota/groups` | List the available quota groups |
| POST | `/admin/quota/groups` | Create a new quota group |
| GET | `/admin/quota/groups/{quotagroup}` | Get information about the quota group |
| DELETE | `/admin/quota/groups/{quotagroup}` | Delete a quota group |
| PUT | `/admin/quota/groups/{quotagroup}/rules/{quotarule}` | Adds a rule to a quota group |
| DELETE | `/admin/quota/groups/{quotagroup}/rules/{quotarule}` | Removes a rule from a quota group |
| GET | `/admin/quota/groups/{quotagroup}/users` | List users in a quota group |
| PUT | `/admin/quota/groups/{quotagroup}/users/{username}` | Add a user to a quota group |
| DELETE | `/admin/quota/groups/{quotagroup}/users/{username}` | Remove a user from a quota group |
| GET | `/admin/quota/rules` | List the available quota rules |
| POST | `/admin/quota/rules` | Create a new quota rule |
| GET | `/admin/quota/rules/{quotarule}` | Get information about a quota rule |
| DELETE | `/admin/quota/rules/{quotarule}` | Deletes a quota rule |
| PATCH | `/admin/quota/rules/{quotarule}` | Change an existing quota rule |
| GET | `/admin/runners/jobs` | Search action jobs according filter conditions |
| GET | `/admin/runners/registration-token` | Get an global actions runner registration token |
| GET | `/admin/unadopted` | List unadopted repositories |
| POST | `/admin/unadopted/{owner}/{repo}` | Adopt unadopted files as a repository |
| DELETE | `/admin/unadopted/{owner}/{repo}` | Delete unadopted files |
| GET | `/admin/users` | Search users according filter conditions |
| POST | `/admin/users` | Create a user account |
| DELETE | `/admin/users/{username}` | Delete user account |
| PATCH | `/admin/users/{username}` | Edit an existing user |
| GET | `/admin/users/{username}/emails` | List all email addresses for a user |
| DELETE | `/admin/users/{username}/emails` | Delete email addresses from a user's account |
| POST | `/admin/users/{username}/keys` | Add an SSH public key to user's account |
| DELETE | `/admin/users/{username}/keys/{id}` | Remove a public key from user's account |
| POST | `/admin/users/{username}/orgs` | Create an organization |
| GET | `/admin/users/{username}/quota` | Get the user's quota info |
| POST | `/admin/users/{username}/quota/groups` | Set the user's quota groups to a given list |
| POST | `/admin/users/{username}/rename` | Rename a user |
| POST | `/admin/users/{username}/repos` | Create a repository on behalf of a user |

## Issue (67)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/repos/issues/search` | Search for issues across repositories |
| GET | `/repos/{owner}/{repo}/issues` | List a repository's issues |
| POST | `/repos/{owner}/{repo}/issues` | Create an issue |
| GET | `/repos/{owner}/{repo}/issues/comments` | List all comments in a repository |
| GET | `/repos/{owner}/{repo}/issues/comments/{id}` | Get a comment |
| DELETE | `/repos/{owner}/{repo}/issues/comments/{id}` | Delete a comment |
| PATCH | `/repos/{owner}/{repo}/issues/comments/{id}` | Edit a comment |
| GET | `/repos/{owner}/{repo}/issues/comments/{id}/assets` | List comment's attachments |
| POST | `/repos/{owner}/{repo}/issues/comments/{id}/assets` | Create a comment attachment |
| GET | `/repos/{owner}/{repo}/issues/comments/{id}/assets/{attachment_id}` | Get a comment attachment |
| DELETE | `/repos/{owner}/{repo}/issues/comments/{id}/assets/{attachment_id}` | Delete a comment attachment |
| PATCH | `/repos/{owner}/{repo}/issues/comments/{id}/assets/{attachment_id}` | Edit a comment attachment |
| GET | `/repos/{owner}/{repo}/issues/comments/{id}/reactions` | Get reactions from a comment |
| POST | `/repos/{owner}/{repo}/issues/comments/{id}/reactions` | Add a reaction to a comment |
| DELETE | `/repos/{owner}/{repo}/issues/comments/{id}/reactions` | Remove a reaction from a comment |
| GET | `/repos/{owner}/{repo}/issues/{index}` | Get an issue |
| DELETE | `/repos/{owner}/{repo}/issues/{index}` | Delete an issue |
| PATCH | `/repos/{owner}/{repo}/issues/{index}` | Edit an issue |
| GET | `/repos/{owner}/{repo}/issues/{index}/assets` | List issue's attachments |
| POST | `/repos/{owner}/{repo}/issues/{index}/assets` | Create an issue attachment |
| GET | `/repos/{owner}/{repo}/issues/{index}/assets/{attachment_id}` | Get an issue attachment |
| DELETE | `/repos/{owner}/{repo}/issues/{index}/assets/{attachment_id}` | Delete an issue attachment |
| PATCH | `/repos/{owner}/{repo}/issues/{index}/assets/{attachment_id}` | Edit an issue attachment |
| GET | `/repos/{owner}/{repo}/issues/{index}/blocks` | List issues that are blocked by this issue |
| POST | `/repos/{owner}/{repo}/issues/{index}/blocks` | Block the issue given in the body |
| DELETE | `/repos/{owner}/{repo}/issues/{index}/blocks` | Unblock the issue given in the body |
| GET | `/repos/{owner}/{repo}/issues/{index}/comments` | List all comments on an issue |
| POST | `/repos/{owner}/{repo}/issues/{index}/comments` | Add a comment to an issue |
| DELETE | `/repos/{owner}/{repo}/issues/{index}/comments/{id}` | Delete a comment (deprecated) |
| PATCH | `/repos/{owner}/{repo}/issues/{index}/comments/{id}` | Edit a comment (deprecated) |
| POST | `/repos/{owner}/{repo}/issues/{index}/deadline` | Set an issue deadline |
| GET | `/repos/{owner}/{repo}/issues/{index}/dependencies` | List an issue's dependencies |
| POST | `/repos/{owner}/{repo}/issues/{index}/dependencies` | Make the issue depend on the issue in the form |
| DELETE | `/repos/{owner}/{repo}/issues/{index}/dependencies` | Remove an issue dependency |
| GET | `/repos/{owner}/{repo}/issues/{index}/labels` | Get an issue's labels |
| PUT | `/repos/{owner}/{repo}/issues/{index}/labels` | Replace an issue's labels |
| POST | `/repos/{owner}/{repo}/issues/{index}/labels` | Add a label to an issue |
| DELETE | `/repos/{owner}/{repo}/issues/{index}/labels` | Remove all labels from an issue |
| DELETE | `/repos/{owner}/{repo}/issues/{index}/labels/{identifier}` | Remove a label from an issue |
| POST | `/repos/{owner}/{repo}/issues/{index}/pin` | Pin an Issue |
| DELETE | `/repos/{owner}/{repo}/issues/{index}/pin` | Unpin an Issue |
| PATCH | `/repos/{owner}/{repo}/issues/{index}/pin/{position}` | Moves the Pin to the given Position |
| GET | `/repos/{owner}/{repo}/issues/{index}/reactions` | Get a list reactions of an issue |
| POST | `/repos/{owner}/{repo}/issues/{index}/reactions` | Add a reaction to an issue |
| DELETE | `/repos/{owner}/{repo}/issues/{index}/reactions` | Remove a reaction from an issue |
| DELETE | `/repos/{owner}/{repo}/issues/{index}/stopwatch/delete` | Delete an issue's existing stopwatch |
| POST | `/repos/{owner}/{repo}/issues/{index}/stopwatch/start` | Start stopwatch on an issue |
| POST | `/repos/{owner}/{repo}/issues/{index}/stopwatch/stop` | Stop an issue's existing stopwatch |
| GET | `/repos/{owner}/{repo}/issues/{index}/subscriptions` | Get users who subscribed on an issue |
| GET | `/repos/{owner}/{repo}/issues/{index}/subscriptions/check` | Check if user is subscribed to an issue |
| PUT | `/repos/{owner}/{repo}/issues/{index}/subscriptions/{user}` | Subscribe user to issue |
| DELETE | `/repos/{owner}/{repo}/issues/{index}/subscriptions/{user}` | Unsubscribe user from issue |
| GET | `/repos/{owner}/{repo}/issues/{index}/timeline` | List all comments and events on an issue |
| GET | `/repos/{owner}/{repo}/issues/{index}/times` | List an issue's tracked times |
| POST | `/repos/{owner}/{repo}/issues/{index}/times` | Add tracked time to an issue |
| DELETE | `/repos/{owner}/{repo}/issues/{index}/times` | Reset a tracked time of an issue |
| DELETE | `/repos/{owner}/{repo}/issues/{index}/times/{id}` | Delete specific tracked time |
| GET | `/repos/{owner}/{repo}/labels` | Get all of a repository's labels |
| POST | `/repos/{owner}/{repo}/labels` | Create a label |
| GET | `/repos/{owner}/{repo}/labels/{id}` | Get a single label |
| DELETE | `/repos/{owner}/{repo}/labels/{id}` | Delete a label |
| PATCH | `/repos/{owner}/{repo}/labels/{id}` | Update a label |
| GET | `/repos/{owner}/{repo}/milestones` | Get all of a repository's opened milestones |
| POST | `/repos/{owner}/{repo}/milestones` | Create a milestone |
| GET | `/repos/{owner}/{repo}/milestones/{id}` | Get a milestone |
| DELETE | `/repos/{owner}/{repo}/milestones/{id}` | Delete a milestone |
| PATCH | `/repos/{owner}/{repo}/milestones/{id}` | Update a milestone |

## Miscellaneous (13)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/gitignore/templates` | Returns a list of all gitignore templates |
| GET | `/gitignore/templates/{name}` | Returns information about a gitignore template |
| GET | `/label/templates` | Returns a list of all label templates |
| GET | `/label/templates/{name}` | Returns all labels in a template |
| GET | `/licenses` | Returns a list of all license templates |
| GET | `/licenses/{name}` | Returns information about a license template |
| POST | `/markdown` | Render a markdown document as HTML |
| POST | `/markdown/raw` | Render raw markdown as HTML |
| POST | `/markup` | Render a markup document as HTML |
| GET | `/nodeinfo` | Returns the nodeinfo of the Forgejo application |
| GET | `/signing-key.gpg` | Get default signing-key.gpg |
| GET | `/signing-key.ssh` | Get default signing-key.ssh |
| GET | `/version` | Returns the version of the running application |

## Notification (7)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/notifications` | List users's notification threads |
| PUT | `/notifications` | Mark notification threads as read, pinned or unread |
| GET | `/notifications/new` | Check if unread notifications exist |
| GET | `/notifications/threads/{id}` | Get notification thread by ID |
| PATCH | `/notifications/threads/{id}` | Mark notification thread as read by ID |
| GET | `/repos/{owner}/{repo}/notifications` | List users's notification threads on a specific repo |
| PUT | `/repos/{owner}/{repo}/notifications` | Mark notification threads as read, pinned or unread on a specific repo |

## Organization (65)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/org/{org}/repos` | Create a repository in an organization (deprecated) |
| GET | `/orgs` | List all organizations |
| POST | `/orgs` | Create an organization |
| GET | `/orgs/{org}` | Get an organization |
| DELETE | `/orgs/{org}` | Delete an organization |
| PATCH | `/orgs/{org}` | Edit an organization |
| GET | `/orgs/{org}/actions/runners/jobs` | Search for organization's action jobs |
| GET | `/orgs/{org}/actions/runners/registration-token` | Get an organization's actions runner registration token |
| GET | `/orgs/{org}/actions/secrets` | List actions secrets of an organization |
| PUT | `/orgs/{org}/actions/secrets/{secretname}` | Create or Update a secret value in an organization |
| DELETE | `/orgs/{org}/actions/secrets/{secretname}` | Delete a secret in an organization |
| GET | `/orgs/{org}/actions/variables` | List variables of an organization |
| GET | `/orgs/{org}/actions/variables/{variablename}` | Get organization's variable by name |
| PUT | `/orgs/{org}/actions/variables/{variablename}` | Update variable in organization |
| POST | `/orgs/{org}/actions/variables/{variablename}` | Create a new variable in organization |
| DELETE | `/orgs/{org}/actions/variables/{variablename}` | Delete organization's variable by name |
| GET | `/orgs/{org}/activities/feeds` | List an organization's activity feeds |
| POST | `/orgs/{org}/avatar` | Update an organization's avatar |
| DELETE | `/orgs/{org}/avatar` | Delete an organization's avatar |
| PUT | `/orgs/{org}/block/{username}` | Blocks a user from the organization |
| GET | `/orgs/{org}/hooks` | List an organization's webhooks |
| POST | `/orgs/{org}/hooks` | Create a hook |
| GET | `/orgs/{org}/hooks/{id}` | Get a hook |
| DELETE | `/orgs/{org}/hooks/{id}` | Delete a hook |
| PATCH | `/orgs/{org}/hooks/{id}` | Update a hook |
| GET | `/orgs/{org}/labels` | List an organization's labels |
| POST | `/orgs/{org}/labels` | Create a label for an organization |
| GET | `/orgs/{org}/labels/{id}` | Get a single label |
| DELETE | `/orgs/{org}/labels/{id}` | Delete a label |
| PATCH | `/orgs/{org}/labels/{id}` | Update a label |
| GET | `/orgs/{org}/list_blocked` | List the organization's blocked users |
| GET | `/orgs/{org}/members` | List an organization's members |
| GET | `/orgs/{org}/members/{username}` | Check if a user is a member of an organization |
| DELETE | `/orgs/{org}/members/{username}` | Remove a member from an organization |
| GET | `/orgs/{org}/public_members` | List an organization's public members |
| GET | `/orgs/{org}/public_members/{username}` | Check if a user is a public member |
| PUT | `/orgs/{org}/public_members/{username}` | Publicize a user's membership |
| DELETE | `/orgs/{org}/public_members/{username}` | Conceal a user's membership |
| GET | `/orgs/{org}/quota` | Get quota information for an organization |
| GET | `/orgs/{org}/quota/artifacts` | List the artifacts affecting the organization's quota |
| GET | `/orgs/{org}/quota/attachments` | List the attachments affecting the organization's quota |
| GET | `/orgs/{org}/quota/check` | Check if the organization is over quota |
| GET | `/orgs/{org}/quota/packages` | List the packages affecting the organization's quota |
| POST | `/orgs/{org}/rename` | Rename an organization |
| GET | `/orgs/{org}/repos` | List an organization's repos |
| POST | `/orgs/{org}/repos` | Create a repository in an organization |
| GET | `/orgs/{org}/teams` | List an organization's teams |
| POST | `/orgs/{org}/teams` | Create a team |
| GET | `/orgs/{org}/teams/search` | Search for teams within an organization |
| PUT | `/orgs/{org}/unblock/{username}` | Unblock a user from the organization |
| GET | `/teams/{id}` | Get a team |
| DELETE | `/teams/{id}` | Delete a team |
| PATCH | `/teams/{id}` | Edit a team |
| GET | `/teams/{id}/activities/feeds` | List a team's activity feeds |
| GET | `/teams/{id}/members` | List a team's members |
| GET | `/teams/{id}/members/{username}` | List a particular member of team |
| PUT | `/teams/{id}/members/{username}` | Add a team member |
| DELETE | `/teams/{id}/members/{username}` | Remove a team member |
| GET | `/teams/{id}/repos` | List a team's repos |
| GET | `/teams/{id}/repos/{org}/{repo}` | List a particular repo of team |
| PUT | `/teams/{id}/repos/{org}/{repo}` | Add a repository to a team |
| DELETE | `/teams/{id}/repos/{org}/{repo}` | Remove a repository from a team |
| GET | `/user/orgs` | List the current user's organizations |
| GET | `/users/{username}/orgs` | List a user's organizations |
| GET | `/users/{username}/orgs/{org}/permissions` | Get user permissions in organization |

## Package (6)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/packages/{owner}` | Gets all packages of an owner |
| POST | `/packages/{owner}/{type}/{name}/-/link/{repo_name}` | Link a package to a repository |
| POST | `/packages/{owner}/{type}/{name}/-/unlink` | Unlink a package from a repository |
| GET | `/packages/{owner}/{type}/{name}/{version}` | Gets a package |
| DELETE | `/packages/{owner}/{type}/{name}/{version}` | Delete a package |
| GET | `/packages/{owner}/{type}/{name}/{version}/files` | Gets all files of a package |

## Repository (184)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/repos/migrate` | Migrate a remote git repository |
| GET | `/repos/search` | Search for repositories |
| GET | `/repos/{owner}/{repo}` | Get a repository |
| DELETE | `/repos/{owner}/{repo}` | Delete a repository |
| PATCH | `/repos/{owner}/{repo}` | Edit a repository's properties |
| GET | `/repos/{owner}/{repo}/actions/runners/jobs` | Search for repository's action jobs |
| GET | `/repos/{owner}/{repo}/actions/runners/registration-token` | Get a repository's actions runner registration token |
| GET | `/repos/{owner}/{repo}/actions/runs` | List a repository's action runs |
| GET | `/repos/{owner}/{repo}/actions/runs/{run_id}` | Get an action run |
| GET | `/repos/{owner}/{repo}/actions/secrets` | List a repo's actions secrets |
| PUT | `/repos/{owner}/{repo}/actions/secrets/{secretname}` | Create or Update a secret value |
| DELETE | `/repos/{owner}/{repo}/actions/secrets/{secretname}` | Delete a secret |
| GET | `/repos/{owner}/{repo}/actions/tasks` | List a repository's action tasks |
| GET | `/repos/{owner}/{repo}/actions/variables` | Get repo-level variables list |
| GET | `/repos/{owner}/{repo}/actions/variables/{variablename}` | Get a repo-level variable |
| PUT | `/repos/{owner}/{repo}/actions/variables/{variablename}` | Update a repo-level variable |
| POST | `/repos/{owner}/{repo}/actions/variables/{variablename}` | Create a repo-level variable |
| DELETE | `/repos/{owner}/{repo}/actions/variables/{variablename}` | Delete a repo-level variable |
| POST | `/repos/{owner}/{repo}/actions/workflows/{workflowfilename}/dispatches` | Dispatches a workflow |
| GET | `/repos/{owner}/{repo}/activities/feeds` | List a repository's activity feeds |
| GET | `/repos/{owner}/{repo}/archive/{archive}` | Get an archive of a repository |
| GET | `/repos/{owner}/{repo}/assignees` | Return all users that can be assigned to issues |
| POST | `/repos/{owner}/{repo}/avatar` | Update a repository's avatar |
| DELETE | `/repos/{owner}/{repo}/avatar` | Delete a repository's avatar |
| GET | `/repos/{owner}/{repo}/branch_protections` | List branch protections |
| POST | `/repos/{owner}/{repo}/branch_protections` | Create a branch protections |
| GET | `/repos/{owner}/{repo}/branch_protections/{name}` | Get a specific branch protection |
| DELETE | `/repos/{owner}/{repo}/branch_protections/{name}` | Delete a specific branch protection |
| PATCH | `/repos/{owner}/{repo}/branch_protections/{name}` | Edit a branch protections |
| GET | `/repos/{owner}/{repo}/branches` | List a repository's branches |
| POST | `/repos/{owner}/{repo}/branches` | Create a branch |
| GET | `/repos/{owner}/{repo}/branches/{branch}` | Retrieve a specific branch |
| DELETE | `/repos/{owner}/{repo}/branches/{branch}` | Delete a specific branch |
| PATCH | `/repos/{owner}/{repo}/branches/{branch}` | Update a branch |
| GET | `/repos/{owner}/{repo}/collaborators` | List a repository's collaborators |
| GET | `/repos/{owner}/{repo}/collaborators/{collaborator}` | Check if a user is a collaborator |
| PUT | `/repos/{owner}/{repo}/collaborators/{collaborator}` | Add a collaborator |
| DELETE | `/repos/{owner}/{repo}/collaborators/{collaborator}` | Delete a collaborator |
| GET | `/repos/{owner}/{repo}/collaborators/{collaborator}/permission` | Get repository permissions |
| GET | `/repos/{owner}/{repo}/commits` | Get a list of all commits |
| GET | `/repos/{owner}/{repo}/commits/{ref}/status` | Get a commit's combined status |
| GET | `/repos/{owner}/{repo}/commits/{ref}/statuses` | Get a commit's statuses |
| GET | `/repos/{owner}/{repo}/commits/{sha}/pull` | Get the pull request of the commit |
| GET | `/repos/{owner}/{repo}/compare/{basehead}` | Get commit comparison information |
| GET | `/repos/{owner}/{repo}/contents` | Gets the metadata of all the entries of the root dir |
| POST | `/repos/{owner}/{repo}/contents` | Modify multiple files in a repository |
| GET | `/repos/{owner}/{repo}/contents/{filepath}` | Gets the metadata and contents of an entry |
| PUT | `/repos/{owner}/{repo}/contents/{filepath}` | Update a file |
| POST | `/repos/{owner}/{repo}/contents/{filepath}` | Create a file |
| DELETE | `/repos/{owner}/{repo}/contents/{filepath}` | Delete a file |
| POST | `/repos/{owner}/{repo}/convert` | Convert a mirror repo to a normal repo |
| POST | `/repos/{owner}/{repo}/diffpatch` | Apply diff patch to repository |
| GET | `/repos/{owner}/{repo}/editorconfig/{filepath}` | Get the EditorConfig definitions |
| GET | `/repos/{owner}/{repo}/flags` | List a repository's flags |
| PUT | `/repos/{owner}/{repo}/flags` | Replace all flags |
| DELETE | `/repos/{owner}/{repo}/flags` | Remove all flags |
| GET | `/repos/{owner}/{repo}/flags/{flag}` | Check if a repository has a given flag |
| PUT | `/repos/{owner}/{repo}/flags/{flag}` | Add a flag |
| DELETE | `/repos/{owner}/{repo}/flags/{flag}` | Remove a flag |
| GET | `/repos/{owner}/{repo}/forks` | List a repository's forks |
| POST | `/repos/{owner}/{repo}/forks` | Fork a repository |
| GET | `/repos/{owner}/{repo}/git/blobs` | Gets multiple blobs |
| GET | `/repos/{owner}/{repo}/git/blobs/{sha}` | Gets the blob |
| GET | `/repos/{owner}/{repo}/git/commits/{sha}` | Get a single commit |
| GET | `/repos/{owner}/{repo}/git/commits/{sha}.{diffType}` | Get a commit's diff or patch |
| GET | `/repos/{owner}/{repo}/git/notes/{sha}` | Get a note corresponding to a single commit |
| POST | `/repos/{owner}/{repo}/git/notes/{sha}` | Set a note |
| DELETE | `/repos/{owner}/{repo}/git/notes/{sha}` | Removes a note |
| GET | `/repos/{owner}/{repo}/git/refs` | Get specified ref or filtered refs |
| GET | `/repos/{owner}/{repo}/git/refs/{ref}` | Get specified ref or filtered refs |
| GET | `/repos/{owner}/{repo}/git/tags/{sha}` | Gets the tag object of an annotated tag |
| GET | `/repos/{owner}/{repo}/git/trees/{sha}` | Gets the tree |
| GET | `/repos/{owner}/{repo}/hooks` | List the hooks |
| POST | `/repos/{owner}/{repo}/hooks` | Create a hook |
| GET | `/repos/{owner}/{repo}/hooks/git` | List the Git hooks |
| GET | `/repos/{owner}/{repo}/hooks/git/{id}` | Get a Git hook |
| DELETE | `/repos/{owner}/{repo}/hooks/git/{id}` | Delete a Git hook |
| PATCH | `/repos/{owner}/{repo}/hooks/git/{id}` | Edit a Git hook |
| GET | `/repos/{owner}/{repo}/hooks/{id}` | Get a hook |
| DELETE | `/repos/{owner}/{repo}/hooks/{id}` | Delete a hook |
| PATCH | `/repos/{owner}/{repo}/hooks/{id}` | Edit a hook |
| POST | `/repos/{owner}/{repo}/hooks/{id}/tests` | Test a push webhook |
| GET | `/repos/{owner}/{repo}/issue_config` | Returns the issue config |
| GET | `/repos/{owner}/{repo}/issue_config/validate` | Returns the validation information |
| GET | `/repos/{owner}/{repo}/issue_templates` | Get available issue templates |
| GET | `/repos/{owner}/{repo}/issues/pinned` | List a repo's pinned issues |
| GET | `/repos/{owner}/{repo}/keys` | List a repository's keys |
| POST | `/repos/{owner}/{repo}/keys` | Add a key |
| GET | `/repos/{owner}/{repo}/keys/{id}` | Get a repository's key by id |
| DELETE | `/repos/{owner}/{repo}/keys/{id}` | Delete a key |
| GET | `/repos/{owner}/{repo}/languages` | Get languages and number of bytes |
| GET | `/repos/{owner}/{repo}/media/{filepath}` | Get a file or LFS object |
| POST | `/repos/{owner}/{repo}/mirror-sync` | Sync a mirrored repository |
| GET | `/repos/{owner}/{repo}/new_pin_allowed` | Returns if new Issue Pins are allowed |
| GET | `/repos/{owner}/{repo}/pulls` | List a repo's pull requests |
| POST | `/repos/{owner}/{repo}/pulls` | Create a pull request |
| GET | `/repos/{owner}/{repo}/pulls/pinned` | List a repo's pinned pull requests |
| GET | `/repos/{owner}/{repo}/pulls/{base}/{head}` | Get a pull request by base and head |
| GET | `/repos/{owner}/{repo}/pulls/{index}` | Get a pull request |
| PATCH | `/repos/{owner}/{repo}/pulls/{index}` | Update a pull request |
| GET | `/repos/{owner}/{repo}/pulls/{index}.{diffType}` | Get a pull request diff or patch |
| GET | `/repos/{owner}/{repo}/pulls/{index}/commits` | Get commits for a pull request |
| GET | `/repos/{owner}/{repo}/pulls/{index}/files` | Get changed files |
| GET | `/repos/{owner}/{repo}/pulls/{index}/merge` | Check if a pull request has been merged |
| POST | `/repos/{owner}/{repo}/pulls/{index}/merge` | Merge a pull request |
| DELETE | `/repos/{owner}/{repo}/pulls/{index}/merge` | Cancel the scheduled auto merge |
| POST | `/repos/{owner}/{repo}/pulls/{index}/requested_reviewers` | Create review requests |
| DELETE | `/repos/{owner}/{repo}/pulls/{index}/requested_reviewers` | Cancel review requests |
| GET | `/repos/{owner}/{repo}/pulls/{index}/reviews` | List all reviews for a pull request |
| POST | `/repos/{owner}/{repo}/pulls/{index}/reviews` | Create a review |
| GET | `/repos/{owner}/{repo}/pulls/{index}/reviews/{id}` | Get a specific review |
| POST | `/repos/{owner}/{repo}/pulls/{index}/reviews/{id}` | Submit a pending review |
| DELETE | `/repos/{owner}/{repo}/pulls/{index}/reviews/{id}` | Delete a specific review |
| GET | `/repos/{owner}/{repo}/pulls/{index}/reviews/{id}/comments` | Get review comments |
| POST | `/repos/{owner}/{repo}/pulls/{index}/reviews/{id}/comments` | Add a new comment |
| GET | `/repos/{owner}/{repo}/pulls/{index}/reviews/{id}/comments/{comment}` | Get a pull review comment |
| DELETE | `/repos/{owner}/{repo}/pulls/{index}/reviews/{id}/comments/{comment}` | Delete a pull review comment |
| POST | `/repos/{owner}/{repo}/pulls/{index}/reviews/{id}/dismissals` | Dismiss a review |
| POST | `/repos/{owner}/{repo}/pulls/{index}/reviews/{id}/undismissals` | Cancel to dismiss a review |
| POST | `/repos/{owner}/{repo}/pulls/{index}/update` | Merge PR's baseBranch into headBranch |
| GET | `/repos/{owner}/{repo}/push_mirrors` | Get all push mirrors |
| POST | `/repos/{owner}/{repo}/push_mirrors` | Set up a new push mirror |
| POST | `/repos/{owner}/{repo}/push_mirrors-sync` | Sync all push mirrored repository |
| GET | `/repos/{owner}/{repo}/push_mirrors/{name}` | Get push mirror by remoteName |
| DELETE | `/repos/{owner}/{repo}/push_mirrors/{name}` | Remove a push mirror |
| GET | `/repos/{owner}/{repo}/raw/{filepath}` | Get a file from a repository |
| GET | `/repos/{owner}/{repo}/releases` | List a repo's releases |
| POST | `/repos/{owner}/{repo}/releases` | Create a release |
| GET | `/repos/{owner}/{repo}/releases/latest` | Gets the most recent release |
| GET | `/repos/{owner}/{repo}/releases/tags/{tag}` | Get a release by tag name |
| DELETE | `/repos/{owner}/{repo}/releases/tags/{tag}` | Delete a release by tag name |
| GET | `/repos/{owner}/{repo}/releases/{id}` | Get a release |
| DELETE | `/repos/{owner}/{repo}/releases/{id}` | Delete a release |
| PATCH | `/repos/{owner}/{repo}/releases/{id}` | Update a release |
| GET | `/repos/{owner}/{repo}/releases/{id}/assets` | List release's attachments |
| POST | `/repos/{owner}/{repo}/releases/{id}/assets` | Create a release attachment |
| GET | `/repos/{owner}/{repo}/releases/{id}/assets/{attachment_id}` | Get a release attachment |
| DELETE | `/repos/{owner}/{repo}/releases/{id}/assets/{attachment_id}` | Delete a release attachment |
| PATCH | `/repos/{owner}/{repo}/releases/{id}/assets/{attachment_id}` | Edit a release attachment |
| GET | `/repos/{owner}/{repo}/reviewers` | Return all users that can be requested to review |
| GET | `/repos/{owner}/{repo}/signing-key.gpg` | Get signing-key.gpg |
| GET | `/repos/{owner}/{repo}/stargazers` | List a repo's stargazers |
| GET | `/repos/{owner}/{repo}/statuses/{sha}` | Get a commit's statuses |
| POST | `/repos/{owner}/{repo}/statuses/{sha}` | Create a commit status |
| GET | `/repos/{owner}/{repo}/subscribers` | List a repo's watchers |
| GET | `/repos/{owner}/{repo}/subscription` | Check if the current user is watching a repo |
| PUT | `/repos/{owner}/{repo}/subscription` | Watch a repo |
| DELETE | `/repos/{owner}/{repo}/subscription` | Unwatch a repo |
| GET | `/repos/{owner}/{repo}/sync_fork` | Gets information about syncing the fork |
| POST | `/repos/{owner}/{repo}/sync_fork` | Syncs the default branch of a fork |
| GET | `/repos/{owner}/{repo}/sync_fork/{branch}` | Gets information about syncing a fork branch |
| POST | `/repos/{owner}/{repo}/sync_fork/{branch}` | Syncs a fork branch |
| GET | `/repos/{owner}/{repo}/tag_protections` | List tag protections |
| POST | `/repos/{owner}/{repo}/tag_protections` | Create a tag protections |
| GET | `/repos/{owner}/{repo}/tag_protections/{id}` | Get a specific tag protection |
| DELETE | `/repos/{owner}/{repo}/tag_protections/{id}` | Delete a specific tag protection |
| PATCH | `/repos/{owner}/{repo}/tag_protections/{id}` | Edit a tag protections |
| GET | `/repos/{owner}/{repo}/tags` | List a repository's tags |
| POST | `/repos/{owner}/{repo}/tags` | Create a new git tag |
| GET | `/repos/{owner}/{repo}/tags/{tag}` | Get the tag by tag name |
| DELETE | `/repos/{owner}/{repo}/tags/{tag}` | Delete a repository's tag by name |
| GET | `/repos/{owner}/{repo}/teams` | List a repository's teams |
| GET | `/repos/{owner}/{repo}/teams/{team}` | Check if a team is assigned to a repository |
| PUT | `/repos/{owner}/{repo}/teams/{team}` | Add a team to a repository |
| DELETE | `/repos/{owner}/{repo}/teams/{team}` | Delete a team from a repository |
| GET | `/repos/{owner}/{repo}/times` | List a repo's tracked times |
| GET | `/repos/{owner}/{repo}/times/{user}` | List a user's tracked times (deprecated) |
| GET | `/repos/{owner}/{repo}/topics` | Get list of topics |
| PUT | `/repos/{owner}/{repo}/topics` | Replace list of topics |
| PUT | `/repos/{owner}/{repo}/topics/{topic}` | Add a topic |
| DELETE | `/repos/{owner}/{repo}/topics/{topic}` | Delete a topic |
| POST | `/repos/{owner}/{repo}/transfer` | Transfer a repo ownership |
| POST | `/repos/{owner}/{repo}/transfer/accept` | Accept a repo transfer |
| POST | `/repos/{owner}/{repo}/transfer/reject` | Reject a repo transfer |
| POST | `/repos/{owner}/{repo}/wiki/new` | Create a wiki page |
| GET | `/repos/{owner}/{repo}/wiki/page/{pageName}` | Get a wiki page |
| DELETE | `/repos/{owner}/{repo}/wiki/page/{pageName}` | Delete a wiki page |
| PATCH | `/repos/{owner}/{repo}/wiki/page/{pageName}` | Edit a wiki page |
| GET | `/repos/{owner}/{repo}/wiki/pages` | Get all wiki pages |
| GET | `/repos/{owner}/{repo}/wiki/revisions/{pageName}` | Get revisions of a wiki page |
| POST | `/repos/{template_owner}/{template_repo}/generate` | Create a repository using a template |
| GET | `/repositories/{id}` | Get a repository by id |
| GET | `/topics/search` | Search for topics by keyword |
| POST | `/user/repos` | Create a repository |

## Settings (4)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/settings/api` | Get instance's global settings for api |
| GET | `/settings/attachment` | Get instance's global settings for Attachment |
| GET | `/settings/repository` | Get instance's global settings for repositories |
| GET | `/settings/ui` | Get instance's global settings for ui |

## User (75)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/user` | Get the authenticated user |
| GET | `/user/actions/runners/jobs` | Search for user's action jobs |
| GET | `/user/actions/runners/registration-token` | Get an user's actions runner registration token |
| PUT | `/user/actions/secrets/{secretname}` | Create or Update a secret value in a user scope |
| DELETE | `/user/actions/secrets/{secretname}` | Delete a secret in a user scope |
| GET | `/user/actions/variables` | Get the user-level list of variables |
| GET | `/user/actions/variables/{variablename}` | Get a user-level variable |
| PUT | `/user/actions/variables/{variablename}` | Update a user-level variable |
| POST | `/user/actions/variables/{variablename}` | Create a user-level variable |
| DELETE | `/user/actions/variables/{variablename}` | Delete a user-level variable |
| GET | `/user/applications/oauth2` | List the authenticated user's oauth2 applications |
| POST | `/user/applications/oauth2` | Creates a new OAuth2 application |
| GET | `/user/applications/oauth2/{id}` | Get an OAuth2 application |
| DELETE | `/user/applications/oauth2/{id}` | Delete an OAuth2 application |
| PATCH | `/user/applications/oauth2/{id}` | Update an OAuth2 application |
| POST | `/user/avatar` | Update avatar of the current user |
| DELETE | `/user/avatar` | Delete avatar of the current user |
| PUT | `/user/block/{username}` | Blocks a user from the doer |
| GET | `/user/emails` | List all email addresses of the current user |
| POST | `/user/emails` | Add an email addresses to the current user's account |
| DELETE | `/user/emails` | Delete email addresses from the current user's account |
| GET | `/user/followers` | List the authenticated user's followers |
| GET | `/user/following` | List the users that the authenticated user is following |
| GET | `/user/following/{username}` | Check whether a user is followed |
| PUT | `/user/following/{username}` | Follow a user |
| DELETE | `/user/following/{username}` | Unfollow a user |
| GET | `/user/gpg_key_token` | Get a Token to verify |
| POST | `/user/gpg_key_verify` | Verify a GPG key |
| GET | `/user/gpg_keys` | List the authenticated user's GPG keys |
| POST | `/user/gpg_keys` | Add a GPG public key |
| GET | `/user/gpg_keys/{id}` | Get a GPG key |
| DELETE | `/user/gpg_keys/{id}` | Remove a GPG public key |
| GET | `/user/hooks` | List the authenticated user's webhooks |
| POST | `/user/hooks` | Create a hook |
| GET | `/user/hooks/{id}` | Get a hook |
| DELETE | `/user/hooks/{id}` | Delete a hook |
| PATCH | `/user/hooks/{id}` | Update a hook |
| GET | `/user/keys` | List the authenticated user's public keys |
| POST | `/user/keys` | Create a public key |
| GET | `/user/keys/{id}` | Get a public key |
| DELETE | `/user/keys/{id}` | Delete a public key |
| GET | `/user/list_blocked` | List the authenticated user's blocked users |
| GET | `/user/quota` | Get quota information |
| GET | `/user/quota/artifacts` | List the artifacts affecting the quota |
| GET | `/user/quota/attachments` | List the attachments affecting the quota |
| GET | `/user/quota/check` | Check if over quota for a given subject |
| GET | `/user/quota/packages` | List the packages affecting the quota |
| GET | `/user/repos` | List the repos that the authenticated user owns |
| POST | `/user/repos` | Create a repository |
| GET | `/user/settings` | Get current user's account settings |
| PATCH | `/user/settings` | Update settings in current user's account |
| GET | `/user/starred` | The repos that the authenticated user has starred |
| GET | `/user/starred/{owner}/{repo}` | Whether the authenticated is starring the repo |
| PUT | `/user/starred/{owner}/{repo}` | Star the given repo |
| DELETE | `/user/starred/{owner}/{repo}` | Unstar the given repo |
| GET | `/user/stopwatches` | Get list of all existing stopwatches |
| GET | `/user/subscriptions` | List repositories watched by the authenticated user |
| GET | `/user/teams` | List all the teams a user belongs to |
| GET | `/user/times` | List the current user's tracked times |
| PUT | `/user/unblock/{username}` | Unblocks a user from the doer |
| GET | `/users/search` | Search for users |
| GET | `/users/{username}` | Get a user |
| GET | `/users/{username}/activities/feeds` | List a user's activity feeds |
| GET | `/users/{username}/followers` | List the given user's followers |
| GET | `/users/{username}/following` | List the users that the given user is following |
| GET | `/users/{username}/following/{target}` | Check if one user is following another |
| GET | `/users/{username}/gpg_keys` | List the given user's GPG keys |
| GET | `/users/{username}/heatmap` | Get a user's heatmap |
| GET | `/users/{username}/keys` | List the given user's public keys |
| GET | `/users/{username}/repos` | List the repos owned by the given user |
| GET | `/users/{username}/starred` | The repos that the given user has starred |
| GET | `/users/{username}/subscriptions` | List the repositories watched by a user |
| GET | `/users/{username}/tokens` | List the specified user's access tokens |
| POST | `/users/{username}/tokens` | Generate an access token |
| DELETE | `/users/{username}/tokens/{token}` | Delete an access token |
