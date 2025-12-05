# Prompt Sets System

A flexible system for managing and generating system prompts, inspired by the `github-mcp-server` toolsets structure. This package allows you to organize prompts into logical groups (prompt sets) and generate dynamic prompts based on templates and handlers.

## Overview

The prompt sets system provides:

- **Organized Prompt Management**: Group related prompts into logical sets
- **Dynamic Prompt Generation**: Use templates with variable substitution or custom handlers
- **Extensibility**: Easily add new prompts and prompt sets
- **Type Safety**: Strongly typed Go structures for prompts and handlers

## Core Concepts

### Prompt

A `Prompt` represents a system prompt template that can be used as a basis for dynamic prompts. It includes:

- **Name**: Unique identifier
- **Description**: Context about the prompt's purpose
- **Content**: The prompt text/template with variable placeholders (`{{variable_name}}`)
- **Variables**: Map of variable names to their descriptions
- **Category**: Groups related prompts
- **Tags**: For searching and filtering

### PromptSet

A `PromptSet` is a collection of related prompts, similar to `github-mcp-server`'s `Toolset`. It provides:

- Prompt storage and retrieval
- Dynamic prompt generation
- Handler support for complex prompt logic

### PromptSetGroup

A `PromptSetGroup` manages multiple prompt sets, allowing you to organize prompts across different categories or domains.

## Usage

### Basic Usage

```go
package main

import (
    "context"
    "fmt"
    "github.com/your-org/gemara-mcp-server/pkg/promptsets"
)

func main() {
    // Get default prompt sets
    group := promptsets.GetDefaultPromptSets()
    
    // Generate a prompt
    ctx := context.Background()
    resp, err := group.GeneratePromptFromGroup(
        ctx,
        "code_generation",  // prompt set name
        "code_generator",   // prompt name
        map[string]interface{}{
            "requirements": "Create a REST API endpoint",
            "language":      "Go",
            "framework":    "Gin",
            "style":        "Clean code",
            "constraints":  "Secure and performant",
        },
    )
    
    if err != nil {
        panic(err)
    }
    
    fmt.Println(resp.Content)
}
```

### Creating Custom Prompts

```go
// Create a new prompt
prompt := &promptsets.Prompt{
    Name:        "my_custom_prompt",
    Description: "A custom prompt for my use case",
    Content:     "You are helping with: {{task}}. Context: {{context}}",
    Variables: map[string]string{
        "task":    "The task to perform",
        "context": "Additional context",
    },
    Category: "custom",
    Tags:     []string{"custom", "specialized"},
}

// Create a prompt set and add the prompt
ps := promptsets.NewPromptSet("my_set", "My custom prompt set")
err := ps.AddPrompt(prompt)
if err != nil {
    panic(err)
}

// Add to a group
group := promptsets.NewPromptSetGroup()
err = group.AddPromptSet(ps)
if err != nil {
    panic(err)
}
```

### Using Dynamic Handlers

For more complex prompt generation logic, use handlers:

```go
// Create a handler
handler := func(ctx context.Context, req promptsets.PromptRequest) (promptsets.PromptResponse, error) {
    // Complex logic to generate prompt
    analysisType := req.Variables["analysis_type"].(string)
    target := req.Variables["target"].(string)
    
    content := fmt.Sprintf(
        "Perform a %s analysis on: %s\n\nProvide detailed insights.",
        analysisType,
        target,
    )
    
    return promptsets.PromptResponse{
        Content: content,
        Metadata: map[string]interface{}{
            "generated_by": "dynamic_handler",
        },
    }, nil
}

// Add handler to prompt set
ps := promptsets.NewPromptSet("analysis", "Analysis prompts")
err := ps.AddHandler("dynamic_analysis", handler)
if err != nil {
    panic(err)
}

// Use the handler
req := promptsets.PromptRequest{
    PromptName: "dynamic_analysis",
    Variables: map[string]interface{}{
        "analysis_type": "security",
        "target":        "authentication system",
    },
}

resp, err := ps.GeneratePrompt(ctx, req)
```

