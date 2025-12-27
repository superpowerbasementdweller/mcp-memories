package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/rocket/mcp-memories/internal/db"
	"github.com/rocket/mcp-memories/internal/mcp"
)

//go:embed templates/*
var templates embed.FS

//go:embed static/*
var static embed.FS

type DashboardData struct {
	Tools      []mcp.ToolDefinition
	Categories map[string][]mcp.ToolDefinition
	Stats      Stats
}

type Stats struct {
	Projects   int
	Memories   int
	Tasks      int
	Bookmarks  int
	Guidelines int
}

func main() {
	// Determine database path
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Failed to get home directory: %v", err)
	}

	dbPath := filepath.Join(homeDir, ".mcp-memory", "memories.db")

	// Open database
	database, err := db.Open(dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer database.Close()

	// Parse templates
	tmpl, err := template.ParseFS(templates, "templates/*.html")
	if err != nil {
		log.Fatalf("Failed to parse templates: %v", err)
	}

	// Static files
	http.Handle("/static/", http.FileServer(http.FS(static)))

	// Dashboard
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tools := mcp.GetToolDefinitions()
		categories := categorizeTools(tools)
		stats := getStats(database)

		data := DashboardData{
			Tools:      tools,
			Categories: categories,
			Stats:      stats,
		}

		if err := tmpl.ExecuteTemplate(w, "index.html", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	// API: Get tools as JSON
	http.HandleFunc("/api/tools", func(w http.ResponseWriter, r *http.Request) {
		tools := mcp.GetToolDefinitions()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tools)
	})

	// API: Get stats
	http.HandleFunc("/api/stats", func(w http.ResponseWriter, r *http.Request) {
		stats := getStats(database)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(stats)
	})

	// API: Restart MCP server
	http.HandleFunc("/api/restart", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var cmd *exec.Cmd
		if runtime.GOOS == "windows" {
			cmd = exec.Command("taskkill", "/IM", "mcp-memories.exe", "/F")
		} else {
			cmd = exec.Command("pkill", "-f", "mcp-memories")
		}

		output, err := cmd.CombinedOutput()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": err == nil || string(output) != "",
			"message": "MCP server killed. It will restart automatically when the IDE reconnects.",
			"output":  string(output),
		})
	})

	// API: CRUD endpoints
	http.HandleFunc("/api/projects", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		projects, _ := database.ListProjects()
		json.NewEncoder(w).Encode(projects)
	})

	http.HandleFunc("/api/memories", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method == "POST" {
			var req struct {
				Content  string   `json:"content"`
				Keywords []string `json:"keywords"`
			}
			json.NewDecoder(r.Body).Decode(&req)
			mem, err := database.CreateMemory(nil, req.Content, req.Keywords)
			if err != nil {
				json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
				return
			}
			json.NewEncoder(w).Encode(mem)
			return
		}
		if r.Method == "DELETE" {
			var req struct {
				ID int64 `json:"id"`
			}
			json.NewDecoder(r.Body).Decode(&req)
			database.DeleteMemory(req.ID)
			json.NewEncoder(w).Encode(map[string]bool{"deleted": true})
			return
		}
		memories, _ := database.SearchMemories(nil, "", nil, 100)
		json.NewEncoder(w).Encode(memories)
	})

	http.HandleFunc("/api/tasks", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method == "POST" {
			var req struct {
				Title       string `json:"title"`
				Description string `json:"description"`
				Priority    int    `json:"priority"`
			}
			json.NewDecoder(r.Body).Decode(&req)
			task, err := database.CreateTask(nil, nil, req.Title, req.Description, req.Priority)
			if err != nil {
				json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
				return
			}
			json.NewEncoder(w).Encode(task)
			return
		}
		if r.Method == "DELETE" {
			var req struct {
				ID int64 `json:"id"`
			}
			json.NewDecoder(r.Body).Decode(&req)
			database.DeleteTask(req.ID)
			json.NewEncoder(w).Encode(map[string]bool{"deleted": true})
			return
		}
		tasks, _ := database.ListTasks(nil, nil, nil)
		json.NewEncoder(w).Encode(tasks)
	})

	http.HandleFunc("/api/guidelines", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method == "POST" {
			var req struct {
				Category string   `json:"category"`
				Title    string   `json:"title"`
				Content  string   `json:"content"`
				Tags     []string `json:"tags"`
				Priority int      `json:"priority"`
			}
			json.NewDecoder(r.Body).Decode(&req)
			g, err := database.CreateGuideline(nil, req.Category, req.Title, req.Content, req.Tags, req.Priority)
			if err != nil {
				json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
				return
			}
			json.NewEncoder(w).Encode(g)
			return
		}
		if r.Method == "DELETE" {
			var req struct {
				ID int64 `json:"id"`
			}
			json.NewDecoder(r.Body).Decode(&req)
			database.DeleteGuideline(req.ID)
			json.NewEncoder(w).Encode(map[string]bool{"deleted": true})
			return
		}
		guidelines, _ := database.ListGuidelines(nil, nil)
		json.NewEncoder(w).Encode(guidelines)
	})

	http.HandleFunc("/api/bookmarks", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method == "POST" {
			var req struct {
				URL           string   `json:"url"`
				Title         string   `json:"title"`
				Excerpt       string   `json:"excerpt"`
				Note          string   `json:"note"`
				DocType       string   `json:"doc_type"`
				PageOrSection string   `json:"page_or_section"`
				Tags          []string `json:"tags"`
			}
			json.NewDecoder(r.Body).Decode(&req)
			b, err := database.CreateBookmark(nil, req.URL, req.Title, req.Excerpt, req.Note, req.DocType, req.PageOrSection, req.Tags)
			if err != nil {
				json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
				return
			}
			json.NewEncoder(w).Encode(b)
			return
		}
		if r.Method == "DELETE" {
			var req struct {
				ID int64 `json:"id"`
			}
			json.NewDecoder(r.Body).Decode(&req)
			database.DeleteBookmark(req.ID)
			json.NewEncoder(w).Encode(map[string]bool{"deleted": true})
			return
		}
		bookmarks, _ := database.ListBookmarks(nil)
		json.NewEncoder(w).Encode(bookmarks)
	})

	port := "8765"
	fmt.Printf("🧠 MCP Memories Dashboard running at http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func categorizeTools(tools []mcp.ToolDefinition) map[string][]mcp.ToolDefinition {
	categories := make(map[string][]mcp.ToolDefinition)

	prefixMap := map[string]string{
		"memory_":    "Memory",
		"task_":      "Task",
		"metadata_":  "Metadata",
		"filetree_":  "Filetree",
		"guideline_": "Guideline",
		"project_":   "Project",
		"bookmark_":  "Bookmark",
	}

	for _, tool := range tools {
		categorized := false
		for prefix, category := range prefixMap {
			if len(tool.Name) >= len(prefix) && tool.Name[:len(prefix)] == prefix {
				categories[category] = append(categories[category], tool)
				categorized = true
				break
			}
		}
		if !categorized {
			categories["Other"] = append(categories["Other"], tool)
		}
	}

	return categories
}

func getStats(database *db.DB) Stats {
	var stats Stats

	database.QueryRow("SELECT COUNT(*) FROM projects").Scan(&stats.Projects)
	database.QueryRow("SELECT COUNT(*) FROM memories").Scan(&stats.Memories)
	database.QueryRow("SELECT COUNT(*) FROM tasks").Scan(&stats.Tasks)
	database.QueryRow("SELECT COUNT(*) FROM bookmarks").Scan(&stats.Bookmarks)
	database.QueryRow("SELECT COUNT(*) FROM guidelines").Scan(&stats.Guidelines)

	return stats
}
