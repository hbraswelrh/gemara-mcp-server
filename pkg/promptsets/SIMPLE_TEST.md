# Simple Testing Instructions

## How It Works - Quick Explanation

### The Flow

1. **User Request** → Chatbot receives: "Create policy for API Security"
2. **Parse Request** → Chatbot extracts scope: "API Security"
3. **Create PromptRequest** → Build request with variables
4. **Generate Prompt** → System substitutes variables into template
5. **Use Prompt** → Send to LLM/agent to execute

### Key Concept: Variable Substitution

The prompt template has placeholders like `{{scope}}`:

```
Scope: {{scope}}
Create policy for: {{scope}}
```

When you provide `scope = "API Security"`, it becomes:

```
Scope: API Security
Create policy for: API Security
```

**Same template, different output based on variables!**

## Quick Test (No Dependencies)

Run the simple demo:

```bash
cd pkg/promptsets
go run quick_demo.go
```

This shows how variable substitution works.

## Full Test (With Package)

### Option 1: Create a Test File

Create `example_usage.go` in the package directory:

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
    userFacingSet, err := group.GetPromptSet("user_facing")
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    
    // Test with API Security
    req := promptsets.PromptRequest{
        PromptName: "create_layer3_policy_with_layer1_mappings",
        Variables: map[string]interface{}{
            "scope":                  "API Security",
            "organization_context":   "Tech company, 1000 employees",
            "risk_appetite":          "Moderate",
            "additional_requirements": "Must support REST and GraphQL",
        },
    }
    
    resp, err := userFacingSet.GeneratePrompt(ctx, req)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    
    fmt.Println("Generated Prompt:")
    fmt.Println("=" + strings.Repeat("=", 80))
    fmt.Println(resp.Content)
    fmt.Println("=" + strings.Repeat("=", 80))
    
    // Verify scope is in the prompt
    if strings.Contains(resp.Content, "API Security") {
        fmt.Println("\n✅ Scope 'API Security' found in prompt!")
    }
}
```

Run it:
```bash
go run example_usage.go
```

### Option 2: Use Go Test

Create `prompts_test.go`:

```go
package promptsets

import (
    "context"
    "strings"
    "testing"
)

func TestCreateLayer3PolicyPrompt(t *testing.T) {
    ctx := context.Background()
    group := GetUserFacingPromptSets()
    userFacingSet, err := group.GetPromptSet("user_facing")
    if err != nil {
        t.Fatalf("Failed to get prompt set: %v", err)
    }
    
    req := PromptRequest{
        PromptName: "create_layer3_policy_with_layer1_mappings",
        Variables: map[string]interface{}{
            "scope":                  "API Security",
            "organization_context":   "Test org",
            "risk_appetite":          "Moderate",
            "additional_requirements": "",
        },
    }
    
    resp, err := userFacingSet.GeneratePrompt(ctx, req)
    if err != nil {
        t.Fatalf("Failed to generate prompt: %v", err)
    }
    
    // Verify scope appears
    if !strings.Contains(resp.Content, "API Security") {
        t.Error("Scope 'API Security' not found in generated prompt")
    }
    
    // Verify Layer 1 mapping instructions
    if !strings.Contains(resp.Content, "Layer 1") {
        t.Error("Layer 1 mapping instructions not found")
    }
    
    // Verify Layer 3 policy instructions
    if !strings.Contains(resp.Content, "Layer 3") {
        t.Error("Layer 3 policy instructions not found")
    }
    
    // Verify schema mention
    if !strings.Contains(strings.ToLower(resp.Content), "schema") {
        t.Error("Schema conformance not mentioned")
    }
}

func TestScopeChange(t *testing.T) {
    ctx := context.Background()
    group := GetUserFacingPromptSets()
    userFacingSet, _ := group.GetPromptSet("user_facing")
    
    scopes := []string{"API Security", "Container Security", "Cloud Security"}
    
    for _, scope := range scopes {
        req := PromptRequest{
            PromptName: "create_layer3_policy_with_layer1_mappings",
            Variables: map[string]interface{}{
                "scope":                  scope,
                "organization_context":   "Test",
                "risk_appetite":          "Test",
                "additional_requirements": "",
            },
        }
        
        resp, err := userFacingSet.GeneratePrompt(ctx, req)
        if err != nil {
            t.Errorf("Failed for scope %s: %v", scope, err)
            continue
        }
        
        if !strings.Contains(resp.Content, scope) {
            t.Errorf("Scope %s not found in prompt", scope)
        }
    }
}
```

Run tests:
```bash
go test -v ./pkg/promptsets
```

## What You Should See

When testing, you should see:

1. ✅ **Prompt Generated**: No errors, prompt content is non-empty
2. ✅ **Scope Included**: Your scope appears in the generated prompt
3. ✅ **Layer 1 Instructions**: Prompt includes gathering Layer 1 mappings
4. ✅ **Layer 3 Instructions**: Prompt includes creating Layer 3 policy
5. ✅ **Schema Mentioned**: Prompt mentions schema conformance

## Testing Different Scopes

Try these scopes:
- "API Security"
- "Container Security"
- "Cloud Infrastructure Security"
- "Data Protection"
- "Network Security"

Each should generate a different prompt (with scope-specific content) but using the same template.

## Integration with Chatbot

To test with a chatbot:

1. Set up your chatbot to call the prompt generation function
2. Parse user messages to extract scope
3. Generate prompt with extracted scope
4. Send generated prompt to LLM
5. LLM executes the instructions in the prompt

Example chatbot flow:
```
User: "Create policy for API Security"
  ↓
Chatbot: Extract scope = "API Security"
  ↓
Chatbot: Generate prompt with scope
  ↓
Chatbot: Send to LLM
  ↓
LLM: Gathers Layer 1 mappings + Creates Layer 3 policy
```

## Troubleshooting

**Error: "prompt not found"**
- Use `GetUserFacingPromptSets()` not `GetDefaultPromptSets()`
- Check prompt name spelling

**Error: "scope is required"**
- Make sure you include `scope` in Variables map
- Scope must be non-empty string

**Scope not in output**
- Check variable substitution is working
- Verify template has `{{scope}}` placeholder

## Next Steps

1. Run `quick_demo.go` to see the concept
2. Create `example_usage.go` to test full system
3. Integrate with your chatbot/agent
4. Test with actual LLM

