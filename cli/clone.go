package cli

import (
	"fmt"
	"log/slog"
	"os"
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

			repoPath := filepath.Join(deps.cfg.WorktreeRootDir, repoName)
			// TODO: ensure direcotry does not exist
			if err := os.MkdirAll(repoPath, 0755); err != nil {
				return fmt.Errorf("could not create repo directory: %w", err)
			}

			if err := deps.git.Clone(url, filepath.Join(repoPath, ".git"), "--bare"); err != nil {
				return err
			}

			// TODO: `fetch = +refs/heads/*:refs/remotes/origin/*` on origin to allow full fetches

			branch, err := deps.git.DefaultBranch(repoPath)
			if err != nil {
				return err
			}

			worktreePath := filepath.Join(repoPath, branch)
			err = deps.git.WorktreeAdd(repoPath, worktreePath, branch)
			if err != nil {
				return err
			}

			err = deps.git.SetOriginFetch(repoPath)
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
