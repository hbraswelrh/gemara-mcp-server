package tools

import "github.com/mark3labs/mcp-go/server"

// ToolKit represents a collection of MCP functionality.
type ToolKit interface {
	Name() string
	Description() string
	Register(mcpServer *server.MCPServer)
}
