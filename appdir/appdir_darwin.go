package appdir

import (
	"errors"
	"os"
	"path/filepath"
)

type dirs struct {
	name string
}

func (d *dirs) UserConfig() (string, error) {
	dir := os.Getenv("HOME")
	if dir == "" {
		return "", errors.New("$HOME is not defined")
	}
	return filepath.Join(dir, "Library", "Application Support", d.name), nil
}

func (d *dirs) UserCache() (string, error) {
	dir := os.Getenv("HOME")
	if dir == "" {
		return "", errors.New("$HOME is not defined")
	}
	return filepath.Join(dir, "Library", "Caches", d.name), nil
}

func (d *dirs) UserLogs() (string, error) {
	dir := os.Getenv("HOME")
	if dir == "" {
		return "", errors.New("$HOME is not defined")
	}
	return filepath.Join(dir, "Library", "Logs", d.name), nil
}

func (d *dirs) UserData() (string, error) {
	dir := os.Getenv("HOME")
	if dir == "" {
		return "", errors.New("$HOME is not defined")
	}
	return filepath.Join(dir, "Library", "Application Support", d.name), nil
}
