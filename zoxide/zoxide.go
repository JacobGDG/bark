package zoxide

import (
	"log/slog"

	"github.com/JacobGDG/bark/shell"
)

type Zoxide interface {
	Add(dir string) error
}

type RealZoxide struct {
	shell shell.Shell
}

func check() error {
	err := shell.Check("zoxide")
	if err != nil {
		return err
	}
	return nil
}

func New(sh shell.Shell) (Zoxide, error) {
	if err := check(); err != nil {
		return nil, err
	}
	return &RealZoxide{shell: sh}, nil
}

// Clone clones url to dest, passing any extra args after the positional arguments.
func (g *RealZoxide) Add(dir string) error {
	_, err := g.shell.Run("zoxide", "add", dir)

	if err == nil {
		slog.Debug("directory added to zoxide", "dir", dir)
	}

	return err
}
