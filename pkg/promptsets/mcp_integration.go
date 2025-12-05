// Package promptsets provides MCP server integration for Gemara prompts.
// Note: This requires the MCP SDK: go get github.com/modelcontextprotocol/go-sdk/mcp
//
//go:build mcp
// +build mcp

package promptsets

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// MCPPromptAdapter adapts our PromptSet system to MCP prompts
type MCPPromptAdapter struct {
	group *PromptSetGroup
}

// NewMCPPromptAdapter creates a new MCP prompt adapter
func NewMCPPromptAdapter(group *PromptSetGroup) *MCPPromptAdapter {
	return &MCPPromptAdapter{
		group: group,
	}
}

// ListMCPPrompts converts PromptSet prompts to MCP prompts
func (a *MCPPromptAdapter) ListMCPPrompts() ([]mcp.Prompt, error) {
	var mcpPrompts []mcp.Prompt

	for _, ps := range a.group.ListPromptSets() {
		for _, prompt := range ps.ListPrompts() {
			// Build argument schema from prompt variables
			properties := make(map[string]mcp.JSONSchema)
			required := make([]string, 0)

			for varName, varDesc := range prompt.Variables {
				properties[varName] = mcp.JSONSchema{
					Type:        "string",
					Description: varDesc,
				}
				required = append(required, varName)
			}

			mcpPrompt := mcp.Prompt{
				Name:        prompt.Name,
				Description: prompt.Description,
				Arguments: []mcp.PromptArgument{
					{
						Name:        "variables",
						Description: "Variables to substitute into the prompt template",
						Required:    false,
					},
				},
			}

			// Add variable arguments if they exist
			if len(properties) > 0 {
				mcpPrompt.Arguments = make([]mcp.PromptArgument, 0, len(properties))
				for varName, schema := range properties {
					mcpPrompt.Arguments = append(mcpPrompt.Arguments, mcp.PromptArgument{
						Name:        varName,
						Description: schema.Description,
						Required:    true,
					})
				}
			}

			mcpPrompts = append(mcpPrompts, mcpPrompt)
		}
	}

	return mcpPrompts, nil
}

// GetMCPPrompt generates an MCP prompt response from a prompt set
func (a *MCPPromptAdapter) GetMCPPrompt(ctx context.Context, promptName string, promptSetName string, variables map[string]interface{}) (mcp.PromptResponse, error) {
	ps, err := a.group.GetPromptSet(promptSetName)
	if err != nil {
		return mcp.PromptResponse{}, fmt.Errorf("prompt set %q not found: %w", promptSetName, err)
	}

	req := PromptRequest{
		PromptName: promptName,
		Variables:  variables,
	}

	resp, err := ps.GeneratePrompt(ctx, req)
	if err != nil {
		return mcp.PromptResponse{}, fmt.Errorf("failed to generate prompt: %w", err)
	}

	return mcp.PromptResponse{
		Content: []mcp.PromptMessage{
			{
				Role:    "user",
				Content: resp.Content,
			},
		},
	}, nil
}

// RegisterPromptsWithMCPServer registers prompts with an MCP server
// This follows the github-mcp-server pattern using ServerPrompt
func RegisterPromptsWithMCPServer(server *mcp.Server, group *PromptSetGroup) error {
	adapter := NewMCPPromptAdapter(group)

	// Get all prompts
	mcpPrompts, err := adapter.ListMCPPrompts()
	if err != nil {
		return fmt.Errorf("failed to list prompts: %w", err)
	}

	// Register each prompt with the server
	for _, mcpPrompt := range mcpPrompts {
		prompt := mcpPrompt // Capture for closure

		handler := func(ctx context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
			// Extract variables from request
			variables := make(map[string]interface{})

			// Parse arguments if provided
			if req.Params.Arguments != nil {
				var args map[string]interface{}
				if err := json.Unmarshal(req.Params.Arguments, &args); err == nil {
					// Extract variables from arguments
					if vars, ok := args["variables"].(map[string]interface{}); ok {
						variables = vars
					} else {
						// If no "variables" key, use all arguments as variables
						variables = args
					}
				}
			}

			// Find the prompt set that contains this prompt
			var promptSetName string
			for _, ps := range group.ListPromptSets() {
				for _, p := range ps.ListPrompts() {
					if p.Name == prompt.Name {
						promptSetName = ps.Name()
						break
					}
				}
				if promptSetName != "" {
					break
				}
			}

			if promptSetName == "" {
				return nil, fmt.Errorf("prompt set not found for prompt %q", prompt.Name)
			}

			// Generate the prompt
			resp, err := adapter.GetMCPPrompt(ctx, prompt.Name, promptSetName, variables)
			if err != nil {
				return nil, err
			}

			return &mcp.GetPromptResult{
				Content: resp.Content,
			}, nil
		}

		server.AddPrompt(&prompt, handler)
	}

	return nil
}

// RegisterGemaraPromptsWithMCPServer registers Gemara-specific prompts with an MCP server
// This is a convenience function that registers the Gemara prompt set
func RegisterGemaraPromptsWithMCPServer(server *mcp.Server) error {
	group := GetGemaraPromptSets()
	return RegisterPromptsWithMCPServer(server, group)
}
