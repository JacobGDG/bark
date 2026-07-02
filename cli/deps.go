package cli

import (
	"github.com/JacobGDG/bark/git"
)

func checkDeps() error {
	err := git.Check()
	if err != nil {
		return err
	}

	return nil
}
