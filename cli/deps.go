package cli

import (
	"github.com/JacobGDG/bark/git"
)

type Deps struct {
	git *git.Client
}

func buildDeps() (*Deps, error) {
	g, err := git.New()
	if err != nil {
		return nil, err
	}

	deps := &Deps{
		git: g,
	}

	return deps, nil
}
