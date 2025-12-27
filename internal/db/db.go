package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"

	"github.com/rocket/mcp-memories/internal/schema"
)

// DB wraps the SQLite database connection
type DB struct {
	*sql.DB
	defaultProjectID int64
}

// Open opens the SQLite database and runs migrations
func Open(dbPath string) (*DB, error) {
	// Ensure directory exists
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("creating db directory: %w", err)
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("opening database: %w", err)
	}

	// Enable WAL mode for better performance
	if _, err := db.Exec("PRAGMA journal_mode=WAL"); err != nil {
		db.Close()
		return nil, fmt.Errorf("enabling WAL mode: %w", err)
	}

	// Enable foreign keys
	if _, err := db.Exec("PRAGMA foreign_keys=ON"); err != nil {
		db.Close()
		return nil, fmt.Errorf("enabling foreign keys: %w", err)
	}

	// Run schema migrations
	if _, err := db.Exec(schema.Schema); err != nil {
		db.Close()
		return nil, fmt.Errorf("running migrations: %w", err)
	}

	return &DB{DB: db, defaultProjectID: 1}, nil
}

// SetDefaultProject sets the default project for operations
func (db *DB) SetDefaultProject(projectID int64) {
	db.defaultProjectID = projectID
}

// DefaultProjectID returns the current default project ID
func (db *DB) DefaultProjectID() int64 {
	return db.defaultProjectID
}

// GetProjectID returns the project ID to use, defaulting to the default project
func (db *DB) GetProjectID(projectID *int64) int64 {
	if projectID != nil && *projectID > 0 {
		return *projectID
	}
	return db.defaultProjectID
}
