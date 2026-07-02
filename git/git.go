package git

import (
	"fmt"
	"log/slog"

	"github.com/JacobGDG/bark/shell"
)

// Client runs git commands against a resolved git binary.
type Client struct {
	git *shell.Client
}

func Check() error {
	err := shell.Check("git")
	if err != nil {
		return err
	}
	return nil
}

// New resolves the git binary and returns a Client.
func New() (*Client, error) {
	g, err := shell.New("git")
	if err != nil {
		return nil, err
	}
	return &Client{git: g}, nil
}

// Clone clones url to dest, passing any extra args after the positional arguments.
func (c *Client) Clone(url, dest string, args ...string) error {
	_, err := c.git.Run(append([]string{"clone", url, dest}, args...)...)

	if err == nil {
		slog.Debug("repository cloned", "url", url, "dest", dest)
	}

	return err
}

// WorktreeAdd adds a worktree at dest from repoPath, passing any extra args after the positional arguments.
func (c *Client) WorktreeAdd(repoPath, dest string, args ...string) error {
	_, err := c.git.Run(append([]string{"-C", repoPath, "worktree", "add", dest}, args...)...)

	if err == nil {
		slog.Debug("worktree added", "repo", repoPath, "dest", dest)
	}

	return err
}

// DefaultBranch returns the default branch name for the repo at repoPath.
func (c *Client) DefaultBranch(repoPath string) (string, error) {
	branch, err := c.git.Run("-C", repoPath, "branch", "--show-current")
	if err != nil {
		return "", fmt.Errorf("could not determine default branch: %w", err)
	}

	slog.Debug("resolved default branch", "branch", branch)

	return branch, nil
}
