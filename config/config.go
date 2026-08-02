package config

import (
	"bytes"
	"errors"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

type configFile struct {
	RootDir string `toml:"root_dir"`
}

type Config struct {
	WorktreeRootDir string
}

func Load() (*Config, error) {
	paths, err := configSearchPaths()
	if err != nil {
		return nil, err
	}

	cfg, err := defaultConfig()
	if err != nil {
		return nil, err
	}

	var cfgData []byte
	var cfgPath string

	for _, path := range paths {
		cfgData, err = os.ReadFile(path)
		if errors.Is(err, os.ErrNotExist) {
			continue
		}
		if err != nil {
			return nil, err
		}

		cfgPath = path
		break
	}

	if cfgData != nil {
		if len(cfgData) != 0 {
			slog.Info("Loading config file.", "configPath", cfgPath)

			var fileCfg configFile
			d := toml.NewDecoder(bytes.NewReader(cfgData))
			d.DisallowUnknownFields()
			if err := d.Decode(&fileCfg); err != nil {
				return nil, err
			}
			if fileCfg.RootDir != "" {
				cfg.WorktreeRootDir = fileCfg.RootDir
			}
		} else {
			slog.Warn("Config file found, but is empty, continuing with defaults", "configPath", cfgPath)
		}
	}

	cfg.log()
	return cfg, nil
}

func (c *Config) log() {
	slog.Debug("config loaded",
		"worktree_root_dir", c.WorktreeRootDir,
	)
}

func defaultConfig() (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	cfg := &Config{
		WorktreeRootDir: filepath.Join(home, "worktrees"),
	}
	return cfg, nil
}

func configSearchPaths() ([]string, error) {
	xdgCfgPath, err := xdgConfigPath()
	if err != nil {
		return nil, err
	}

	locations := []string{
		xdgCfgPath,
	}

	return locations, nil
}

func xdgConfigPath() (string, error) {
	if v := os.Getenv("XDG_CONFIG_HOME"); v != "" {
		return filepath.Join(v, "bark", "bark.toml"), nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, ".config", "bark", "bark.toml"), nil
}
