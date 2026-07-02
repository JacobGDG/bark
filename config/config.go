package config

import (
	"log/slog"
	"os"
	"path/filepath"
)

type Config struct {
	WorktreeRootDir string
}

func Load() (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	cfg := &Config{
		WorktreeRootDir: filepath.Join(home, "worktrees"),
	}
	cfg.log()
	return cfg, nil
}

func (c *Config) log() {
	slog.Debug("config loaded", "worktree_root_dir", c.WorktreeRootDir)
}

func (c *Config) DepRequired(dep string) bool {
	switch dep {
	case "git":
		return false
	default:
		return false
	}
}
