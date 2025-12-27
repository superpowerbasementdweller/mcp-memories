package db

import (
	"database/sql"
)

// FileAnnotation represents a file/directory annotation
type FileAnnotation struct {
	ID        int64  `json:"id"`
	ProjectID int64  `json:"project_id"`
	Path      string `json:"path"`
	Note      string `json:"note"`
	IsDir     bool   `json:"is_dir"`
}

// AnnotateFile adds or updates a note on a file path
func (db *DB) AnnotateFile(projectID *int64, path, note string, isDir bool) (*FileAnnotation, error) {
	pid := db.GetProjectID(projectID)

	_, err := db.Exec(
		"INSERT INTO filetree (project_id, path, note, is_dir) VALUES (?, ?, ?, ?) ON CONFLICT(project_id, path) DO UPDATE SET note = ?, is_dir = ?",
		pid, path, note, isDir, note, isDir,
	)
	if err != nil {
		return nil, err
	}

	return db.GetFileAnnotation(projectID, path)
}

// GetFileAnnotation gets an annotation for a specific path
func (db *DB) GetFileAnnotation(projectID *int64, path string) (*FileAnnotation, error) {
	pid := db.GetProjectID(projectID)

	f := &FileAnnotation{}
	err := db.QueryRow(
		"SELECT id, project_id, path, note, is_dir FROM filetree WHERE project_id = ? AND path = ?",
		pid, path,
	).Scan(&f.ID, &f.ProjectID, &f.Path, &f.Note, &f.IsDir)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return f, nil
}

// ListFileAnnotations lists all annotations for a project
func (db *DB) ListFileAnnotations(projectID *int64) ([]FileAnnotation, error) {
	pid := db.GetProjectID(projectID)

	rows, err := db.Query(
		"SELECT id, project_id, path, note, is_dir FROM filetree WHERE project_id = ? ORDER BY path",
		pid,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []FileAnnotation
	for rows.Next() {
		var f FileAnnotation
		if err := rows.Scan(&f.ID, &f.ProjectID, &f.Path, &f.Note, &f.IsDir); err != nil {
			return nil, err
		}
		items = append(items, f)
	}
	return items, rows.Err()
}

// DeleteFileAnnotation deletes an annotation
func (db *DB) DeleteFileAnnotation(projectID *int64, path string) error {
	pid := db.GetProjectID(projectID)
	_, err := db.Exec("DELETE FROM filetree WHERE project_id = ? AND path = ?", pid, path)
	return err
}
