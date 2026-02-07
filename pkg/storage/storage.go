// Package storage defines the storage interface for persisting and retrieving pages.
// It provides a common interface that can be implemented by different storage backends.
package storage

import (
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	"go_link_storage/pkg/lib/e"
	"io"
)

// Storage defines the interface for page storage operations.
// Implementations should provide persistent storage for pages associated with users.
type Storage interface {
	// Save stores a page in the storage.
	Save(ctx context.Context, p *Page) error
	// PickRandom retrieves a random page for the given user.
	PickRandom(ctx context.Context, userName string) (*Page, error)
	// Remove deletes a page from the storage.
	Remove(ctx context.Context, p *Page) error
	// Exists checks if a page already exists in the storage.
	Exists(ctx context.Context, p *Page) (bool, error)
}

// ErrNoSavedPages is returned when attempting to pick a random page
// but no pages are saved for the user.
var ErrNoSavedPages = errors.New("no saved pages")

// Page represents a saved web page with its URL and associated username.
type Page struct {
	URL      string // The URL of the page
	UserName string // The username of the user who saved the page
}

// Hash calculates a SHA1 hash of the page based on its URL and username.
// This hash is used as a unique identifier for the page in storage.
func (p Page) Hash() (string, error) {
	h := sha1.New()

	if _, err := io.WriteString(h, p.URL); err != nil {
		return "", e.Wrap("cannot calculate hash", err)
	}

	if _, err := io.WriteString(h, p.UserName); err != nil {
		return "", e.Wrap("cannot calculate hash", err)
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