## Available Prompt Sets

The package includes several default prompt sets:

### Code Generation (`code_generation`)
- `code_generator`: Generate code based on requirements
- `code_reviewer`: Review and improve code

### Documentation (`documentation`)
- `documentation_generator`: Generate technical documentation

### Analysis (`analysis`)
- `code_analyzer`: Analyze code structure and patterns
- `security_analyzer`: Perform security analysis

### Testing (`testing`)
- `test_generator`: Generate test cases

### General (`general`)
- `task_executor`: General-purpose task execution
- `problem_solver`: Problem-solving and debugging

### Gemara GRC Model (`gemara`)

The Gemara prompt set provides system prompts that are aware of the [OpenSSF Gemara Project](https://github.com/ossf/gemara) (GRC Engineering Model for Automated Risk Assessment) layered architecture. This enables deterministic usage of Gemara layers in dynamic prompts.

#### Gemara Architecture

The Gemara model organizes governance activities into 6 layers:

| Layer | Name        | Description                                            |
| ----- | ----------- | ------------------------------------------------------ |
| 1     | Guidance    | High-level guidance on cybersecurity measures          |
| 2     | Controls    | Technology-specific, threat-informed security controls |
| 3     | Policy      | Risk-informed guidance tailored to an organization     |
| 4     | Evaluation  | Inspection of code, configurations, and deployments    |
| 5     | Enforcement | Prevention or remediation based on assessment findings |
| 6     | Audit       | Review of organizational policy and conformance        |

#### Gemara Prompts

- `gemara_architecture_explainer`: Explains the complete Gemara layered architecture
- `gemara_layer_context`: Provides context for a specific Gemara layer with validation
- `gemara_layer_1_guidance`: System prompt for Layer 1 (Guidance) work
- `gemara_layer_2_controls`: System prompt for Layer 2 (Controls) work
- `gemara_layer_3_policy`: System prompt for Layer 3 (Policy) work
- `gemara_layer_4_evaluation`: System prompt for Layer 4 (Evaluation) work
- `gemara_layer_5_enforcement`: System prompt for Layer 5 (Enforcement) work
- `gemara_layer_6_audit`: System prompt for Layer 6 (Audit) work
- `gemara_layer_validator`: Validates layer references and ensures proper usage

#### Using Gemara Prompts

```go
// Get Gemara prompt sets
group := promptsets.GetGemaraPromptSets()
gemaraSet, _ := group.GetPromptSet("gemara")

// Use Layer 1 (Guidance) prompt
req := promptsets.PromptRequest{
    PromptName: "gemara_layer_1_guidance",
    Variables: map[string]interface{}{
        "task":    "Create guidance for secure API development",
        "subject": "REST API Security",
    },
}

resp, err := gemaraSet.GeneratePrompt(ctx, req)
```

#### Deterministic Layer References

Layer references are validated to ensure deterministic usage. Valid formats include:
- `"Layer 1"` through `"Layer 6"` (case-insensitive)
- Numbers `"1"` through `"6"`
- Layer names: `"Guidance"`, `"Controls"`, `"Policy"`, `"Evaluation"`, `"Enforcement"`, `"Audit"`

```go
// Validate layer reference
layerNum, err := promptsets.ValidateLayerReference("Layer 1")
if err != nil {
    // Handle invalid reference
}

// Get layer information
layer, _ := promptsets.GetLayer(layerNum)
fmt.Printf("Layer %d: %s\n", layer.Number, layer.Name)
```

#### Layer Context with Validation

The `gemara_layer_context` prompt automatically validates layer references and provides appropriate context:

```go
req := promptsets.PromptRequest{
    PromptName: "gemara_layer_context",
    Variables: map[string]interface{}{
        "layer_reference": "Layer 2",  // Automatically validated
        "task":            "Generate security controls",
        "subject":         "Kubernetes security",
    },
}

resp, err := gemaraSet.GeneratePrompt(ctx, req)
// resp.Content includes validated layer context
// resp.Metadata includes layer_number and layer_name
```

#### Layer-Specific Artifacts

Each layer has expected artifacts that should be associated with it:

- **Layer 1 (Guidance)**: Guidance frameworks, industry standards, high-level cybersecurity rules
- **Layer 2 (Controls)**: Technology-specific controls, threat-informed controls, control catalogs
- **Layer 3 (Policy)**: Organizational policies, risk-informed governance rules
- **Layer 4 (Evaluation)**: Assessment results, evaluation reports, code inspection results
- **Layer 5 (Enforcement)**: Enforcement actions, remediation plans, prevention mechanisms
- **Layer 6 (Audit)**: Audit reports, conformance reviews, policy effectiveness assessments

The system prompts ensure that users reference the correct layer and that artifacts align with layer expectations.

## Variable Substitution

Prompts support variable substitution using `{{variable_name}}` syntax:

```go
prompt := &promptsets.Prompt{
    Content: "Hello {{name}}, you are working on {{task}}",
    Variables: map[string]string{
        "name": "The person's name",
        "task": "The current task",
    },
}

req := promptsets.PromptRequest{
    PromptName: "my_prompt",
    Variables: map[string]interface{}{
        "name": "Alice",
        "task": "building an API",
    },
}

// Generates: "Hello Alice, you are working on building an API"
```

## Best Practices

1. **Organize by Domain**: Group related prompts into the same prompt set
2. **Use Descriptive Names**: Choose clear, descriptive names for prompts and sets
3. **Document Variables**: Always document what each variable represents
4. **Use Handlers for Complexity**: Use handlers when simple variable substitution isn't enough
5. **Tag Your Prompts**: Use tags to make prompts searchable and filterable
6. **Validate Input**: Always validate variables before using them in handlers

## Integration with MCP Server

To integrate with an MCP server, you can expose prompts as MCP resources or tools:

```go
// Example: Expose prompts as MCP prompts
func (s *Server) ListPrompts(ctx context.Context) ([]mcp.Prompt, error) {
    group := promptsets.GetDefaultPromptSets()
    var mcpPrompts []mcp.Prompt
    
    for _, ps := range group.ListPromptSets() {
        for _, prompt := range ps.ListPrompts() {
            mcpPrompts = append(mcpPrompts, mcp.Prompt{
                Name:        prompt.Name,
                Description: prompt.Description,
            })
        }
    }
    
    return mcpPrompts, nil
}

func (s *Server) GetPrompt(ctx context.Context, req mcp.PromptRequest) (mcp.PromptResponse, error) {
    group := promptsets.GetDefaultPromptSets()
    
    // Extract prompt set and prompt name from request
    promptSetName := req.Params["prompt_set"].(string)
    promptName := req.Name
    
    resp, err := group.GeneratePromptFromGroup(
        ctx,
        promptSetName,
        promptName,
        req.Params,
    )
    
    if err != nil {
        return mcp.PromptResponse{}, err
    }
    
    return mcp.PromptResponse{
        Content: []mcp.PromptMessage{
            {
                Role:    "user",
                Content: resp.Content,
            },
        },
    }, nil
}
```

## Examples

See `examples.go` for complete usage examples including:
- Basic prompt generation
- Custom prompt creation
- Dynamic handler usage
- Prompt set management

See `gemara_examples.go` for Gemara-specific examples including:
- Gemara architecture explanation
- Layer-specific prompt usage
- Layer reference validation
- Deterministic layer usage patterns

## Contributing

When adding new prompts:

1. Add them to the appropriate category in `SystemPromptDefinitions`
2. Ensure all variables are documented
3. Include appropriate tags
4. Test variable substitution
5. Update this README if adding a new category

## License

Apache License 2.0

