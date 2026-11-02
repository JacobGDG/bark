package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

//TODO:
// 2 types of clone, reference and dev. Clone to different places
// clone and checkout a specific version or tag. Useful for reference
// after clone switch open new worktree as a tmux session

func NewCloneCommand(deps *Deps) *cobra.Command {
	var stay bool
	var worktreeRootDir string

	cmd := &cobra.Command{
		Use:   "clone <url>",
		Short: "Clone repo into worktrees root",
		Long: `Clone a repo into the worktrees root directory, in a --bare repo and default branch worktree like so:

- <worktree-root>/<repo>/
  - .git/   a --bare clone
  - main/   worktree for the remote's default branch (named after it)

Side effects:
- If zoxide is installed, the <repo> and default branch directories are added to its database.
- Unless --stay is passed, a tmux session named <repo>/<branch> is created and the
  current tmux client is switched to it, then tmux.startup_command is sent to it.
  Pass --stay when scripting or running non-interactively.
		`,
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
