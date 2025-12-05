# Testing Guide for User-Facing Prompts

## Quick Start

### Option 1: Run the Test File

```bash
cd /home/hbraswel/GIT/PSCE/Baklava/developer/gem-mcp/gemara-mcp-server/pkg/promptsets
go run test_user_facing.go prompts.go promptsets.go gemara.go
```

This will run comprehensive tests showing:
- Available prompts
- Prompt generation with different scopes
- Variable substitution
- Layer 1 mapping instructions
- Full prompt output

### Option 2: Use Go Test

Create a test file `prompts_test.go`:

```go
package promptsets

import (
    "context"
    "testing"
)

func TestUserFacingPrompts(t *testing.T) {
    ctx := context.Background()
    group := GetUserFacingPromptSets()
    userFacingSet, err := group.GetPromptSet("user_facing")
    if err != nil {
        t.Fatalf("Failed to get user-facing prompt set: %v", err)
    }

    req := PromptRequest{
        PromptName: "create_layer3_policy_with_layer1_mappings",
        Variables: map[string]interface{}{
            "scope":                  "API Security",
            "organization_context":   "Test organization",
            "risk_appetite":          "Moderate",
            "additional_requirements": "",
        },
    }

    resp, err := userFacingSet.GeneratePrompt(ctx, req)
    if err != nil {
        t.Fatalf("Failed to generate prompt: %v", err)
    }

    if resp.Content == "" {
        t.Error("Generated prompt content is empty")
    }

    // Verify scope is in the prompt
    if !contains(resp.Content, "API Security") {
        t.Error("Scope not found in generated prompt")
    }
}
```

Run with:
```bash
go test -v ./pkg/promptsets
```

## Manual Testing

### Step 1: Create a Test Program

Create `test_manual.go`:

```go
package main

import (
    "context"
    "fmt"
    "github.com/your-org/gemara-mcp-server/pkg/promptsets"
)

func main() {
    ctx := context.Background()
    
    // Get user-facing prompts
    group := promptsets.GetUserFacingPromptSets()
    userFacingSet, _ := group.GetPromptSet("user_facing")
    
    // Test with your scope
    req := promptsets.PromptRequest{
        PromptName: "create_layer3_policy_with_layer1_mappings",
        Variables: map[string]interface{}{
            "scope":                  "Your Scope Here",
            "organization_context":   "Your organization",
            "risk_appetite":          "Your risk appetite",
            "additional_requirements": "Your requirements",
        },
    }
    
    resp, err := userFacingSet.GeneratePrompt(ctx, req)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    
    fmt.Println("Generated Prompt:")
    fmt.Println("=" + string(make([]byte, 80)) + "=")
    fmt.Println(resp.Content)
    fmt.Println("=" + string(make([]byte, 80)) + "=")
}
```

### Step 2: Run It

```bash
go run test_manual.go
```

## Testing Different Scopes

### Test Script

```go
package main

import (
    "context"
    "fmt"
    "github.com/your-org/gemara-mcp-server/pkg/promptsets"
)

func main() {
    ctx := context.Background()
    group := promptsets.GetUserFacingPromptSets()
    userFacingSet, _ := group.GetPromptSet("user_facing")
    
    scopes := []string{
        "API Security",
        "Container Security",
        "Cloud Infrastructure Security",
        "Data Protection",
    }
    
    for _, scope := range scopes {
        fmt.Printf("\n=== Testing Scope: %s ===\n\n", scope)
        
        req := promptsets.PromptRequest{
            PromptName: "create_layer3_policy_with_layer1_mappings",
            Variables: map[string]interface{}{
                "scope":                  scope,
                "organization_context":   "Test org",
                "risk_appetite":          "Moderate",
                "additional_requirements": "",
            },
        }
        
        resp, err := userFacingSet.GeneratePrompt(ctx, req)
        if err != nil {
            fmt.Printf("Error: %v\n", err)
            continue
        }
        
        // Verify scope appears
        if contains(resp.Content, scope) {
            fmt.Printf("✅ Scope '%s' correctly included\n", scope)
        } else {
            fmt.Printf("❌ Scope '%s' NOT found\n", scope)
        }
        
        fmt.Printf("Prompt length: %d characters\n", len(resp.Content))
    }
}

func contains(s, substr string) bool {
    // Simple contains check
    for i := 0; i <= len(s)-len(substr); i++ {
        if s[i:i+len(substr)] == substr {
            return true
        }
    }
    return false
}
```

