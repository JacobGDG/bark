package cli

import (
	"fmt"
	"log/slog"
)

// worktreeSetup holds the post-worktree-add bookkeeping shared by clone and
// add: linking the branch upstream, indexing the new directories with
// zoxide, and opening a tmux session for the new worktree.
type worktreeSetup struct {
	repoName     string
	repoPath     string
	branch       string
	worktreePath string
	setUpstream  bool
	stay         bool
}

func finishWorktreeSetup(deps *Deps, s worktreeSetup) {
	if s.setUpstream {
		if err := deps.git.SetBranchUpstream(s.repoPath, s.branch); err != nil {
			slog.Warn("failed to set branch upstream", "repo", s.repoPath, "branch", s.branch, "err", err)
		}
	}

	if deps.zoxide != nil {
		for _, dir := range []string{s.worktreePath, s.repoPath} {
			if err := deps.zoxide.Add(dir); err != nil {
				slog.Warn("Attempted to add dir to zoxide and failed", "dir", dir, "message", err)
			}
		}
	}

	if s.stay || deps.tmux == nil {
		return
	}

	session := fmt.Sprintf("%s/%s", s.repoName, s.branch)
	if err := deps.tmux.NewSession(session, s.worktreePath); err != nil {
		slog.Error("failed to create tmux session", "session", session, "err", err)
		return
	}
	if err := deps.tmux.SwitchClient(session); err != nil {
		slog.Error("failed to switch to tmux session", "session", session, "err", err)
		return
	}
	if deps.cfg.TmuxStartupCommand != "" {
		if err := deps.tmux.SendKeys(session, deps.cfg.TmuxStartupCommand); err != nil {
			slog.Warn("startup command failed", "session", session, "err", err)
		}
	}
}

func resolveWorktreeRoot(deps *Deps, override string) string {
	if override != "" {
		return override
	}
	return deps.cfg.WorktreeRootDir
}
