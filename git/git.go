package git

import (
	"fmt"
	"log/slog"
	"path"
	"strings"

	"github.com/JacobGDG/bark/shell"
)

// Client runs git commands against a resolved git binary.
type Client struct {
	git *shell.Client
}

// New resolves the git binary and returns a Client.
func New() (*Client, error) {
	g, err := shell.New("git")
	if err != nil {
		return nil, fmt.Errorf("git client: %w", err)
	}
	return &Client{git: g}, nil
}

// Clone clones url to dest, passing any extra args after the positional arguments.
func (c *Client) Clone(url, dest string, args ...string) error {
	_, err := c.git.Run(append([]string{"clone", url, dest}, args...)...)

	slog.Debug("Repository cloned", "url", url)
	slog.Debug("Repository cloned", "destination", dest)

	return err
}

// WorktreeAdd adds a worktree at dest from repoPath, passing any extra args after the positional arguments.
func (c *Client) WorktreeAdd(repoPath, dest string, args ...string) error {
	_, err := c.git.Run(append([]string{"-C", repoPath, "worktree", "add", dest}, args...)...)

	slog.Debug("Worktree added", "repoPath", repoPath)
	slog.Debug("Worktree added", "destination", dest)

	return err
}

// SetRemoteHead resolves and sets refs/remotes/origin/HEAD for the repo at repoPath.
func (c *Client) SetRemoteHead(repoPath string) error {
	_, err := c.git.Run("-C", repoPath, "remote", "set-head", "origin", "--auto")

	slog.Debug("Remote head set", "repoPath", repoPath)

	return err
}

// DefaultBranch returns the default branch name for the repo at repoPath.
func (c *Client) DefaultBranch(repoPath string) (string, error) {
	ref, err := c.git.Run("-C", repoPath, "symbolic-ref", "refs/remotes/origin/HEAD")
	if err != nil {
		return "", fmt.Errorf("could not determine default branch: %w", err)
	}

	slog.Debug("Default branch set", "repoPath", repoPath)

	return path.Base(strings.TrimSpace(ref)), nil
}
