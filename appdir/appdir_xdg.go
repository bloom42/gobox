// +build !darwin,!windows

package appdir

import (
	"errors"
	"os"
	"path/filepath"
)

type dirs struct {
	name string
}

func (d *dirs) joinWithHome(xdgEnv string, paths ...string) (string, error) {
	baseDir := os.Getenv(xdgEnv)
	if baseDir == "" {
		home := os.Getenv("HOME")
		if dir == "" {
			return "", fmt.Errorf("neither $%s nor $HOME are defined", xdgEnv)
		}
		baseDir = filepath.Join(home, paths...)
	}

	// handle flatpak case https://docs.flatpak.org/en/latest/conventions.html#xdg-base-directories
	if strings.Contains(baseDir, d.name) {
		return baseDir, nil
	}
	return filepath.Join(baseDir, d.name), nil
}

func (d *dirs) UserConfig() (string, error) {
	return d.joinWithHome("XDG_CONFIG_HOME", ".config")
}

func (d *dirs) UserCache() (string, error) {
	return d.joinWithHome("XDG_CACHE_HOME", ".cache")
}

func (d *dirs) UserLogs() (string, error) {
	return d.joinWithHome("XDG_STATE_HOME", ".local", "state")
}

func (d *dirs) UserData() (string, error) {
	return d.joinWithHome("XDG_DATA_HOME", ".local", "share")
}
