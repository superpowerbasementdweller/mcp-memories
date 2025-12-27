package mcp

import (
	"errors"
	"fmt"

	"github.com/rocket/mcp-memories/internal/db"
)

var ErrUnknownTool = errors.New("unknown tool")

// HandleToolCall routes a tool call to the appropriate handler
func HandleToolCall(database *db.DB, name string, args map[string]interface{}) (interface{}, error) {
	switch name {
	// Memory tools
	case "memory_store":
		return handleMemoryStore(database, args)
	case "memory_search":
		return handleMemorySearch(database, args)
	case "memory_delete":
		return handleMemoryDelete(database, args)

	// Task tools
	case "task_create":
		return handleTaskCreate(database, args)
	case "task_update":
		return handleTaskUpdate(database, args)
	case "task_list":
		return handleTaskList(database, args)
	case "task_delete":
		return handleTaskDelete(database, args)

	// Metadata tools
	case "metadata_set":
		return handleMetadataSet(database, args)
	case "metadata_get":
		return handleMetadataGet(database, args)
	case "metadata_list":
		return handleMetadataList(database, args)
	case "metadata_delete":
		return handleMetadataDelete(database, args)

	// Filetree tools
	case "filetree_annotate":
		return handleFiletreeAnnotate(database, args)
	case "filetree_get":
		return handleFiletreeGet(database, args)
	case "filetree_delete":
		return handleFiletreeDelete(database, args)

	// Guideline tools
	case "guideline_create":
		return handleGuidelineCreate(database, args)
	case "guideline_update":
		return handleGuidelineUpdate(database, args)
	case "guideline_list":
		return handleGuidelineList(database, args)
	case "guideline_search":
		return handleGuidelineSearch(database, args)
	case "guideline_get":
		return handleGuidelineGet(database, args)
	case "guideline_delete":
		return handleGuidelineDelete(database, args)

	// Project tools
	case "project_create":
		return handleProjectCreate(database, args)
	case "project_list":
		return handleProjectList(database, args)
	case "project_set_default":
		return handleProjectSetDefault(database, args)

	// Bookmark tools
	case "bookmark_create":
		return handleBookmarkCreate(database, args)
	case "bookmark_search":
		return handleBookmarkSearch(database, args)
	case "bookmark_list":
		return handleBookmarkList(database, args)
	case "bookmark_delete":
		return handleBookmarkDelete(database, args)

	default:
		return nil, fmt.Errorf("%w: %s", ErrUnknownTool, name)
	}
}

// Helper functions for argument extraction
func getString(args map[string]interface{}, key string) string {
	if v, ok := args[key].(string); ok {
		return v
	}
	return ""
}

func getStringPtr(args map[string]interface{}, key string) *string {
	if v, ok := args[key].(string); ok {
		return &v
	}
	return nil
}

func getInt(args map[string]interface{}, key string) int {
	if v, ok := args[key].(float64); ok {
		return int(v)
	}
	return 0
}

func getIntPtr(args map[string]interface{}, key string) *int {
	if v, ok := args[key].(float64); ok {
		i := int(v)
		return &i
	}
	return nil
}

func getInt64(args map[string]interface{}, key string) int64 {
	if v, ok := args[key].(float64); ok {
		return int64(v)
	}
	return 0
}

func getInt64Ptr(args map[string]interface{}, key string) *int64 {
	if v, ok := args[key].(float64); ok {
		i := int64(v)
		return &i
	}
	return nil
}

func getBool(args map[string]interface{}, key string) bool {
	if v, ok := args[key].(bool); ok {
		return v
	}
	return false
}

func getStringArray(args map[string]interface{}, key string) []string {
	if v, ok := args[key].([]interface{}); ok {
		result := make([]string, 0, len(v))
		for _, item := range v {
			if s, ok := item.(string); ok {
				result = append(result, s)
			}
		}
		return result
	}
	return nil
}

func getStringArrayPtr(args map[string]interface{}, key string) *[]string {
	if arr := getStringArray(args, key); arr != nil {
		return &arr
	}
	return nil
}

func getProjectID(database *db.DB, args map[string]interface{}) *int64 {
	if slug := getString(args, "project"); slug != "" {
		if p, err := database.GetOrCreateProject(slug); err == nil {
			return &p.ID
		}
	}
	return nil
}

// Memory handlers
func handleMemoryStore(database *db.DB, args map[string]interface{}) (interface{}, error) {
	content := getString(args, "content")
	if content == "" {
		return nil, fmt.Errorf("content is required")
	}
	keywords := getStringArray(args, "keywords")
	projectID := getProjectID(database, args)
	return database.CreateMemory(projectID, content, keywords)
}

