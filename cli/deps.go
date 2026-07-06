package cli

import (
	"github.com/JacobGDG/bark/config"
	"github.com/JacobGDG/bark/git"
)

type Deps struct {
	cfg *config.Config
	git git.Git
}

func createDeps(cfg *config.Config) (*Deps, error) {
	g, err := git.New()
	if err != nil {
		return nil, err
	}

	return &Deps{
		cfg: cfg,
		git: g,
	}, nil
}
