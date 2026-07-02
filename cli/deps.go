package cli

import (
	"github.com/JacobGDG/bark/config"
	"github.com/JacobGDG/bark/git"
)

func checkDeps(cfg *config.Config) error {
	err := git.Check()
	if err != nil && cfg.DepRequired("git") {
		return err
	}

	return nil
}
