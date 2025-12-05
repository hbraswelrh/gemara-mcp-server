package main

import (
	"fmt"
	"strings"
)

// Simple test that doesn't require the full package
// This demonstrates the concept

func main() {
	fmt.Println("=" + strings.Repeat("=", 80))
	fmt.Println("Quick Test: How User-Facing Prompts Work")
	fmt.Println("=" + strings.Repeat("=", 80))
	fmt.Println()

	// Simulate the prompt template
	template := `You are a GRC expert. Create a Layer 3 policy for: {{scope}}

Gather all Layer 1 guidance mappings for: {{scope}}
Create Layer 3 policy conforming to Gemara schema for: {{scope}}`

	// Simulate variable substitution
	variables := map[string]interface{}{
		"scope": "API Security",
	}

	// Simple substitution
	result := template
	for key, value := range variables {
		placeholder := fmt.Sprintf("{{%s}}", key)
		result = strings.ReplaceAll(result, placeholder, fmt.Sprintf("%v", value))
	}

	fmt.Println("1. Original Template:")
	fmt.Println(template)
	fmt.Println()

	fmt.Println("2. After Variable Substitution (scope = 'API Security'):")
	fmt.Println(result)
	fmt.Println()

	// Change scope
	variables["scope"] = "Container Security"
	result2 := template
	for key, value := range variables {
		placeholder := fmt.Sprintf("{{%s}}", key)
		result2 = strings.ReplaceAll(result2, placeholder, fmt.Sprintf("%v", value))
	}

	fmt.Println("3. After Changing Scope to 'Container Security':")
	fmt.Println(result2)
	fmt.Println()

	fmt.Println("=" + strings.Repeat("=", 80))
	fmt.Println("Key Points:")
	fmt.Println("=" + strings.Repeat("=", 80))
	fmt.Println("✅ Same template, different output based on variables")
	fmt.Println("✅ Scope can be changed dynamically")
	fmt.Println("✅ No code changes needed for new scopes")
	fmt.Println("✅ Perfect for chatbot interfaces")
	fmt.Println()
	fmt.Println("To test the full system, use:")
	fmt.Println("  go test ./pkg/promptsets")
	fmt.Println("or see TESTING_GUIDE.md")
}

