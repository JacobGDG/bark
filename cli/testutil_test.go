package cli

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

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
