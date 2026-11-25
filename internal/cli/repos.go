package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

type repoJSON struct {
	Name      string         `json:"name"`
	Path      string         `json:"path"`
	Worktrees []worktreeJSON `json:"worktrees"`
}

type worktreeJSON struct {
	Branch string `json:"branch"`
	Path   string `json:"path"`
}

func NewReposCommand(deps *Deps) *cobra.Command {
	var asJSON bool
	var worktreeRootDir string

	cmd := &cobra.Command{
		Use:   "repos",
		Short: "List repos in the worktree root",
		Long: `Print the name of each bark-managed repo in the worktree root, one per line.
--json adds each repo's path and worktrees.

Full manual: man bark-repos`,
		Example: `  bark repos | fzf
  bark repos --json | jq -r '.[].worktrees[].path'`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			root := resolveWorktreeRoot(deps, worktreeRootDir)

			names, err := listRepos(root)
			if err != nil {
				return err
			}

			out := cmd.OutOrStdout()
			if !asJSON {
				for _, name := range names {
					fmt.Fprintln(out, name)
				}
				return nil
			}

			repos := []repoJSON{}
			for _, name := range names {
				repoPath := filepath.Join(root, name)
				worktrees, err := deps.git.Worktrees(repoPath)
				if err != nil {
					return err
				}

				repo := repoJSON{Name: name, Path: repoPath, Worktrees: []worktreeJSON{}}
				for _, wt := range worktrees {
					repo.Worktrees = append(repo.Worktrees, worktreeJSON{Branch: wt.Branch, Path: wt.Path})
				}
				repos = append(repos, repo)
			}

			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(repos)
		},
	}

	cmd.Flags().BoolVar(&asJSON, "json", false, "Print repos and their worktrees as JSON")
	cmd.Flags().StringVar(&worktreeRootDir, "worktree-root-dir", "", "List repos under this directory instead of the configured worktree root")
	return cmd
}

func listRepos(root string) ([]string, error) {
	entries, err := os.ReadDir(root)
	// A missing root just means nothing has been cloned yet.
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var names []string
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		if validateRepoDir(root, entry.Name()) == nil {
			names = append(names, entry.Name())
		}
	}
	return names, nil
}
