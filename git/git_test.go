package git

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/JacobGDG/bark/shell"
)

type mockShell struct {
	gotName string
	gotArgs []string
	err     *shell.CmdError
}

func (m *mockShell) Run(name string, args ...string) (string, *shell.CmdError) {
	m.gotName = name
	m.gotArgs = args
	return "", m.err
}

func TestRemoteExists(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		shellErr *shell.CmdError
		wantName string
		wantArgs []string
		wantErr  bool
	}{
		{
			name:     "reachable repository",
			url:      "git@github.com:JacobGDG/bark.git",
			shellErr: nil,
			wantName: "git",
			wantArgs: []string{"ls-remote", "git@github.com:JacobGDG/bark.git", "HEAD"},
			wantErr:  false,
		},
		{
			name:     "unreachable repository returns error",
			url:      "git@github.com:JacobGDG/thisdoesnotexist.git",
			shellErr: &shell.CmdError{Stderr: "ERROR: Repository not found.", ExitCode: 128},
			wantName: "git",
			wantArgs: []string{"ls-remote", "git@github.com:JacobGDG/thisdoesnotexist.git", "HEAD"},
			wantErr:  true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			sh := &mockShell{err: test.shellErr}
			g := &RealGit{shell: sh}

			err := g.RemoteExists(test.url)

			assert.Equal(t, test.wantName, sh.gotName)
			assert.Equal(t, test.wantArgs, sh.gotArgs)

			if test.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
