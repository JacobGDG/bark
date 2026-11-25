package cli

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

//TODO:
// .env files or node_modules wont be in git but will want to be in any given worktree

const (
	onExistsError     = "error"
	onExistsOpen      = "open"
	onExistsIncrement = "increment"
)

func NewAddCommand(deps *Deps) *cobra.Command {
	var stay bool
	var printPath bool
	var worktreeRootDir string
	var onExists string
	var from string

	cmd := &cobra.Command{
		Use:   "add <repo|.> <branch>",
		Short: "Add a new worktree for an existing repo",
		Long: `Fetch origin, then create a worktree for <branch> at
<worktree-root>/<repo>/<branch> and open it in a tmux session unless --stay.
New branches start from origin's default branch, or --from. Pass "." as
<repo> to use the current repo.

--on-exists decides what happens when <branch> already has a worktree:
error (exit 2), open it, or increment to the first free <branch>-N.

Full manual: man bark-add`,
		Example: `  bark add bark my-feature --stay
  bark add . my-feature --stay   # from anywhere inside <worktree-root>/bark/
  bark add bark JIR-123-fix-bug --on-exists=increment --print-path`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			repoArg := args[0]
			branch := args[1]

			switch onExists {
			case onExistsError, onExistsOpen, onExistsIncrement:
			default:
				return fmt.Errorf("invalid --on-exists %q: want error, open or increment", onExists)
			}

			root := resolveWorktreeRoot(deps, worktreeRootDir)

			repoName, err := resolveRepo(root, repoArg)
			if err != nil {
				return err
			}

			repoPath := filepath.Join(root, repoName)

			// Fetch first so remote branches, and the base for new ones, are
			// current. A failed fetch (e.g. offline) isn't fatal.
			if err := deps.git.Fetch(repoPath); err != nil {
				slog.Warn("fetch failed, using existing remote refs", "repo", repoPath, "err", err)
			}

			worktrees, err := deps.git.Worktrees(repoPath)
			if err != nil {
				return err
			}

			if onExists == onExistsIncrement {
				branch, err = nextFreeBranch(deps, repoPath, branch, worktrees)
				if err != nil {
					return err
				}
			}

			state, err := lookupBranch(deps, repoPath, branch, worktrees)
			if err != nil {
				return err
			}

			if state.worktree != nil {
				if onExists != onExistsOpen {
					return &ExitCodeError{
						Code: ExitWorktreeExists,
						Err:  fmt.Errorf("worktree already exists: %s", state.worktree.Path),
					}
				}

				finishWorktreeSetup(deps, worktreeSetup{
					repoName:     repoName,
					repoPath:     repoPath,
					branch:       branch,
					worktreePath: state.worktree.Path,
					stay:         stay,
				})
				if printPath {
					fmt.Fprintln(cmd.OutOrStdout(), state.worktree.Path)
				}
				return nil
			}

			worktreePath := filepath.Join(repoPath, branch)
			if _, err := os.Stat(worktreePath); err == nil {
				return fmt.Errorf("%s already exists but isn't a worktree for %s", worktreePath, branch)
			}

			if from != "" && state.exists() {
				slog.Warn("--from ignored: branch already exists", "branch", branch, "from", from)
			}

			switch {
			case state.local:
				err = deps.git.WorktreeAdd(repoPath, worktreePath, branch)
			case state.remote:
				err = deps.git.WorktreeAdd(repoPath, worktreePath, "-b", branch, "origin/"+branch)
			default:
				base := from
				if base == "" {
					if base, err = defaultBaseRef(deps, repoPath); err != nil {
						return err
					}
				}
				// --no-track: the new branch must not adopt its base (e.g.
				// origin/main) as upstream.
				err = deps.git.WorktreeAdd(repoPath, worktreePath, "--no-track", "-b", branch, base)
			}
			if err != nil {
				return err
			}

			finishWorktreeSetup(deps, worktreeSetup{
				repoName:     repoName,
				repoPath:     repoPath,
				branch:       branch,
				worktreePath: worktreePath,
				setUpstream:  !state.local && state.remote,
				stay:         stay,
			})

			if printPath {
				fmt.Fprintln(cmd.OutOrStdout(), worktreePath)
			}
			return nil
		},
	}

	cmd.Flags().BoolVar(&stay, "stay", false, "Skip tmux session creation after add")
	cmd.Flags().BoolVar(&printPath, "print-path", false, "Print the worktree path to stdout")
	cmd.Flags().StringVar(&worktreeRootDir, "worktree-root-dir", "", "Override the worktree root directory for this add")
	cmd.Flags().StringVar(&onExists, "on-exists", onExistsError, "When the branch already has a worktree: error, open or increment")
	cmd.Flags().StringVar(&from, "from", "", "Start a new branch from this ref instead of origin's default branch")
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
