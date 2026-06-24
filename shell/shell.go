package shell

import (
	"fmt"
	"os/exec"
	"strings"
)

// Client runs a single resolved binary.
type Client struct {
	binPath string
}

// New resolves binaryName via PATH and returns a Client, or an error if not found.
func New(binName string) (*Client, error) {
	binaryPath, err := exec.LookPath(binName)
	if err != nil {
		return nil, fmt.Errorf("%s not found in PATH: %w", binName, err)
	}
	return &Client{binPath: binaryPath}, nil
}

// Capture executes the binary with args and returns stdout on success.
func (c *Client) Run(args ...string) (string, error) {
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
