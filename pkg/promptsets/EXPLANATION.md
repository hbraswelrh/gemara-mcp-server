# How User-Facing Prompts Work - Complete Explanation

## Overview

User-facing prompts allow users to interact with the Gemara prompt system via chatbot/agent interfaces. The key feature is **dynamic scope customization** - users can change the scope (e.g., "API Security" → "Container Security") without creating new prompts.

## How It Works

### 1. The Prompt Template

The prompt is stored as a template with placeholders:

```
You are a GRC expert. Create a Layer 3 policy for: {{scope}}

Gather all Layer 1 guidance mappings for: {{scope}}
Create Layer 3 policy conforming to Gemara schema for: {{scope}}
```

### 2. User Makes Request

**Via Chatbot:**
```
User: "Create an OpenSSF Gemara Project layer 3 policy that conforms to 
      the schema and is scoped for API Security and gather all gemara 
      layer 1 guidance mappings that apply to the specified scope."
```

### 3. System Processes Request

```go
// Chatbot extracts scope: "API Security"
req := PromptRequest{
    PromptName: "create_layer3_policy_with_layer1_mappings",
    Variables: map[string]interface{}{
        "scope": "API Security",  // ← User's scope
        "organization_context": "...",
        "risk_appetite": "...",
    },
}
```

### 4. Variable Substitution

The system replaces `{{scope}}` with "API Security":

```
You are a GRC expert. Create a Layer 3 policy for: API Security

Gather all Layer 1 guidance mappings for: API Security
Create Layer 3 policy conforming to Gemara schema for: API Security
```

### 5. Generated Prompt Sent to LLM

The LLM receives the complete prompt with all variables substituted and:
- Understands it needs to gather Layer 1 guidance for "API Security"
- Creates a Layer 3 policy for "API Security"
- Ensures schema conformance
- Tailors to organization context

### 6. User Changes Scope

**User says:** "Actually, make it for Container Security"

```go
// Just change the scope variable!
req.Variables["scope"] = "Container Security"
```

The same template generates a new prompt, now focused on "Container Security".

## Key Benefits

1. **One Template, Many Scopes**: Same prompt template works for any scope
2. **No Code Changes**: Add new scopes without modifying code
3. **Chatbot-Friendly**: Perfect for conversational interfaces
4. **Deterministic**: Same scope always produces same prompt structure
5. **Comprehensive**: Combines Layer 1 mapping + Layer 3 policy generation

## Testing

### Quick Demo (No Dependencies)

```bash
cd pkg/promptsets
go run quick_demo.go
```

This shows variable substitution in action.

### Full Test

Create a test file:

```go
package main

import (
    "context"
    "fmt"
    "strings"
    "github.com/your-org/gemara-mcp-server/pkg/promptsets"
)

func main() {
    ctx := context.Background()
    group := promptsets.GetUserFacingPromptSets()
    userFacingSet, _ := group.GetPromptSet("user_facing")
    
    req := promptsets.PromptRequest{
        PromptName: "create_layer3_policy_with_layer1_mappings",
        Variables: map[string]interface{}{
            "scope":                  "API Security",
            "organization_context":   "Tech company",
            "risk_appetite":          "Moderate",
            "additional_requirements": "",
        },
    }
    
    resp, _ := userFacingSet.GeneratePrompt(ctx, req)
    
    // Verify scope is in the prompt
    if strings.Contains(resp.Content, "API Security") {
        fmt.Println("✅ Scope correctly included!")
    }
    
    fmt.Println(resp.Content)
}
```

### Go Test

```bash
go test -v ./pkg/promptsets
```

## Architecture

```
┌──────────────┐
│   User       │
│  (Chatbot)   │
└──────┬───────┘
       │ "Create policy for API Security"
       ▼
┌──────────────────┐
│  Parse Request   │
│  Extract Scope   │
└──────┬───────────┘
       │ scope = "API Security"
       ▼
┌──────────────────┐
│  PromptRequest   │
│  Variables:      │
│  - scope         │
│  - org_context   │
│  - risk_appetite │
└──────┬───────────┘
       │
       ▼
┌──────────────────┐
│  PromptSet       │
│  GeneratePrompt  │
└──────┬───────────┘
       │ Variable substitution
       │ {{scope}} → "API Security"
       ▼
┌──────────────────┐
│  PromptResponse  │
│  (Final Prompt)  │
└──────┬───────────┘
       │
       ▼
┌──────────────────┐
│  LLM/Agent       │
│  Executes Prompt │
└──────────────────┘
```

## What the Generated Prompt Contains

The `create_layer3_policy_with_layer1_mappings` prompt includes:

1. **Instructions to gather Layer 1 guidance mappings**
   - Identify relevant Layer 1 sources (NIST, ISO 27001, etc.)
   - Map controls/requirements to the scope
   - Document mappings

2. **Instructions to create Layer 3 policy**
   - Generate policy document
   - Conform to Gemara Layer 3 schema
   - Tailor to organization and risk appetite
   - Scope for the specified domain

3. **Schema conformance requirements**
   - Ensure proper structure
   - Include required fields
   - Follow Gemara schema format

4. **Scope-specific context**
   - All instructions reference the provided scope
   - Policy is tailored to that specific domain

## Example: Changing Scope

### Initial Request
```go
Variables: {
    "scope": "API Security"
}
```

Generated prompt focuses on API Security.

### User Changes Scope
```go
Variables: {
    "scope": "Container Security"  // ← Changed
}
```

Same template, but now all instructions reference "Container Security" instead.

## Integration Points

### With Chatbot
1. User sends message
2. Chatbot parses to extract scope
3. Chatbot calls `GeneratePrompt()` with scope
4. Chatbot sends generated prompt to LLM
5. LLM executes instructions

### With MCP Server
1. Register prompts with MCP server
2. Client requests prompt via MCP protocol
3. Server generates prompt with provided variables
4. Returns complete prompt to client
5. Client uses prompt with LLM

### With Agent Framework
1. Agent receives user request
2. Agent calls prompt generation function
3. Agent receives generated prompt
4. Agent uses prompt as system message
5. Agent executes task based on prompt

## Verification Checklist

When testing, verify:

- ✅ Prompt generates without errors
- ✅ Scope variable is substituted correctly
- ✅ Prompt contains Layer 1 mapping instructions
- ✅ Prompt contains Layer 3 policy instructions
- ✅ Prompt mentions schema conformance
- ✅ Different scopes produce different prompts
- ✅ Same scope produces consistent prompts

## Files Reference

- **`prompts.go`**: Contains the prompt templates
- **`promptsets.go`**: Core prompt management system
- **`gemara.go`**: Gemara layer definitions
- **`quick_demo.go`**: Simple demonstration
- **`test_user_facing.go`**: Comprehensive test
- **`chatbot_examples.go`**: Chatbot integration examples

## Next Steps

1. **Run the demo**: `go run quick_demo.go`
2. **Test the system**: Create your own test file
3. **Integrate**: Connect to your chatbot/agent
4. **Customize**: Modify prompts for your needs

## Summary

The system works by:
1. Storing prompts as templates with variables
2. Substituting variables when generating prompts
3. Allowing users to change variables (especially scope) dynamically
4. Generating complete, ready-to-use prompts for LLMs

This enables flexible, chatbot-friendly prompt generation without code changes for new scopes or requirements.

