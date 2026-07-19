package cli

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/JacobGDG/bark/config"
)

func TestNewConfigCommand(t *testing.T) {
	tests := []struct {
		name    string
		arg     string
		wantOut string
		wantErr bool
	}{
		{
			name:    "root key prints WorktreeRootDir",
			arg:     "root",
			wantOut: "/test/worktrees\n",
			wantErr: false,
		},
		{
			name:    "unknown key returns error",
			arg:     "unknown",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{WorktreeRootDir: "/test/worktrees"}
			cmd := NewConfigCommand(cfg)

			var out bytes.Buffer
			// set outpu to be saved to out, rather than stdout
			cmd.SetOut(&out)

			// given we error, we want to hide this error for cleaner testing
			cmd.SetErr(&bytes.Buffer{})
			cmd.SetArgs([]string{tt.arg})

			err := cmd.Execute()

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantOut, out.String())
			}
		})
	}
}