## Testing with Chatbot Interface

### Simulated Chatbot

```go
package main

import (
    "context"
    "fmt"
    "github.com/your-org/gemara-mcp-server/pkg/promptsets"
)

func simulateChatbot() {
    ctx := context.Background()
    group := promptsets.GetUserFacingPromptSets()
    userFacingSet, _ := group.GetPromptSet("user_facing")
    
    // Simulate user conversation
    conversations := []struct {
        userMessage string
        scope       string
    }{
        {
            userMessage: "Create a policy for API Security",
            scope:       "API Security",
        },
        {
            userMessage: "Change it to Container Security",
            scope:       "Container Security",
        },
    }
    
    for _, conv := range conversations {
        fmt.Printf("User: %s\n", conv.userMessage)
        
        req := promptsets.PromptRequest{
            PromptName: "create_layer3_policy_with_layer1_mappings",
            Variables: map[string]interface{}{
                "scope":                  conv.scope,
                "organization_context":   "User's organization",
                "risk_appetite":          "Moderate",
                "additional_requirements": "",
            },
        }
        
        resp, err := userFacingSet.GeneratePrompt(ctx, req)
        if err != nil {
            fmt.Printf("Error: %v\n\n", err)
            continue
        }
        
        fmt.Printf("Chatbot: Generated prompt for '%s' (%d chars)\n\n", 
            conv.scope, len(resp.Content))
    }
}

func main() {
    simulateChatbot()
}
```

## What to Verify

When testing, verify:

1. ✅ **Prompt Generation**: Prompt is generated without errors
2. ✅ **Scope Substitution**: The `{{scope}}` variable is replaced with actual scope
3. ✅ **Layer 1 Instructions**: Prompt includes instructions to gather Layer 1 mappings
4. ✅ **Layer 3 Instructions**: Prompt includes instructions to create Layer 3 policy
5. ✅ **Schema Conformance**: Prompt mentions schema conformance
6. ✅ **Different Scopes**: Same prompt works with different scope values
7. ✅ **Metadata**: Response includes useful metadata (scope, user_facing flag)

## Expected Output

When you run the tests, you should see:

```
Testing User-Facing Prompts
================================================================================

TEST 1: List Available User-Facing Prompts
--------------------------------------------------------------------------------
1. create_layer3_policy_with_layer1_mappings
   Description: Create a Gemara Layer 3 policy conforming to the schema...
   Variables: scope organization_context risk_appetite additional_requirements

TEST 2: Generate Prompt for API Security
--------------------------------------------------------------------------------
✅ Prompt generated successfully
   Scope: API Security
   Prompt length: 2847 characters

   First 300 characters:
   You are a GRC expert working with the OpenSSF Gemara Project. Your task is to create a comprehensive Layer 3 (Policy) document that:

1. Conforms to the Gemara Layer 3 schema
2. Is scoped for: API Security
...
```

## Integration Testing

### With MCP Server

If you have an MCP server set up:

```go
// Register prompts
group := promptsets.GetAllPromptSets()
promptsets.RegisterPromptsWithMCPServer(server, group)

// Test via MCP protocol
// (Use MCP client to call prompts)
```

### With LLM/Agent

1. Generate the prompt using the system
2. Send the generated prompt to your LLM/agent
3. Verify the LLM understands it needs to:
   - Gather Layer 1 guidance mappings
   - Create Layer 3 policy
   - Ensure schema conformance
   - Focus on the specified scope

## Troubleshooting

### Error: "prompt not found"
- Make sure you're using `GetUserFacingPromptSets()` or `GetAllPromptSets()`
- Verify the prompt name is correct: `"create_layer3_policy_with_layer1_mappings"`

### Error: "scope is required"
- Ensure you provide the `scope` variable in `Variables` map
- Scope must be a non-empty string

### Scope not appearing in prompt
- Check that variable substitution is working
- Verify the prompt template contains `{{scope}}`

### Prompt too long/short
- This is normal - prompts are comprehensive
- Length varies based on variable content
- Typical length: 2000-4000 characters

## Next Steps

1. Run the test file to see it in action
2. Modify scopes to test different domains
3. Integrate with your chatbot/agent interface
4. Test with actual LLM to verify prompt effectiveness

