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
		value     string
		wantParts []string
		wantErr   bool
	}{
		{
			name:  "man prints a roff page for bark",
			value: "man",
			wantParts: []string{
				`.TH "BARK" "1" "Jan 1970" "bark vtest" "User Commands"`,
				".SH NAME",
				"<worktree-root>/", // placeholders survive the Markdown renderer
				"bark-add(1)",
			},
		},
		{
			name:    "unknown value returns error",
			value:   "html",
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
			cmd.SetArgs([]string{"--generate", tt.value})

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
