//go:build mcp
// +build mcp

package promptsets

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// ExampleMCPServerSetup demonstrates how to set up an MCP server with Gemara prompts
// This follows the github-mcp-server pattern
func ExampleMCPServerSetup() {
	// Create MCP server
	server := mcp.NewServer(
		mcp.Options{
			Name:    "gemara-mcp-server",
			Version: "0.1.0",
		},
	)

	// Register Gemara prompts with the server
	// This automatically registers all Gemara prompts including:
	// - gemara_layer1_to_layer3_policy (the key prompt for Layer 1 → Layer 3 generation)
	// - gemara_layer_1_guidance
	// - gemara_layer_2_controls
	// - gemara_layer_3_policy
	// - gemara_layer_context
	// - gemara_architecture_explainer
	// - etc.
	if err := RegisterGemaraPromptsWithMCPServer(server); err != nil {
		log.Fatalf("Failed to register Gemara prompts: %v", err)
	}

	log.Println("Gemara MCP server initialized with prompts")
	log.Println("Registered prompts:")

	// List registered prompts
	group := GetGemaraPromptSets()
	gemaraSet, _ := group.GetPromptSet("gemara")
	for _, prompt := range gemaraSet.ListPrompts() {
		log.Printf("  - %s: %s", prompt.Name, prompt.Description)
	}
}

// ExampleLayer1ToLayer3Usage demonstrates using the Layer 1 → Layer 3 policy prompt
func ExampleLayer1ToLayer3Usage() {
	ctx := context.Background()
	group := GetGemaraPromptSets()
	gemaraSet, _ := group.GetPromptSet("gemara")

	// Example: Generate Layer 3 policy from Layer 1 guidance
	// The scope variable allows users to change what domain/area the policy covers
	req := PromptRequest{
		PromptName: "gemara_layer1_to_layer3_policy",
		Variables: map[string]interface{}{
			"scope":                  "Cloud Infrastructure Security", // User can change this scope
			"layer1_guidance_source": "NIST Cybersecurity Framework",
			"layer1_guidance_content": `
The NIST Cybersecurity Framework provides guidance on:
- Identify: Develop organizational understanding to manage cybersecurity risk
- Protect: Develop and implement safeguards
- Detect: Develop and implement activities to identify cybersecurity events
- Respond: Develop and implement activities to take action
- Recover: Develop and implement activities to maintain resilience
			`,
			"organization_context":    "Financial services organization, 5000+ employees, subject to SOX and PCI DSS",
			"risk_appetite":           "Low risk tolerance, compliance-focused",
			"additional_requirements": "Must align with existing cloud security standards and support multi-cloud strategy",
		},
	}

	resp, err := gemaraSet.GeneratePrompt(ctx, req)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Println("Generated Layer 3 Policy Prompt:")
	fmt.Println("=" + string(make([]byte, 80)) + "=")
	fmt.Println(resp.Content)
	fmt.Println("=" + string(make([]byte, 80)) + "=")
}

// ExampleMCPServerHandler shows how to handle MCP prompt requests
// This is what you would use in your actual MCP server implementation
func ExampleMCPServerHandler(ctx context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	group := GetGemaraPromptSets()
	adapter := NewMCPPromptAdapter(group)

	// Extract variables from the request
	variables := make(map[string]interface{})

	if req.Params.Arguments != nil {
		var args map[string]interface{}
		if err := json.Unmarshal(req.Params.Arguments, &args); err == nil {
			// Support both direct variables and nested "variables" key
			if vars, ok := args["variables"].(map[string]interface{}); ok {
				variables = vars
			} else {
				variables = args
			}
		}
	}

	// Find which prompt set contains this prompt
	var promptSetName string
	for _, ps := range group.ListPromptSets() {
		for _, p := range ps.ListPrompts() {
			if p.Name == req.Params.Name {
				promptSetName = ps.Name()
				break
			}
		}
		if promptSetName != "" {
			break
		}
	}

	if promptSetName == "" {
		return nil, fmt.Errorf("prompt %q not found", req.Params.Name)
	}

	// Generate the prompt
	resp, err := adapter.GetMCPPrompt(ctx, req.Params.Name, promptSetName, variables)
	if err != nil {
		return nil, err
	}

	return &mcp.GetPromptResult{
		Content: resp.Content,
	}, nil
}

