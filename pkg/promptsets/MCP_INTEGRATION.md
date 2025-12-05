# MCP Server Integration Guide

This guide explains how to integrate Gemara prompts with an MCP (Model Context Protocol) server, following the pattern used by `github-mcp-server`.

## Overview

The integration allows you to:
- Register Gemara prompts as MCP prompts
- Expose prompts that support variable scope changes
- Enable Layer 1 → Layer 3 policy generation with configurable scope

## Key Prompt: Layer 1 → Layer 3 Policy Generation

The `gemara_layer1_to_layer3_policy` prompt is designed to generate Layer 3 (Policy) documents from Layer 1 (Guidance) with a configurable scope variable.

### Variables

- **`scope`**: The domain/area the policy covers (e.g., "Cloud Infrastructure Security", "API Security")
- **`layer1_guidance_source`**: Source of Layer 1 guidance (e.g., "NIST Cybersecurity Framework")
- **`layer1_guidance_content`**: The specific guidance content
- **`organization_context`**: Organization-specific context
- **`risk_appetite`**: Organization's risk tolerance
- **`additional_requirements`**: Any additional constraints

## Setup

### 1. Install Dependencies

```bash
go get github.com/modelcontextprotocol/go-sdk/mcp
```

### 2. Register Prompts with MCP Server

```go
package main

import (
    "log"
    "github.com/modelcontextprotocol/go-sdk/mcp"
    "github.com/your-org/gemara-mcp-server/pkg/promptsets"
)

func main() {
    // Create MCP server
    server, err := mcp.NewServer(
        mcp.Options{
            Name:    "gemara-mcp-server",
            Version: "0.1.0",
        },
    )
    if err != nil {
        log.Fatalf("Failed to create server: %v", err)
    }

    // Register all Gemara prompts
    if err := promptsets.RegisterGemaraPromptsWithMCPServer(server); err != nil {
        log.Fatalf("Failed to register prompts: %v", err)
    }

    // Your server setup continues...
}
```

### 3. Selective Prompt Registration

If you only want to expose specific prompts:

```go
// Register only the Layer 1 → Layer 3 policy prompt
selectedPrompts := []string{
    "gemara_layer1_to_layer3_policy",
    "gemara_layer_context",
    "gemara_layer_validator",
}

// Use ExampleSelectivePromptRegistration from mcp_server_example.go
```

## Usage Examples

### Example 1: Generate Policy with Different Scopes

```go
ctx := context.Background()
group := promptsets.GetGemaraPromptSets()
gemaraSet, _ := group.GetPromptSet("gemara")

// Change scope to generate different policies
req := promptsets.PromptRequest{
    PromptName: "gemara_layer1_to_layer3_policy",
    Variables: map[string]interface{}{
        "scope": "API Security", // Change this to change scope
        "layer1_guidance_source": "NIST Cybersecurity Framework",
        "layer1_guidance_content": "NIST guidance content...",
        "organization_context": "Tech company, 1000 employees",
        "risk_appetite": "Moderate",
        "additional_requirements": "Must support microservices",
    },
}

resp, err := gemaraSet.GeneratePrompt(ctx, req)
// resp.Content contains the generated Layer 3 policy prompt
```

### Example 2: Using via MCP Protocol

When a client requests the prompt via MCP:

```json
{
  "method": "prompts/get",
  "params": {
    "name": "gemara_layer1_to_layer3_policy",
    "arguments": {
      "scope": "Cloud Infrastructure Security",
      "layer1_guidance_source": "ISO 27001",
      "layer1_guidance_content": "ISO 27001 requirements...",
      "organization_context": "Financial services, 5000+ employees",
      "risk_appetite": "Low risk tolerance",
      "additional_requirements": "Multi-cloud strategy"
    }
  }
}
```

The server will:
1. Validate the prompt exists
2. Generate the prompt with variable substitution
3. Return the complete prompt content

## Available Prompts

When you register all Gemara prompts, the following are available:

1. **`gemara_layer1_to_layer3_policy`** - Generate Layer 3 policy from Layer 1 guidance (with scope variable)
2. **`gemara_architecture_explainer`** - Explains the complete Gemara architecture
3. **`gemara_layer_context`** - Provides context for a specific layer
4. **`gemara_layer_1_guidance`** - System prompt for Layer 1 work
5. **`gemara_layer_2_controls`** - System prompt for Layer 2 work
6. **`gemara_layer_3_policy`** - System prompt for Layer 3 work
7. **`gemara_layer_4_evaluation`** - System prompt for Layer 4 work
8. **`gemara_layer_5_enforcement`** - System prompt for Layer 5 work
9. **`gemara_layer_6_audit`** - System prompt for Layer 6 work
10. **`gemara_layer_validator`** - Validates layer references

## Benefits of This Approach

1. **Modularity**: Prompts are organized into logical sets
2. **Flexibility**: Users can change scope via variables
3. **Deterministic**: Layer references are validated
4. **Consistent**: Follows github-mcp-server patterns
5. **Extensible**: Easy to add new prompts

## Pattern Comparison with github-mcp-server

The integration follows the same pattern as `github-mcp-server`:

- **Toolsets → PromptSets**: Organized collections of related functionality
- **ServerTool → ServerPrompt**: Wrapper for MCP registration
- **ToolsetGroup → PromptSetGroup**: Management of multiple sets
- **Enable/Disable**: Can selectively register prompts

This ensures consistency with established MCP server patterns.

## Next Steps

1. Implement your MCP server using the examples
2. Test with different scope values
3. Customize prompts for your organization's needs
4. Add additional prompts as needed

See `mcp_server_example.go` for complete working examples.