func handleMemorySearch(database *db.DB, args map[string]interface{}) (interface{}, error) {
	query := getString(args, "query")
	keywords := getStringArray(args, "keywords")
	projectID := getProjectID(database, args)
	limit := getInt(args, "limit")
	if limit == 0 {
		limit = 20
	}
	return database.SearchMemories(projectID, query, keywords, limit)
}

func handleMemoryDelete(database *db.DB, args map[string]interface{}) (interface{}, error) {
	id := getInt64(args, "id")
	if id == 0 {
		return nil, fmt.Errorf("id is required")
	}
	if err := database.DeleteMemory(id); err != nil {
		return nil, err
	}
	return map[string]interface{}{"deleted": true, "id": id}, nil
}

// Task handlers
func handleTaskCreate(database *db.DB, args map[string]interface{}) (interface{}, error) {
	title := getString(args, "title")
	if title == "" {
		return nil, fmt.Errorf("title is required")
	}
	description := getString(args, "description")
	parentID := getInt64Ptr(args, "parent_id")
	priority := getInt(args, "priority")
	projectID := getProjectID(database, args)
	return database.CreateTask(projectID, parentID, title, description, priority)
}

func handleTaskUpdate(database *db.DB, args map[string]interface{}) (interface{}, error) {
	id := getInt64(args, "id")
	if id == 0 {
		return nil, fmt.Errorf("id is required")
	}
	return database.UpdateTask(id, getStringPtr(args, "title"), getStringPtr(args, "description"), getStringPtr(args, "status"), getIntPtr(args, "priority"))
}

func handleTaskList(database *db.DB, args map[string]interface{}) (interface{}, error) {
	projectID := getProjectID(database, args)
	status := getStringPtr(args, "status")
	parentID := getInt64Ptr(args, "parent_id")
	return database.ListTasks(projectID, status, parentID)
}

func handleTaskDelete(database *db.DB, args map[string]interface{}) (interface{}, error) {
	id := getInt64(args, "id")
	if id == 0 {
		return nil, fmt.Errorf("id is required")
	}
	if err := database.DeleteTask(id); err != nil {
		return nil, err
	}
	return map[string]interface{}{"deleted": true, "id": id}, nil
}

// Metadata handlers
func handleMetadataSet(database *db.DB, args map[string]interface{}) (interface{}, error) {
	key := getString(args, "key")
	if key == "" {
		return nil, fmt.Errorf("key is required")
	}
	value := getString(args, "value")
	projectID := getProjectID(database, args)
	return database.SetMetadata(projectID, key, value)
}

func handleMetadataGet(database *db.DB, args map[string]interface{}) (interface{}, error) {
	key := getString(args, "key")
	if key == "" {
		return nil, fmt.Errorf("key is required")
	}
	projectID := getProjectID(database, args)
	m, err := database.GetMetadata(projectID, key)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return map[string]interface{}{"key": key, "value": nil}, nil
	}
	return m, nil
}

func handleMetadataList(database *db.DB, args map[string]interface{}) (interface{}, error) {
	projectID := getProjectID(database, args)
	return database.ListMetadata(projectID)
}

func handleMetadataDelete(database *db.DB, args map[string]interface{}) (interface{}, error) {
	key := getString(args, "key")
	if key == "" {
		return nil, fmt.Errorf("key is required")
	}
	projectID := getProjectID(database, args)
	if err := database.DeleteMetadata(projectID, key); err != nil {
		return nil, err
	}
	return map[string]interface{}{"deleted": true, "key": key}, nil
}

// Filetree handlers
func handleFiletreeAnnotate(database *db.DB, args map[string]interface{}) (interface{}, error) {
	path := getString(args, "path")
	if path == "" {
		return nil, fmt.Errorf("path is required")
	}
	note := getString(args, "note")
	if note == "" {
		return nil, fmt.Errorf("note is required")
	}
	isDir := getBool(args, "is_dir")
	projectID := getProjectID(database, args)
	return database.AnnotateFile(projectID, path, note, isDir)
}

func handleFiletreeGet(database *db.DB, args map[string]interface{}) (interface{}, error) {
	projectID := getProjectID(database, args)
	if path := getString(args, "path"); path != "" {
		return database.GetFileAnnotation(projectID, path)
	}
	return database.ListFileAnnotations(projectID)
}

func handleFiletreeDelete(database *db.DB, args map[string]interface{}) (interface{}, error) {
	path := getString(args, "path")
	if path == "" {
		return nil, fmt.Errorf("path is required")
	}
	projectID := getProjectID(database, args)
	if err := database.DeleteFileAnnotation(projectID, path); err != nil {
		return nil, err
	}
	return map[string]interface{}{"deleted": true, "path": path}, nil
}

