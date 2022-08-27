package content

import (
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"

	"golang.org/x/net/webdav"
)

//go:embed templates/*
var templateFS embed.FS

// ErrNotADirectory is returned when something exists at a user directory location that is not a directory.
var ErrNotADirectory = errors.New("user directory is not a directory")

// ErrAlreadyExists is returned when a user directory can't be created because it already exists.
var ErrAlreadyExists = errors.New("user directory already exists")

// userDirFormat defines the format for naming user directories.
const userDirFormat = "%s/~%s"

// Manager provides functions to manage user-owned content.
type Manager struct {
	Dir string
}

// UserDir returns the path to a user directory.
func (m *Manager) UserDir(username string) string {
	return fmt.Sprintf(userDirFormat, m.Dir, username)
}

// Exists returns true if a user directory exists for the given user.
func (m *Manager) Exists(username string) (bool, error) {
	dir := m.UserDir(username)
	stat, err := os.Stat(dir)
	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if !stat.IsDir() {
		return false, ErrNotADirectory
	}
	return true, nil
}

// Create creates a user directory, failing if it already exists.
func (m *Manager) Create(username string) error {
	exists, err := m.Exists(username)
	if err != nil {
		return err
	}
	if exists {
		return ErrAlreadyExists
	}
	err = os.Mkdir(m.UserDir(username), 0700)
	if err != nil {
		return err
	}
	data, err := templateFS.ReadFile(path.Join("templates", "index.gmi"))
	if err != nil {
		return err
	}
	err = os.WriteFile(path.Join(m.UserDir(username), "index.gmi"), data, 0600)
	if err != nil {
		return err
	}
	return nil
}

// Rename renames a user directory.
func (m *Manager) Rename(oldUsername string, newUsername string) error {
	return os.Rename(m.UserDir(oldUsername), m.UserDir(newUsername))
}

// Delete deletes a user directory.
func (m *Manager) Delete(username string) error {
	return os.RemoveAll(m.UserDir(username))
}

// WebDAVDirFor returns a webdav.Dir for a user directory.
func (m *Manager) WebDAVDirFor(username string) webdav.Dir {
	return webdav.Dir(m.UserDir(username))
}
