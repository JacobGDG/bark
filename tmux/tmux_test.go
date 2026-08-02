package tmux

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

func TestNewSession(t *testing.T) {
	tests := []struct {
		name      string
		session   string
		dir       string
		inSession bool
		shellErr  *shell.CmdError
		wantName  string
		wantArgs  []string
		wantErr   bool
	}{
		{
			name:      "success",
			session:   "bark:main",
			dir:       "/some/worktrees/bark/main",
			inSession: true,
			shellErr:  nil,
			wantName:  "tmux",
			wantArgs:  []string{"new-session", "-d", "-s", "bark:main", "-c", "/some/worktrees/bark/main"},
			wantErr:   false,
		},
		{
			name:      "shell error is returned",
			session:   "bark:main",
			dir:       "/some/worktrees/bark/main",
			inSession: true,
			shellErr:  &shell.CmdError{Stderr: "duplicate session: bark:main", ExitCode: 1},
			wantName:  "tmux",
			wantArgs:  []string{"new-session", "-d", "-s", "bark:main", "-c", "/some/worktrees/bark/main"},
			wantErr:   true,
		},
		{
			name:      "not in a tmux session",
			session:   "bark:main",
			dir:       "/some/worktrees/bark/main",
			inSession: false,
			shellErr:  nil,
			wantName:  "",
			wantArgs:  nil,
			wantErr:   true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.inSession {
				t.Setenv("TMUX", "/tmp/tmux-1000/default,1234,0")
			} else {
				t.Setenv("TMUX", "")
			}

			sh := &mockShell{err: test.shellErr}
			tm := &RealTmux{shell: sh}

			err := tm.NewSession(test.session, test.dir)

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

func TestSendKeys(t *testing.T) {
	tests := []struct {
		name      string
		session   string
		cmd       string
		inSession bool
		shellErr  *shell.CmdError
		wantName  string
		wantArgs  []string
		wantErr   bool
	}{
		{
			name:      "success",
			session:   "bark/main",
			cmd:       "tmuxifier load-window vimsplit && tmux move-window -t 0 && tmux kill-window -t 1",
			inSession: true,
			shellErr:  nil,
			wantName:  "tmux",
			wantArgs:  []string{"send-keys", "-t", "bark/main", "tmuxifier load-window vimsplit && tmux move-window -t 0 && tmux kill-window -t 1", "Enter"},
			wantErr:   false,
		},
		{
			name:      "shell error is returned",
			session:   "bark/main",
			cmd:       "some-command",
			inSession: true,
			shellErr:  &shell.CmdError{Stderr: "no session", ExitCode: 1},
			wantName:  "tmux",
			wantArgs:  []string{"send-keys", "-t", "bark/main", "some-command", "Enter"},
			wantErr:   true,
		},
		{
			name:      "not in a tmux session",
			session:   "bark/main",
			cmd:       "some-command",
			inSession: false,
			shellErr:  nil,
			wantName:  "",
			wantArgs:  nil,
			wantErr:   true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.inSession {
				t.Setenv("TMUX", "/tmp/tmux-1000/default,1234,0")
			} else {
				t.Setenv("TMUX", "")
			}

			sh := &mockShell{err: test.shellErr}
			tm := &RealTmux{shell: sh}

			err := tm.SendKeys(test.session, test.cmd)

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

func TestSwitchClient(t *testing.T) {
	tests := []struct {
		name      string
		session   string
		inSession bool
		shellErr  *shell.CmdError
		wantName  string
		wantArgs  []string
		wantErr   bool
	}{
		{
			name:      "success",
			session:   "bark:main",
			inSession: true,
			shellErr:  nil,
			wantName:  "tmux",
			wantArgs:  []string{"switch-client", "-t", "bark:main"},
			wantErr:   false,
		},
		{
			name:      "shell error is returned",
			session:   "bark:main",
			inSession: true,
			shellErr:  &shell.CmdError{Stderr: "no client found", ExitCode: 1},
			wantName:  "tmux",
			wantArgs:  []string{"switch-client", "-t", "bark:main"},
			wantErr:   true,
		},
		{
			name:      "not in a tmux session",
			session:   "bark:main",
			inSession: false,
			shellErr:  nil,
			wantName:  "",
			wantArgs:  nil,
			wantErr:   true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.inSession {
				t.Setenv("TMUX", "/tmp/tmux-1000/default,1234,0")
			} else {
				t.Setenv("TMUX", "")
			}

			sh := &mockShell{err: test.shellErr}
			tm := &RealTmux{shell: sh}

			err := tm.SwitchClient(test.session)

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
