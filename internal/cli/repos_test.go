package cli

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRepos(t *testing.T) {
	deps, root, _, repoPath := clonedRepo(t)
	require.NoError(t, os.MkdirAll(filepath.Join(root, "not-a-repo"), 0755))
	require.NoError(t, os.WriteFile(filepath.Join(root, "file"), nil, 0644))

	t.Run("plain lists names", func(t *testing.T) {
		out, err := runCmd(t, NewReposCommand(deps))
		require.NoError(t, err)
		assert.Equal(t, "source\n", out)
	})

	t.Run("json includes worktrees", func(t *testing.T) {
		out, err := runCmd(t, NewReposCommand(deps), "--json")
		require.NoError(t, err)

		var repos []repoJSON
		require.NoError(t, json.Unmarshal([]byte(out), &repos))
		require.Len(t, repos, 1)
		assert.Equal(t, "source", repos[0].Name)
		assert.Equal(t, repoPath, repos[0].Path)
		require.Len(t, repos[0].Worktrees, 1)
		assert.Equal(t, "main", repos[0].Worktrees[0].Branch)
	})

	t.Run("missing root lists nothing", func(t *testing.T) {
		missing := filepath.Join(t.TempDir(), "nope")

		out, err := runCmd(t, NewReposCommand(deps), "--worktree-root-dir", missing)
		require.NoError(t, err)
		assert.Empty(t, out)

		out, err = runCmd(t, NewReposCommand(deps), "--worktree-root-dir", missing, "--json")
		require.NoError(t, err)
		assert.Equal(t, "[]\n", out)
	})
}
