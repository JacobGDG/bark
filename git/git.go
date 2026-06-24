package git

import (
	"fmt"
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

// Clone creates a bare clone of url at dest.
func (c *Client) Clone(url, dest string) error {
	_, err := c.git.Run("clone", "--bare", url, dest)

	return err
}

// WorktreeAdd adds a worktree for branch at the given path inside repoPath.
func (c *Client) WorktreeAdd(repoPath, branch string) error {
	_, err := c.git.Run("-C", repoPath, "worktree", "add", branch)

	return err
}

// DefaultBranch returns the default branch name for the repo at repoPath.
func (c *Client) DefaultBranch(repoPath string) (string, error) {
	ref, err := c.git.Run("-C", repoPath, "symbolic-ref", "refs/remotes/origin/HEAD")
	if err != nil {
		return "", fmt.Errorf("could not determine default branch: %w", err)
	}
	return path.Base(strings.TrimSpace(ref)), nil
}
