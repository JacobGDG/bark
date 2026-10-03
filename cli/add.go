package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

//TODO:
// .env files or node_modules wont be in git but will want to be in any given worktree

func NewAddCommand(deps *Deps) *cobra.Command {
	var stay bool
	var worktreeRootDir string

	cmd := &cobra.Command{
		Use:   "add <repo|.> <branch>",
		Short: "Add a new worktree for an existing repo",
		Long: `Create a new worktree for <branch> inside an existing cloned repo:

- Pass "." as <repo> to detect the repo from the current working directory
  (it must be somewhere inside <worktree-root>/<repo>/...).
- Otherwise <repo> must exactly match a directory under the worktree root.
- If <branch> already exists locally, a worktree is attached to it.
- If <branch> exists on origin but not locally, a local branch is created tracking it.
- Otherwise, a brand-new branch is created off the repo's current HEAD.

If inside a tmux session, it will create a new session and switch to it.
		`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			repoArg := args[0]
			branch := args[1]

			root := resolveWorktreeRoot(deps, worktreeRootDir)

			repoName, err := resolveRepo(root, repoArg)
			if err != nil {
				return err
			}

			repoPath := filepath.Join(root, repoName)
			worktreePath := filepath.Join(repoPath, branch)

			if _, err := os.Stat(worktreePath); err == nil {
				return fmt.Errorf("worktree already exists: %s", worktreePath)
			}

			localExists, err := deps.git.LocalBranchExists(repoPath, branch)
			if err != nil {
				return err
			}

			remoteExists := false
			if !localExists {
				remoteExists, err = deps.git.RemoteBranchExists(repoPath, branch)
				if err != nil {
					return err
				}
			}

			switch {
			case localExists:
				err = deps.git.WorktreeAdd(repoPath, worktreePath, branch)
			case remoteExists:
				err = deps.git.WorktreeAdd(repoPath, worktreePath, "-b", branch, "origin/"+branch)
			default:
				err = deps.git.WorktreeAdd(repoPath, worktreePath, "-b", branch)
			}
			if err != nil {
				return err
			}

			finishWorktreeSetup(deps, worktreeSetup{
				repoName:     repoName,
				repoPath:     repoPath,
				branch:       branch,
				worktreePath: worktreePath,
				setUpstream:  remoteExists,
				stay:         stay,
			})

			return nil
		},
	}

	cmd.Flags().BoolVar(&stay, "stay", false, "Skip tmux session creation after add")
	cmd.Flags().StringVar(&worktreeRootDir, "worktree-root-dir", "", "Override the worktree root directory for this add")
	return cmd
}

func resolveRepo(root, repoArg string) (string, error) {
	repo := repoArg

	if repo == "." {
		cwd, err := os.Getwd()
		if err != nil {
			return "", err
		}

		repo, err = repoNameFromCwd(root, cwd)
		if err != nil {
			return "", err
		}
	}

	if err := validateRepoDir(root, repo); err != nil {
		return "", err
	}

	return repo, nil
}

func repoNameFromCwd(root, cwd string) (string, error) {
	rel, err := filepath.Rel(root, cwd)
	if err != nil {
		return "", err
	}

	if rel == "." || strings.HasPrefix(rel, "..") {
		return "", fmt.Errorf("current directory is not inside worktree root %s; pass a repo name instead of \".\"", root)
	}

	parts := strings.Split(rel, string(filepath.Separator))
	return parts[0], nil
}

func validateRepoDir(root, repo string) error {
	repoPath := filepath.Join(root, repo)

	info, err := os.Stat(repoPath)
	if err != nil || !info.IsDir() {
		return fmt.Errorf("no repo found at %s", repoPath)
	}

	if _, err := os.Stat(filepath.Join(repoPath, ".git")); err != nil {
		return fmt.Errorf("%s does not look like a bark-managed repo (missing .git)", repoPath)
	}

	return nil
}
