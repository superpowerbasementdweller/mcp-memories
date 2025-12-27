package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/rocket/mcp-memories/internal/db"
	"github.com/rocket/mcp-memories/internal/mcp"
)

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

	// Create and run MCP server
	server := mcp.NewServer(database)
	if err := server.Run(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
