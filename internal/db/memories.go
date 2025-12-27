package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Memory represents a stored memory
type Memory struct {
	ID        int64     `json:"id"`
	ProjectID int64     `json:"project_id"`
	Content   string    `json:"content"`
	Keywords  []string  `json:"keywords"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateMemory creates a new memory
func (db *DB) CreateMemory(projectID *int64, content string, keywords []string) (*Memory, error) {
	pid := db.GetProjectID(projectID)

	keywordsJSON, err := json.Marshal(keywords)
	if err != nil {
		return nil, fmt.Errorf("marshaling keywords: %w", err)
	}

	result, err := db.Exec(
		"INSERT INTO memories (project_id, content, keywords) VALUES (?, ?, ?)",
		pid, content, string(keywordsJSON),
	)
	if err != nil {
		return nil, fmt.Errorf("creating memory: %w", err)
	}

	id, _ := result.LastInsertId()
	return db.GetMemory(id)
}

// GetMemory gets a memory by ID
func (db *DB) GetMemory(id int64) (*Memory, error) {
	m := &Memory{}
	var keywordsJSON sql.NullString
	err := db.QueryRow(
		"SELECT id, project_id, content, keywords, created_at, updated_at FROM memories WHERE id = ?",
		id,
	).Scan(&m.ID, &m.ProjectID, &m.Content, &keywordsJSON, &m.CreatedAt, &m.UpdatedAt)
	if err != nil {
		return nil, err
	}
	if keywordsJSON.Valid {
		json.Unmarshal([]byte(keywordsJSON.String), &m.Keywords)
	}
	return m, nil
}

// SearchMemories searches memories by content and/or keywords
func (db *DB) SearchMemories(projectID *int64, query string, keywords []string, limit int) ([]Memory, error) {
	pid := db.GetProjectID(projectID)

	var conditions []string
	var args []interface{}

	conditions = append(conditions, "project_id = ?")
	args = append(args, pid)

	if query != "" {
		conditions = append(conditions, "content LIKE ?")
		args = append(args, "%"+query+"%")
	}

	for _, kw := range keywords {
		conditions = append(conditions, "keywords LIKE ?")
		args = append(args, "%\""+kw+"\"%")
	}

	sqlQuery := fmt.Sprintf(
		"SELECT id, project_id, content, keywords, created_at, updated_at FROM memories WHERE %s ORDER BY updated_at DESC",
		strings.Join(conditions, " AND "),
	)

	if limit > 0 {
		sqlQuery += fmt.Sprintf(" LIMIT %d", limit)
	}

	rows, err := db.Query(sqlQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var memories []Memory
	for rows.Next() {
		var m Memory
		var keywordsJSON sql.NullString
		if err := rows.Scan(&m.ID, &m.ProjectID, &m.Content, &keywordsJSON, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, err
		}
		if keywordsJSON.Valid {
			json.Unmarshal([]byte(keywordsJSON.String), &m.Keywords)
		}
		memories = append(memories, m)
	}
	return memories, rows.Err()
}

// DeleteMemory deletes a memory by ID
func (db *DB) DeleteMemory(id int64) error {
	_, err := db.Exec("DELETE FROM memories WHERE id = ?", id)
	return err
}
