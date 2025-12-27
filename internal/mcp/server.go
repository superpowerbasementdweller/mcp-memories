package mcp

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"log"

	"github.com/rocket/mcp-memories/internal/db"
)

// Server implements the MCP protocol over stdio
type Server struct {
	db     *db.DB
	reader *bufio.Reader
	writer io.Writer
	mu     sync.Mutex
	logger *log.Logger
}

const maxMessageBytes = 8 * 1024 * 1024

// NewServer creates a new MCP server
func NewServer(database *db.DB) *Server {
	// Setup logging to file
	homeDir, _ := os.UserHomeDir()
	logDir := filepath.Join(homeDir, ".mcp-memory")
	_ = os.MkdirAll(logDir, 0755)
	logPath := filepath.Join(logDir, "server.log")
	f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	var logger *log.Logger
	if err == nil {
		logger = log.New(f, "[MCP] ", log.LstdFlags)
	} else {
		logger = log.New(os.Stderr, "[MCP] ", log.LstdFlags)
	}

	return &Server{
		db:     database,
		reader: bufio.NewReaderSize(os.Stdin, 64*1024),
		writer: os.Stdout,
		logger: logger,
	}
}

// Request represents a JSON-RPC 2.0 request
type Request struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

// Response represents a JSON-RPC 2.0 response
type Response struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id,omitempty"`
	Result  interface{} `json:"result,omitempty"`
	Error   *Error      `json:"error,omitempty"`
}

// Error represents a JSON-RPC 2.0 error
type Error struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Standard JSON-RPC error codes
const (
	ParseError     = -32700
	InvalidRequest = -32600
	MethodNotFound = -32601
	InvalidParams  = -32602
	InternalError  = -32603
)

// Run starts the server's main loop
func (s *Server) Run() error {
	s.logger.Println("Server started")
	defer func() {
		if r := recover(); r != nil {
			s.logger.Printf("Panic recovered: %v", r)
		}
	}()

	for {
		line, err := readLineLimited(s.reader, maxMessageBytes)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			s.logger.Printf("Read error: %v", err)
			return fmt.Errorf("reading input: %w", err)
		}

		line = bytesTrimSpaceCRLF(line)
		if len(line) == 0 {
			continue
		}

		var req Request
		if err := json.Unmarshal(line, &req); err != nil {
			s.logger.Printf("Parse error: %v", err)
			s.sendError(nil, ParseError, "Parse error", err.Error())
			continue
		}

		// Handle request in a goroutine to allow concurrent processing?
		// For now, keep it synchronous to avoid race conditions with DB if not thread safe (DB is thread safe mostly but let's be safe)
		// But we should recover individually per request too
		func() {
			defer func() {
				if r := recover(); r != nil {
					s.logger.Printf("Panic handling request: %v", r)
					s.sendError(req.ID, InternalError, "Internal error", fmt.Sprintf("Panic: %v", r))
				}
			}()
			s.handleRequest(&req)
		}()
	}
}

func bytesTrimSpaceCRLF(b []byte) []byte {
	// json.Unmarshal allows whitespace, but trimming avoids surprises with CRLF and empty lines.
	s := strings.TrimSpace(string(b))
	return []byte(s)
}

func readLineLimited(r *bufio.Reader, max int) ([]byte, error) {
	var out []byte
	for {
		frag, err := r.ReadSlice('\n')
		if err == nil {
			out = append(out, frag...)
			if len(out) > max {
				return nil, fmt.Errorf("message too large: %d bytes", len(out))
			}
			return out, nil
		}
		if errors.Is(err, bufio.ErrBufferFull) {
			out = append(out, frag...)
			if len(out) > max {
				// Drain until end-of-line to re-sync.
				for {
					_, drainErr := r.ReadSlice('\n')
					if drainErr == nil {
						break
					}
					if !errors.Is(drainErr, bufio.ErrBufferFull) {
						break
					}
				}
				return nil, fmt.Errorf("message too large: exceeds %d bytes", max)
			}
			continue
		}
		if err == io.EOF {
			// Return any final partial line for processing.
			if len(frag) > 0 {
				out = append(out, frag...)
				if len(out) > max {
					return nil, fmt.Errorf("message too large: exceeds %d bytes", max)
				}
				return out, nil
			}
			return nil, io.EOF
		}
		return nil, err
	}
}

