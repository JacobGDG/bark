package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// runBark runs the full bark command tree with args and returns its stdout.
func runBark(t *testing.T, args ...string) (*cobra.Command, string, error) {
	t.Helper()
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	t.Setenv("SOURCE_DATE_EPOCH", "0")

	cmd, err := NewRootCommand("vtest")
	require.NoError(t, err)

	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs(args)

	err = cmd.Execute()
	return cmd, out.String(), err
}

func TestGenerateFlag(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		wantParts []string
		wantErr   bool
	}{
		{
			name: "man prints a roff page for bark",
			args: []string{"--generate", "man"},
			wantParts: []string{
				`.TH "BARK" "1" "Jan 1970" "bark vtest" "User Commands"`,
				".SH NAME",
				"<worktree-root>/", // placeholders survive the Markdown renderer
				"bark-add(1)",
			},
		},
		{
			name: "subcommand man page uses the full manual, not --help text",
			args: []string{"add", "--generate", "man"},
			wantParts: []string{
				`.TH "BARK-ADD" "1"`,
				"<repo|.> <branch>", // SYNOPSIS keeps the args from Use
				"brand-new branch is created off the repo's current HEAD",
			},
		},
		{
			name: "every subcommand can generate a man page",
			args: []string{"config", "--generate", "man"},
			wantParts: []string{
				`.TH "BARK-CONFIG" "1"`,
				"tmux.startup_command  command sent to new tmux sessions",
			},
		},
		{
			name:      "--help points at the man page",
			args:      []string{"add", "--help"},
			wantParts: []string{"Full manual: man bark-add"},
		},
		{
			name:    "unknown value returns error",
			args:    []string{"--generate", "html"},
			wantErr: true,
		},
		{
			name:    "argument validation still applies without --generate",
			args:    []string{"add"},
			wantErr: true,
		},
		{
			name:    "unknown subcommands are still rejected",
			args:    []string{"nope"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, out, err := runBark(t, tt.args...)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			for _, part := range tt.wantParts {
				assert.Contains(t, out, part)
			}
		})
	}
}

func TestGenerateManPages(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)

	cmd, _, err := runBark(t, "--generate", "man-pages")
	require.NoError(t, err)

	// Every command gets a page without being listed anywhere.
	walkCommands(cmd, func(c *cobra.Command) {
		if c.IsAvailableCommand() {
			assert.FileExists(t, filepath.Join(dir, manPageName(c)+".1"))
		}
	})

	// Pages render the same way as --generate man (covered above); just
	// check the manual text made it in.
	page, err := os.ReadFile(filepath.Join(dir, "bark-add.1"))
	require.NoError(t, err)
	assert.Contains(t, string(page), "brand-new branch is created off the repo's current HEAD")
}
