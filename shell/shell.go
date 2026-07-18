package shell

import (
	"os/exec"
	"strings"
)

// Shell centralises all exec.Command calls into a single place.
type Shell interface {
	// Run executes name with args and returns trimmed stdout on success.
	Run(name string, args ...string) (string, *CmdError)
}

type RealShell struct{}

func Check(binName string) error {
	_, err := exec.LookPath(binName)
	if err != nil {
		return err
	}
	return nil
}

func New() (Shell, error) {
	return &RealShell{}, nil
}

func (c *RealShell) Run(name string, args ...string) (string, *CmdError) {
	cmd := exec.Command(name, args...)
	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", &CmdError{
			Stderr:   strings.TrimSpace(stderr.String()),
			ExitCode: cmd.ProcessState.ExitCode(),
			Err:      err,
		}
	}
	return strings.TrimSpace(stdout.String()), nil
}