func (s *Server) handleRequest(req *Request) {
	s.logger.Printf("Handling request: %s", req.Method)
	if req.Method != "notifications/initialized" {
		if err := validateRequestID(req.ID); err != nil {
			s.sendError(nil, InvalidRequest, "Invalid request", err.Error())
			return
		}
	}

	switch req.Method {
	case "initialize":
		s.handleInitialize(req)
	case "tools/list":
		s.handleToolsList(req)
	case "tools/call":
		s.handleToolsCall(req)
	case "notifications/initialized":
		// Acknowledgment, no response needed
	default:
		s.logger.Printf("Method not found: %s", req.Method)
		s.sendError(req.ID, MethodNotFound, "Method not found", req.Method)
	}
}

func validateRequestID(id interface{}) error {
	if id == nil {
		return fmt.Errorf("missing id")
	}
	switch v := id.(type) {
	case string:
		if v == "" {
			return fmt.Errorf("id must not be empty")
		}
		return nil
	case float64:
		return nil
	case int:
		return nil
	case int64:
		return nil
	default:
		return fmt.Errorf("id must be string or number")
	}
}

func (s *Server) handleInitialize(req *Request) {
	result := map[string]interface{}{
		"protocolVersion": "2024-11-05",
		"capabilities": map[string]interface{}{
			"tools": map[string]interface{}{"listChanged": false},
		},
		"serverInfo": map[string]interface{}{
			"name":    "mcp-memories",
			"version": "1.0.0",
		},
	}
	s.sendResult(req.ID, result)
}

func (s *Server) handleToolsList(req *Request) {
	s.sendResult(req.ID, map[string]interface{}{
		"tools": GetToolDefinitions(),
	})
}

func (s *Server) handleToolsCall(req *Request) {
	var params struct {
		Name      string                 `json:"name"`
		Arguments map[string]interface{} `json:"arguments"`
	}

	if err := json.Unmarshal(req.Params, &params); err != nil {
		s.logger.Printf("Invalid params for tool call: %v", err)
		s.sendError(req.ID, InvalidParams, "Invalid params", err.Error())
		return
	}
	if params.Name == "" {
		s.sendError(req.ID, InvalidParams, "Invalid params", "tool name is required")
		return
	}

	s.logger.Printf("Calling tool: %s", params.Name)
	result, err := HandleToolCall(s.db, params.Name, params.Arguments)
	if err != nil {
		s.logger.Printf("Tool error: %v", err)
		if errors.Is(err, ErrUnknownTool) {
			s.sendError(req.ID, InvalidParams, "Unknown tool", err.Error())
			return
		}
		s.sendResult(req.ID, map[string]interface{}{
			"content": []map[string]interface{}{
				{
					"type": "text",
					"text": fmt.Sprintf("Error: %v", err),
				},
			},
			"isError": true,
		})
		return
	}

	// Format result as text content
	resultJSON, _ := json.Marshal(result) // keep compact
	s.sendResult(req.ID, map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": string(resultJSON),
			},
		},
		"isError": false,
	})
}

func (s *Server) sendResult(id interface{}, result interface{}) {
	s.send(&Response{
		JSONRPC: "2.0",
		ID:      id,
		Result:  result,
	})
}

func (s *Server) sendError(id interface{}, code int, message string, data interface{}) {
	s.send(&Response{
		JSONRPC: "2.0",
		ID:      id,
		Error: &Error{
			Code:    code,
			Message: message,
			Data:    data,
		},
	})
}

func (s *Server) send(resp *Response) {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := json.Marshal(resp)
	if err != nil {
		s.logger.Printf("Error marshaling response: %v", err)
		return
	}
	fmt.Fprintf(s.writer, "%s\n", data)
}
