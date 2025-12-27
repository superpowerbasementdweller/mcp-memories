package db

import (
	"database/sql"
)

// Metadata represents a key-value pair
type Metadata struct {
	ID        int64  `json:"id"`
	ProjectID int64  `json:"project_id"`
	Key       string `json:"key"`
	Value     string `json:"value"`
}

// SetMetadata sets a metadata key-value pair
func (db *DB) SetMetadata(projectID *int64, key, value string) (*Metadata, error) {
	pid := db.GetProjectID(projectID)

	_, err := db.Exec(
		"INSERT INTO metadata (project_id, key, value) VALUES (?, ?, ?) ON CONFLICT(project_id, key) DO UPDATE SET value = ?",
		pid, key, value, value,
	)
	if err != nil {
		return nil, err
	}

	return db.GetMetadata(projectID, key)
}

// GetMetadata gets a metadata value by key
func (db *DB) GetMetadata(projectID *int64, key string) (*Metadata, error) {
	pid := db.GetProjectID(projectID)

	m := &Metadata{}
	err := db.QueryRow(
		"SELECT id, project_id, key, value FROM metadata WHERE project_id = ? AND key = ?",
		pid, key,
	).Scan(&m.ID, &m.ProjectID, &m.Key, &m.Value)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return m, nil
}

// ListMetadata lists all metadata for a project
func (db *DB) ListMetadata(projectID *int64) ([]Metadata, error) {
	pid := db.GetProjectID(projectID)

	rows, err := db.Query(
		"SELECT id, project_id, key, value FROM metadata WHERE project_id = ? ORDER BY key",
		pid,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []Metadata
	for rows.Next() {
		var m Metadata
		if err := rows.Scan(&m.ID, &m.ProjectID, &m.Key, &m.Value); err != nil {
			return nil, err
		}
		items = append(items, m)
	}
	return items, rows.Err()
}

// DeleteMetadata deletes a metadata key
func (db *DB) DeleteMetadata(projectID *int64, key string) error {
	pid := db.GetProjectID(projectID)
	_, err := db.Exec("DELETE FROM metadata WHERE project_id = ? AND key = ?", pid, key)
	return err
}
