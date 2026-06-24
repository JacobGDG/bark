package cli

import (
	"github.com/spf13/cobra"
)

func NewRootCommand(version string) *cobra.Command {
	rootCommand := &cobra.Command{
		Use:     "bark",
		Version: version,
		Short:   "Git Worktree manager",
	}

	return rootCommand
}
