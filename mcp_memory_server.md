# MCP Memory Server

A local, database-backed MCP server for persistent context across AI coding sessions. Stores memories, tasks, metadata, and file annotations in SQLite, scoped per-project but portable via `~/.mcp-memory/`.

## Database Schema

```sql
-- Project namespacing (default = "global")
CREATE TABLE projects (
    id INTEGER PRIMARY KEY,
    slug TEXT UNIQUE NOT NULL,
    name TEXT,
    root_path TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Freeform memories with keyword tagging
CREATE TABLE memories (
    id INTEGER PRIMARY KEY,
    project_id INTEGER REFERENCES projects(id),
    content TEXT NOT NULL,
    keywords TEXT, -- JSON array of strings
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Key-value metadata per project
CREATE TABLE metadata (
    id INTEGER PRIMARY KEY,
    project_id INTEGER REFERENCES projects(id),
    key TEXT NOT NULL,
    value TEXT,
    UNIQUE(project_id, key)
);

-- Hierarchical tasks with status
CREATE TABLE tasks (
    id INTEGER PRIMARY KEY,
    project_id INTEGER REFERENCES projects(id),
    parent_id INTEGER REFERENCES tasks(id),
    title TEXT NOT NULL,
    description TEXT,
    status TEXT DEFAULT 'todo', -- todo, in_progress, done, blocked
    priority INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- File/directory annotations
CREATE TABLE filetree (
    id INTEGER PRIMARY KEY,
    project_id INTEGER REFERENCES projects(id),
    path TEXT NOT NULL,
    note TEXT,
    is_dir BOOLEAN DEFAULT FALSE,
    UNIQUE(project_id, path)
);

-- Guidelines/How-tos for knowledge transfer between models
CREATE TABLE guidelines (
    id INTEGER PRIMARY KEY,
    project_id INTEGER REFERENCES projects(id),
    category TEXT NOT NULL,    -- e.g., "coding_style", "architecture", "workflow", "debugging"
    title TEXT NOT NULL,
    content TEXT NOT NULL,     -- Markdown content with step-by-step instructions
    tags TEXT,                 -- JSON array for searchability
    priority INTEGER DEFAULT 0, -- Higher = more important to follow
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(project_id, category, title)
);

CREATE INDEX idx_memories_project ON memories(project_id);
CREATE INDEX idx_guidelines_project_category ON guidelines(project_id, category);
CREATE INDEX idx_memories_keywords ON memories(keywords);
CREATE INDEX idx_tasks_project_status ON tasks(project_id, status);
CREATE INDEX idx_tasks_parent ON tasks(parent_id);
CREATE INDEX idx_filetree_project ON filetree(project_id);
```

---

## MCP Tools

### Memory Tools
| Tool | Parameters | Description |
|------|------------|-------------|
| `memory_store` | `content`, `keywords[]`, `project?` | Store a new memory |
| `memory_search` | `query?`, `keywords[]?`, `project?`, `limit?` | Search memories by content/keywords |
| `memory_delete` | `id` | Delete a memory |

### Task Tools
| Tool | Parameters | Description |
|------|------------|-------------|
| `task_create` | `title`, `description?`, `parent_id?`, `priority?`, `project?` | Create a task |
| `task_update` | `id`, `status?`, `title?`, `description?`, `priority?` | Update a task |
| `task_list` | `project?`, `status?`, `parent_id?` | List tasks with filters |
| `task_delete` | `id` | Delete a task and its subtasks |

### Metadata Tools
| Tool | Parameters | Description |
|------|------------|-------------|
| `metadata_set` | `key`, `value`, `project?` | Set a metadata key |
| `metadata_get` | `key`, `project?` | Get a metadata value |
| `metadata_list` | `project?` | List all metadata for a project |
| `metadata_delete` | `key`, `project?` | Delete a metadata key |

### Filetree Tools
| Tool | Parameters | Description |
|------|------------|-------------|
| `filetree_annotate` | `path`, `note`, `project?` | Add/update a note on a path |
| `filetree_get` | `path?`, `project?` | Get annotations (all or for path) |
| `filetree_delete` | `path`, `project?` | Delete an annotation |

### Guidelines Tools
| Tool | Parameters | Description |
|------|------------|-------------|
| `guideline_create` | `category`, `title`, `content`, `tags[]?`, `priority?`, `project?` | Create a guideline |
| `guideline_update` | `id`, `content?`, `tags[]?`, `priority?` | Update a guideline |
| `guideline_list` | `category?`, `project?` | List guidelines by category |
| `guideline_search` | `query`, `category?`, `project?` | Search guidelines by content/tags |
| `guideline_get` | `id` | Get a specific guideline with full content |
| `guideline_delete` | `id` | Delete a guideline |

### Project Tools
| Tool | Parameters | Description |
|------|------------|-------------|
| `project_create` | `slug`, `name?`, `root_path?` | Create a project |
| `project_list` | | List all projects |
| `project_set_default` | `slug` | Set the default project for this session |

---

## Proposed Changes

### [NEW] Project Structure

```
c:\Users\rocket\Documents\code\go\mcp_memories\
├── go.mod
├── go.sum
├── cmd/
│   └── mcp-memories/
│       └── main.go           # Entry point, stdio MCP server
├── internal/
│   ├── db/
│   │   ├── db.go             # SQLite connection, migrations
│   │   ├── memories.go       # Memory CRUD
│   │   ├── tasks.go          # Task CRUD
│   │   ├── metadata.go       # Metadata CRUD
│   │   ├── filetree.go       # Filetree CRUD
│   │   ├── guidelines.go     # Guidelines CRUD
│   │   └── projects.go       # Project CRUD
│   ├── mcp/
│   │   ├── server.go         # MCP JSON-RPC handler
│   │   ├── tools.go          # Tool definitions
│   │   └── handlers.go       # Tool implementations
│   └── schema/
│       └── migrations.go     # SQL schema embedded
└── README.md
```

---

### [NEW] [main.go](file:///c:/Users/rocket/Documents/code/go/mcp_memories/cmd/mcp-memories/main.go)
Entry point that:
- Initializes SQLite database at `~/.mcp-memory/memories.db`
- Starts MCP server on stdio (stdin/stdout JSON-RPC)
- Handles graceful shutdown

---

### [NEW] [db.go](file:///c:/Users/rocket/Documents/code/go/mcp_memories/internal/db/db.go)
Database layer:
- Opens SQLite with WAL mode for performance
- Runs migrations on startup
- Provides `*sql.DB` wrapper with helper methods

---

### [NEW] [server.go](file:///c:/Users/rocket/Documents/code/go/mcp_memories/internal/mcp/server.go)
MCP protocol implementation:
- JSON-RPC 2.0 over stdio
- Handles `initialize`, `tools/list`, `tools/call`
- Routes tool calls to handlers

---

## Dependencies

| Package | Purpose |
|---------|---------|
| `modernc.org/sqlite` | Pure Go SQLite driver (no CGO) |
| Standard library only otherwise | Keep it simple |

---

## Verification Plan

### Manual Testing
1. Build the server: `go build -o mcp-memories.exe ./cmd/mcp-memories`
2. Test via stdin/stdout with JSON-RPC requests
3. Verify database file created at `~/.mcp-memory/memories.db`

### Integration Testing
1. Configure in MCP client settings
2. Verify tools appear in tool list
3. Test each tool category (memory, task, metadata, filetree)

---

## Installation

After building, the binary can be placed anywhere. MCP client configuration:

```json
{
  "mcp-memories": {
    "command": "C:\\Users\\rocket\\.mcp-memory\\mcp-memories.exe"
  }
}
```
