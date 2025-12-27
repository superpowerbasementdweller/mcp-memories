# MCP Memories

A local, database-backed MCP (Model Context Protocol) server for persistent context across AI coding sessions. Stores memories, tasks, metadata, file annotations, guidelines, and bookmarks in SQLite.

## Features

- **Memories**: Store and search freeform notes with keyword tagging
- **Tasks**: Hierarchical task management with subtasks and status tracking
- **Metadata**: Key-value store for project configuration and context
- **Filetree**: Annotate files and directories with notes
- **Guidelines**: Document patterns, how-tos, and conventions for knowledge transfer
- **Bookmarks**: Save references to external docs, PDFs, images, and URLs
- **Projects**: Namespace data per-project while keeping everything in one portable database
- **Dashboard**: Web UI to view data and restart the MCP server

## Installation

### Build from source
```powershell
cd c:\Users\rocket\Documents\code\go\mcp_memories
go build -o mcp-memories.exe ./cmd/mcp-memories
go build -o mcp-dashboard.exe ./cmd/dashboard
```

### Install
Copy the binaries to `~/.mcp-memory/`:
```powershell
Copy-Item .\mcp-memories.exe "$env:USERPROFILE\.mcp-memory\mcp-memories.exe"
Copy-Item .\mcp-dashboard.exe "$env:USERPROFILE\.mcp-memory\mcp-dashboard.exe"
```

## Configuration

Add to your MCP client configuration:

```json
{
  "mcpServers": {
    "mcp-memories": {
      "command": "C:\\Users\\rocket\\.mcp-memory\\mcp-memories.exe"
    }
  }
}
```

## Web Dashboard

Run the dashboard to view stored data and manage the MCP server:

```powershell
& "$env:USERPROFILE\.mcp-memory\mcp-dashboard.exe"
```

Then open **http://localhost:8765** in your browser.

**Dashboard features:**
- 📊 Stats overview (projects, memories, tasks, guidelines, bookmarks)
- 🔧 All 27 tools organized by category
- 📋 Data browser with tabs to view stored data
- 🔄 Restart button to kill the MCP server

## Available Tools (27 total)

### Memory Tools (3)
| Tool | Description |
|------|-------------|
| `memory_store` | Store a new memory with optional keywords |
| `memory_search` | Search memories by content and/or keywords |
| `memory_delete` | Delete a memory by ID |

### Task Tools (4)
| Tool | Description |
|------|-------------|
| `task_create` | Create a task with optional parent for subtasks |
| `task_update` | Update status, title, description, or priority |
| `task_list` | List tasks with filters (project, status, parent) |
| `task_delete` | Delete a task and its subtasks |

### Metadata Tools (4)
| Tool | Description |
|------|-------------|
| `metadata_set` | Set a key-value pair |
| `metadata_get` | Get a value by key |
| `metadata_list` | List all metadata for a project |
| `metadata_delete` | Delete a key |

### Filetree Tools (3)
| Tool | Description |
|------|-------------|
| `filetree_annotate` | Add/update a note on a file or directory |
| `filetree_get` | Get annotations (all or for specific path) |
| `filetree_delete` | Delete an annotation |

### Guideline Tools (6)
| Tool | Description |
|------|-------------|
| `guideline_create` | Create a guideline with category, title, content |
| `guideline_update` | Update content, tags, or priority |
| `guideline_list` | List guidelines by category |
| `guideline_search` | Search by content, title, or tags |
| `guideline_get` | Get full guideline content |
| `guideline_delete` | Delete a guideline |

### Bookmark Tools (4)
| Tool | Description |
|------|-------------|
| `bookmark_create` | Create a bookmark for docs, PDFs, images, URLs |
| `bookmark_search` | Search by query and/or tags |
| `bookmark_list` | List all bookmarks for a project |
| `bookmark_delete` | Delete a bookmark by ID |

### Project Tools (3)
| Tool | Description |
|------|-------------|
| `project_create` | Create a new project namespace |
| `project_list` | List all projects |
| `project_set_default` | Set the default project for operations |

## Database Location

All data is stored in a single SQLite database at:
```
~/.mcp-memory/memories.db
```

## Example Usage

### Store a memory
```json
{
  "name": "memory_store",
  "arguments": {
    "content": "The Buh project uses SDL3 with GPU-accelerated vector graphics",
    "keywords": ["buh", "sdl3", "graphics", "architecture"],
    "project": "odin-buh"
  }
}
```

### Create a bookmark for a PDF
```json
{
  "name": "bookmark_create",
  "arguments": {
    "url": "c:\\docs\\sdl3_gpu_api.pdf",
    "title": "SDL3 GPU API Reference",
    "excerpt": "SDL_GPUGraphicsPipeline requires vertex and fragment shaders",
    "note": "Key reference for understanding the rendering pipeline",
    "doc_type": "pdf",
    "page_or_section": "Page 42",
    "tags": ["sdl3", "gpu", "reference"],
    "project": "odin-buh"
  }
}
```

### Create a guideline for knowledge transfer
```json
{
  "name": "guideline_create",
  "arguments": {
    "category": "coding_style",
    "title": "Odin procedure naming",
    "content": "## Odin Naming Conventions\n\n1. Use `snake_case` for procedures\n2. Use `PascalCase` for types\n3. Prefix private procedures with `_`",
    "tags": ["odin", "style", "naming"],
    "priority": 10,
    "project": "odin-buh"
  }
}
```

## License

MIT
