package promptsets

import (
	"fmt"
	"strings"
)

// GemaraLayer represents a layer in the Gemara GRC Engineering Model
type GemaraLayer struct {
	Number      int
	Name        string
	Description string
	Artifacts   []string
	Examples    []string
}

// GemaraLayers defines all 6 layers of the Gemara model
var GemaraLayers = map[int]*GemaraLayer{
	1: {
		Number:      1,
		Name:        "Guidance",
		Description: "High-level guidance on cybersecurity measures. Activities in this layer provide high-level rules pertaining to cybersecurity measures. Guidance is typically developed by industry groups, government agencies, or international standards bodies.",
		Artifacts: []string{
			"Guidance frameworks",
			"Industry standards",
			"International standards",
			"High-level cybersecurity rules",
			"Guideline mappings (crosswalk references)",
		},
		Examples: []string{
			"NIST Cybersecurity Framework",
			"ISO 27001",
			"PCI DSS",
			"HIPAA",
			"GDPR",
			"CRA",
		},
	},
	2: {
		Number:      2,
		Name:        "Controls",
		Description: "Technology-specific, threat-informed security controls. Controls are the specific guardrails that organizations put in place to protect their information systems. They are typically informed by the best practices and industry standards which are produced in Layer 1.",
		Artifacts: []string{
			"Technology-specific controls",
			"Threat-informed controls",
			"Control mappings to threats",
			"Control mappings to Layer 1 guidance",
			"Control catalogs",
		},
		Examples: []string{
			"CIS Benchmarks",
			"FINOS Common Cloud Controls",
			"Open Source Project Security (OSPS) Baseline",
		},
	},
	3: {
		Number:      3,
		Name:        "Policy",
		Description: "Risk-informed guidance tailored to an organization. Activities in this layer provide risk-informed governance rules that — while based on best practices and industry standards — are tailored to an organization. Policies cannot be properly developed without consideration for organization-specific risk appetite and risk-acceptance.",
		Artifacts: []string{
			"Organizational policies",
			"Risk-informed governance rules",
			"Policy documents",
			"Risk assessment documents",
		},
		Examples: []string{
			"Organization-specific security policies",
			"Risk-tailored governance rules",
			"Custom policy documents",
		},
	},
	4: {
		Number:      4,
		Name:        "Evaluation",
		Description: "Inspection of code, configurations, and deployments. Activities in this layer provide inspection of code, configurations, and deployments. Evaluation activities may be built based on outputs from layers 2 or 3. While automated assessments are often developed by vendors or industry groups, robust evaluation should be informed by organizational policies in order to custom-tailor the assessment to the needs of the compliance program.",
		Artifacts: []string{
			"Assessment results",
			"Evaluation reports",
			"Code inspection results",
			"Configuration analysis",
			"Deployment evaluations",
			"Control evaluation mappings",
		},
		Examples: []string{
			"Automated security scans",
			"Compliance assessments",
			"Code review results",
			"Configuration audits",
		},
	},
	5: {
		Number:      5,
		Name:        "Enforcement",
		Description: "Prevention or remediation based on assessment findings. Activities in this layer provide prevention or remediation. These enforcement actions should be guided by Layer 3 policies and based on assessment findings from Layer 4 evaluations. This layer ensures that the organization is complying with policy when evidence of noncompliance is found.",
		Artifacts: []string{
			"Enforcement actions",
			"Remediation plans",
			"Prevention mechanisms",
			"Policy compliance actions",
			"Automated enforcement rules",
		},
		Examples: []string{
			"Blocking non-compliant deployments",
			"Automated remediation scripts",
			"Policy enforcement gates",
			"Compliance remediation workflows",
		},
	},
	6: {
		Number:      6,
		Name:        "Audit",
		Description: "Review of organizational policy and conformance. Activities in this layer provide a review of organizational policy and conformance. Audits consider information from all of the lower layers. These activities are typically performed by internal or external auditors to ensure that the organization has designed and enforced effective policies based on the organization's requirements.",
		Artifacts: []string{
			"Audit reports",
			"Conformance reviews",
			"Policy effectiveness assessments",
			"Compliance audit findings",
			"Cross-layer audit analysis",
		},
		Examples: []string{
			"Internal compliance audits",
			"External security audits",
			"Policy effectiveness reviews",
			"Conformance assessments",
		},
	},
}

// ValidateLayerReference validates and normalizes a layer reference
// Returns the layer number and an error if invalid
func ValidateLayerReference(layerRef string) (int, error) {
	// Normalize the input
	normalized := strings.ToLower(strings.TrimSpace(layerRef))

	// Try to parse as number
	var layerNum int
	if _, err := fmt.Sscanf(normalized, "layer %d", &layerNum); err == nil {
		if layerNum >= 1 && layerNum <= 6 {
			return layerNum, nil
		}
		return 0, fmt.Errorf("layer number must be between 1 and 6, got %d", layerNum)
	}

	// Try direct number
	if _, err := fmt.Sscanf(normalized, "%d", &layerNum); err == nil {
		if layerNum >= 1 && layerNum <= 6 {
			return layerNum, nil
		}
		return 0, fmt.Errorf("layer number must be between 1 and 6, got %d", layerNum)
	}

	// Try layer name matching
	for num, layer := range GemaraLayers {
		if strings.ToLower(layer.Name) == normalized {
			return num, nil
		}
	}

	return 0, fmt.Errorf("invalid layer reference: %q. Must be 'layer 1' through 'layer 6', a number 1-6, or a layer name (Guidance, Controls, Policy, Evaluation, Enforcement, Audit)", layerRef)
}

// GetLayer returns a GemaraLayer by number
func GetLayer(layerNum int) (*GemaraLayer, error) {
	layer, exists := GemaraLayers[layerNum]
	if !exists {
		return nil, fmt.Errorf("layer %d does not exist. Valid layers are 1-6", layerNum)
	}
	return layer, nil
}

// GetLayerByReference returns a GemaraLayer by reference string
func GetLayerByReference(layerRef string) (*GemaraLayer, error) {
	layerNum, err := ValidateLayerReference(layerRef)
	if err != nil {
		return nil, err
	}
	return GetLayer(layerNum)
}

// FormatLayerContext formats layer information for use in prompts
func FormatLayerContext(layerNum int) (string, error) {
	layer, err := GetLayer(layerNum)
	if err != nil {
		return "", err
	}

	var artifactsList strings.Builder
	for i, artifact := range layer.Artifacts {
		if i > 0 {
			artifactsList.WriteString(", ")
		}
		artifactsList.WriteString(artifact)
	}

	var examplesList strings.Builder
	for i, example := range layer.Examples {
		if i > 0 {
			examplesList.WriteString(", ")
		}
		examplesList.WriteString(example)
	}

	return fmt.Sprintf(`Gemara Layer %d: %s

Description: %s

Expected Artifacts: %s

Examples: %s`,
		layer.Number,
		layer.Name,
		layer.Description,
		artifactsList.String(),
		examplesList.String(),
	), nil
}

// GetAllLayersContext returns formatted context for all layers
func GetAllLayersContext() string {
	var context strings.Builder
	context.WriteString("Gemara GRC Engineering Model - Layered Architecture\n\n")
	context.WriteString("The Gemara model organizes governance activities into 6 layers:\n\n")

	for i := 1; i <= 6; i++ {
		layerContext, _ := FormatLayerContext(i)
		context.WriteString(layerContext)
		context.WriteString("\n\n")
	}

	return context.String()
}
