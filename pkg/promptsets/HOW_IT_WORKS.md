# How User-Facing Prompts Work

## Architecture Overview

```
┌─────────────────┐
│  User/Chatbot   │
│   Interface     │
└────────┬────────┘
         │
         │ "Create policy for API Security"
         ▼
┌─────────────────┐
│  Prompt Request │
│  (with scope)   │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  PromptSet      │
│  GeneratePrompt │
└────────┬────────┘
         │
         │ Variable substitution
         │ + Handler logic (if any)
         ▼
┌─────────────────┐
│  PromptResponse │
│  (Final prompt)  │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  LLM/Agent      │
│  (Uses prompt)   │
└─────────────────┘
```

## Step-by-Step Flow

### 1. User Makes Request via Chatbot

**User says:** "Create an OpenSSF Gemara Project layer 3 policy that conforms to the schema and is scoped for API Security and gather all gemara layer 1 guidance mappings that apply to the specified scope."

### 2. Chatbot Parses Request

The chatbot extracts:
- **Intent**: Create Layer 3 policy with Layer 1 mappings
- **Scope**: "API Security"
- **Other context**: Organization info, risk appetite, etc.

### 3. Chatbot Creates PromptRequest

```go
req := PromptRequest{
    PromptName: "create_layer3_policy_with_layer1_mappings",
    Variables: map[string]interface{}{
        "scope":                  "API Security",
        "organization_context":   "Technology company, 1000 employees",
        "risk_appetite":          "Moderate risk tolerance",
        "additional_requirements": "Must support REST and GraphQL APIs",
    },
}
```

### 4. PromptSet Generates Final Prompt

The system:
1. Loads the prompt template
2. Substitutes variables (`{{scope}}` → "API Security")
3. Applies any handlers (if defined)
4. Returns the complete prompt

### 5. Generated Prompt is Used

The final prompt is sent to the LLM/agent, which:
- Understands it needs to gather Layer 1 guidance mappings for "API Security"
- Creates a Layer 3 policy conforming to Gemara schema
- Ensures the policy is scoped for "API Security"

### 6. User Changes Scope

**User says:** "Actually, make it for Container Security instead"

The chatbot simply changes the scope variable:

```go
req.Variables["scope"] = "Container Security" // Just change this!
```

The same prompt template is used, but now all instructions reference "Container Security" instead of "API Security".

## Key Components

### Prompt Template

The prompt template contains placeholders like `{{scope}}` that get replaced:

```
Scope: {{scope}}
The scope defines the domain, technology, or area this policy will cover...
```

### Variable Substitution

When `scope = "API Security"`, the template becomes:

```
Scope: API Security
The scope defines the domain, technology, or area this policy will cover...
```

### Dynamic Behavior

- **Same prompt template** for all scopes
- **Different output** based on variable values
- **No code changes** needed to support new scopes

## Testing

See `test_user_facing.go` for complete test examples.

