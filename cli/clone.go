package cli

import (
	"fmt"
	"log/slog"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func NewCloneCommand(deps *Deps) *cobra.Command {
	return &cobra.Command{
		Use:   "clone <url>",
		Short: "Clone repo into worktrees root",
		Long: `Clone a repo into the worktrees root directory, in a --bare repo and default branch worktree like so:

- <worktree-root>/<repo>/
  - .git/   a --bare clone
  - main/   the default branch
		`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			url := args[0]

			repoName, err := repoNameFromURL(url)
			if err != nil {
				return err
			}

			// --config writes the refspec before the initial transfer, so
			// the remote-tracking refs exist without a separate fetch.
			repoPath := filepath.Join(deps.cfg.WorktreeRootDir, repoName)
			err = deps.git.Clone(url, filepath.Join(repoPath, ".git"),
				"--bare",
				"--config", "remote.origin.fetch=+refs/heads/*:refs/remotes/origin/*")
			if err != nil {
				return err
			}

			branch, err := deps.git.DefaultBranch(repoPath)
			if err != nil {
				return err
			}

			if err := deps.git.SetBranchUpstream(repoPath, branch); err != nil {
				return err
			}

			worktreePath := filepath.Join(repoPath, branch)
			err = deps.git.WorktreeAdd(repoPath, worktreePath, branch)
			if err != nil {
				return err
			}

			if deps.zoxide != nil {
				for _, dir := range []string{worktreePath, repoPath} {
					if err := deps.zoxide.Add(dir); err != nil {
						slog.Warn("Attempted to add dir to zoxide and failed", "dir", dir, "message", err)
					}
				}
			}

			return nil
		},
	}
}

func repoNameFromURL(rawURL string) (string, error) {
	parts := strings.Split(rawURL, "/")
	length := len(parts)
	lastItem := parts[length-1]

	if length < 2 || lastItem == "" {
		return "", fmt.Errorf("could not derive repo name from URL: %s", rawURL)
	}

	name := strings.TrimSuffix(lastItem, ".git")
	if name == "" || name == "." {
		return "", fmt.Errorf("could not derive repo name from URL: %s", rawURL)
	}
	return name, nil
}
