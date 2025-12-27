package schema

// Schema contains the SQL to initialize the database
const Schema = `
-- Project namespacing (default = "global")
CREATE TABLE IF NOT EXISTS projects (
    id INTEGER PRIMARY KEY,
    slug TEXT UNIQUE NOT NULL,
    name TEXT,
    root_path TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Freeform memories with keyword tagging
CREATE TABLE IF NOT EXISTS memories (
    id INTEGER PRIMARY KEY,
    project_id INTEGER REFERENCES projects(id),
    content TEXT NOT NULL,
    keywords TEXT, -- JSON array of strings
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Key-value metadata per project
CREATE TABLE IF NOT EXISTS metadata (
    id INTEGER PRIMARY KEY,
    project_id INTEGER REFERENCES projects(id),
    key TEXT NOT NULL,
    value TEXT,
    UNIQUE(project_id, key)
);

-- Hierarchical tasks with status
CREATE TABLE IF NOT EXISTS tasks (
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
CREATE TABLE IF NOT EXISTS filetree (
    id INTEGER PRIMARY KEY,
    project_id INTEGER REFERENCES projects(id),
    path TEXT NOT NULL,
    note TEXT,
    is_dir BOOLEAN DEFAULT FALSE,
    UNIQUE(project_id, path)
);

-- Guidelines/How-tos for knowledge transfer between models
CREATE TABLE IF NOT EXISTS guidelines (
    id INTEGER PRIMARY KEY,
    project_id INTEGER REFERENCES projects(id),
    category TEXT NOT NULL,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    tags TEXT, -- JSON array for searchability
    priority INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(project_id, category, title)
);

-- Bookmarks for external documents, PDFs, images, URLs
CREATE TABLE IF NOT EXISTS bookmarks (
    id INTEGER PRIMARY KEY,
    project_id INTEGER REFERENCES projects(id),
    url TEXT NOT NULL,           -- file path or URL
    title TEXT NOT NULL,
    excerpt TEXT,                -- relevant quote or description
    note TEXT,                   -- why this is useful
    doc_type TEXT,               -- pdf, image, url, markdown, etc.
    page_or_section TEXT,        -- page number, section, or anchor
    tags TEXT,                   -- JSON array
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Insert default global project
INSERT OR IGNORE INTO projects (id, slug, name) VALUES (1, 'global', 'Global');

-- Indexes
CREATE INDEX IF NOT EXISTS idx_memories_project ON memories(project_id);
CREATE INDEX IF NOT EXISTS idx_memories_keywords ON memories(keywords);
CREATE INDEX IF NOT EXISTS idx_tasks_project_status ON tasks(project_id, status);
CREATE INDEX IF NOT EXISTS idx_tasks_parent ON tasks(parent_id);
CREATE INDEX IF NOT EXISTS idx_filetree_project ON filetree(project_id);
CREATE INDEX IF NOT EXISTS idx_guidelines_project_category ON guidelines(project_id, category);
CREATE INDEX IF NOT EXISTS idx_metadata_project ON metadata(project_id);
CREATE INDEX IF NOT EXISTS idx_bookmarks_project ON bookmarks(project_id);
CREATE INDEX IF NOT EXISTS idx_bookmarks_tags ON bookmarks(tags);
`
