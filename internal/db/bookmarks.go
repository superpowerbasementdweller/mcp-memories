package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Bookmark represents a reference to an external document
type Bookmark struct {
	ID            int64     `json:"id"`
	ProjectID     int64     `json:"project_id"`
	URL           string    `json:"url"`
	Title         string    `json:"title"`
	Excerpt       string    `json:"excerpt,omitempty"`
	Note          string    `json:"note,omitempty"`
	DocType       string    `json:"doc_type,omitempty"`
	PageOrSection string    `json:"page_or_section,omitempty"`
	Tags          []string  `json:"tags,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}

// CreateBookmark creates a new bookmark
func (db *DB) CreateBookmark(projectID *int64, url, title, excerpt, note, docType, pageOrSection string, tags []string) (*Bookmark, error) {
	pid := db.GetProjectID(projectID)

	tagsJSON, err := json.Marshal(tags)
	if err != nil {
		return nil, fmt.Errorf("marshaling tags: %w", err)
	}

	result, err := db.Exec(
		"INSERT INTO bookmarks (project_id, url, title, excerpt, note, doc_type, page_or_section, tags) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		pid, url, title, excerpt, note, docType, pageOrSection, string(tagsJSON),
	)
	if err != nil {
		return nil, fmt.Errorf("creating bookmark: %w", err)
	}

	id, _ := result.LastInsertId()
	return db.GetBookmark(id)
}

// GetBookmark gets a bookmark by ID
func (db *DB) GetBookmark(id int64) (*Bookmark, error) {
	b := &Bookmark{}
	var excerpt, note, docType, pageOrSection, tagsJSON sql.NullString
	err := db.QueryRow(
		"SELECT id, project_id, url, title, excerpt, note, doc_type, page_or_section, tags, created_at FROM bookmarks WHERE id = ?",
		id,
	).Scan(&b.ID, &b.ProjectID, &b.URL, &b.Title, &excerpt, &note, &docType, &pageOrSection, &tagsJSON, &b.CreatedAt)
	if err != nil {
		return nil, err
	}
	b.Excerpt = excerpt.String
	b.Note = note.String
	b.DocType = docType.String
	b.PageOrSection = pageOrSection.String
	if tagsJSON.Valid {
		json.Unmarshal([]byte(tagsJSON.String), &b.Tags)
	}
	return b, nil
}

// SearchBookmarks searches bookmarks by query and/or tags
func (db *DB) SearchBookmarks(projectID *int64, query string, tags []string, docType *string) ([]Bookmark, error) {
	pid := db.GetProjectID(projectID)

	var conditions []string
	var args []interface{}

	conditions = append(conditions, "project_id = ?")
	args = append(args, pid)

	if query != "" {
		conditions = append(conditions, "(title LIKE ? OR excerpt LIKE ? OR note LIKE ? OR url LIKE ?)")
		likeQuery := "%" + query + "%"
		args = append(args, likeQuery, likeQuery, likeQuery, likeQuery)
	}

	for _, tag := range tags {
		conditions = append(conditions, "tags LIKE ?")
		args = append(args, "%\""+tag+"\"%")
	}

	if docType != nil && *docType != "" {
		conditions = append(conditions, "doc_type = ?")
		args = append(args, *docType)
	}

	sqlQuery := fmt.Sprintf(
		"SELECT id, project_id, url, title, excerpt, note, doc_type, page_or_section, tags, created_at FROM bookmarks WHERE %s ORDER BY created_at DESC",
		strings.Join(conditions, " AND "),
	)

	rows, err := db.Query(sqlQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookmarks []Bookmark
	for rows.Next() {
		var b Bookmark
		var excerpt, note, docType, pageOrSection, tagsJSON sql.NullString
		if err := rows.Scan(&b.ID, &b.ProjectID, &b.URL, &b.Title, &excerpt, &note, &docType, &pageOrSection, &tagsJSON, &b.CreatedAt); err != nil {
			return nil, err
		}
		b.Excerpt = excerpt.String
		b.Note = note.String
		b.DocType = docType.String
		b.PageOrSection = pageOrSection.String
		if tagsJSON.Valid {
			json.Unmarshal([]byte(tagsJSON.String), &b.Tags)
		}
		bookmarks = append(bookmarks, b)
	}
	return bookmarks, rows.Err()
}

// ListBookmarks lists all bookmarks for a project
func (db *DB) ListBookmarks(projectID *int64) ([]Bookmark, error) {
	return db.SearchBookmarks(projectID, "", nil, nil)
}

// DeleteBookmark deletes a bookmark by ID
func (db *DB) DeleteBookmark(id int64) error {
	_, err := db.Exec("DELETE FROM bookmarks WHERE id = ?", id)
	return err
}
