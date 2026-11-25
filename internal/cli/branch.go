package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/JacobGDG/bark/internal/git"
)

type branchState struct {
	local    bool
	remote   bool
	worktree *git.Worktree
}

func (b branchState) exists() bool { return b.local || b.remote }

func lookupBranch(deps *Deps, repoPath, branch string, worktrees []git.Worktree) (branchState, error) {
	var state branchState
	var err error

	state.worktree = findWorktree(worktrees, branch)

	if state.local, err = deps.git.LocalBranchExists(repoPath, branch); err != nil {
		return state, err
	}
	if state.remote, err = deps.git.RemoteBranchExists(repoPath, branch); err != nil {
		return state, err
	}
	return state, nil
}

func findWorktree(worktrees []git.Worktree, branch string) *git.Worktree {
	for i := range worktrees {
		if worktrees[i].Branch == branch {
			return &worktrees[i]
		}
	}
	return nil
}

func defaultBaseRef(deps *Deps, repoPath string) (string, error) {
	def, err := deps.git.DefaultBranch(repoPath)
	if err != nil {
		return "", err
	}
	remote, err := deps.git.RemoteBranchExists(repoPath, def)
	if err != nil {
		return "", err
	}
	if remote {
		return "origin/" + def, nil
	}
	return def, nil
}

const maxIncrement = 1000

func nextFreeBranch(deps *Deps, repoPath, branch string, worktrees []git.Worktree) (string, error) {
	for i := 0; i <= maxIncrement; i++ {
		candidate := branch
		if i > 0 {
			candidate = fmt.Sprintf("%s-%d", branch, i)
		}

		state, err := lookupBranch(deps, repoPath, candidate, worktrees)
		if err != nil {
			return "", err
		}
		if state.exists() || state.worktree != nil {
			continue
		}
		if _, err := os.Stat(filepath.Join(repoPath, candidate)); err == nil {
			continue
		}
		return candidate, nil
	}
	return "", fmt.Errorf("no free branch name found for %s after %d attempts", branch, maxIncrement)
}
