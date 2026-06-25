package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/JacobGDG/bark/config"
	"github.com/JacobGDG/bark/git"
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
			url := args[0]

			repoName, err := git.RepoNameFromURL(url)
			if err != nil {
				return err
			}

			repoPath := filepath.Join(cfg.WorktreeRootDir, repoName)
			bareClone := filepath.Join(repoPath, ".git")

			if err := os.MkdirAll(repoPath, 0755); err != nil {
				return fmt.Errorf("could not create repo directory: %w", err)
			}

			g, err := git.New()
			if err != nil {
				return err
			}

			if err := g.Clone(url, bareClone, "--bare"); err != nil {
				return err
			}

			if err := g.SetRemoteHead(bareClone); err != nil {
				return err
			}

			branch, err := g.DefaultBranch(bareClone)
			if err != nil {
				return err
			}

			worktreePath := filepath.Join(repoPath, branch)
			return g.WorktreeAdd(repoPath, worktreePath, "-b", branch, "origin/"+branch)
		},
	}
}
