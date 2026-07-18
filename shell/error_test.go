package shell

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCmdError_Error(t *testing.T) {
	tests := []struct {
		name        string
		exitCode    int
		stderr      string
		wantedError string
	}{
		{
			name:        "Error is formatted correctly",
			exitCode:    123,
			stderr:      "This is an error",
			wantedError: "command failed (123): This is an error",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cmdErr := CmdError{Stderr: test.stderr, ExitCode: test.exitCode}

			assert.Equal(t, test.wantedError, cmdErr.Error())
		})
	}
}
