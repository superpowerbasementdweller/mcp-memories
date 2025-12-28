package mcp

import (
	"testing"

	"github.com/rocket/mcp-memories/internal/db"
)

// TestAllTools is a comprehensive integration test that exercises all 24 MCP tools
// by storing data and then recalling it to verify functionality.
func TestAllTools(t *testing.T) {
	// Setup: Create in-memory database for test isolation
	database, err := db.Open(":memory:")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer database.Close()

	// Track IDs from create operations for use in recall/delete operations
	var (
		memoryID    int64
		taskID      int64
		guidelineID int64
		bookmarkID  int64
	)

	// ========================================
	// PROJECT TOOLS (3 tools)
	// ========================================
	t.Run("project_create", func(t *testing.T) {
		result, err := HandleToolCall(database, "project_create", map[string]interface{}{
			"slug":      "test-project",
			"name":      "Test Project",
			"root_path": "/path/to/project",
		})
		if err != nil {
			t.Fatalf("project_create failed: %v", err)
		}
		if result == nil {
			t.Error("project_create returned nil")
		}
		t.Logf("project_create result: %+v", result)
	})

	t.Run("project_list", func(t *testing.T) {
		result, err := HandleToolCall(database, "project_list", map[string]interface{}{})
		if err != nil {
			t.Fatalf("project_list failed: %v", err)
		}
		projects, ok := result.([]db.Project)
		if !ok {
			t.Fatalf("unexpected result type: %T", result)
		}
		// Should have at least the default project and our test project
		if len(projects) < 1 {
			t.Error("project_list returned no projects")
		}
		t.Logf("project_list found %d projects", len(projects))
	})

	t.Run("project_set_default", func(t *testing.T) {
		result, err := HandleToolCall(database, "project_set_default", map[string]interface{}{
			"slug": "test-project",
		})
		if err != nil {
			t.Fatalf("project_set_default failed: %v", err)
		}
		if result == nil {
			t.Error("project_set_default returned nil")
		}
		t.Logf("project_set_default result: %+v", result)
	})

	// ========================================
	// MEMORY TOOLS (3 tools)
	// ========================================
	t.Run("memory_store", func(t *testing.T) {
		result, err := HandleToolCall(database, "memory_store", map[string]interface{}{
			"content":  "This is a test memory about Go programming",
			"keywords": []interface{}{"go", "programming", "test"},
		})
		if err != nil {
			t.Fatalf("memory_store failed: %v", err)
		}
		memory, ok := result.(*db.Memory)
		if !ok {
			t.Fatalf("unexpected result type: %T", result)
		}
		memoryID = memory.ID
		if memory.Content != "This is a test memory about Go programming" {
			t.Errorf("memory content mismatch: %s", memory.Content)
		}
		t.Logf("memory_store created memory ID: %d", memoryID)
	})

	t.Run("memory_search", func(t *testing.T) {
		result, err := HandleToolCall(database, "memory_search", map[string]interface{}{
			"query":    "Go programming",
			"keywords": []interface{}{"go"},
		})
		if err != nil {
			t.Fatalf("memory_search failed: %v", err)
		}
		memories, ok := result.([]db.Memory)
		if !ok {
			t.Fatalf("unexpected result type: %T", result)
		}
		if len(memories) == 0 {
			t.Error("memory_search returned no results")
		} else {
			t.Logf("memory_search found %d memories", len(memories))
		}
	})

	t.Run("memory_delete", func(t *testing.T) {
		result, err := HandleToolCall(database, "memory_delete", map[string]interface{}{
			"id": float64(memoryID), // JSON numbers are float64
		})
		if err != nil {
			t.Fatalf("memory_delete failed: %v", err)
		}
		resultMap, ok := result.(map[string]interface{})
		if !ok || resultMap["deleted"] != true {
			t.Error("memory_delete did not return deleted: true")
		}
		t.Logf("memory_delete succeeded for ID: %d", memoryID)
	})

	// ========================================
	// TASK TOOLS (4 tools)
	// ========================================
	t.Run("task_create", func(t *testing.T) {
		result, err := HandleToolCall(database, "task_create", map[string]interface{}{
			"title":       "Test Task",
			"description": "A test task for integration testing",
			"priority":    float64(1),
		})
		if err != nil {
			t.Fatalf("task_create failed: %v", err)
		}
		task, ok := result.(*db.Task)
		if !ok {
			t.Fatalf("unexpected result type: %T", result)
		}
		taskID = task.ID
		t.Logf("task_create created task ID: %d", taskID)
	})

	t.Run("task_list", func(t *testing.T) {
		result, err := HandleToolCall(database, "task_list", map[string]interface{}{})
		if err != nil {
			t.Fatalf("task_list failed: %v", err)
		}
		tasks, ok := result.([]db.Task)
		if !ok {
			t.Fatalf("unexpected result type: %T", result)
		}
		if len(tasks) == 0 {
			t.Error("task_list returned no tasks")
		}
		t.Logf("task_list found %d tasks", len(tasks))
	})

	t.Run("task_update", func(t *testing.T) {
		result, err := HandleToolCall(database, "task_update", map[string]interface{}{
			"id":     float64(taskID),
			"status": "in_progress",
			"title":  "Updated Test Task",
		})
		if err != nil {
			t.Fatalf("task_update failed: %v", err)
		}
		task, ok := result.(*db.Task)
		if !ok {
			t.Fatalf("unexpected result type: %T", result)
		}
		if task.Status != "in_progress" {
			t.Errorf("task status not updated: %s", task.Status)
		}
		t.Logf("task_update succeeded: status=%s", task.Status)
	})

	t.Run("task_delete", func(t *testing.T) {
		result, err := HandleToolCall(database, "task_delete", map[string]interface{}{
			"id": float64(taskID),
		})
		if err != nil {
			t.Fatalf("task_delete failed: %v", err)
		}
		resultMap, ok := result.(map[string]interface{})
		if !ok || resultMap["deleted"] != true {
			t.Error("task_delete did not return deleted: true")
		}
		t.Logf("task_delete succeeded for ID: %d", taskID)
	})

	// ========================================
	// METADATA TOOLS (4 tools)
	// ========================================
	t.Run("metadata_set", func(t *testing.T) {
		result, err := HandleToolCall(database, "metadata_set", map[string]interface{}{
			"key":   "test_key",
			"value": "test_value",
		})
		if err != nil {
			t.Fatalf("metadata_set failed: %v", err)
		}
		if result == nil {
			t.Error("metadata_set returned nil")
		}
		t.Logf("metadata_set result: %+v", result)
	})

	t.Run("metadata_get", func(t *testing.T) {
		result, err := HandleToolCall(database, "metadata_get", map[string]interface{}{
			"key": "test_key",
		})
		if err != nil {
			t.Fatalf("metadata_get failed: %v", err)
		}
		meta, ok := result.(*db.Metadata)
		if !ok {
			t.Fatalf("unexpected result type: %T", result)
		}
		if meta.Value != "test_value" {
			t.Errorf("metadata value mismatch: %s", meta.Value)
		}
		t.Logf("metadata_get found: key=%s, value=%s", meta.Key, meta.Value)
	})

	t.Run("metadata_list", func(t *testing.T) {
		result, err := HandleToolCall(database, "metadata_list", map[string]interface{}{})
		if err != nil {
			t.Fatalf("metadata_list failed: %v", err)
		}
		metadata, ok := result.([]db.Metadata)
		if !ok {
			t.Fatalf("unexpected result type: %T", result)
		}
		if len(metadata) == 0 {
			t.Error("metadata_list returned no metadata")
		}
		t.Logf("metadata_list found %d items", len(metadata))
	})

	t.Run("metadata_delete", func(t *testing.T) {
		result, err := HandleToolCall(database, "metadata_delete", map[string]interface{}{
			"key": "test_key",
		})
		if err != nil {
			t.Fatalf("metadata_delete failed: %v", err)
		}
		resultMap, ok := result.(map[string]interface{})
		if !ok || resultMap["deleted"] != true {
			t.Error("metadata_delete did not return deleted: true")
		}
		t.Logf("metadata_delete succeeded for key: test_key")
	})

	// ========================================
	// FILETREE TOOLS (3 tools)
	// ========================================
	t.Run("filetree_annotate", func(t *testing.T) {
		result, err := HandleToolCall(database, "filetree_annotate", map[string]interface{}{
			"path":   "/src/main.go",
			"note":   "Main entry point for the application",
			"is_dir": false,
		})
		if err != nil {
			t.Fatalf("filetree_annotate failed: %v", err)
		}
		if result == nil {
			t.Error("filetree_annotate returned nil")
		}
		t.Logf("filetree_annotate result: %+v", result)
	})

	t.Run("filetree_get", func(t *testing.T) {
		result, err := HandleToolCall(database, "filetree_get", map[string]interface{}{
			"path": "/src/main.go",
		})
		if err != nil {
			t.Fatalf("filetree_get failed: %v", err)
		}
		annotation, ok := result.(*db.FileAnnotation)
		if !ok {
			t.Fatalf("unexpected result type: %T", result)
		}
		if annotation.Note != "Main entry point for the application" {
			t.Errorf("annotation note mismatch: %s", annotation.Note)
		}
		t.Logf("filetree_get found: path=%s, note=%s", annotation.Path, annotation.Note)
	})

	t.Run("filetree_delete", func(t *testing.T) {
		result, err := HandleToolCall(database, "filetree_delete", map[string]interface{}{
			"path": "/src/main.go",
		})
		if err != nil {
			t.Fatalf("filetree_delete failed: %v", err)
		}
		resultMap, ok := result.(map[string]interface{})
		if !ok || resultMap["deleted"] != true {
			t.Error("filetree_delete did not return deleted: true")
		}
		t.Logf("filetree_delete succeeded for path: /src/main.go")
	})

	// ========================================
	// GUIDELINE TOOLS (6 tools)
	// ========================================
	t.Run("guideline_create", func(t *testing.T) {
		result, err := HandleToolCall(database, "guideline_create", map[string]interface{}{
			"category": "coding_style",
			"title":    "Go Error Handling",
			"content":  "Always handle errors explicitly. Never ignore returned errors.",
			"tags":     []interface{}{"go", "errors", "best-practices"},
			"priority": float64(10),
		})
		if err != nil {
			t.Fatalf("guideline_create failed: %v", err)
		}
		guideline, ok := result.(*db.Guideline)
		if !ok {
			t.Fatalf("unexpected result type: %T", result)
		}
		guidelineID = guideline.ID
		t.Logf("guideline_create created guideline ID: %d", guidelineID)
	})

	t.Run("guideline_get", func(t *testing.T) {
		result, err := HandleToolCall(database, "guideline_get", map[string]interface{}{
			"id": float64(guidelineID),
		})
		if err != nil {
			t.Fatalf("guideline_get failed: %v", err)
		}
		guideline, ok := result.(*db.Guideline)
		if !ok {
			t.Fatalf("unexpected result type: %T", result)
		}
		if guideline.Title != "Go Error Handling" {
			t.Errorf("guideline title mismatch: %s", guideline.Title)
		}
		t.Logf("guideline_get found: title=%s", guideline.Title)
	})

	t.Run("guideline_update", func(t *testing.T) {
		result, err := HandleToolCall(database, "guideline_update", map[string]interface{}{
			"id":       float64(guidelineID),
			"content":  "Always handle errors explicitly. Never ignore returned errors. Use errors.Is and errors.As for error checking.",
			"priority": float64(20),
		})
		if err != nil {
			t.Fatalf("guideline_update failed: %v", err)
		}
		guideline, ok := result.(*db.Guideline)
		if !ok {
			t.Fatalf("unexpected result type: %T", result)
		}
		if guideline.Priority != 20 {
			t.Errorf("guideline priority not updated: %d", guideline.Priority)
		}
		t.Logf("guideline_update succeeded: priority=%d", guideline.Priority)
	})

	t.Run("guideline_list", func(t *testing.T) {
		result, err := HandleToolCall(database, "guideline_list", map[string]interface{}{
			"category": "coding_style",
		})
		if err != nil {
			t.Fatalf("guideline_list failed: %v", err)
		}
		guidelines, ok := result.([]db.Guideline)
		if !ok {
			t.Fatalf("unexpected result type: %T", result)
		}
		if len(guidelines) == 0 {
			t.Error("guideline_list returned no guidelines")
		}
		t.Logf("guideline_list found %d guidelines", len(guidelines))
	})

	t.Run("guideline_search", func(t *testing.T) {
		result, err := HandleToolCall(database, "guideline_search", map[string]interface{}{
			"query": "error handling",
		})
		if err != nil {
			t.Fatalf("guideline_search failed: %v", err)
		}
		guidelines, ok := result.([]db.Guideline)
		if !ok {
			t.Fatalf("unexpected result type: %T", result)
		}
		if len(guidelines) == 0 {
			t.Error("guideline_search returned no results")
		}
		t.Logf("guideline_search found %d guidelines", len(guidelines))
	})

	t.Run("guideline_delete", func(t *testing.T) {
		result, err := HandleToolCall(database, "guideline_delete", map[string]interface{}{
			"id": float64(guidelineID),
		})
		if err != nil {
			t.Fatalf("guideline_delete failed: %v", err)
		}
		resultMap, ok := result.(map[string]interface{})
		if !ok || resultMap["deleted"] != true {
			t.Error("guideline_delete did not return deleted: true")
		}
		t.Logf("guideline_delete succeeded for ID: %d", guidelineID)
	})

	// ========================================
	// BOOKMARK TOOLS (4 tools)
	// ========================================
	t.Run("bookmark_create", func(t *testing.T) {
		result, err := HandleToolCall(database, "bookmark_create", map[string]interface{}{
			"url":             "https://go.dev/doc/effective_go",
			"title":           "Effective Go",
			"excerpt":         "This document gives tips for writing clear, idiomatic Go code.",
			"note":            "Great reference for Go best practices",
			"doc_type":        "url",
			"page_or_section": "Introduction",
			"tags":            []interface{}{"go", "documentation", "guide"},
		})
		if err != nil {
			t.Fatalf("bookmark_create failed: %v", err)
		}
		bookmark, ok := result.(*db.Bookmark)
		if !ok {
			t.Fatalf("unexpected result type: %T", result)
		}
		bookmarkID = bookmark.ID
		t.Logf("bookmark_create created bookmark ID: %d", bookmarkID)
	})

	t.Run("bookmark_search", func(t *testing.T) {
		result, err := HandleToolCall(database, "bookmark_search", map[string]interface{}{
			"query": "Effective Go",
			"tags":  []interface{}{"go"},
		})
		if err != nil {
			t.Fatalf("bookmark_search failed: %v", err)
		}
		bookmarks, ok := result.([]db.Bookmark)
		if !ok {
			t.Fatalf("unexpected result type: %T", result)
		}
		if len(bookmarks) == 0 {
			t.Error("bookmark_search returned no results")
		}
		t.Logf("bookmark_search found %d bookmarks", len(bookmarks))
	})

	t.Run("bookmark_list", func(t *testing.T) {
		result, err := HandleToolCall(database, "bookmark_list", map[string]interface{}{})
		if err != nil {
			t.Fatalf("bookmark_list failed: %v", err)
		}
		bookmarks, ok := result.([]db.Bookmark)
		if !ok {
			t.Fatalf("unexpected result type: %T", result)
		}
		if len(bookmarks) == 0 {
			t.Error("bookmark_list returned no bookmarks")
		}
		t.Logf("bookmark_list found %d bookmarks", len(bookmarks))
	})

	t.Run("bookmark_delete", func(t *testing.T) {
		result, err := HandleToolCall(database, "bookmark_delete", map[string]interface{}{
			"id": float64(bookmarkID),
		})
		if err != nil {
			t.Fatalf("bookmark_delete failed: %v", err)
		}
		resultMap, ok := result.(map[string]interface{})
		if !ok || resultMap["deleted"] != true {
			t.Error("bookmark_delete did not return deleted: true")
		}
		t.Logf("bookmark_delete succeeded for ID: %d", bookmarkID)
	})
}
