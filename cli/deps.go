package cli

import (
	"github.com/JacobGDG/bark/config"
	"github.com/JacobGDG/bark/git"
	"github.com/JacobGDG/bark/shell"
	"github.com/JacobGDG/bark/tmux"
	"github.com/JacobGDG/bark/zoxide"
)

// Deps is a container for all external dependencies needed for Bark to run.
// Constructing this early means we can fail fast if any required commands are not available.
type Deps struct {
	cfg    *config.Config
	git    git.Git
	zoxide zoxide.Zoxide
	tmux   tmux.Tmux
}

func createDeps(cfg *config.Config) (*Deps, error) {
	sh, err := shell.New()
	if err != nil {
		return nil, err
	}

	g, err := git.New(sh)
	if err != nil {
		return nil, err
	}

	z, _ := zoxide.New(sh)
	t, _ := tmux.New(sh)

	return &Deps{
		cfg:    cfg,
		git:    g,
		zoxide: z,
		tmux:   t,
	}, nil
}
