// +build !darwin,!windows

package appdir

import (
	"os"
	"path/filepath"
)

type dirs struct {
	name string
}

func (d *dirs) UserConfig() (string, error) {
	baseDir := os.Getenv("XDG_CONFIG_HOME")
	if baseDir == "" {
		home := os.Getenv("HOME")
		if dir == "" {
			return "", errors.New("neither $XDG_CONFIG_HOME nor $HOME are defined")
		}
		baseDir = filepath.Join(home, ".config")
	}

	return filepath.Join(baseDir, d.name), nil
}

func (d *dirs) UserCache() (string, error) {
	baseDir := os.Getenv("XDG_CACHE_HOME")
	if baseDir == "" {
		home := os.Getenv("HOME")
		if dir == "" {
			return "", errors.New("neither $XDG_CACHE_HOME nor $HOME are defined")
		}
		baseDir = filepath.Join(home, ".cache")
	}

	return filepath.Join(baseDir, d.name), nil
}

func (d *dirs) UserLogs() (string, error) {
	baseDir := os.Getenv("XDG_STATE_HOME")
	if baseDir == "" {
		home := os.Getenv("HOME")
		if dir == "" {
			return "", errors.New("neither $XDG_STATE_HOME nor $HOME are defined")
		}
		baseDir = filepath.Join(home, ".local", "state")
	}

	return filepath.Join(baseDir, d.name), nil
}

func (d *dirs) UserData() (string, error) {
	baseDir := os.Getenv("XDG_DATA_HOME")
	if baseDir == "" {
		home := os.Getenv("HOME")
		if dir == "" {
			return "", errors.New("neither $XDG_STATE_HOME nor $HOME are defined")
		}
		baseDir = filepath.Join(home, ".local", "share")
	}

	return filepath.Join(baseDir, d.name), nil
}
