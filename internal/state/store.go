package state

import (
	"database/sql"
	"errors"
	"time"

	_ "modernc.org/sqlite"
)

// Store manages persistent state for discovered LinkedIn profiles.
// This state layer is intentionally minimal and read-only oriented.
type Store struct {
	db *sql.DB
}

// NewStore opens (or creates) the SQLite database and ensures schema.
func NewStore(path string) (*Store, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	if err := migrate(db); err != nil {
		return nil, err
	}

	return &Store{db: db}, nil
}

// migrate ensures the required schema exists.
// The schema tracks only discovery and visit state.
func migrate(db *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS profiles (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		profile_url TEXT UNIQUE NOT NULL,
		visited BOOLEAN NOT NULL DEFAULT 0,
		created_at DATETIME NOT NULL,
		last_visited_at DATETIME
	);
	`
	_, err := db.Exec(schema)
	return err
}

// EnsureProfile inserts a profile URL if it does not already exist.
func (s *Store) EnsureProfile(profileURL string) error {
	_, err := s.db.Exec(
		`INSERT OR IGNORE INTO profiles (profile_url, created_at)
		 VALUES (?, ?)`,
		profileURL,
		time.Now(),
	)
	return err
}

// MarkVisited marks a profile as visited.
func (s *Store) MarkVisited(profileURL string) error {
	_, err := s.db.Exec(
		`UPDATE profiles
		 SET visited = 1, last_visited_at = ?
		 WHERE profile_url = ?`,
		time.Now(),
		profileURL,
	)
	return err
}

// GetPendingProfiles returns profiles that have not yet been visited.
func (s *Store) GetPendingProfiles(limit int) ([]ProfileState, error) {
	rows, err := s.db.Query(
		`SELECT profile_url, visited, created_at, last_visited_at
		 FROM profiles
		 WHERE visited = 0
		 ORDER BY created_at
		 LIMIT ?`,
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []ProfileState
	for rows.Next() {
		var p ProfileState
		if err := rows.Scan(
			&p.ProfileURL,
			&p.Visited,
			&p.CreatedAt,
			&p.LastVisitedAt,
		); err != nil {
			return nil, err
		}
		results = append(results, p)
	}

	return results, nil
}

// Close closes the underlying database connection.
func (s *Store) Close() error {
	if s.db == nil {
		return errors.New("store not initialized")
	}
	return s.db.Close()
}
