package cli

import (
	"fmt"

	"github.com/JacobGDG/bark/config"
	"github.com/spf13/cobra"
)

func NewCloneCommand(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "clone",
		Short: "Clone repo into worktrees root",
		Long: `Clone a repo into the worktrees root directory, in a --bare repo and default branch worktree like so:
			
- <worktree-root>/<repo>/
  - .git/   a --bare clone
  - main/   the default branch
				`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(cfg.WorktreeRootDir)
		},
	}
}
