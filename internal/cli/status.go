package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"path/filepath"
	"strconv"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"
)

// Pointer fields are null in JSON when they don't apply, e.g. ahead/behind
// for a branch with no remote copy.
type branchStatus struct {
	Repo         string     `json:"repo"`
	Branch       string     `json:"branch"`
	Exists       bool       `json:"exists"`
	Local        bool       `json:"local"`
	Remote       bool       `json:"remote"`
	Worktree     *string    `json:"worktree"`
	ChangedFiles *int       `json:"changed_files"`
	Ahead        *int       `json:"ahead"`
	Behind       *int       `json:"behind"`
	LastCommit   *time.Time `json:"last_commit"`
	Merged       *bool      `json:"merged"`
}

func NewStatusCommand(deps *Deps) *cobra.Command {
	var asJSON bool
	var fetch bool
	var worktreeRootDir string

	cmd := &cobra.Command{
		Use:   "status <repo|.> <branch>",
		Short: "Show where a branch exists and its state",
		Long: `Show whether <branch> exists locally, on origin and as a worktree, plus
uncommitted changes, ahead/behind origin, last commit and whether it's merged
into the default branch. Uses the last fetch unless --fetch. Pass "." as
<repo> to use the current repo.

Full manual: man bark-status`,
		Example: `  bark status bark my-feature
  bark status . my-feature --fetch --json | jq .exists`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			root := resolveWorktreeRoot(deps, worktreeRootDir)

			repoName, err := resolveRepo(root, args[0])
			if err != nil {
				return err
			}
			repoPath := filepath.Join(root, repoName)

			if fetch {
				if err := deps.git.Fetch(repoPath); err != nil {
					slog.Warn("fetch failed, using existing remote refs", "repo", repoPath, "err", err)
				}
			}

			status, err := getBranchStatus(deps, repoName, repoPath, args[1])
			if err != nil {
				return err
			}

			if asJSON {
				enc := json.NewEncoder(cmd.OutOrStdout())
				enc.SetIndent("", "  ")
				return enc.Encode(status)
			}
			return printBranchStatus(cmd.OutOrStdout(), status)
		},
	}

	cmd.Flags().BoolVar(&asJSON, "json", false, "Print the status as JSON")
	cmd.Flags().BoolVar(&fetch, "fetch", false, "Fetch origin before checking")
	cmd.Flags().StringVar(&worktreeRootDir, "worktree-root-dir", "", "Look for the repo under this directory instead of the configured worktree root")
	return cmd
}

func getBranchStatus(deps *Deps, repoName, repoPath, branch string) (branchStatus, error) {
	status := branchStatus{Repo: repoName, Branch: branch}

	worktrees, err := deps.git.Worktrees(repoPath)
	if err != nil {
		return status, err
	}
	state, err := lookupBranch(deps, repoPath, branch, worktrees)
	if err != nil {
		return status, err
	}

	status.Exists = state.exists()
	status.Local = state.local
	status.Remote = state.remote

	if state.worktree != nil {
		path := state.worktree.Path
		status.Worktree = &path
		changed, err := deps.git.ChangedFiles(path)
		if err != nil {
			return status, err
		}
		status.ChangedFiles = &changed
	}

	if !status.Exists {
		return status, nil
	}

	remoteRef := "refs/remotes/origin/" + branch
	tip := remoteRef
	if state.local {
		tip = "refs/heads/" + branch
	}

	if state.local && state.remote {
		ahead, behind, err := deps.git.AheadBehind(repoPath, tip, remoteRef)
		if err != nil {
			return status, err
		}
		status.Ahead, status.Behind = &ahead, &behind
	}

	lastCommit, err := deps.git.CommitTime(repoPath, tip)
	if err != nil {
		return status, err
	}
	status.LastCommit = &lastCommit

	base, err := defaultBaseRef(deps, repoPath)
	if err != nil {
		return status, err
	}
	merged, err := deps.git.IsAncestor(repoPath, tip, base)
	if err != nil {
		return status, err
	}
	status.Merged = &merged

	return status, nil
}

func printBranchStatus(w io.Writer, s branchStatus) error {
	yesNo := func(b bool) string {
		if b {
			return "yes"
		}
		return "no"
	}

	worktree := "-"
	if s.Worktree != nil {
		worktree = *s.Worktree
	}
	changed := "-"
	if s.ChangedFiles != nil {
		changed = strconv.Itoa(*s.ChangedFiles)
	}
	ahead, behind := "-", "-"
	if s.Ahead != nil {
		ahead, behind = strconv.Itoa(*s.Ahead), strconv.Itoa(*s.Behind)
	}
	lastCommit := "-"
	if s.LastCommit != nil {
		lastCommit = s.LastCommit.Format(time.RFC3339)
	}
	merged := "-"
	if s.Merged != nil {
		merged = yesNo(*s.Merged)
	}

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	rows := [][2]string{
		{"repo", s.Repo},
		{"branch", s.Branch},
		{"exists", yesNo(s.Exists)},
		{"local", yesNo(s.Local)},
		{"remote", yesNo(s.Remote)},
		{"worktree", worktree},
		{"changed_files", changed},
		{"ahead", ahead},
		{"behind", behind},
		{"last_commit", lastCommit},
		{"merged", merged},
	}
	for _, row := range rows {
		fmt.Fprintf(tw, "%s:\t%s\n", row[0], row[1])
	}
	return tw.Flush()
}
