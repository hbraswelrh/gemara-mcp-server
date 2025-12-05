package promptsets

import (
	"context"
	"fmt"
)

// ExampleGemaraUsage demonstrates how to use Gemara-aware prompts
func ExampleGemaraUsage() {
	// Get Gemara prompt sets
	group := GetGemaraPromptSets()
	ctx := context.Background()

	// Example 1: Get architecture explanation
	fmt.Println("=== Example 1: Gemara Architecture Explanation ===")

	archReq := PromptRequest{
		PromptName: "gemara_architecture_explainer",
		Variables:  map[string]interface{}{},
	}

	gemaraSet, _ := group.GetPromptSet("gemara")
	archResp, err := gemaraSet.GeneratePrompt(ctx, archReq)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println(archResp.Content)
		fmt.Println()
	}

	// Example 2: Work with a specific layer (Layer 1 - Guidance)
	fmt.Println("=== Example 2: Layer 1 (Guidance) Context ===")

	layer1Req := PromptRequest{
		PromptName: "gemara_layer_context",
		Variables: map[string]interface{}{
			"layer_reference":    "Layer 1",
			"task":               "Create guidance for secure API development",
			"subject":            "REST API Security",
			"additional_context": "Focus on authentication and authorization",
		},
	}

	layer1Resp, err := gemaraSet.GeneratePrompt(ctx, layer1Req)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println(layer1Resp.Content)
		fmt.Printf("Metadata: %+v\n\n", layer1Resp.Metadata)
	}

	// Example 3: Use Layer 2 specific prompt
	fmt.Println("=== Example 3: Layer 2 (Controls) Prompt ===")

	layer2Req := PromptRequest{
		PromptName: "gemara_layer_2_controls",
		Variables: map[string]interface{}{
			"task":       "Generate security controls for Kubernetes",
			"subject":    "Kubernetes cluster security",
			"technology": "Kubernetes",
			"threats":    "Unauthorized access, data exfiltration, privilege escalation",
		},
	}

	layer2Resp, err := gemaraSet.GeneratePrompt(ctx, layer2Req)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println(layer2Resp.Content)
		fmt.Println()
	}

	// Example 4: Validate layer reference
	fmt.Println("=== Example 4: Layer Reference Validation ===")

	validReq := PromptRequest{
		PromptName: "gemara_layer_validator",
		Variables: map[string]interface{}{
			"layer_reference": "layer 3",
			"task":            "Create organizational security policy",
			"subject":         "Cloud infrastructure security",
		},
	}

	validResp, err := gemaraSet.GeneratePrompt(ctx, validReq)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println(validResp.Content)
		fmt.Printf("Valid: %v\n", validResp.Metadata["valid"])
		fmt.Println()
	}

	// Example 5: Invalid layer reference (demonstrates validation)
	fmt.Println("=== Example 5: Invalid Layer Reference (Validation Error) ===")

	invalidReq := PromptRequest{
		PromptName: "gemara_layer_validator",
		Variables: map[string]interface{}{
			"layer_reference": "layer 10", // Invalid
			"task":            "Some task",
			"subject":         "Some subject",
		},
	}

	invalidResp, err := gemaraSet.GeneratePrompt(ctx, invalidReq)
	if err != nil {
		fmt.Printf("Validation caught error: %v\n", err)
	} else {
		fmt.Println(invalidResp.Content)
		fmt.Println()
	}

	// Example 6: Layer 4 Evaluation
	fmt.Println("=== Example 6: Layer 4 (Evaluation) Prompt ===")

	layer4Req := PromptRequest{
		PromptName: "gemara_layer_4_evaluation",
		Variables: map[string]interface{}{
			"task":             "Evaluate Kubernetes cluster security configuration",
			"subject":          "Kubernetes RBAC and network policies",
			"evaluation_type":  "Configuration audit",
			"related_controls": "CIS Kubernetes Benchmark controls",
		},
	}

	layer4Resp, err := gemaraSet.GeneratePrompt(ctx, layer4Req)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println(layer4Resp.Content)
		fmt.Println()
	}
}

// ExampleDeterministicLayerUsage demonstrates deterministic layer usage
func ExampleDeterministicLayerUsage() {
	fmt.Println("=== Deterministic Layer Reference Examples ===")

	testCases := []string{
		"Layer 1",
		"layer 2",
		"LAYER 3",
		"4",
		"5",
		"6",
		"Guidance",
		"Controls",
		"Policy",
		"Evaluation",
		"Enforcement",
		"Audit",
	}

	for _, testCase := range testCases {
		layerNum, err := ValidateLayerReference(testCase)
		if err != nil {
			fmt.Printf("❌ '%s' -> Error: %v\n", testCase, err)
		} else {
			layer, _ := GetLayer(layerNum)
			fmt.Printf("✅ '%s' -> Layer %d: %s\n", testCase, layerNum, layer.Name)
		}
	}

	fmt.Println()
	fmt.Println("=== Invalid Layer References ===")

	invalidCases := []string{
		"layer 0",
		"layer 7",
		"layer 10",
		"invalid",
		"",
	}

	for _, testCase := range invalidCases {
		_, err := ValidateLayerReference(testCase)
		if err != nil {
			fmt.Printf("✅ Correctly rejected '%s': %v\n", testCase, err)
		} else {
			fmt.Printf("❌ Should have rejected '%s'\n", testCase)
		}
	}
}
