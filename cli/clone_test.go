package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRepoNameFromURL(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:    "Expected https url",
			input:   "github.com/JacobGDG/bark.git",
			want:    "bark",
			wantErr: false,
		},
		{
			name:    "Expected ssh url",
			input:   "git@github.com:JacobGDG/bark.git",
			want:    "bark",
			wantErr: false,
		},
		{
			name:    "Expected wiki url",
			input:   "git@github.com:JacobGDG/bark.wiki.git",
			want:    "bark.wiki",
			wantErr: false,
		},
		{
			name:    "Unexpected empty path url",
			input:   "git@github.com:JacobGDG/",
			want:    "",
			wantErr: true,
		},
		{
			name:    "Unexpected string",
			input:   "asdf",
			want:    "",
			wantErr: true,
		},
		{
			name:    "Unexpected empty string",
			input:   "",
			want:    "",
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := repoNameFromURL(test.input)

			if !test.wantErr {
				require.NoError(t, err)

				if assert.NotNil(t, got) {
					assert.Equal(t, test.want, got)
				}
			} else {
				assert.Error(t, err)
			}
		})
	}
}
