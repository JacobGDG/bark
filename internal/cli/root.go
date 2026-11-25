package cli

import (
	"github.com/JacobGDG/bark/internal/config"
	"github.com/spf13/cobra"
)

// NewRootCommand constructs the top-level "bark" command and all its subcommands.
func NewRootCommand(version string) (*cobra.Command, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}
	deps, err := createDeps(cfg)
	if err != nil {
		return nil, err
	}

	rootCommand := &cobra.Command{
		Use:     "bark",
		Version: version,
		Short:   "Git Worktree manager",
		Long: `Git Worktree manager: one bare clone per repo, one worktree per branch,
under ~/worktrees by default.

Full manual: man bark
Each command has its own: man bark-<command>
Config file format: man 5 bark`,
		// Runs after argument validation, so usage is still shown for bad
		// arguments but not for runtime failures, keeping stderr readable in
		// scripts.
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			cmd.SilenceUsage = true
		},
	}

	rootCommand.AddCommand(
		NewCloneCommand(deps),
		NewAddCommand(deps),
		NewReposCommand(deps),
		NewStatusCommand(deps),

		NewConfigCommand(cfg),
		// switch - start or switch to a session focussed on a given worktree
		// purge - Remove all worktrees whos branches have been merged
		//			 - --all do this for all pulled repos
		// update - git pull the current repo and all branches
		//				- --all for all repos
		// tui - a UI for managing worktrees. depending on location called it will change behaviour
		//     - if in worktree root it will show all, with keymaps to call clone, switch ect
		//		 - if in repo root, show worktrees for given repo
		//		 - --tmux allows opening a popup for either of the above
	)

	return rootCommand, nil
}
