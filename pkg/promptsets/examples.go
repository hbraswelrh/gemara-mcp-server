package promptsets

import (
	"context"
	"fmt"
)

// ExampleUsage demonstrates how to use the prompt system
func ExampleUsage() {
	// Create a new prompt set group with default prompts
	group := GetDefaultPromptSets()

	// Example 1: Generate a code generation prompt
	ctx := context.Background()

	codeGenReq := PromptRequest{
		PromptName: "code_generator",
		Variables: map[string]interface{}{
			"requirements": "Create a REST API endpoint for user authentication",
			"language":     "Go",
			"framework":    "Gin",
			"style":        "Clean code with proper error handling",
			"constraints":  "Must be secure and follow OWASP guidelines",
		},
	}

	codeGenSet, _ := group.GetPromptSet("code_generation")
	codeGenResp, err := codeGenSet.GeneratePrompt(ctx, codeGenReq)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Println("Generated Code Generation Prompt:")
	fmt.Println(codeGenResp.Content)
	fmt.Println()

	// Example 2: Generate a documentation prompt
	docReq := PromptRequest{
		PromptName: "documentation_generator",
		Variables: map[string]interface{}{
			"subject":  "REST API Authentication",
			"doc_type": "API Documentation",
			"audience": "Developers",
			"format":   "Markdown",
			"content":  "API endpoints for user authentication",
		},
	}

	docSet, _ := group.GetPromptSet("documentation")
	docResp, err := docSet.GeneratePrompt(ctx, docReq)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Println("Generated Documentation Prompt:")
	fmt.Println(docResp.Content)
	fmt.Println()

	// Example 3: Create a custom prompt set with a dynamic handler
	customSet := NewPromptSet("custom", "Custom prompt set for specialized tasks")

	// Add a custom prompt
	customPrompt := &Prompt{
		Name:        "custom_task",
		Description: "A custom prompt for specialized tasks",
		Content:     "You are helping with: {{task}}. Context: {{context}}",
		Variables: map[string]string{
			"task":    "The task to perform",
			"context": "Additional context",
		},
		Category: "custom",
		Tags:     []string{"custom", "specialized"},
	}

	_ = customSet.AddPrompt(customPrompt)

	// Add a dynamic handler for more complex prompt generation
	_ = customSet.AddHandler("dynamic_analysis", func(ctx context.Context, req PromptRequest) (PromptResponse, error) {
		// This handler can perform complex logic to generate prompts
		analysisType := req.Variables["analysis_type"].(string)
		target := req.Variables["target"].(string)

		content := fmt.Sprintf(
			"Perform a %s analysis on: %s\n\nProvide detailed insights and recommendations.",
			analysisType,
			target,
		)

		return PromptResponse{
			Content: content,
			Metadata: map[string]interface{}{
				"generated_by":  "dynamic_handler",
				"analysis_type": analysisType,
			},
		}, nil
	})

	_ = group.AddPromptSet(customSet)

	// Use the dynamic handler
	dynamicReq := PromptRequest{
		PromptName: "dynamic_analysis",
		Variables: map[string]interface{}{
			"analysis_type": "security",
			"target":        "authentication system",
		},
	}

	dynamicResp, err := customSet.GeneratePrompt(ctx, dynamicReq)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Println("Generated Dynamic Prompt:")
	fmt.Println(dynamicResp.Content)
	fmt.Printf("Metadata: %+v\n", dynamicResp.Metadata)
}

// CreateCustomPromptSet demonstrates creating a custom prompt set
func CreateCustomPromptSet(name, description string, prompts []*Prompt) (*PromptSet, error) {
	ps := NewPromptSet(name, description)

	for _, prompt := range prompts {
		if err := ps.AddPrompt(prompt); err != nil {
			return nil, fmt.Errorf("failed to add prompt %q: %w", prompt.Name, err)
		}
	}

	return ps, nil
}

// CreatePromptSetWithHandlers demonstrates creating a prompt set with handlers
func CreatePromptSetWithHandlers(name, description string, handlers map[string]PromptHandler) (*PromptSet, error) {
	ps := NewPromptSet(name, description)

	for handlerName, handler := range handlers {
		if err := ps.AddHandler(handlerName, handler); err != nil {
			return nil, fmt.Errorf("failed to add handler %q: %w", handlerName, err)
		}
	}

	return ps, nil
}