// Guideline handlers
func handleGuidelineCreate(database *db.DB, args map[string]interface{}) (interface{}, error) {
	category := getString(args, "category")
	if category == "" {
		return nil, fmt.Errorf("category is required")
	}
	title := getString(args, "title")
	if title == "" {
		return nil, fmt.Errorf("title is required")
	}
	content := getString(args, "content")
	if content == "" {
		return nil, fmt.Errorf("content is required")
	}
	tags := getStringArray(args, "tags")
	priority := getInt(args, "priority")
	projectID := getProjectID(database, args)
	return database.CreateGuideline(projectID, category, title, content, tags, priority)
}

func handleGuidelineUpdate(database *db.DB, args map[string]interface{}) (interface{}, error) {
	id := getInt64(args, "id")
	if id == 0 {
		return nil, fmt.Errorf("id is required")
	}
	return database.UpdateGuideline(id, getStringPtr(args, "content"), getStringArrayPtr(args, "tags"), getIntPtr(args, "priority"))
}

func handleGuidelineList(database *db.DB, args map[string]interface{}) (interface{}, error) {
	projectID := getProjectID(database, args)
	category := getStringPtr(args, "category")
	return database.ListGuidelines(projectID, category)
}

func handleGuidelineSearch(database *db.DB, args map[string]interface{}) (interface{}, error) {
	query := getString(args, "query")
	if query == "" {
		return nil, fmt.Errorf("query is required")
	}
	projectID := getProjectID(database, args)
	category := getStringPtr(args, "category")
	return database.SearchGuidelines(projectID, query, category)
}

func handleGuidelineGet(database *db.DB, args map[string]interface{}) (interface{}, error) {
	id := getInt64(args, "id")
	if id == 0 {
		return nil, fmt.Errorf("id is required")
	}
	return database.GetGuideline(id)
}

func handleGuidelineDelete(database *db.DB, args map[string]interface{}) (interface{}, error) {
	id := getInt64(args, "id")
	if id == 0 {
		return nil, fmt.Errorf("id is required")
	}
	if err := database.DeleteGuideline(id); err != nil {
		return nil, err
	}
	return map[string]interface{}{"deleted": true, "id": id}, nil
}

// Project handlers
func handleProjectCreate(database *db.DB, args map[string]interface{}) (interface{}, error) {
	slug := getString(args, "slug")
	if slug == "" {
		return nil, fmt.Errorf("slug is required")
	}
	name := getString(args, "name")
	rootPath := getString(args, "root_path")
	return database.CreateProject(slug, name, rootPath)
}

func handleProjectList(database *db.DB, args map[string]interface{}) (interface{}, error) {
	return database.ListProjects()
}

func handleProjectSetDefault(database *db.DB, args map[string]interface{}) (interface{}, error) {
	slug := getString(args, "slug")
	if slug == "" {
		return nil, fmt.Errorf("slug is required")
	}
	p, err := database.GetOrCreateProject(slug)
	if err != nil {
		return nil, err
	}
	database.SetDefaultProject(p.ID)
	return map[string]interface{}{"default_project": p}, nil
}

// Bookmark handlers
func handleBookmarkCreate(database *db.DB, args map[string]interface{}) (interface{}, error) {
	url := getString(args, "url")
	if url == "" {
		return nil, fmt.Errorf("url is required")
	}
	title := getString(args, "title")
	if title == "" {
		return nil, fmt.Errorf("title is required")
	}
	excerpt := getString(args, "excerpt")
	note := getString(args, "note")
	docType := getString(args, "doc_type")
	pageOrSection := getString(args, "page_or_section")
	tags := getStringArray(args, "tags")
	projectID := getProjectID(database, args)
	return database.CreateBookmark(projectID, url, title, excerpt, note, docType, pageOrSection, tags)
}

func handleBookmarkSearch(database *db.DB, args map[string]interface{}) (interface{}, error) {
	query := getString(args, "query")
	tags := getStringArray(args, "tags")
	docType := getStringPtr(args, "doc_type")
	projectID := getProjectID(database, args)
	return database.SearchBookmarks(projectID, query, tags, docType)
}

func handleBookmarkList(database *db.DB, args map[string]interface{}) (interface{}, error) {
	projectID := getProjectID(database, args)
	return database.ListBookmarks(projectID)
}

func handleBookmarkDelete(database *db.DB, args map[string]interface{}) (interface{}, error) {
	id := getInt64(args, "id")
	if id == 0 {
		return nil, fmt.Errorf("id is required")
	}
	if err := database.DeleteBookmark(id); err != nil {
		return nil, err
	}
	return map[string]interface{}{"deleted": true, "id": id}, nil
}
