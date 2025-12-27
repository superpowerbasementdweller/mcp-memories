package db

import (
	"database/sql"
	"fmt"
	"time"
)

// Project represents a project namespace
type Project struct {
	ID        int64     `json:"id"`
	Slug      string    `json:"slug"`
	Name      string    `json:"name,omitempty"`
	RootPath  string    `json:"root_path,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateProject creates a new project
func (db *DB) CreateProject(slug, name, rootPath string) (*Project, error) {
	result, err := db.Exec(
		"INSERT INTO projects (slug, name, root_path) VALUES (?, ?, ?)",
		slug, name, rootPath,
	)
	if err != nil {
		return nil, fmt.Errorf("creating project: %w", err)
	}

	id, _ := result.LastInsertId()
	return db.GetProjectByID(id)
}

// GetProjectByID gets a project by ID
func (db *DB) GetProjectByID(id int64) (*Project, error) {
	p := &Project{}
	var name, rootPath sql.NullString
	err := db.QueryRow(
		"SELECT id, slug, name, root_path, created_at FROM projects WHERE id = ?",
		id,
	).Scan(&p.ID, &p.Slug, &name, &rootPath, &p.CreatedAt)
	if err != nil {
		return nil, err
	}
	p.Name = name.String
	p.RootPath = rootPath.String
	return p, nil
}

// GetProjectBySlug gets a project by slug
func (db *DB) GetProjectBySlug(slug string) (*Project, error) {
	p := &Project{}
	var name, rootPath sql.NullString
	err := db.QueryRow(
		"SELECT id, slug, name, root_path, created_at FROM projects WHERE slug = ?",
		slug,
	).Scan(&p.ID, &p.Slug, &name, &rootPath, &p.CreatedAt)
	if err != nil {
		return nil, err
	}
	p.Name = name.String
	p.RootPath = rootPath.String
	return p, nil
}

// ListProjects lists all projects
func (db *DB) ListProjects() ([]Project, error) {
	rows, err := db.Query("SELECT id, slug, name, root_path, created_at FROM projects ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []Project
	for rows.Next() {
		var p Project
		var name, rootPath sql.NullString
		if err := rows.Scan(&p.ID, &p.Slug, &name, &rootPath, &p.CreatedAt); err != nil {
			return nil, err
		}
		p.Name = name.String
		p.RootPath = rootPath.String
		projects = append(projects, p)
	}
	return projects, rows.Err()
}

// GetOrCreateProject gets a project by slug, creating it if it doesn't exist
func (db *DB) GetOrCreateProject(slug string) (*Project, error) {
	p, err := db.GetProjectBySlug(slug)
	if err == nil {
		return p, nil
	}
	if err != sql.ErrNoRows {
		return nil, err
	}
	return db.CreateProject(slug, "", "")
}
