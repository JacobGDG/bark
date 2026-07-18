package zoxide

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockShell struct {
	gotName string
	gotArgs []string
	err     error
}

func (m *mockShell) Run(name string, args ...string) (string, error) {
	m.gotName = name
	m.gotArgs = args
	return "", m.err
}

func TestAdd(t *testing.T) {
	tests := []struct {
		name     string
		dir      string
		shellErr error
		wantName string
		wantArgs []string
		wantErr  bool
	}{
		{
			name:     "success",
			dir:      "/some/dir",
			shellErr: nil,
			wantName: "zoxide",
			wantArgs: []string{"add", "/some/dir"},
			wantErr:  false,
		},
		{
			name:     "shell error is returned",
			dir:      "/some/dir",
			shellErr: errors.New("zoxide failed"),
			wantName: "zoxide",
			wantArgs: []string{"add", "/some/dir"},
			wantErr:  true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			sh := &mockShell{err: test.shellErr}
			z := &RealZoxide{shell: sh}

			err := z.Add(test.dir)

			assert.Equal(t, test.wantName, sh.gotName)
			assert.Equal(t, test.wantArgs, sh.gotArgs)

			if test.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
