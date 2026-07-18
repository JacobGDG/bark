package git

import (
	"fmt"
	"log/slog"

	"github.com/JacobGDG/bark/shell"
)

type Git interface {
	Clone(url, dest string, args ...string) error
	WorktreeAdd(repoPath, dest string, args ...string) error
	DefaultBranch(repoPath string) (string, error)
	SetBranchUpstream(repoPath, branch string) error
}

type RealGit struct {
	shell shell.Shell
}

func check() error {
	err := shell.Check("git")
	if err != nil {
		return err
	}
	return nil
}

func New(sh shell.Shell) (Git, error) {
	if err := check(); err != nil {
		return nil, err
	}
	return &RealGit{shell: sh}, nil
}

// Clone clones url to dest, passing any extra args after the positional arguments.
func (g *RealGit) Clone(url, dest string, args ...string) error {
	_, err := g.shell.Run("git", append([]string{"clone", url, dest}, args...)...)

	if err == nil {
		slog.Debug("repository cloned", "url", url, "dest", dest)
	}

	return err
}

// WorktreeAdd adds a worktree at dest from repoPath, passing any extra args after the positional arguments.
func (g *RealGit) WorktreeAdd(repoPath, dest string, args ...string) error {
	_, err := g.shell.Run("git", append([]string{"-C", repoPath, "worktree", "add", dest}, args...)...)

	if err == nil {
		slog.Debug("worktree added", "repo", repoPath, "dest", dest)
	}

	return err
}

// DefaultBranch returns the default branch name for the repo at repoPath.
func (g *RealGit) DefaultBranch(repoPath string) (string, error) {
	branch, err := g.shell.Run("git", "-C", repoPath, "branch", "--show-current")
	if err != nil {
		return "", fmt.Errorf("could not determine default branch: %w", err)
	}

	slog.Debug("resolved default branch", "branch", branch)

	return branch, nil
}

func (g *RealGit) SetBranchUpstream(repoPath, branch string) error {
	upstream := fmt.Sprintf("origin/%s", branch)
	_, err := g.shell.Run("git", "-C", repoPath, "branch", "--set-upstream-to", upstream, branch)
	if err != nil {
		return fmt.Errorf("could not set upstream for %s: %w", branch, err)
	}

	slog.Debug("branch upstream set", "branch", branch, "upstream", upstream)

	return nil
}
