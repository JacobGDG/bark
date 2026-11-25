package cli

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"

	"github.com/JacobGDG/bark/internal/config"
	"github.com/JacobGDG/bark/internal/git"
	"github.com/JacobGDG/bark/internal/shell"
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

func commitFile(t *testing.T, dir, name, content string) {
	t.Helper()

	require.NoError(t, os.WriteFile(filepath.Join(dir, name), []byte(content), 0644))
	runGit(t, dir, "add", name)
	runGit(t, dir, "-c", "user.email=test@example.com", "-c", "user.name=bark test", "commit", "-m", name, "--no-verify")
}

func gitOutput(t *testing.T, dir string, args ...string) string {
	t.Helper()

	out, err := exec.Command("git", append([]string{"-C", dir}, args...)...).CombinedOutput()
	require.NoError(t, err, "git %v: %s", args, out)
	return strings.TrimSpace(string(out))
}

func clonedRepo(t *testing.T) (deps *Deps, root, source, repoPath string) {
	t.Helper()

	source = newLocalSourceRepo(t)
	root = t.TempDir()

	sh, err := shell.New()
	require.NoError(t, err)
	g, err := git.New(sh)
	require.NoError(t, err)

	deps = &Deps{
		cfg: &config.Config{WorktreeRootDir: root},
		git: g,
	}

	cloneCmd := NewCloneCommand(deps)
	cloneCmd.SetArgs([]string{source, "--stay"})
	require.NoError(t, cloneCmd.Execute())

	return deps, root, source, filepath.Join(root, "source")
}

func runCmd(t *testing.T, cmd *cobra.Command, args ...string) (string, error) {
	t.Helper()

	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(io.Discard)
	cmd.SetArgs(args)
	err := cmd.Execute()
	return out.String(), err
}
