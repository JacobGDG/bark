package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func NewCloneCommand(deps *Deps) *cobra.Command {
	var stay bool
	var worktreeRootDir string

	cmd := &cobra.Command{
		Use:   "clone <url>",
		Short: "Clone repo into worktrees root",
		Long: `Bare-clone <url> into <worktree-root>/<repo>/ with a worktree for the default
branch, then open it in a tmux session unless --stay.

Full manual: man bark-clone`,
		Example: `  bark clone git@github.com:JacobGDG/bark.git --stay
  cd "$(bark config root)/bark/main"`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			url := args[0]

			if err := deps.git.RemoteExists(url); err != nil {
				return err
			}

			repoName, err := repoNameFromURL(url)
			if err != nil {
				return err
			}

			root := resolveWorktreeRoot(deps, worktreeRootDir)

			repoPath := filepath.Join(root, repoName)

			if _, err := os.Stat(repoPath); err == nil {
				return fmt.Errorf("repo clone location already exists: %s", repoPath)
			}

			// --config writes the refspec before the initial transfer, so
			// the remote-tracking refs exist without a separate fetch.
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

			worktreePath := filepath.Join(repoPath, branch)
			err = deps.git.WorktreeAdd(repoPath, worktreePath, branch)
			if err != nil {
				return err
			}

			finishWorktreeSetup(deps, worktreeSetup{
				repoName:     repoName,
				repoPath:     repoPath,
				branch:       branch,
				worktreePath: worktreePath,
				setUpstream:  true,
				stay:         stay,
			})

			return nil
		},
	}

	cmd.Flags().BoolVar(&stay, "stay", false, "Skip tmux session creation after clone")
	cmd.Flags().StringVar(&worktreeRootDir, "worktree-root-dir", "", "Override the worktree root directory for this clone")
	return cmd
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
