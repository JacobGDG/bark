package git

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/JacobGDG/bark/internal/shell"
)

type Git interface {
	Clone(url, dest string, args ...string) error
	WorktreeAdd(repoPath, dest string, args ...string) error
	DefaultBranch(repoPath string) (string, error)
	// SetBranchUpstream ensures a given local branch is correctly linked to a branch in origin,
	// without this any interactions with origin would not know where to point.
	SetBranchUpstream(repoPath, branch string) error
	// RemoteExists checks that url points to a reachable repository, without
	// transferring any objects.
	RemoteExists(url string) error
	// LocalBranchExists reports whether branch exists as a local branch in repoPath.
	LocalBranchExists(repoPath, branch string) (bool, error)
	// RemoteBranchExists reports whether branch exists on origin in repoPath.
	RemoteBranchExists(repoPath, branch string) (bool, error)
	Fetch(repoPath string) error
	// Worktrees lists the worktrees of repoPath, excluding the bare repo itself.
	Worktrees(repoPath string) ([]Worktree, error)
	CommitTime(repoPath, ref string) (time.Time, error)
	// AheadBehind counts the commits on ref that aren't on base (ahead), and
	// the commits on base that aren't on ref (behind).
	AheadBehind(repoPath, ref, base string) (ahead, behind int, err error)
	IsAncestor(repoPath, ref, base string) (bool, error)
	// ChangedFiles counts uncommitted changes, including untracked files, in a worktree.
	ChangedFiles(worktreePath string) (int, error)
}

type Worktree struct {
	Path string
	// Branch is the short branch name, or empty for a detached HEAD.
	Branch string
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

func (g *RealGit) Clone(url, dest string, args ...string) error {
	_, err := g.shell.Run("git", append([]string{"clone", url, dest}, args...)...)
	if err != nil {
		return fmt.Errorf("failed to clone %s: %w", url, err)
	}
	slog.Debug("repository cloned", "url", url, "dest", dest)
	return nil
}

func (g *RealGit) WorktreeAdd(repoPath, dest string, args ...string) error {
	_, err := g.shell.Run("git", append([]string{"-C", repoPath, "worktree", "add", dest}, args...)...)
	if err != nil {
		return fmt.Errorf("failed to add worktree at %s: %w", dest, err)
	}
	slog.Debug("worktree added", "repo", repoPath, "dest", dest)
	return nil
}

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

func (g *RealGit) RemoteExists(url string) error {
	_, err := g.shell.Run("git", "ls-remote", url, "HEAD")
	if err != nil {
		return fmt.Errorf("remote repository not found or inaccessible %s: %w", url, err)
	}

	slog.Debug("remote repository verified", "url", url)

	return nil
}

func (g *RealGit) LocalBranchExists(repoPath, branch string) (bool, error) {
	_, err := g.shell.Run("git", "-C", repoPath, "show-ref", "--verify", "--quiet", "refs/heads/"+branch)
	if err == nil {
		slog.Debug("local branch found", "repo", repoPath, "branch", branch)
		return true, nil
	}
	if err.ExitCode == 1 {
		slog.Debug("local branch not found", "repo", repoPath, "branch", branch)
		return false, nil
	}
	return false, fmt.Errorf("could not check local branch %s: %w", branch, err)
}

func (g *RealGit) RemoteBranchExists(repoPath, branch string) (bool, error) {
	_, err := g.shell.Run("git", "-C", repoPath, "show-ref", "--verify", "--quiet", "refs/remotes/origin/"+branch)
	if err == nil {
		slog.Debug("remote branch found", "repo", repoPath, "branch", branch)
		return true, nil
	}
	if err.ExitCode == 1 {
		slog.Debug("remote branch not found", "repo", repoPath, "branch", branch)
		return false, nil
	}
	return false, fmt.Errorf("could not check remote branch %s: %w", branch, err)
}

func (g *RealGit) Fetch(repoPath string) error {
	_, err := g.shell.Run("git", "-C", repoPath, "fetch", "--prune", "origin")
	if err != nil {
		return fmt.Errorf("failed to fetch origin: %w", err)
	}
	slog.Debug("fetched origin", "repo", repoPath)
	return nil
}

func (g *RealGit) Worktrees(repoPath string) ([]Worktree, error) {
	out, err := g.shell.Run("git", "-C", repoPath, "worktree", "list", "--porcelain")
	if err != nil {
		return nil, fmt.Errorf("could not list worktrees: %w", err)
	}
	return parseWorktrees(out), nil
}

func parseWorktrees(out string) []Worktree {
	var worktrees []Worktree
	for _, record := range strings.Split(out, "\n\n") {
		var wt Worktree
		bare := false
		for _, line := range strings.Split(record, "\n") {
			key, value, _ := strings.Cut(line, " ")
			switch key {
			case "worktree":
				wt.Path = value
			case "branch":
				wt.Branch = strings.TrimPrefix(value, "refs/heads/")
			case "bare":
				bare = true
			}
		}
		if wt.Path != "" && !bare {
			worktrees = append(worktrees, wt)
		}
	}
	return worktrees
}

func (g *RealGit) CommitTime(repoPath, ref string) (time.Time, error) {
	out, err := g.shell.Run("git", "-C", repoPath, "log", "-1", "--format=%cI", ref, "--")
	if err != nil {
		return time.Time{}, fmt.Errorf("could not read commit time of %s: %w", ref, err)
	}
	t, perr := time.Parse(time.RFC3339, out)
	if perr != nil {
		return time.Time{}, fmt.Errorf("could not parse commit time %q: %w", out, perr)
	}
	return t, nil
}

func (g *RealGit) AheadBehind(repoPath, ref, base string) (int, int, error) {
	out, err := g.shell.Run("git", "-C", repoPath, "rev-list", "--left-right", "--count", ref+"..."+base)
	if err != nil {
		return 0, 0, fmt.Errorf("could not compare %s with %s: %w", ref, base, err)
	}
	fields := strings.Fields(out)
	if len(fields) != 2 {
		return 0, 0, fmt.Errorf("unexpected rev-list output: %q", out)
	}
	ahead, aerr := strconv.Atoi(fields[0])
	behind, berr := strconv.Atoi(fields[1])
	if aerr != nil || berr != nil {
		return 0, 0, fmt.Errorf("unexpected rev-list output: %q", out)
	}
	return ahead, behind, nil
}

func (g *RealGit) IsAncestor(repoPath, ref, base string) (bool, error) {
	_, err := g.shell.Run("git", "-C", repoPath, "merge-base", "--is-ancestor", ref, base)
	if err == nil {
		return true, nil
	}
	if err.ExitCode == 1 {
		return false, nil
	}
	return false, fmt.Errorf("could not check whether %s is merged into %s: %w", ref, base, err)
}

func (g *RealGit) ChangedFiles(worktreePath string) (int, error) {
	out, err := g.shell.Run("git", "-C", worktreePath, "status", "--porcelain")
	if err != nil {
		return 0, fmt.Errorf("could not read status of %s: %w", worktreePath, err)
	}
	if out == "" {
		return 0, nil
	}
	return len(strings.Split(out, "\n")), nil
}
