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
        t.Fatalf("Failed to get user-facing prompt set: %v", err)
    }
    
    req := PromptRequest{
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
        t.Fatalf("Failed to generate prompt: %v", err)
    }
    
    // Verify prompt is not empty
    if resp.Content == "" {
        t.Error("Generated prompt content is empty")
    }
    
    // Verify scope appears in the prompt
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
    contentLower := strings.ToLower(resp.Content)
    if !strings.Contains(contentLower, "schema") {
        t.Error("Schema conformance not mentioned in prompt")
    }
}

func TestScopeChange(t *testing.T) {
    ctx := context.Background()
    group := GetUserFacingPromptSets()
    userFacingSet, err := group.GetPromptSet("user_facing")
    if err != nil {
        t.Fatalf("Failed to get user-facing prompt set: %v", err)
    }
    
    scopes := []string{
        "API Security",
        "Container Security",
        "Cloud Infrastructure Security",
        "Data Protection and Privacy",
    }
    
    for _, scope := range scopes {
        req := PromptRequest{
            PromptName: "create_layer3_policy_with_layer1_mappings",
            Variables: map[string]interface{}{
                "scope":                  scope,
                "organization_context":   "Test organization",
                "risk_appetite":          "Moderate",
                "additional_requirements": "",
            },
        }
        
        resp, err := userFacingSet.GeneratePrompt(ctx, req)
        if err != nil {
            t.Errorf("Failed to generate prompt for scope %q: %v", scope, err)
            continue
        }
        
        if !strings.Contains(resp.Content, scope) {
            t.Errorf("Scope %q not found in generated prompt", scope)
        }
        
        // Verify metadata
        if resp.Metadata["scope"] != scope {
            t.Errorf("Expected scope %q in metadata, got %v", scope, resp.Metadata["scope"])
        }
    }
}

func TestUserFacingPromptSetExists(t *testing.T) {
    group := GetUserFacingPromptSets()
    userFacingSet, err := group.GetPromptSet("user_facing")
    if err != nil {
        t.Fatalf("User-facing prompt set not found: %v", err)
    }
    
    prompts := userFacingSet.ListPrompts()
    if len(prompts) == 0 {
        t.Error("No prompts found in user-facing prompt set")
    }
    
    // Verify the main prompt exists
    found := false
    for _, prompt := range prompts {
        if prompt.Name == "create_layer3_policy_with_layer1_mappings" {
            found = true
            break
        }
    }
    
    if !found {
        t.Error("Main prompt 'create_layer3_policy_with_layer1_mappings' not found")
    }
}

