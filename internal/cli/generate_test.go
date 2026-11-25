package cli

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateFlag(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	t.Setenv("SOURCE_DATE_EPOCH", "0")

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
			cmd, err := NewRootCommand("vtest")
			require.NoError(t, err)

			var out bytes.Buffer
			cmd.SetOut(&out)
			cmd.SetErr(&bytes.Buffer{})
			cmd.SetArgs(tt.args)

			err = cmd.Execute()

			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			for _, part := range tt.wantParts {
				assert.Contains(t, out.String(), part)
			}
		})
	}
}