// ExampleSelectivePromptRegistration shows how to register only specific prompts
// This is useful if you want to expose only certain prompts via MCP
func ExampleSelectivePromptRegistration(server *mcp.Server) error {
	group := GetGemaraPromptSets()
	gemaraSet, _ := group.GetPromptSet("gemara")

	// Select only the prompts you want to expose
	selectedPrompts := []string{
		"gemara_layer1_to_layer3_policy", // The key prompt for Layer 1 → Layer 3
		"gemara_layer_context",           // For getting layer context
		"gemara_layer_validator",         // For validating layer references
	}

	for _, promptName := range selectedPrompts {
		prompt, err := gemaraSet.GetPrompt(promptName)
		if err != nil {
			return fmt.Errorf("prompt %q not found: %w", promptName, err)
		}

		// Convert to MCP prompt
		mcpPrompt := mcp.Prompt{
			Name:        prompt.Name,
			Description: prompt.Description,
		}

		// Add arguments based on prompt variables
		for varName, varDesc := range prompt.Variables {
			mcpPrompt.Arguments = append(mcpPrompt.Arguments, mcp.PromptArgument{
				Name:        varName,
				Description: varDesc,
				Required:    true,
			})
		}

		// Create handler
		handler := func(ctx context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
			return ExampleMCPServerHandler(ctx, req)
		}

		server.AddPrompt(&mcpPrompt, handler)
	}

	return nil
}

// RunExampleServer demonstrates a complete MCP server setup
func RunExampleServer() {
	// This would typically be in your main.go or server.go
	server := mcp.NewServer(
		mcp.Options{
			Name:    "gemara-mcp-server",
			Version: "0.1.0",
		},
	)

	// Option 1: Register all Gemara prompts
	if err := RegisterGemaraPromptsWithMCPServer(server); err != nil {
		log.Fatalf("Failed to register prompts: %v", err)
	}

	// Option 2: Register only specific prompts (uncomment to use)
	// if err := ExampleSelectivePromptRegistration(server); err != nil {
	//     log.Fatalf("Failed to register selected prompts: %v", err)
	// }

	// Run the server (this is pseudocode - actual implementation depends on your server setup)
	log.Println("MCP server ready")
	log.Println("Registered prompts can be accessed via MCP protocol")

	// In a real implementation, you would:
	// - Set up stdio transport or other transport
	// - Handle incoming requests
	// - Process prompt requests using the handlers
}

// ExampleVariableScopeChange demonstrates how users can change scope
func ExampleVariableScopeChange() {
	ctx := context.Background()
	group := GetGemaraPromptSets()
	gemaraSet, _ := group.GetPromptSet("gemara")

	// Different scopes for the same Layer 1 → Layer 3 prompt
	scopes := []string{
		"Cloud Infrastructure Security",
		"API Security",
		"Container Security",
		"Data Protection and Privacy",
		"Network Security",
	}

	for _, scope := range scopes {
		req := PromptRequest{
			PromptName: "gemara_layer1_to_layer3_policy",
			Variables: map[string]interface{}{
				"scope":                   scope, // Scope changes for each iteration
				"layer1_guidance_source":  "ISO 27001",
				"layer1_guidance_content": "ISO 27001 information security management system requirements",
				"organization_context":    "Technology company, 1000 employees",
				"risk_appetite":           "Moderate risk tolerance",
				"additional_requirements": "Must support agile development practices",
			},
		}

		resp, err := gemaraSet.GeneratePrompt(ctx, req)
		if err != nil {
			log.Printf("Error for scope %q: %v", scope, err)
			continue
		}

		fmt.Printf("\n=== Scope: %s ===\n", scope)
		fmt.Printf("Prompt length: %d characters\n", len(resp.Content))
		fmt.Printf("Metadata: %+v\n", resp.Metadata)
	}
}
