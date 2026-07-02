package cli

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/JacobGDG/bark/config"
	"github.com/spf13/cobra"
)

func NewCloneCommand(cfg *config.Config) *cobra.Command {
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
			deps, err := buildDeps()
			if err != nil {
				return err
			}

			url := args[0]

			repoName, err := repoNameFromURL(url)
			if err != nil {
				return err
			}

			repoPath := filepath.Join(cfg.WorktreeRootDir, repoName)
			// TODO: ensure direcotry does not exist
			if err := os.MkdirAll(repoPath, 0755); err != nil {
				return fmt.Errorf("could not create repo directory: %w", err)
			}

			if err := deps.git.Clone(url, filepath.Join(repoPath, ".git"), "--bare"); err != nil {
				return err
			}

			branch, err := deps.git.DefaultBranch(repoPath)
			if err != nil {
				return err
			}

			worktreePath := filepath.Join(repoPath, branch)
			return deps.git.WorktreeAdd(repoPath, worktreePath, branch)
		},
	}
}

func repoNameFromURL(rawURL string) (string, error) {
	name := strings.TrimSuffix(path.Base(rawURL), ".git")
	if name == "" || name == "." {
		return "", fmt.Errorf("could not derive repo name from URL: %s", rawURL)
	}
	return name, nil
}
