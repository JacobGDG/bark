package shell

import (
	"os/exec"
	"strings"
)

type Shell interface {
	Run(args ...string) (string, error)
}

// Client runs a single resolved binary.
type RealShell struct {
	binPath string
}

// simply check the bin exists in path, do nothing
func Check(binName string) error {
	_, err := exec.LookPath(binName)
	if err != nil {
		return err
	}
	return nil
}

// New resolves binaryName via PATH and returns a Client, or an error if not found.
func NewShell(binName string) (Shell, error) {
	binPath, err := exec.LookPath(binName)
	if err != nil {
		return nil, err
	}
	return &RealShell{binPath}, nil
}

// Capture executes the binary with args and returns stdout on success.
func (c *RealShell) Run(args ...string) (string, error) {
	cmd := exec.Command(c.binPath, args...)
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
