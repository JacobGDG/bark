package config

import (
	"bytes"
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	home, err := os.UserHomeDir()
	require.NoError(t, err)
	defaultRootDir := filepath.Join(home, "worktrees")

	tests := []struct {
		name        string
		fileContent *string // nil = file does not exist
		wantRootDir string
		wantErr     bool
	}{
		{
			name:        "no file found uses defaults",
			fileContent: nil,
			wantRootDir: defaultRootDir,
		},
		{
			name:        "empty file uses defaults",
			fileContent: new(""),
			wantRootDir: defaultRootDir,
		},
		{
			name:        "root_dir absent uses default",
			fileContent: new("# bark config\n"),
			wantRootDir: defaultRootDir,
		},
		{
			name:        "root_dir set is used",
			fileContent: new(`root_dir = "/custom/path"`),
			wantRootDir: "/custom/path",
		},
		{
			name:        "invalid path accepted without error",
			fileContent: new(`root_dir = "/nonexistent/does/not/exist"`),
			wantRootDir: "/nonexistent/does/not/exist",
		},
		{
			name:        "malformed TOML returns error",
			fileContent: new("root_dir = ["),
			wantErr:     true,
		},
		{
			name:        "unknown field returns error",
			fileContent: new(`unknown_key = "value"`),
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			configDir := t.TempDir()
			t.Setenv("XDG_CONFIG_HOME", configDir)

			if tt.fileContent != nil {
				barkDir := filepath.Join(configDir, "bark")
				require.NoError(t, os.MkdirAll(barkDir, 0755))
				require.NoError(t, os.WriteFile(
					filepath.Join(barkDir, "config.toml"),
					[]byte(*tt.fileContent),
					0644,
				))
			}

			cfg, err := Load()
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantRootDir, cfg.WorktreeRootDir)
		})
	}
}

func TestLoad_EmptyFileLogsWarning(t *testing.T) {
	var buf bytes.Buffer
	original := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelWarn})))
	defer slog.SetDefault(original)

	configDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", configDir)
	barkDir := filepath.Join(configDir, "bark")
	require.NoError(t, os.MkdirAll(barkDir, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(barkDir, "config.toml"), []byte(""), 0644))

	_, err := Load()
	require.NoError(t, err)
	assert.Contains(t, buf.String(), "WARN")
}
