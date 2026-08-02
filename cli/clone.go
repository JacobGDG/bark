package cli

import (
	"fmt"
	"log/slog"
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

	cmd := &cobra.Command{
		Use:   "clone <url>",
		Short: "Clone repo into worktrees root",
		Long: `Clone a repo into the worktrees root directory, in a --bare repo and default branch worktree like so:

- <worktree-root>/<repo>/
  - .git/   a --bare clone
  - main/   the default branch

If zoxide is available, it will also add the <repo> and main directories to its database.
If inside a tmux session, it will create a new session and switch to it.
		`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			url := args[0]

			repoName, err := repoNameFromURL(url)
			if err != nil {
				return err
			}

			repoPath := filepath.Join(deps.cfg.WorktreeRootDir, repoName)

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

			if err := deps.git.SetBranchUpstream(repoPath, branch); err != nil {
				return err
			}

			if deps.zoxide != nil {
				for _, dir := range []string{worktreePath, repoPath} {
					if err := deps.zoxide.Add(dir); err != nil {
						slog.Warn("Attempted to add dir to zoxide and failed", "dir", dir, "message", err)
					}
				}
			}

			if !stay && deps.tmux != nil {
				session := fmt.Sprintf("%s/%s", repoName, branch)
				if err := deps.tmux.NewSession(session, worktreePath); err != nil {
					slog.Error("failed to create tmux session", "session", session, "err", err)
				} else if err := deps.tmux.SwitchClient(session); err != nil {
					slog.Error("failed to switch to tmux session", "session", session, "err", err)
				} else if deps.cfg.TmuxStartupCommand != "" {
					if err := deps.tmux.SendKeys(session, deps.cfg.TmuxStartupCommand); err != nil {
						slog.Warn("startup command failed", "session", session, "err", err)
					}
				}
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&stay, "stay", false, "Skip tmux session creation after clone")
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
