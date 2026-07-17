package cli

import (
	"github.com/JacobGDG/bark/config"
	"github.com/JacobGDG/bark/git"
	"github.com/JacobGDG/bark/shell"
)

type Deps struct {
	cfg *config.Config
	git git.Git
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

	return &Deps{
		cfg: cfg,
		git: g,
	}, nil
}
