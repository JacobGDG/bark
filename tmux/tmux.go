package tmux

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/JacobGDG/bark/shell"
)

type Tmux interface {
	NewSession(name, dir string) error
	SwitchClient(name string) error
}

type RealTmux struct {
	shell shell.Shell
}

func check() error {
	return shell.Check("tmux")
}

func inTmuxSession() bool {
	return os.Getenv("TMUX") != ""
}

func New(sh shell.Shell) (Tmux, error) {
	if err := check(); err != nil {
		return nil, err
	}
	return &RealTmux{shell: sh}, nil
}

func (t *RealTmux) NewSession(name, dir string) error {
	if !inTmuxSession() {
		return fmt.Errorf("not in a tmux session")
	}
	_, err := t.shell.Run("tmux", "new-session", "-d", "-s", name, "-c", dir)
	if err != nil {
		return fmt.Errorf("failed to create tmux session %s: %w", name, err)
	}
	slog.Debug("tmux session created", "session", name, "dir", dir)
	return nil
}

func (t *RealTmux) SwitchClient(name string) error {
	if !inTmuxSession() {
		return fmt.Errorf("not in a tmux session")
	}
	_, err := t.shell.Run("tmux", "switch-client", "-t", name)
	if err != nil {
		return fmt.Errorf("failed to switch to tmux session %s: %w", name, err)
	}
	slog.Debug("switched to tmux session", "session", name)
	return nil
}
