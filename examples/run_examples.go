package main

import (
	"fmt"
	"os"

	"github.com/your-org/gemara-mcp-server/pkg/promptsets"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run run_examples.go <example_name>")
		fmt.Println("\nAvailable examples:")
		fmt.Println("  user_facing     - Test user-facing prompts with different scopes")
		fmt.Println("  chatbot         - Test chatbot integration examples")
		fmt.Println("  gemara          - Test Gemara layer examples")
		fmt.Println("  all             - Run all examples")
		os.Exit(1)
	}

	example := os.Args[1]

	switch example {
	case "user_facing":
		fmt.Println("=== Running User-Facing Prompt Example ===")
		fmt.Println()
		promptsets.ExampleUserFacingPrompt()

	case "chatbot":
		fmt.Println("=== Running Chatbot Examples ===")
		fmt.Println()
		promptsets.ExampleChatbotUsage()
		fmt.Println()
		promptsets.ExampleChatbotInterface()

	case "gemara":
		fmt.Println("=== Running Gemara Examples ===")
		fmt.Println()
		promptsets.ExampleGemaraUsage()
		fmt.Println()
		promptsets.ExampleDeterministicLayerUsage()

	case "all":
		fmt.Println("=== Running All Examples ===")
		fmt.Println()
		
		fmt.Println("\n--- User-Facing Prompts ---")
		promptsets.ExampleUserFacingPrompt()
		
		fmt.Println("\n--- Chatbot Examples ---")
		promptsets.ExampleChatbotUsage()
		
		fmt.Println("\n--- Gemara Examples ---")
		promptsets.ExampleGemaraUsage()

	default:
		fmt.Printf("Unknown example: %s\n", example)
		fmt.Println("Use: user_facing, chatbot, gemara, or all")
		os.Exit(1)
	}
}

