package db

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

// Task represents a task
type Task struct {
	ID          int64     `json:"id"`
	ProjectID   int64     `json:"project_id"`
	ParentID    *int64    `json:"parent_id,omitempty"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	Status      string    `json:"status"`
	Priority    int       `json:"priority"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateTask creates a new task
func (db *DB) CreateTask(projectID *int64, parentID *int64, title, description string, priority int) (*Task, error) {
	pid := db.GetProjectID(projectID)

	result, err := db.Exec(
		"INSERT INTO tasks (project_id, parent_id, title, description, priority) VALUES (?, ?, ?, ?, ?)",
		pid, parentID, title, description, priority,
	)
	if err != nil {
		return nil, fmt.Errorf("creating task: %w", err)
	}

	id, _ := result.LastInsertId()
	return db.GetTask(id)
}

// GetTask gets a task by ID
func (db *DB) GetTask(id int64) (*Task, error) {
	t := &Task{}
	var parentID sql.NullInt64
	var description sql.NullString
	err := db.QueryRow(
		"SELECT id, project_id, parent_id, title, description, status, priority, created_at, updated_at FROM tasks WHERE id = ?",
		id,
	).Scan(&t.ID, &t.ProjectID, &parentID, &t.Title, &description, &t.Status, &t.Priority, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return nil, err
	}
	if parentID.Valid {
		t.ParentID = &parentID.Int64
	}
	t.Description = description.String
	return t, nil
}

// UpdateTask updates a task
func (db *DB) UpdateTask(id int64, title, description, status *string, priority *int) (*Task, error) {
	var sets []string
	var args []interface{}

	if title != nil {
		sets = append(sets, "title = ?")
		args = append(args, *title)
	}
	if description != nil {
		sets = append(sets, "description = ?")
		args = append(args, *description)
	}
	if status != nil {
		sets = append(sets, "status = ?")
		args = append(args, *status)
	}
	if priority != nil {
		sets = append(sets, "priority = ?")
		args = append(args, *priority)
	}

	if len(sets) == 0 {
		return db.GetTask(id)
	}

	sets = append(sets, "updated_at = CURRENT_TIMESTAMP")
	args = append(args, id)

	_, err := db.Exec(
		fmt.Sprintf("UPDATE tasks SET %s WHERE id = ?", strings.Join(sets, ", ")),
		args...,
	)
	if err != nil {
		return nil, fmt.Errorf("updating task: %w", err)
	}

	return db.GetTask(id)
}

// ListTasks lists tasks with optional filters
func (db *DB) ListTasks(projectID *int64, status *string, parentID *int64) ([]Task, error) {
	pid := db.GetProjectID(projectID)

	var conditions []string
	var args []interface{}

	conditions = append(conditions, "project_id = ?")
	args = append(args, pid)

	if status != nil {
		conditions = append(conditions, "status = ?")
		args = append(args, *status)
	}

	if parentID != nil {
		if *parentID == 0 {
			conditions = append(conditions, "parent_id IS NULL")
		} else {
			conditions = append(conditions, "parent_id = ?")
			args = append(args, *parentID)
		}
	}

	query := fmt.Sprintf(
		"SELECT id, project_id, parent_id, title, description, status, priority, created_at, updated_at FROM tasks WHERE %s ORDER BY priority DESC, created_at",
		strings.Join(conditions, " AND "),
	)

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var t Task
		var parentID sql.NullInt64
		var description sql.NullString
		if err := rows.Scan(&t.ID, &t.ProjectID, &parentID, &t.Title, &description, &t.Status, &t.Priority, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		if parentID.Valid {
			t.ParentID = &parentID.Int64
		}
		t.Description = description.String
		tasks = append(tasks, t)
	}
	return tasks, rows.Err()
}

// DeleteTask deletes a task and its subtasks
func (db *DB) DeleteTask(id int64) error {
	// First delete all subtasks recursively
	_, err := db.Exec("DELETE FROM tasks WHERE parent_id = ?", id)
	if err != nil {
		return err
	}
	_, err = db.Exec("DELETE FROM tasks WHERE id = ?", id)
	return err
}
