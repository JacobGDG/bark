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

func TestLocalBranchExists(t *testing.T) {
	tests := []struct {
		name     string
		shellErr *shell.CmdError
		wantArgs []string
		want     bool
		wantErr  bool
	}{
		{
			name:     "branch exists",
			shellErr: nil,
			wantArgs: []string{"-C", "/repo", "show-ref", "--verify", "--quiet", "refs/heads/main"},
			want:     true,
			wantErr:  false,
		},
		{
			name:     "branch does not exist",
			shellErr: &shell.CmdError{Stderr: "", ExitCode: 1},
			wantArgs: []string{"-C", "/repo", "show-ref", "--verify", "--quiet", "refs/heads/main"},
			want:     false,
			wantErr:  false,
		},
		{
			name:     "unexpected failure",
			shellErr: &shell.CmdError{Stderr: "fatal: not a git repository", ExitCode: 128},
			wantArgs: []string{"-C", "/repo", "show-ref", "--verify", "--quiet", "refs/heads/main"},
			want:     false,
			wantErr:  true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			sh := &mockShell{err: test.shellErr}
			g := &RealGit{shell: sh}

			got, err := g.LocalBranchExists("/repo", "main")

			assert.Equal(t, test.wantArgs, sh.gotArgs)
			assert.Equal(t, test.want, got)

			if test.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestRemoteBranchExists(t *testing.T) {
	tests := []struct {
		name     string
		shellErr *shell.CmdError
		wantArgs []string
		want     bool
		wantErr  bool
	}{
		{
			name:     "branch exists",
			shellErr: nil,
			wantArgs: []string{"-C", "/repo", "show-ref", "--verify", "--quiet", "refs/remotes/origin/main"},
			want:     true,
			wantErr:  false,
		},
		{
			name:     "branch does not exist",
			shellErr: &shell.CmdError{Stderr: "", ExitCode: 1},
			wantArgs: []string{"-C", "/repo", "show-ref", "--verify", "--quiet", "refs/remotes/origin/main"},
			want:     false,
			wantErr:  false,
		},
		{
			name:     "unexpected failure",
			shellErr: &shell.CmdError{Stderr: "fatal: not a git repository", ExitCode: 128},
			wantArgs: []string{"-C", "/repo", "show-ref", "--verify", "--quiet", "refs/remotes/origin/main"},
			want:     false,
			wantErr:  true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			sh := &mockShell{err: test.shellErr}
			g := &RealGit{shell: sh}

			got, err := g.RemoteBranchExists("/repo", "main")

			assert.Equal(t, test.wantArgs, sh.gotArgs)
			assert.Equal(t, test.want, got)

			if test.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
