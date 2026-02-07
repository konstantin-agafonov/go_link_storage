// Package sqlite provides a SQLite implementation of the storage.Storage interface.
package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"go_link_storage/pkg/storage"

	_ "modernc.org/sqlite"
)

// Storage implements the storage.Storage interface using SQLite database.
type Storage struct {
	db *sql.DB // SQLite database connection
}

// New creates a new SQLite storage instance.
// It opens a connection to the database at the given path and verifies connectivity.
func New(path string) (*Storage, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("can't open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("cannot connect to db: %w", err)
	}

	return &Storage{db: db}, nil
}

// Save stores a page in the SQLite database.
func (s *Storage) Save(ctx context.Context, p *storage.Page) error {
	q := `INSERT INTO pages (url, user_name) VALUES (?, ?);`

	if _, err := s.db.ExecContext(ctx, q, p.URL, p.UserName); err != nil {
		return fmt.Errorf("cannot save page: %w", err)
	}

	return nil
}

// PickRandom retrieves a random page for the given user from the database.
func (s *Storage) PickRandom(ctx context.Context, userName string) (*storage.Page, error) {
	q := `SELECT url FROM pages WHERE user_name = ? ORDER BY RANDOM() LIMIT 1;`

	var url string

	err := s.db.QueryRowContext(ctx, q, userName).Scan(&url)
	if err == sql.ErrNoRows {
		return nil, storage.ErrNoSavedPages
	}

	if err != nil {
		return nil, fmt.Errorf("cannot select url: %w", err)
	}

	return &storage.Page{
		URL:      url,
		UserName: userName,
	}, nil
}

// Remove deletes a page from the SQLite database.
func (s *Storage) Remove(ctx context.Context, p *storage.Page) error {
	q := `DELETE FROM pages WHERE url = ? AND user_name = ?;`

	if _, err := s.db.ExecContext(ctx, q, p.URL, p.UserName); err != nil {
		return fmt.Errorf("cannot remove page: %w", err)
	}

	return nil
}

// Exists checks if a page already exists in the SQLite database.
func (s *Storage) Exists(ctx context.Context, p *storage.Page) (bool, error) {
	q := `SELECT COUNT() FROM pages WHERE url = ? AND user_name = ?;`

	var count int

	if err := s.db.QueryRowContext(ctx, q, p.URL, p.UserName).Scan(&count); err != nil {
		return false, fmt.Errorf("cannot select url: %w", err)
	}

	return count > 0, nil
}

func (s *Storage) Init(ctx context.Context) error {
	q := `CREATE TABLE IF NOT EXISTS pages (url TEXT, user_name TEXT);`

	_, err := s.db.ExecContext(ctx, q)
	if err != nil {
		return fmt.Errorf("cannot create table: %w", err)
	}

	return nil
}
