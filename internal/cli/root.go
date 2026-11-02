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
		Long: `Git Worktree manager.

bark clones each repo as a bare git repo and checks out one git worktree per
branch, all under a single worktree root:

  <worktree-root>/
    <repo>/
      .git/      bare clone (no files checked out here)
      main/      worktree for the default branch
      <branch>/  one worktree per branch added with "bark add"

Each <branch> directory is a normal git worktree: cd into it and run git as
usual. Branch names containing "/" become nested directories
(feature/x -> <repo>/feature/x/).

Worktree root:
  Defaults to ~/worktrees. Print the current value with:

    bark config root

Config file:
  $XDG_CONFIG_HOME/bark/bark.toml (default ~/.config/bark/bark.toml), TOML.
  Unknown keys are an error.

    root_dir = "/path/to/worktrees"

    [tmux]
    startup_command = "nvim ."

Side effects:
  "clone" and "add" register the new directories with zoxide (if installed)
  and, unless --stay is passed, create a tmux session named <repo>/<branch> and
  switch the current tmux client to it. Pass --stay when scripting or running
  non-interactively.`,
	}

	rootCommand.AddCommand(
		NewCloneCommand(deps),
		NewAddCommand(deps),

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
		// status - fetch stats about a current/given repo
		//        - if outdated, age, local changes, no. worktrees etc
		//				- likely used to inform other commands
	)

	return rootCommand, nil
}
