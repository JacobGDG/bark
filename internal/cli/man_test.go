package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// manDir holds the hand-written man pages, relative to this package.
const manDir = "../../man"

// TestManPagesMatchCommands keeps man/ in step with the command tree, since
// the pages are written by hand: every command has a page that documents
// each of its flags, every page has a command, and --help points at it.
func TestManPagesMatchCommands(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	root, err := NewRootCommand("vtest")
	require.NoError(t, err)

	commandPages := map[string]bool{}
	var check func(*cobra.Command)
	check = func(cmd *cobra.Command) {
		name := strings.ReplaceAll(cmd.CommandPath(), " ", "-")
		commandPages[name+".1.md"] = true

		page, err := os.ReadFile(filepath.Join(manDir, name+".1.md"))
		if assert.NoError(t, err, "missing man page for %q", cmd.CommandPath()) {
			cmd.LocalNonPersistentFlags().VisitAll(func(f *pflag.Flag) {
				assert.Contains(t, string(page), "**--"+f.Name+"**",
					"man/%s.1.md doesn't document --%s", name, f.Name)
			})
		}
		assert.Contains(t, cmd.Long, "man "+name, "%s --help doesn't point at its man page", cmd.CommandPath())

		for _, c := range cmd.Commands() {
			if c.IsAvailableCommand() {
				check(c)
			}
		}
	}
	check(root)

	pages, err := filepath.Glob(filepath.Join(manDir, "*.1.md"))
	require.NoError(t, err)
	for _, page := range pages {
		assert.True(t, commandPages[filepath.Base(page)], "%s has no matching command", page)
	}
}
