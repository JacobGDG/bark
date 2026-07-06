package cli

import (
	"testing"
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repoNameFromURL(tt.input)

			if err != nil && !tt.wantErr {
				t.Errorf("expected '%s' to return '%s' but got error '%s'", tt.input, tt.want, err)
				return
			}

			if err == nil && tt.wantErr {
				t.Errorf("expected '%s' to return an Err but got '%s'", tt.input, got)
				return
			}

			if got != tt.want {
				t.Errorf("expected '%s' to return '%s' but got '%s'", tt.input, tt.want, got)
				return
			}
		})
	}
}
