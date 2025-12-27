package mcp

// ToolDefinition represents an MCP tool definition
type ToolDefinition struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	InputSchema map[string]interface{} `json:"inputSchema"`
}

// GetToolDefinitions returns all available tool definitions
func GetToolDefinitions() []ToolDefinition {
	return []ToolDefinition{
		// Memory tools
		{
			Name:        "memory_store",
			Description: "Store a new memory with optional keywords for later retrieval",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"content":  map[string]interface{}{"type": "string", "description": "The content to remember"},
					"keywords": map[string]interface{}{"type": "array", "items": map[string]interface{}{"type": "string"}, "description": "Keywords for categorization and search"},
					"project":  map[string]interface{}{"type": "string", "description": "Project slug (optional, defaults to current project)"},
				},
				"required": []string{"content"},
			},
		},
		{
			Name:        "memory_search",
			Description: "Search memories by content and/or keywords",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"query":    map[string]interface{}{"type": "string", "description": "Text to search for in content"},
					"keywords": map[string]interface{}{"type": "array", "items": map[string]interface{}{"type": "string"}, "description": "Keywords to filter by"},
					"project":  map[string]interface{}{"type": "string", "description": "Project slug (optional)"},
					"limit":    map[string]interface{}{"type": "integer", "description": "Maximum results to return"},
				},
			},
		},
		{
			Name:        "memory_delete",
			Description: "Delete a memory by ID",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"id": map[string]interface{}{"type": "integer", "description": "Memory ID to delete"},
				},
				"required": []string{"id"},
			},
		},

		// Task tools
		{
			Name:        "task_create",
			Description: "Create a new task with optional parent for subtasks",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"title":       map[string]interface{}{"type": "string", "description": "Task title"},
					"description": map[string]interface{}{"type": "string", "description": "Detailed description"},
					"parent_id":   map[string]interface{}{"type": "integer", "description": "Parent task ID for subtasks"},
					"priority":    map[string]interface{}{"type": "integer", "description": "Priority (higher = more important)"},
					"project":     map[string]interface{}{"type": "string", "description": "Project slug (optional)"},
				},
				"required": []string{"title"},
			},
		},
		{
			Name:        "task_update",
			Description: "Update a task's status, title, description, or priority",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"id":          map[string]interface{}{"type": "integer", "description": "Task ID"},
					"status":      map[string]interface{}{"type": "string", "enum": []string{"todo", "in_progress", "done", "blocked"}, "description": "Task status"},
					"title":       map[string]interface{}{"type": "string", "description": "New title"},
					"description": map[string]interface{}{"type": "string", "description": "New description"},
					"priority":    map[string]interface{}{"type": "integer", "description": "New priority"},
				},
				"required": []string{"id"},
			},
		},
		{
			Name:        "task_list",
			Description: "List tasks with optional filters",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"project":   map[string]interface{}{"type": "string", "description": "Project slug (optional)"},
					"status":    map[string]interface{}{"type": "string", "enum": []string{"todo", "in_progress", "done", "blocked"}, "description": "Filter by status"},
					"parent_id": map[string]interface{}{"type": "integer", "description": "Filter by parent (0 for root tasks)"},
				},
			},
		},
		{
			Name:        "task_delete",
			Description: "Delete a task and its subtasks",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"id": map[string]interface{}{"type": "integer", "description": "Task ID to delete"},
				},
				"required": []string{"id"},
			},
		},

		// Metadata tools
		{
			Name:        "metadata_set",
			Description: "Set a key-value metadata pair for a project",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"key":     map[string]interface{}{"type": "string", "description": "Metadata key"},
					"value":   map[string]interface{}{"type": "string", "description": "Metadata value"},
					"project": map[string]interface{}{"type": "string", "description": "Project slug (optional)"},
				},
				"required": []string{"key", "value"},
			},
		},
		{
			Name:        "metadata_get",
			Description: "Get a metadata value by key",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"key":     map[string]interface{}{"type": "string", "description": "Metadata key"},
					"project": map[string]interface{}{"type": "string", "description": "Project slug (optional)"},
				},
				"required": []string{"key"},
			},
		},
		{
			Name:        "metadata_list",
			Description: "List all metadata for a project",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"project": map[string]interface{}{"type": "string", "description": "Project slug (optional)"},
				},
			},
		},
		{
			Name:        "metadata_delete",
			Description: "Delete a metadata key",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"key":     map[string]interface{}{"type": "string", "description": "Metadata key to delete"},
					"project": map[string]interface{}{"type": "string", "description": "Project slug (optional)"},
				},
				"required": []string{"key"},
			},
		},

		// Filetree tools
		{
			Name:        "filetree_annotate",
			Description: "Add or update a note on a file or directory path",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"path":    map[string]interface{}{"type": "string", "description": "File or directory path"},
					"note":    map[string]interface{}{"type": "string", "description": "Annotation note"},
					"is_dir":  map[string]interface{}{"type": "boolean", "description": "Whether path is a directory"},
					"project": map[string]interface{}{"type": "string", "description": "Project slug (optional)"},
				},
				"required": []string{"path", "note"},
			},
		},
		{
			Name:        "filetree_get",
			Description: "Get file annotations for a project or specific path",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"path":    map[string]interface{}{"type": "string", "description": "Specific path (optional, returns all if omitted)"},
					"project": map[string]interface{}{"type": "string", "description": "Project slug (optional)"},
				},
			},
		},
		{
			Name:        "filetree_delete",
			Description: "Delete a file annotation",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"path":    map[string]interface{}{"type": "string", "description": "Path to delete annotation for"},
					"project": map[string]interface{}{"type": "string", "description": "Project slug (optional)"},
				},
				"required": []string{"path"},
			},
		},

		// Guideline tools
		{
			Name:        "guideline_create",
			Description: "Create a new guideline or how-to document for knowledge transfer",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"category": map[string]interface{}{"type": "string", "description": "Category (e.g., coding_style, architecture, workflow, debugging)"},
					"title":    map[string]interface{}{"type": "string", "description": "Guideline title"},
					"content":  map[string]interface{}{"type": "string", "description": "Markdown content with instructions"},
					"tags":     map[string]interface{}{"type": "array", "items": map[string]interface{}{"type": "string"}, "description": "Tags for searchability"},
					"priority": map[string]interface{}{"type": "integer", "description": "Priority (higher = more important)"},
					"project":  map[string]interface{}{"type": "string", "description": "Project slug (optional)"},
				},
				"required": []string{"category", "title", "content"},
			},
		},
		{
			Name:        "guideline_update",
			Description: "Update a guideline's content, tags, or priority",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"id":       map[string]interface{}{"type": "integer", "description": "Guideline ID"},
					"content":  map[string]interface{}{"type": "string", "description": "New content"},
					"tags":     map[string]interface{}{"type": "array", "items": map[string]interface{}{"type": "string"}, "description": "New tags"},
					"priority": map[string]interface{}{"type": "integer", "description": "New priority"},
				},
				"required": []string{"id"},
			},
		},
		{
			Name:        "guideline_list",
			Description: "List guidelines, optionally filtered by category",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"category": map[string]interface{}{"type": "string", "description": "Filter by category"},
					"project":  map[string]interface{}{"type": "string", "description": "Project slug (optional)"},
				},
			},
		},
		{
			Name:        "guideline_search",
			Description: "Search guidelines by content, title, or tags",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"query":    map[string]interface{}{"type": "string", "description": "Search query"},
					"category": map[string]interface{}{"type": "string", "description": "Filter by category"},
					"project":  map[string]interface{}{"type": "string", "description": "Project slug (optional)"},
				},
				"required": []string{"query"},
			},
		},
		{
			Name:        "guideline_get",
			Description: "Get a specific guideline with full content",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"id": map[string]interface{}{"type": "integer", "description": "Guideline ID"},
				},
				"required": []string{"id"},
			},
		},
		{
			Name:        "guideline_delete",
			Description: "Delete a guideline",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"id": map[string]interface{}{"type": "integer", "description": "Guideline ID to delete"},
				},
				"required": []string{"id"},
			},
		},

		// Project tools
		{
			Name:        "project_create",
			Description: "Create a new project namespace",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"slug":      map[string]interface{}{"type": "string", "description": "Unique project identifier"},
					"name":      map[string]interface{}{"type": "string", "description": "Human-readable name"},
					"root_path": map[string]interface{}{"type": "string", "description": "Project root directory path"},
				},
				"required": []string{"slug"},
			},
		},
		{
			Name:        "project_list",
			Description: "List all projects",
			InputSchema: map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{},
			},
		},
		{
			Name:        "project_set_default",
			Description: "Set the default project for this session",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"slug": map[string]interface{}{"type": "string", "description": "Project slug to set as default"},
				},
				"required": []string{"slug"},
			},
		},

		// Bookmark tools
		{
			Name:        "bookmark_create",
			Description: "Create a bookmark for an external document, PDF, image, or URL with notes",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"url":             map[string]interface{}{"type": "string", "description": "File path or URL to bookmark"},
					"title":           map[string]interface{}{"type": "string", "description": "Descriptive title"},
					"excerpt":         map[string]interface{}{"type": "string", "description": "Relevant quote or key information from the document"},
					"note":            map[string]interface{}{"type": "string", "description": "Why this is useful, what to remember"},
					"doc_type":        map[string]interface{}{"type": "string", "description": "Document type (pdf, image, url, markdown, etc.)"},
					"page_or_section": map[string]interface{}{"type": "string", "description": "Page number, section name, or anchor"},
					"tags":            map[string]interface{}{"type": "array", "items": map[string]interface{}{"type": "string"}, "description": "Tags for searchability"},
					"project":         map[string]interface{}{"type": "string", "description": "Project slug (optional)"},
				},
				"required": []string{"url", "title"},
			},
		},
		{
			Name:        "bookmark_search",
			Description: "Search bookmarks by query and/or tags",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"query":    map[string]interface{}{"type": "string", "description": "Search in title, excerpt, note, or URL"},
					"tags":     map[string]interface{}{"type": "array", "items": map[string]interface{}{"type": "string"}, "description": "Filter by tags"},
					"doc_type": map[string]interface{}{"type": "string", "description": "Filter by document type"},
					"project":  map[string]interface{}{"type": "string", "description": "Project slug (optional)"},
				},
			},
		},
		{
			Name:        "bookmark_list",
			Description: "List all bookmarks for a project",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"project": map[string]interface{}{"type": "string", "description": "Project slug (optional)"},
				},
			},
		},
		{
			Name:        "bookmark_delete",
			Description: "Delete a bookmark by ID",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"id": map[string]interface{}{"type": "integer", "description": "Bookmark ID to delete"},
				},
				"required": []string{"id"},
			},
		},
	}
}
