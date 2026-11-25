package cli

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/JacobGDG/bark/internal/config"
	"github.com/JacobGDG/bark/internal/git"
	"github.com/JacobGDG/bark/internal/shell"
)

func newRepoDir(t *testing.T, root, name string) string {
	t.Helper()

	repoPath := filepath.Join(root, name)
	require.NoError(t, os.MkdirAll(filepath.Join(repoPath, ".git"), 0755))
	return repoPath
}

func TestRepoNameFromCwd(t *testing.T) {
	root := t.TempDir()
	newRepoDir(t, root, "myrepo")

	tests := []struct {
		name    string
		cwd     string
		want    string
		wantErr bool
	}{
		{
			name: "cwd at repo root",
			cwd:  filepath.Join(root, "myrepo"),
			want: "myrepo",
		},
		{
			name: "cwd inside a worktree",
			cwd:  filepath.Join(root, "myrepo", "main", "pkg"),
			want: "myrepo",
		},
		{
			name:    "cwd is root itself",
			cwd:     root,
			wantErr: true,
		},
		{
			name:    "cwd outside root",
			cwd:     filepath.Dir(root),
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := repoNameFromCwd(root, test.cwd)

			if test.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, test.want, got)
		})
	}
}

func TestValidateRepoDir(t *testing.T) {
	root := t.TempDir()
	newRepoDir(t, root, "myrepo")
	require.NoError(t, os.MkdirAll(filepath.Join(root, "not-a-repo"), 0755))

	tests := []struct {
		name    string
		repo    string
		wantErr bool
	}{
		{
			name: "valid bark-managed repo",
			repo: "myrepo",
		},
		{
			name:    "directory without .git",
			repo:    "not-a-repo",
			wantErr: true,
		},
		{
			name:    "missing directory",
			repo:    "nope",
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := validateRepoDir(root, test.repo)

			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestResolveRepo(t *testing.T) {
	root := t.TempDir()
	newRepoDir(t, root, "myrepo")

	t.Run("explicit repo name", func(t *testing.T) {
		got, err := resolveRepo(root, "myrepo")
		require.NoError(t, err)
		assert.Equal(t, "myrepo", got)
	})

	t.Run("explicit repo name that does not exist", func(t *testing.T) {
		_, err := resolveRepo(root, "nope")
		assert.Error(t, err)
	})

	t.Run("dot detects repo from cwd", func(t *testing.T) {
		cwd, err := os.Getwd()
		require.NoError(t, err)
		t.Chdir(filepath.Join(root, "myrepo"))
		t.Cleanup(func() { _ = os.Chdir(cwd) })

		got, err := resolveRepo(root, ".")
		require.NoError(t, err)
		assert.Equal(t, "myrepo", got)
	})
}

func TestAdd_E2E(t *testing.T) {
	sourceRepo := newLocalSourceRepo(t)

	runGit(t, sourceRepo, "checkout", "-b", "feature-remote")
	require.NoError(t, os.WriteFile(filepath.Join(sourceRepo, "feature.txt"), []byte("feature"), 0644))
	runGit(t, sourceRepo, "add", "feature.txt")
	runGit(t, sourceRepo, "-c", "user.email=test@example.com", "-c", "user.name=bark test", "commit", "-m", "feature", "--no-verify")
	runGit(t, sourceRepo, "checkout", "main")

	worktreeRoot := t.TempDir()

	sh, err := shell.New()
	require.NoError(t, err)
	g, err := git.New(sh)
	require.NoError(t, err)

	deps := &Deps{
		cfg: &config.Config{WorktreeRootDir: filepath.Join(t.TempDir(), "should-not-be-used")},
		git: g,
	}

	cloneCmd := NewCloneCommand(deps)
	cloneCmd.SetArgs([]string{sourceRepo, "--worktree-root-dir", worktreeRoot, "--stay"})
	require.NoError(t, cloneCmd.Execute())

	repoPath := filepath.Join(worktreeRoot, "source")

	// simulate a local branch that already exists but has no worktree yet
	runGit(t, repoPath, "branch", "local-branch", "main")

	t.Run("attaches existing local branch", func(t *testing.T) {
		addCmd := NewAddCommand(deps)
		addCmd.SetArgs([]string{"source", "local-branch", "--worktree-root-dir", worktreeRoot, "--stay"})
		require.NoError(t, addCmd.Execute())

		worktreePath := filepath.Join(repoPath, "local-branch")
		assert.FileExists(t, filepath.Join(worktreePath, ".git"))
		assert.FileExists(t, filepath.Join(worktreePath, "README.md"))
	})

	t.Run("creates local branch tracking remote", func(t *testing.T) {
		addCmd := NewAddCommand(deps)
		addCmd.SetArgs([]string{"source", "feature-remote", "--worktree-root-dir", worktreeRoot, "--stay"})
		require.NoError(t, addCmd.Execute())

		worktreePath := filepath.Join(repoPath, "feature-remote")
		assert.FileExists(t, filepath.Join(worktreePath, "feature.txt"))
	})

	t.Run("creates brand new branch", func(t *testing.T) {
		addCmd := NewAddCommand(deps)
		addCmd.SetArgs([]string{"source", "feature-new", "--worktree-root-dir", worktreeRoot, "--stay"})
		require.NoError(t, addCmd.Execute())

		worktreePath := filepath.Join(repoPath, "feature-new")
		assert.FileExists(t, filepath.Join(worktreePath, "README.md"))
	})

	t.Run("errors when worktree already exists", func(t *testing.T) {
		addCmd := NewAddCommand(deps)
		addCmd.SetArgs([]string{"source", "feature-new", "--worktree-root-dir", worktreeRoot, "--stay"})
		assert.Error(t, addCmd.Execute())
	})

	t.Run("errors on unknown repo", func(t *testing.T) {
		addCmd := NewAddCommand(deps)
		addCmd.SetArgs([]string{"not-a-repo", "something", "--worktree-root-dir", worktreeRoot, "--stay"})
		assert.Error(t, addCmd.Execute())
	})
}

func TestAdd_Scripting(t *testing.T) {
	t.Run("fetches branches pushed after clone", func(t *testing.T) {
		deps, _, source, repoPath := clonedRepo(t)
		runGit(t, source, "checkout", "-b", "pushed-later")
		commitFile(t, source, "later.txt", "later")

		_, err := runCmd(t, NewAddCommand(deps), "source", "pushed-later", "--stay")
		require.NoError(t, err)
		assert.FileExists(t, filepath.Join(repoPath, "pushed-later", "later.txt"))
	})

	t.Run("new branch starts from fresh origin default without tracking it", func(t *testing.T) {
		deps, _, source, repoPath := clonedRepo(t)
		commitFile(t, source, "upstream.txt", "new on main")

		_, err := runCmd(t, NewAddCommand(deps), "source", "fresh", "--stay")
		require.NoError(t, err)

		worktree := filepath.Join(repoPath, "fresh")
		assert.FileExists(t, filepath.Join(worktree, "upstream.txt"))
		_, upstreamErr := exec.Command("git", "-C", worktree, "rev-parse", "--abbrev-ref", "@{upstream}").Output()
		assert.Error(t, upstreamErr, "new branch should have no upstream")
	})

	t.Run("--from sets the base of a new branch", func(t *testing.T) {
		deps, _, source, repoPath := clonedRepo(t)
		runGit(t, source, "checkout", "-b", "base")
		commitFile(t, source, "base.txt", "base")

		_, err := runCmd(t, NewAddCommand(deps), "source", "on-base", "--from", "origin/base", "--stay")
		require.NoError(t, err)
		assert.FileExists(t, filepath.Join(repoPath, "on-base", "base.txt"))
	})

	t.Run("--print-path prints only the worktree path", func(t *testing.T) {
		deps, _, _, repoPath := clonedRepo(t)

		out, err := runCmd(t, NewAddCommand(deps), "source", "printed", "--stay", "--print-path")
		require.NoError(t, err)
		assert.Equal(t, filepath.Join(repoPath, "printed")+"\n", out)
	})

	t.Run("--on-exists=error exits with ExitWorktreeExists", func(t *testing.T) {
		deps, _, _, _ := clonedRepo(t)

		_, err := runCmd(t, NewAddCommand(deps), "source", "main", "--stay")
		var exitErr *ExitCodeError
		require.ErrorAs(t, err, &exitErr)
		assert.Equal(t, ExitWorktreeExists, exitErr.Code)
	})

	t.Run("--on-exists=open returns the existing worktree", func(t *testing.T) {
		deps, _, _, repoPath := clonedRepo(t)

		out, err := runCmd(t, NewAddCommand(deps), "source", "main", "--stay", "--on-exists=open", "--print-path")
		require.NoError(t, err)
		assert.Equal(t, filepath.Join(repoPath, "main")+"\n", out)
	})

	t.Run("--on-exists=increment skips names in use anywhere", func(t *testing.T) {
		deps, _, source, repoPath := clonedRepo(t)
		// JIR-1 has a worktree, JIR-1-1 is only local, JIR-1-2 is only on origin.
		_, err := runCmd(t, NewAddCommand(deps), "source", "JIR-1", "--stay")
		require.NoError(t, err)
		runGit(t, repoPath, "branch", "JIR-1-1", "main")
		runGit(t, source, "branch", "JIR-1-2", "main")

		out, err := runCmd(t, NewAddCommand(deps), "source", "JIR-1", "--stay", "--on-exists=increment", "--print-path")
		require.NoError(t, err)
		assert.Equal(t, filepath.Join(repoPath, "JIR-1-3")+"\n", out)
	})

	t.Run("--on-exists=increment keeps a free name as is", func(t *testing.T) {
		deps, _, _, repoPath := clonedRepo(t)

		out, err := runCmd(t, NewAddCommand(deps), "source", "unused", "--stay", "--on-exists=increment", "--print-path")
		require.NoError(t, err)
		assert.Equal(t, filepath.Join(repoPath, "unused")+"\n", out)
	})

	t.Run("rejects unknown --on-exists", func(t *testing.T) {
		deps, _, _, _ := clonedRepo(t)

		_, err := runCmd(t, NewAddCommand(deps), "source", "x", "--stay", "--on-exists=maybe")
		assert.ErrorContains(t, err, "invalid --on-exists")
	})
}
