package promptsets

import (
	"context"
	"fmt"
	"log"
)

// ExampleChatbotUsage demonstrates how user-facing prompts work with chatbot interfaces
func ExampleChatbotUsage() {
	ctx := context.Background()
	group := GetUserFacingPromptSets()
	userFacingSet, _ := group.GetPromptSet("user_facing")

	// Example 1: User wants to create a Layer 3 policy with Layer 1 mappings
	// The user can change the scope via the chatbot interface
	fmt.Println("=== Example 1: User Creates Policy with Custom Scope ===")

	// User says: "Create a policy for API Security"
	req1 := PromptRequest{
		PromptName: "create_layer3_policy_with_layer1_mappings",
		Variables: map[string]interface{}{
			"scope":                   "API Security", // User changed this via chatbot
			"organization_context":    "Technology company, 1000 employees, microservices architecture",
			"risk_appetite":           "Moderate risk tolerance",
			"additional_requirements": "Must support REST and GraphQL APIs",
		},
	}

	resp1, err := userFacingSet.GeneratePrompt(ctx, req1)
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Printf("Generated prompt for scope: %s\n", resp1.Metadata["scope"])
		fmt.Printf("Prompt length: %d characters\n\n", len(resp1.Content))
		fmt.Println("First 500 characters of prompt:")
		if len(resp1.Content) > 500 {
			fmt.Println(resp1.Content[:500] + "...")
		} else {
			fmt.Println(resp1.Content)
		}
		fmt.Println()
	}

	// Example 2: User changes scope to a different domain
	fmt.Println("=== Example 2: User Changes Scope ===")

	// User says: "Actually, make it for Container Security instead"
	req2 := PromptRequest{
		PromptName: "create_layer3_policy_with_layer1_mappings",
		Variables: map[string]interface{}{
			"scope":                   "Container Security", // User changed scope
			"organization_context":    "Same organization",
			"risk_appetite":           "Same risk appetite",
			"additional_requirements": "Focus on Kubernetes and Docker",
		},
	}

	resp2, err := userFacingSet.GeneratePrompt(ctx, req2)
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Printf("Generated prompt for NEW scope: %s\n", resp2.Metadata["scope"])
		fmt.Printf("Prompt automatically adjusted for new scope\n\n")
	}

	// Example 3: Interactive scope customization
	fmt.Println("=== Example 3: Interactive Scope Customization ===")

	userRequest := "I want to create a policy but I'm not sure what scope to use. Can you help?"
	currentScope := "Cloud Infrastructure Security"

	req3 := PromptRequest{
		PromptName: "customize_policy_scope",
		Variables: map[string]interface{}{
			"scope":        currentScope,
			"user_request": userRequest,
		},
	}

	resp3, err := userFacingSet.GeneratePrompt(ctx, req3)
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Println("Chatbot response prompt:")
		fmt.Println(resp3.Content)
		fmt.Println()
	}
}

// ExampleChatbotInterface shows how a chatbot would use these prompts
func ExampleChatbotInterface() {
	ctx := context.Background()
	group := GetUserFacingPromptSets()
	userFacingSet, _ := group.GetPromptSet("user_facing")

	// Simulate a chatbot conversation
	conversations := []struct {
		userMessage string
		scope       string
	}{
		{
			userMessage: "Create an OpenSSF Gemara Project layer 3 policy that conforms to the schema and is scoped for Cloud Infrastructure Security and gather all gemara layer 1 guidance mappings that apply to the specified scope.",
			scope:       "Cloud Infrastructure Security",
		},
		{
			userMessage: "Now do the same but for API Security",
			scope:       "API Security",
		},
		{
			userMessage: "What about Container Security?",
			scope:       "Container Security",
		},
	}

	fmt.Println("=== Simulated Chatbot Conversation ===")

	for i, conv := range conversations {
		fmt.Printf("User: %s\n\n", conv.userMessage)

		req := PromptRequest{
			PromptName: "create_layer3_policy_with_layer1_mappings",
			Variables: map[string]interface{}{
				"scope":                   conv.scope,
				"organization_context":    "Technology company",
				"risk_appetite":           "Moderate",
				"additional_requirements": "",
			},
		}

		resp, err := userFacingSet.GeneratePrompt(ctx, req)
		if err != nil {
			fmt.Printf("Error: %v\n\n", err)
			continue
		}

		fmt.Printf("Chatbot: Generated prompt for scope '%s'\n", conv.scope)
		fmt.Printf("Prompt includes:\n")
		fmt.Printf("  - Layer 3 policy generation instructions\n")
		fmt.Printf("  - Layer 1 guidance mapping gathering\n")
		fmt.Printf("  - Schema conformance requirements\n")
		fmt.Printf("  - Scope-specific context: %s\n", conv.scope)
		fmt.Printf("\n---\n\n")

		if i == 0 {
			fmt.Println("(Showing full prompt for first request)")
			fmt.Println(resp.Content)
			fmt.Println()
			fmt.Println("---")
			fmt.Println()
		}
	}
}

