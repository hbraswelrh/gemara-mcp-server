package promptsets

import (
	"context"
	"fmt"
	"strings"
)

// ExampleUserFacingPrompt demonstrates how to use user-facing prompts
// This can be run as a test or used as an example
func ExampleUserFacingPrompt() {
	ctx := context.Background()

	// Get user-facing prompts
	group := GetUserFacingPromptSets()
	userFacingSet, err := group.GetPromptSet("user_facing")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Test with different scopes
	scopes := []string{
		"API Security",
		"Container Security",
		"Cloud Infrastructure Security",
		"Data Protection and Privacy",
		"Network Security",
		"Identity and Access Management",
		"Incident Response",
		"Secure Software Development",
	}

	for _, scope := range scopes {
		req := PromptRequest{
			PromptName: "create_layer3_policy_with_layer1_mappings",
			Variables: map[string]interface{}{
				"scope":                   scope,
				"organization_context":    "Tech company, 1000 employees",
				"risk_appetite":           "Moderate",
				"additional_requirements": "Must support REST and GraphQL",
			},
		}

		resp, err := userFacingSet.GeneratePrompt(ctx, req)
		if err != nil {
			fmt.Printf("Error for scope %q: %v\n", scope, err)
			continue
		}

		// Verify scope is in the prompt
		if strings.Contains(resp.Content, scope) {
			fmt.Printf("✅ Scope '%s' correctly included!\n", scope)
		} else {
			fmt.Printf("❌ Scope '%s' NOT found in prompt\n", scope)
		}
	}

	// Show full prompt for first scope
	req := PromptRequest{
		PromptName: "create_layer3_policy_with_layer1_mappings",
		Variables: map[string]interface{}{
			"scope":                   "API Security",
			"organization_context":    "Tech company, 1000 employees",
			"risk_appetite":           "Moderate",
			"additional_requirements": "Must support REST and GraphQL",
		},
	}

	resp, err := userFacingSet.GeneratePrompt(ctx, req)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Print the generated prompt
	fmt.Println("\nGenerated Prompt for 'API Security':")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println(resp.Content)
	fmt.Println(strings.Repeat("=", 80))
}
