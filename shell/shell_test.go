package shell

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCheck(t *testing.T) {
	tests := []struct {
		name    string
		bin     string
		wantErr bool
	}{
		{
			name:    "binary exists",
			bin:     "echo",
			wantErr: false,
		},
		{
			name:    "binary not found",
			bin:     "definitelynotabinary",
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := Check(test.bin)

			if test.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestRealShell_Run(t *testing.T) {
	tests := []struct {
		name         string
		bin          string
		args         []string
		wantOut      string
		wantExitCode int
		wantStderr   string
		wantErr      bool
	}{
		{
			name:    "success — returns stdout",
			bin:     "echo",
			args:    []string{"hello"},
			wantOut: "hello",
		},
		{
			name:    "trims trailing newline",
			bin:     "echo",
			args:    []string{"hello\n"},
			wantOut: "hello",
		},
		{
			name:         "failure — returns CmdError with exit code and stderr",
			bin:          "sh",
			args:         []string{"-c", "echo 'something failed' >&2; exit 42"},
			wantExitCode: 42,
			wantStderr:   "something failed",
			wantErr:      true,
		},
		{
			name:         "binary not found",
			bin:          "definitelynotabinary",
			args:         nil,
			wantExitCode: -1,
			wantErr:      true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			sh := &RealShell{}
			out, err := sh.Run(test.bin, test.args...)

			assert.Equal(t, test.wantOut, out)

			if test.wantErr {
				require.NotNil(t, err)
				if test.wantExitCode != 0 {
					assert.Equal(t, test.wantExitCode, err.ExitCode)
				}
				if test.wantStderr != "" {
					assert.Equal(t, test.wantStderr, err.Stderr)
				}
			} else {
				require.Nil(t, err)
			}
		})
	}
}
