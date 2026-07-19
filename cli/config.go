package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/JacobGDG/bark/config"
)

func NewConfigCommand(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "config <key>",
		Short: "Get current config values",

		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			arg := args[0]

			switch arg {
			case "root":
				// By using the custom writer we can test the output
				fmt.Fprintln(cmd.OutOrStdout(), cfg.WorktreeRootDir)
			default:
				return fmt.Errorf("Config key not found: %s", arg)
			}

			return nil
		},
	}
}
