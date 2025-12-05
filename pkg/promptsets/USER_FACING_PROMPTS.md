# User-Facing Prompts for Chatbot/Agent Interfaces

This document describes the user-facing prompts that can be accessed and modified via chatbot/agent interfaces.

## Overview

User-facing prompts are designed to be:
- **Accessible**: Can be called directly by users via chatbot interfaces
- **Customizable**: Variables (especially `scope`) can be changed dynamically
- **Interactive**: Support conversational workflows
- **Comprehensive**: Combine multiple operations (e.g., Layer 1 mapping + Layer 3 policy generation)

## Key Prompt: `create_layer3_policy_with_layer1_mappings`

This is the main prompt that addresses your requirement:

> "Create an OpenSSF Gemara Project layer 3 policy that conforms to the schema and is scoped for [scope] and gather all gemara layer 1 guidance mappings that apply to the specified scope."

### Usage

```go
req := PromptRequest{
    PromptName: "create_layer3_policy_with_layer1_mappings",
    Variables: map[string]interface{}{
        "scope": "API Security", // ← User can change this via chatbot
        "organization_context": "Technology company, 1000 employees",
        "risk_appetite": "Moderate risk tolerance",
        "additional_requirements": "Must support REST and GraphQL",
    },
}

resp, err := userFacingSet.GeneratePrompt(ctx, req)
```

### What It Does

1. **Gathers Layer 1 Guidance Mappings**: 
   - Identifies all relevant Layer 1 guidance sources (NIST, ISO 27001, PCI DSS, etc.)
   - Maps specific controls/requirements to the specified scope
   - Documents how Layer 1 guidance informs the Layer 3 policy

2. **Creates Layer 3 Policy**:
   - Generates a complete Layer 3 policy document
   - Conforms to Gemara Layer 3 schema
   - Tailored to the organization and risk appetite
   - Scoped for the specified domain

3. **Schema Conformance**:
   - Ensures the policy follows Gemara Layer 3 schema structure
   - Includes all required schema fields

### Variables

- **`scope`** (required): The domain/area for the policy. **This is the key variable users can change via chatbot interface.**
  - Examples: "Cloud Infrastructure Security", "API Security", "Container Security", "Data Protection", etc.
  
- **`organization_context`**: Information about the organization
  
- **`risk_appetite`**: Organization's risk tolerance
  
- **`additional_requirements`**: Any specific constraints or requirements

## Available User-Facing Prompts

### 1. `create_layer3_policy_with_layer1_mappings`
**Primary prompt for creating Layer 3 policies with Layer 1 guidance mapping**

- Combines Layer 1 mapping gathering with Layer 3 policy generation
- Supports dynamic scope changes
- Ensures schema conformance

### 2. `analyze_layer1_guidance_for_scope`
**Analyzes Layer 1 guidance sources for a specific scope**

- Identifies applicable Layer 1 guidance sources
- Extracts relevant controls/requirements
- Documents mappings and cross-references

### 3. `generate_layer3_policy_from_guidance`
**Generates Layer 3 policy from provided Layer 1 guidance**

- Takes Layer 1 guidance mappings as input
- Creates schema-conformant Layer 3 policy
- Tailors to organization context

### 4. `customize_policy_scope`
**Interactive prompt for scope customization**

- Helps users choose appropriate scope
- Provides guidance on available scopes
- Supports conversational scope selection

## Chatbot Integration

### Example Conversation Flow

```
User: "Create an OpenSSF Gemara Project layer 3 policy that conforms to the schema and is scoped for API Security and gather all gemara layer 1 guidance mappings that apply to the specified scope."

Chatbot: [Uses create_layer3_policy_with_layer1_mappings with scope="API Security"]
         [Generates prompt that includes both Layer 1 mapping gathering and Layer 3 policy creation]

User: "Actually, make it for Container Security instead"

Chatbot: [Uses same prompt with scope="Container Security"]
         [Automatically adjusts all instructions for new scope]
```

### Changing Scope Dynamically

The `scope` variable can be changed without creating a new prompt:

```go
// First request
req1 := PromptRequest{
    PromptName: "create_layer3_policy_with_layer1_mappings",
    Variables: map[string]interface{}{
        "scope": "API Security",
        // ... other variables
    },
}

// User changes scope
req2 := PromptRequest{
    PromptName: "create_layer3_policy_with_layer1_mappings", // Same prompt
    Variables: map[string]interface{}{
        "scope": "Container Security", // ← Changed scope
        // ... same other variables
    },
}
```

## Getting User-Facing Prompts

```go
// Get user-facing prompt sets
group := promptsets.GetUserFacingPromptSets()
userFacingSet, _ := group.GetPromptSet("user_facing")

// List available prompts
for _, prompt := range userFacingSet.ListPrompts() {
    fmt.Printf("%s: %s\n", prompt.Name, prompt.Description)
}

// Generate a prompt
req := PromptRequest{
    PromptName: "create_layer3_policy_with_layer1_mappings",
    Variables: map[string]interface{}{
        "scope": "Your Scope Here",
        // ... other variables
    },
}

resp, err := userFacingSet.GeneratePrompt(ctx, req)
```

## Integration with MCP Server

User-facing prompts can be registered with an MCP server just like other prompts:

```go
// Get all prompts (including user-facing)
group := promptsets.GetAllPromptSets()

// Register with MCP server
RegisterPromptsWithMCPServer(server, group)
```

Users can then access these prompts via the MCP protocol, and the chatbot/agent interface can:
1. List available prompts
2. Show prompt descriptions and variables
3. Allow users to set variable values (especially `scope`)
4. Generate customized prompts on demand

## Examples

See `chatbot_examples.go` for complete examples including:
- Basic usage with scope changes
- Interactive chatbot conversations
- Dynamic scope customization
- Integration patterns

## Benefits

1. **User-Friendly**: Prompts are designed for direct user interaction
2. **Flexible**: Scope and other variables can be changed dynamically
3. **Comprehensive**: Combines multiple operations (mapping + generation)
4. **Schema-Aware**: Ensures Gemara schema conformance
5. **Interactive**: Supports conversational workflows

## Next Steps

1. Integrate with your chatbot/agent interface
2. Expose prompts via MCP server
3. Allow users to customize scope and other variables
4. Use the generated prompts with your LLM/agent system

