package shell

import (
	"fmt"
)

type CmdError struct {
	Stderr   string
	ExitCode int
	Err      error
}

func (e *CmdError) Error() string {
	return fmt.Sprintf("command failed (%d): %s", e.ExitCode, e.Stderr)

}
