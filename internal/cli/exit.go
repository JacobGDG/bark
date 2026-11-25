package cli

const (
	ExitError = 1
	// bark add with --on-exists=error found an existing worktree.
	ExitWorktreeExists = 2
)

type ExitCodeError struct {
	Code int
	Err  error
}

func (e *ExitCodeError) Error() string { return e.Err.Error() }
func (e *ExitCodeError) Unwrap() error { return e.Err }
