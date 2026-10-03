package cli

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/JacobGDG/bark/config"
	"github.com/JacobGDG/bark/git"
	"github.com/JacobGDG/bark/shell"
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

// newLocalSourceRepo creates a throwaway git repo with one commit on "main",
// so it can be used as a local clone source without touching the network.
func newLocalSourceRepo(t *testing.T) string {
	t.Helper()

	dir := filepath.Join(t.TempDir(), "source")
	runGit(t, "", "init", "--initial-branch=main", dir)
	require.NoError(t, os.WriteFile(filepath.Join(dir, "README.md"), []byte("hello"), 0644))
	runGit(t, dir, "add", "README.md")
	runGit(t, dir, "-c", "user.email=test@example.com", "-c", "user.name=bark test", "commit", "-m", "init", "--no-verify")

	return dir
}

func runGit(t *testing.T, dir string, args ...string) {
	t.Helper()

	fullArgs := args
	if dir != "" {
		fullArgs = append([]string{"-C", dir}, args...)
	}

	cmd := exec.Command("git", fullArgs...)
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, "git %v: %s", args, out)
}

func TestClone_E2E(t *testing.T) {
	sourceRepo := newLocalSourceRepo(t)
	worktreeRoot := t.TempDir()

	sh, err := shell.New()
	require.NoError(t, err)
	g, err := git.New(sh)
	require.NoError(t, err)

	deps := &Deps{
		cfg: &config.Config{WorktreeRootDir: filepath.Join(t.TempDir(), "should-not-be-used")},
		git: g,
	}

	cmd := NewCloneCommand(deps)
	cmd.SetArgs([]string{sourceRepo, "--worktree-root-dir", worktreeRoot, "--stay"})
	require.NoError(t, cmd.Execute())

	repoPath := filepath.Join(worktreeRoot, "source")
	assert.FileExists(t, filepath.Join(repoPath, ".git", "HEAD"))

	worktreePath := filepath.Join(repoPath, "main")
	assert.FileExists(t, filepath.Join(worktreePath, ".git"))
	assert.FileExists(t, filepath.Join(worktreePath, "README.md"))

	assert.NoDirExists(t, deps.cfg.WorktreeRootDir)
}
