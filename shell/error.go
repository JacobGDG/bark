package shell

type CmdError struct {
	Stderr   string
	ExitCode int
	Err      error
}

func (e *CmdError) Error() string {
	return e.Stderr
}
