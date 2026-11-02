package shell

import (
	"fmt"
)

// CmdError wraps a command failure with its stderr output and exit code.
type CmdError struct {
	Stderr   string
	ExitCode int
	Err      error
}

// Error presents CmdError data in a standard format keeping context
func (e *CmdError) Error() string {
	return fmt.Sprintf("command failed (%d): %s", e.ExitCode, e.Stderr)

}
