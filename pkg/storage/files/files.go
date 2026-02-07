// Package files provides a file-based implementation of the storage.Storage interface.
// Pages are stored as individual files using gob encoding, organized by username.
package files

import (
	"encoding/gob"
	"errors"
	"fmt"
	"go_link_storage/pkg/lib/e"
	"go_link_storage/pkg/storage"
	"math/rand"
	"os"
	"path/filepath"
)

// Storage implements the storage.Storage interface using the file system.
type Storage struct {
	basePath string // Base directory path for storing files
}

const defaultPerm = 0774 // Default file permissions for created directories

// New creates a new file-based storage instance with the given base path.
func New(basePath string) Storage {
	return Storage{basePath: basePath}
}

// Save stores a page as a file in the file system.
// The file is encoded using gob and stored in a directory named after the username.
func (s Storage) Save(page *storage.Page) (err error) {
	defer func() { err = e.WrapIfErr("cannot save page", err) }()

	fPath := filepath.Join(s.basePath, page.UserName)

	if err := os.MkdirAll(fPath, defaultPerm); err != nil {
		return err
	}

	fName, err := fileName(page)
	if err != nil {
		return err
	}

	fPath = filepath.Join(fPath, fName)

	file, err := os.Create(fPath)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	if err := gob.NewEncoder(file).Encode(page); err != nil {
		return err
	}

	return nil
}

// PickRandom selects and returns a random page from the files stored for the given user.
func (s Storage) PickRandom(userName string) (page *storage.Page, err error) {
	defer func() { err = e.WrapIfErr("cannot pick page", err) }()

	path := filepath.Join(s.basePath, userName)

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, storage.ErrNoSavedPages
	}

	n := rand.Intn(len(files))

	file := files[n]

	return s.decodePage(filepath.Join(path, file.Name()))
}

// Remove deletes the file associated with the given page.
func (s Storage) Remove(p *storage.Page) error {
	fileName, err := fileName(p)
	if err != nil {
		return e.Wrap("cannot remove page", err)
	}

	path := filepath.Join(s.basePath, p.UserName, fileName)

	if err := os.Remove(path); err != nil {
		msg := fmt.Sprintf("cannot remove file %s", path)
		return e.Wrap(msg, err)
	}

	return nil
}

// Exists checks if a file exists for the given page.
func (s Storage) Exists(p *storage.Page) (bool, error) {
	fileName, err := fileName(p)
	if err != nil {
		return false, e.Wrap("cannot check if file exists", err)
	}

	path := filepath.Join(s.basePath, p.UserName, fileName)

	switch _, err = os.Stat(path); {
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	case err != nil:
		msg := fmt.Sprintf("cannot check if file %s exists", path)
		return false, e.Wrap(msg, err)
	}

	return true, nil
}

// decodePage reads and decodes a page from a file using gob.
func (s Storage) decodePage(filePath string) (*storage.Page, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, e.Wrap("cannot decode page", err)
	}
	defer func() { _ = f.Close() }()

	var p storage.Page

	if err := gob.NewDecoder(f).Decode(&p); err != nil {
		return nil, e.Wrap("cannot decode page", err)
	}

	return &p, nil
}

// fileName generates a filename for a page based on its hash.
func fileName(p *storage.Page) (string, error) {
	return p.Hash()
}
