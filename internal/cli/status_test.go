package cli

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func statusJSON(t *testing.T, deps *Deps, args ...string) branchStatus {
	t.Helper()

	out, err := runCmd(t, NewStatusCommand(deps), append(args, "--json")...)
	require.NoError(t, err)

	var status branchStatus
	require.NoError(t, json.Unmarshal([]byte(out), &status))
	return status
}

func TestStatus(t *testing.T) {
	deps, _, source, repoPath := clonedRepo(t)

	t.Run("missing branch", func(t *testing.T) {
		s := statusJSON(t, deps, "source", "nope")
		assert.False(t, s.Exists)
		assert.Nil(t, s.Worktree)
		assert.Nil(t, s.LastCommit)
		assert.Nil(t, s.Merged)
	})

	t.Run("worktree with changes, ahead of origin, unmerged", func(t *testing.T) {
		_, err := runCmd(t, NewAddCommand(deps), "source", "main", "--stay", "--on-exists=open")
		require.NoError(t, err)
		worktree := filepath.Join(repoPath, "main")
		commitFile(t, worktree, "local.txt", "local")
		require.NoError(t, os.WriteFile(filepath.Join(worktree, "untracked.txt"), nil, 0644))

		s := statusJSON(t, deps, "source", "main")
		assert.True(t, s.Exists)
		assert.True(t, s.Local)
		assert.True(t, s.Remote)
		require.NotNil(t, s.Worktree)
		assert.Equal(t, worktree, *s.Worktree)
		require.NotNil(t, s.ChangedFiles)
		assert.Equal(t, 1, *s.ChangedFiles)
		require.NotNil(t, s.Ahead)
		assert.Equal(t, 1, *s.Ahead)
		assert.Equal(t, 0, *s.Behind)
		assert.NotNil(t, s.LastCommit)
		require.NotNil(t, s.Merged)
		assert.False(t, *s.Merged, "local commit isn't on origin/main")
	})

	t.Run("remote-only branch needs --fetch", func(t *testing.T) {
		runGit(t, source, "branch", "remote-only", "main")

		assert.False(t, statusJSON(t, deps, "source", "remote-only").Exists)

		s := statusJSON(t, deps, "source", "remote-only", "--fetch")
		assert.True(t, s.Exists)
		assert.False(t, s.Local)
		assert.True(t, s.Remote)
		assert.Nil(t, s.Ahead, "ahead/behind needs both copies")
		require.NotNil(t, s.Merged)
		assert.True(t, *s.Merged)
	})

	t.Run("plain output", func(t *testing.T) {
		out, err := runCmd(t, NewStatusCommand(deps), "source", "nope")
		require.NoError(t, err)
		assert.Contains(t, out, "exists:         no\n")
		assert.Contains(t, out, "worktree:       -\n")
	})
}