// ExampleVariableScopeChangeViaChatbot demonstrates how users can change scope dynamically
func ExampleVariableScopeChangeViaChatbot() {
	ctx := context.Background()
	group := GetUserFacingPromptSets()
	userFacingSet, _ := group.GetPromptSet("user_facing")

	// Different scopes a user might request via chatbot
	userScopes := []string{
		"Cloud Infrastructure Security",
		"API Security",
		"Container Security",
		"Data Protection and Privacy",
		"Network Security",
		"Identity and Access Management",
		"Incident Response",
		"Secure Software Development",
		"Supply Chain Security",
	}

	fmt.Println("=== User Can Change Scope Dynamically ===")
	fmt.Println("The same prompt can be used with different scopes:")

	for _, scope := range userScopes {
		req := PromptRequest{
			PromptName: "create_layer3_policy_with_layer1_mappings",
			Variables: map[string]interface{}{
				"scope":                   scope,
				"organization_context":    "Example organization",
				"risk_appetite":           "Moderate",
				"additional_requirements": "",
			},
		}

		resp, err := userFacingSet.GeneratePrompt(ctx, req)
		if err != nil {
			fmt.Printf("❌ Error for scope '%s': %v\n", scope, err)
			continue
		}

		fmt.Printf("✅ Scope: %s\n", scope)
		fmt.Printf("   Prompt generated successfully\n")
		fmt.Printf("   Metadata: scope=%v, user_facing=%v\n", resp.Metadata["scope"], resp.Metadata["user_facing"])
		fmt.Println()
	}
}

// ExampleChatbotIntegration shows how to integrate with an actual chatbot
func ExampleChatbotIntegration() {
	ctx := context.Background()
	group := GetUserFacingPromptSets()
	userFacingSet, _ := group.GetPromptSet("user_facing")

	// This is how a chatbot would handle user input
	handleUserRequest := func(userMessage string, currentScope string) (string, error) {
		// Parse user message to extract scope if mentioned
		// For this example, we'll use a simple approach

		// Check if user wants to create a policy
		if contains(userMessage, "create") && contains(userMessage, "policy") {
			// Extract scope from message or use current scope
			scope := extractScope(userMessage, currentScope)

			req := PromptRequest{
				PromptName: "create_layer3_policy_with_layer1_mappings",
				Variables: map[string]interface{}{
					"scope":                   scope,
					"organization_context":    "User's organization",
					"risk_appetite":           "To be determined",
					"additional_requirements": "",
				},
			}

			resp, err := userFacingSet.GeneratePrompt(ctx, req)
			if err != nil {
				return "", err
			}

			return resp.Content, nil
		}

		// Handle scope customization
		if contains(userMessage, "scope") || contains(userMessage, "change") {
			req := PromptRequest{
				PromptName: "customize_policy_scope",
				Variables: map[string]interface{}{
					"scope":        currentScope,
					"user_request": userMessage,
				},
			}

			resp, err := userFacingSet.GeneratePrompt(ctx, req)
			if err != nil {
				return "", err
			}

			return resp.Content, nil
		}

		return "I can help you create Gemara Layer 3 policies. What scope would you like?", nil
	}

	// Example usage
	userMessages := []string{
		"Create a policy for API Security",
		"Change scope to Container Security",
		"What scopes are available?",
	}

	fmt.Println("=== Chatbot Integration Example ===")

	for _, msg := range userMessages {
		fmt.Printf("User: %s\n", msg)
		response, err := handleUserRequest(msg, "Cloud Infrastructure Security")
		if err != nil {
			fmt.Printf("Error: %v\n\n", err)
			continue
		}
		fmt.Printf("Chatbot: [Generated prompt - %d chars]\n\n", len(response))
	}
}

// Helper functions for chatbot integration
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		(len(s) > len(substr) &&
			(s[:len(substr)] == substr ||
				s[len(s)-len(substr):] == substr ||
				containsMiddle(s, substr))))
}

func containsMiddle(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func extractScope(message, defaultScope string) string {
	// Simple scope extraction - in production, use NLP or more sophisticated parsing
	scopes := []string{
		"Cloud Infrastructure Security",
		"API Security",
		"Container Security",
		"Data Protection",
		"Network Security",
		"Identity and Access Management",
		"Incident Response",
		"Secure Software Development",
		"Supply Chain Security",
	}

	messageLower := toLower(message)
	for _, scope := range scopes {
		if contains(messageLower, toLower(scope)) {
			return scope
		}
	}

	return defaultScope
}

func toLower(s string) string {
	result := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		if s[i] >= 'A' && s[i] <= 'Z' {
			result[i] = s[i] + 32
		} else {
			result[i] = s[i]
		}
	}
	return string(result)
}
