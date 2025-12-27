package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Guideline represents a how-to or pattern documentation
type Guideline struct {
	ID        int64     `json:"id"`
	ProjectID int64     `json:"project_id"`
	Category  string    `json:"category"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Tags      []string  `json:"tags"`
	Priority  int       `json:"priority"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateGuideline creates a new guideline
func (db *DB) CreateGuideline(projectID *int64, category, title, content string, tags []string, priority int) (*Guideline, error) {
	pid := db.GetProjectID(projectID)

	tagsJSON, err := json.Marshal(tags)
	if err != nil {
		return nil, fmt.Errorf("marshaling tags: %w", err)
	}

	result, err := db.Exec(
		"INSERT INTO guidelines (project_id, category, title, content, tags, priority) VALUES (?, ?, ?, ?, ?, ?)",
		pid, category, title, content, string(tagsJSON), priority,
	)
	if err != nil {
		return nil, fmt.Errorf("creating guideline: %w", err)
	}

	id, _ := result.LastInsertId()
	return db.GetGuideline(id)
}

// GetGuideline gets a guideline by ID
func (db *DB) GetGuideline(id int64) (*Guideline, error) {
	g := &Guideline{}
	var tagsJSON sql.NullString
	err := db.QueryRow(
		"SELECT id, project_id, category, title, content, tags, priority, created_at, updated_at FROM guidelines WHERE id = ?",
		id,
	).Scan(&g.ID, &g.ProjectID, &g.Category, &g.Title, &g.Content, &tagsJSON, &g.Priority, &g.CreatedAt, &g.UpdatedAt)
	if err != nil {
		return nil, err
	}
	if tagsJSON.Valid {
		json.Unmarshal([]byte(tagsJSON.String), &g.Tags)
	}
	return g, nil
}

// UpdateGuideline updates a guideline
func (db *DB) UpdateGuideline(id int64, content *string, tags *[]string, priority *int) (*Guideline, error) {
	var sets []string
	var args []interface{}

	if content != nil {
		sets = append(sets, "content = ?")
		args = append(args, *content)
	}
	if tags != nil {
		tagsJSON, _ := json.Marshal(*tags)
		sets = append(sets, "tags = ?")
		args = append(args, string(tagsJSON))
	}
	if priority != nil {
		sets = append(sets, "priority = ?")
		args = append(args, *priority)
	}

	if len(sets) == 0 {
		return db.GetGuideline(id)
	}

	sets = append(sets, "updated_at = CURRENT_TIMESTAMP")
	args = append(args, id)

	_, err := db.Exec(
		fmt.Sprintf("UPDATE guidelines SET %s WHERE id = ?", strings.Join(sets, ", ")),
		args...,
	)
	if err != nil {
		return nil, fmt.Errorf("updating guideline: %w", err)
	}

	return db.GetGuideline(id)
}

// ListGuidelines lists guidelines with optional category filter
func (db *DB) ListGuidelines(projectID *int64, category *string) ([]Guideline, error) {
	pid := db.GetProjectID(projectID)

	var conditions []string
	var args []interface{}

	conditions = append(conditions, "project_id = ?")
	args = append(args, pid)

	if category != nil {
		conditions = append(conditions, "category = ?")
		args = append(args, *category)
	}

	query := fmt.Sprintf(
		"SELECT id, project_id, category, title, content, tags, priority, created_at, updated_at FROM guidelines WHERE %s ORDER BY priority DESC, category, title",
		strings.Join(conditions, " AND "),
	)

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var guidelines []Guideline
	for rows.Next() {
		var g Guideline
		var tagsJSON sql.NullString
		if err := rows.Scan(&g.ID, &g.ProjectID, &g.Category, &g.Title, &g.Content, &tagsJSON, &g.Priority, &g.CreatedAt, &g.UpdatedAt); err != nil {
			return nil, err
		}
		if tagsJSON.Valid {
			json.Unmarshal([]byte(tagsJSON.String), &g.Tags)
		}
		guidelines = append(guidelines, g)
	}
	return guidelines, rows.Err()
}

// SearchGuidelines searches guidelines by content and tags
func (db *DB) SearchGuidelines(projectID *int64, query string, category *string) ([]Guideline, error) {
	pid := db.GetProjectID(projectID)

	var conditions []string
	var args []interface{}

	conditions = append(conditions, "project_id = ?")
	args = append(args, pid)

	if query != "" {
		conditions = append(conditions, "(content LIKE ? OR title LIKE ? OR tags LIKE ?)")
		likeQuery := "%" + query + "%"
		args = append(args, likeQuery, likeQuery, likeQuery)
	}

	if category != nil {
		conditions = append(conditions, "category = ?")
		args = append(args, *category)
	}

	sqlQuery := fmt.Sprintf(
		"SELECT id, project_id, category, title, content, tags, priority, created_at, updated_at FROM guidelines WHERE %s ORDER BY priority DESC, updated_at DESC",
		strings.Join(conditions, " AND "),
	)

	rows, err := db.Query(sqlQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var guidelines []Guideline
	for rows.Next() {
		var g Guideline
		var tagsJSON sql.NullString
		if err := rows.Scan(&g.ID, &g.ProjectID, &g.Category, &g.Title, &g.Content, &tagsJSON, &g.Priority, &g.CreatedAt, &g.UpdatedAt); err != nil {
			return nil, err
		}
		if tagsJSON.Valid {
			json.Unmarshal([]byte(tagsJSON.String), &g.Tags)
		}
		guidelines = append(guidelines, g)
	}
	return guidelines, rows.Err()
}

// DeleteGuideline deletes a guideline
func (db *DB) DeleteGuideline(id int64) error {
	_, err := db.Exec("DELETE FROM guidelines WHERE id = ?", id)
	return err
}
