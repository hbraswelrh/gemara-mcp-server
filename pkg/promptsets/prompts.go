package promptsets

import (
	"context"
	"fmt"
	"strings"
)

// SystemPromptDefinitions contains predefined system prompts organized by category
var SystemPromptDefinitions = map[string][]*Prompt{
	"code_generation": {
		{
			Name:        "code_generator",
			Description: "A system prompt for generating code based on requirements",
			Content: `You are an expert software developer. Your task is to generate high-quality, production-ready code based on the following requirements:

Requirements:
{{requirements}}

Context:
- Programming Language: {{language}}
- Framework/Library: {{framework}}
- Code Style: {{style}}
- Additional Constraints: {{constraints}}

Please generate code that:
1. Follows best practices and design patterns
2. Includes proper error handling
3. Has clear, meaningful variable and function names
4. Includes appropriate comments and documentation
5. Is efficient and maintainable
6. Adheres to the specified style guide

Generate the code now:`,
			Variables: map[string]string{
				"requirements": "The functional requirements for the code",
				"language":     "The programming language to use",
				"framework":    "The framework or library to use (optional)",
				"style":        "The coding style to follow",
				"constraints":  "Any additional constraints or requirements",
			},
			Category: "code_generation",
			Tags:     []string{"code", "generation", "development"},
		},
		{
			Name:        "code_reviewer",
			Description: "A system prompt for reviewing and improving code",
			Content: `You are an expert code reviewer. Analyze the following code and provide comprehensive feedback:

Code to Review:
{{code}}

Review Criteria:
- Code Quality: {{quality_focus}}
- Performance: {{performance_focus}}
- Security: {{security_focus}}
- Maintainability: {{maintainability_focus}}

Please provide:
1. Overall assessment
2. Specific issues found (with line numbers if applicable)
3. Suggestions for improvement
4. Best practices recommendations
5. Security concerns (if any)
6. Performance optimizations (if applicable)

Review the code:`,
			Variables: map[string]string{
				"code":                  "The code to review",
				"quality_focus":         "Areas to focus on for code quality",
				"performance_focus":     "Performance considerations",
				"security_focus":        "Security considerations",
				"maintainability_focus": "Maintainability considerations",
			},
			Category: "code_generation",
			Tags:     []string{"code", "review", "quality"},
		},
	},
	"documentation": {
		{
			Name:        "documentation_generator",
			Description: "A system prompt for generating documentation",
			Content: `You are a technical writer specializing in software documentation. Generate comprehensive documentation for the following:

Subject: {{subject}}
Type: {{doc_type}}
Audience: {{audience}}
Format: {{format}}

Content to Document:
{{content}}

Please create documentation that:
1. Is clear and easy to understand for the target audience
2. Includes examples where appropriate
3. Follows the specified format and style
4. Covers all important aspects
5. Is well-organized with proper headings and structure

Generate the documentation:`,
			Variables: map[string]string{
				"subject":  "The subject to document",
				"doc_type": "Type of documentation (API, user guide, README, etc.)",
				"audience": "Target audience (developers, users, etc.)",
				"format":   "Documentation format (markdown, HTML, etc.)",
				"content":  "The content or code to document",
			},
			Category: "documentation",
			Tags:     []string{"documentation", "writing", "technical"},
		},
	},
	"analysis": {
		{
			Name:        "code_analyzer",
			Description: "A system prompt for analyzing code structure and patterns",
			Content: `You are a code analysis expert. Analyze the following codebase and provide insights:

Codebase:
{{codebase}}

Analysis Focus:
- Architecture: {{architecture_focus}}
- Patterns: {{patterns_focus}}
- Dependencies: {{dependencies_focus}}
- Complexity: {{complexity_focus}}

Please provide:
1. Architecture overview
2. Design patterns identified
3. Dependency analysis
4. Complexity metrics and observations
5. Potential improvements
6. Risk assessment

Analyze the codebase:`,
			Variables: map[string]string{
				"codebase":           "The code or codebase to analyze",
				"architecture_focus": "Focus areas for architecture analysis",
				"patterns_focus":     "Focus areas for pattern analysis",
				"dependencies_focus": "Focus areas for dependency analysis",
				"complexity_focus":   "Focus areas for complexity analysis",
			},
			Category: "analysis",
			Tags:     []string{"analysis", "code", "architecture"},
		},
		{
			Name:        "security_analyzer",
			Description: "A system prompt for security analysis",
			Content: `You are a security expert. Perform a security analysis of the following:

Target: {{target}}
Type: {{target_type}}
Context: {{context}}

Security Focus Areas:
- Vulnerabilities: {{vulnerability_focus}}
- Best Practices: {{best_practices_focus}}
- Compliance: {{compliance_focus}}
- Threat Model: {{threat_model_focus}}

Please analyze and provide:
1. Security vulnerabilities identified
2. Risk assessment for each vulnerability
3. Recommendations for remediation
4. Security best practices to follow
5. Compliance considerations
6. Threat model analysis

Perform security analysis:`,
			Variables: map[string]string{
				"target":               "The code, system, or component to analyze",
				"target_type":          "Type of target (code, API, system, etc.)",
				"context":              "Additional context about the target",
				"vulnerability_focus":  "Areas to focus on for vulnerabilities",
				"best_practices_focus": "Security best practices to check",
				"compliance_focus":     "Compliance requirements to check",
				"threat_model_focus":   "Threat modeling considerations",
			},
			Category: "analysis",
			Tags:     []string{"security", "analysis", "vulnerability"},
		},
	},
	"testing": {
		{
			Name:        "test_generator",
			Description: "A system prompt for generating test cases",
			Content: `You are a testing expert. Generate comprehensive test cases for the following:

Code/Functionality to Test:
{{code}}

Testing Requirements:
- Test Type: {{test_type}}
- Coverage Goals: {{coverage_goals}}
- Framework: {{framework}}
- Test Style: {{test_style}}

Please generate:
1. Unit tests for individual components
2. Integration tests where applicable
3. Edge cases and boundary conditions
4. Error handling tests
5. Performance tests (if applicable)
6. Test data and fixtures

Generate the test cases:`,
			Variables: map[string]string{
				"code":           "The code or functionality to test",
				"test_type":      "Type of tests (unit, integration, e2e, etc.)",
				"coverage_goals": "Code coverage goals",
				"framework":      "Testing framework to use",
				"test_style":     "Test style (TDD, BDD, etc.)",
			},
			Category: "testing",
			Tags:     []string{"testing", "test", "quality"},
		},
	},
	"general": {
		{
			Name:        "task_executor",
			Description: "A general-purpose system prompt for task execution",
			Content: `You are an AI assistant helping with the following task:

Task: {{task}}
Context: {{context}}
Constraints: {{constraints}}
Expected Output: {{expected_output}}

Please:
1. Understand the task requirements
2. Consider the provided context
3. Work within the specified constraints
4. Deliver the expected output format
5. Provide clear explanations of your approach

Execute the task:`,
			Variables: map[string]string{
				"task":            "The task to execute",
				"context":         "Additional context for the task",
				"constraints":     "Constraints to work within",
				"expected_output": "Expected format or type of output",
			},
			Category: "general",
			Tags:     []string{"general", "task", "assistant"},
		},
		{
			Name:        "problem_solver",
			Description: "A system prompt for problem-solving and debugging",
			Content: `You are a problem-solving expert. Help solve the following problem:

Problem Description:
{{problem}}

Current State:
{{current_state}}

Attempted Solutions:
{{attempted_solutions}}

Error Messages/Issues:
{{errors}}

Please:
1. Analyze the problem thoroughly
2. Identify root causes
3. Suggest multiple solution approaches
4. Provide step-by-step resolution
5. Explain why the solution works
6. Suggest preventive measures

Solve the problem:`,
			Variables: map[string]string{
				"problem":             "Description of the problem",
				"current_state":       "Current state of the system/code",
				"attempted_solutions": "Solutions already attempted",
				"errors":              "Error messages or issues encountered",
			},
			Category: "general",
			Tags:     []string{"problem", "debugging", "solving"},
		},
	},
	"gemara": {
		{
			Name:        "gemara_architecture_explainer",
			Description: "Explains the Gemara GRC Engineering Model layered architecture",
			Content: `You are an expert in the OpenSSF Gemara Project (GRC Engineering Model for Automated Risk Assessment). 

{{gemara_layers_context}}

Key Principles:
1. Each layer builds upon the lower layers
2. Higher-level use cases may bring multiple lower layers together
3. Layers have specific artifacts and expected outputs
4. Layer references must be deterministic and validated

When working with Gemara layers:
- Always validate layer references (Layer 1-6, or layer names)
- Ensure artifacts match the expected outputs for the specified layer
- Understand the relationships between layers
- Use appropriate layer context for the task at hand

Answer questions about the Gemara architecture:`,
			Variables: map[string]string{
				"gemara_layers_context": "The complete Gemara layers context (automatically populated)",
			},
			Category: "gemara",
			Tags:     []string{"gemara", "architecture", "grc", "compliance"},
		},
		{
			Name:        "gemara_layer_context",
			Description: "Provides context for a specific Gemara layer",
			Content: `You are working with the OpenSSF Gemara Project. The following context is for {{layer_reference}}:

{{layer_context}}

Task: {{task}}

When working with {{layer_reference}}:
1. Ensure all artifacts align with the expected outputs for this layer
2. Reference appropriate lower layers when needed
3. Follow the layer-specific guidelines and examples
4. Maintain consistency with Gemara model principles

{{additional_context}}

Proceed with the task:`,
			Variables: map[string]string{
				"layer_reference":    "The Gemara layer reference (e.g., 'Layer 1', 'layer 2', 'Guidance', etc.)",
				"layer_context":      "The formatted context for the specified layer (automatically populated)",
				"task":               "The specific task to perform",
				"additional_context": "Any additional context or requirements",
			},
			Category: "gemara",
			Tags:     []string{"gemara", "layer", "context"},
		},
		{
			Name:        "gemara_layer_1_guidance",
			Description: "System prompt for working with Gemara Layer 1 (Guidance)",
			Content: `You are working with Gemara Layer 1: Guidance.

Layer 1 provides high-level guidance on cybersecurity measures. Activities in this layer produce high-level rules pertaining to cybersecurity measures. Guidance is typically developed by industry groups, government agencies, or international standards bodies.

Expected Artifacts:
- Guidance frameworks
- Industry standards
- International standards
- High-level cybersecurity rules
- Guideline mappings (crosswalk references)

Examples: NIST Cybersecurity Framework, ISO 27001, PCI DSS, HIPAA, GDPR, CRA

Task: {{task}}

Subject: {{subject}}

When working with Layer 1:
1. Focus on high-level, abstract guidance
2. Reference industry standards and frameworks
3. Create mappings between different guidance documents when applicable
4. Ensure guidance is technology-agnostic and broadly applicable
5. Layer 1 controls are high-level and abstract

Generate Layer 1 guidance content:`,
			Variables: map[string]string{
				"task":    "The specific task related to Layer 1 Guidance",
				"subject": "The subject matter for the guidance",
			},
			Category: "gemara",
			Tags:     []string{"gemara", "layer1", "guidance", "standards"},
		},
		{
			Name:        "gemara_layer_2_controls",
			Description: "System prompt for working with Gemara Layer 2 (Controls)",
			Content: `You are working with Gemara Layer 2: Controls.

Layer 2 produces technology-specific, threat-informed security controls. Controls are the specific guardrails that organizations put in place to protect their information systems. They are typically informed by the best practices and industry standards which are produced in Layer 1.

Expected Artifacts:
- Technology-specific controls
- Threat-informed controls
- Control mappings to threats
- Control mappings to Layer 1 guidance
- Control catalogs

Examples: CIS Benchmarks, FINOS Common Cloud Controls, Open Source Project Security (OSPS) Baseline

Task: {{task}}

Subject: {{subject}}

Technology: {{technology}}

Threats: {{threats}}

When working with Layer 2:
1. Develop technology-specific controls (not abstract)
2. Map controls to threats and vulnerabilities
3. Reference Layer 1 guidance that informs the controls
4. Ensure controls are actionable and implementable
5. Consider technology capabilities and threats to those capabilities

Generate Layer 2 control content:`,
			Variables: map[string]string{
				"task":       "The specific task related to Layer 2 Controls",
				"subject":    "The subject matter for the controls",
				"technology": "The specific technology the controls apply to",
				"threats":    "The threats the controls mitigate",
			},
			Category: "gemara",
			Tags:     []string{"gemara", "layer2", "controls", "security"},
		},
		{
			Name:        "gemara_layer_3_policy",
			Description: "System prompt for working with Gemara Layer 3 (Policy)",
			Content: `You are working with Gemara Layer 3: Policy.

Layer 3 provides risk-informed governance rules that — while based on best practices and industry standards — are tailored to an organization. Policies cannot be properly developed without consideration for organization-specific risk appetite and risk-acceptance.

Expected Artifacts:
- Organizational policies
- Risk-informed governance rules
- Policy documents
- Risk assessment documents

Task: {{task}}

Subject: {{subject}}

Organization Context: {{organization_context}}

Risk Appetite: {{risk_appetite}}

When working with Layer 3:
1. Tailor policies to the specific organization
2. Consider organizational risk appetite and risk acceptance
3. Base policies on Layer 1 guidance and Layer 2 controls
4. Ensure policies are risk-informed, not just compliance-driven
5. Policies should be actionable and enforceable

Generate Layer 3 policy content:`,
			Variables: map[string]string{
				"task":                 "The specific task related to Layer 3 Policy",
				"subject":              "The subject matter for the policy",
				"organization_context": "Context about the organization",
				"risk_appetite":        "The organization's risk appetite",
			},
			Category: "gemara",
			Tags:     []string{"gemara", "layer3", "policy", "governance"},
		},
		{
			Name:        "gemara_layer_4_evaluation",
			Description: "System prompt for working with Gemara Layer 4 (Evaluation)",
			Content: `You are working with Gemara Layer 4: Evaluation.

Layer 4 provides inspection of code, configurations, and deployments. Evaluation activities may be built based on outputs from layers 2 or 3. While automated assessments are often developed by vendors or industry groups, robust evaluation should be informed by organizational policies in order to custom-tailor the assessment to the needs of the compliance program.

Expected Artifacts:
- Assessment results
- Evaluation reports
- Code inspection results
- Configuration analysis
- Deployment evaluations
- Control evaluation mappings

Examples: Automated security scans, compliance assessments, code review results, configuration audits

Task: {{task}}

Subject: {{subject}}

Evaluation Type: {{evaluation_type}}

Related Controls: {{related_controls}}

When working with Layer 4:
1. Inspect code, configurations, and deployments
2. Map evaluations to Layer 2 controls or Layer 3 policies
3. Provide actionable findings and recommendations
4. Ensure evaluations are tailored to organizational needs
5. Consider both automated and manual evaluation methods

Generate Layer 4 evaluation content:`,
			Variables: map[string]string{
				"task":             "The specific task related to Layer 4 Evaluation",
				"subject":          "The subject matter to evaluate",
				"evaluation_type":  "Type of evaluation (code review, config audit, security scan, etc.)",
				"related_controls": "Related Layer 2 controls or Layer 3 policies",
			},
			Category: "gemara",
			Tags:     []string{"gemara", "layer4", "evaluation", "assessment"},
		},
		{
			Name:        "gemara_layer_5_enforcement",
			Description: "System prompt for working with Gemara Layer 5 (Enforcement)",
			Content: `You are working with Gemara Layer 5: Enforcement.

Layer 5 provides prevention or remediation. These enforcement actions should be guided by Layer 3 policies and based on assessment findings from Layer 4 evaluations. This layer ensures that the organization is complying with policy when evidence of noncompliance is found.

Expected Artifacts:
- Enforcement actions
- Remediation plans
- Prevention mechanisms
- Policy compliance actions
- Automated enforcement rules

Examples: Blocking non-compliant deployments, automated remediation scripts, policy enforcement gates, compliance remediation workflows

Task: {{task}}

Subject: {{subject}}

Evaluation Findings: {{evaluation_findings}}

Related Policies: {{related_policies}}

When working with Layer 5:
1. Base enforcement on Layer 4 evaluation findings
2. Guide actions by Layer 3 policies
3. Provide both prevention and remediation mechanisms
4. Ensure enforcement is automated where possible
5. Document enforcement actions and outcomes

Generate Layer 5 enforcement content:`,
			Variables: map[string]string{
				"task":                "The specific task related to Layer 5 Enforcement",
				"subject":             "The subject matter for enforcement",
				"evaluation_findings": "Findings from Layer 4 evaluations",
				"related_policies":    "Related Layer 3 policies",
			},
			Category: "gemara",
			Tags:     []string{"gemara", "layer5", "enforcement", "remediation"},
		},
		{
			Name:        "gemara_layer_6_audit",
			Description: "System prompt for working with Gemara Layer 6 (Audit)",
			Content: `You are working with Gemara Layer 6: Audit.

Layer 6 provides a review of organizational policy and conformance. Audits consider information from all of the lower layers. These activities are typically performed by internal or external auditors to ensure that the organization has designed and enforced effective policies based on the organization's requirements.

Expected Artifacts:
- Audit reports
- Conformance reviews
- Policy effectiveness assessments
- Compliance audit findings
- Cross-layer audit analysis

Examples: Internal compliance audits, external security audits, policy effectiveness reviews, conformance assessments

Task: {{task}}

Subject: {{subject}}

Audit Scope: {{audit_scope}}

Layers to Review: {{layers_to_review}}

When working with Layer 6:
1. Review information from all relevant lower layers
2. Assess policy effectiveness and conformance
3. Provide comprehensive audit findings
4. Consider both design and enforcement aspects
5. Ensure audits are thorough and objective

Generate Layer 6 audit content:`,
			Variables: map[string]string{
				"task":             "The specific task related to Layer 6 Audit",
				"subject":          "The subject matter for the audit",
				"audit_scope":      "The scope of the audit",
				"layers_to_review": "Which lower layers to review (1-5)",
			},
			Category: "gemara",
			Tags:     []string{"gemara", "layer6", "audit", "conformance"},
		},
		{
			Name:        "gemara_layer1_to_layer3_policy",
			Description: "Generates Layer 3 policy from Layer 1 guidance with configurable scope",
			Content: `You are a GRC expert working with the Gemara model. Your task is to generate a Layer 3 (Policy) document based on Layer 1 (Guidance) input.

Gemara Context:
- Layer 1 (Guidance): High-level guidance on cybersecurity measures from industry groups, government agencies, or international standards bodies
- Layer 3 (Policy): Risk-informed governance rules tailored to an organization, based on Layer 1 guidance and Layer 2 controls

Policy Generation Task:
Generate a Layer 3 policy document for the specified scope based on the provided Layer 1 guidance.

Scope: {{scope}}
The scope defines the domain, technology, or area this policy will cover (e.g., "Cloud Infrastructure Security", "API Security", "Container Security", "Data Protection").

Layer 1 Guidance Source: {{layer1_guidance_source}}
The source of Layer 1 guidance (e.g., "NIST Cybersecurity Framework", "ISO 27001", "PCI DSS", "CIS Benchmarks").

Layer 1 Guidance Content:
{{layer1_guidance_content}}
The specific Layer 1 guidance, standards, or requirements that should inform this policy.

Organization Context: {{organization_context}}
Information about the organization that will help tailor the policy (e.g., industry, size, regulatory requirements, existing security posture).

Risk Appetite: {{risk_appetite}}
The organization's risk tolerance and risk acceptance criteria (e.g., "Low risk tolerance", "Moderate risk appetite", "Risk-averse").

Additional Requirements: {{additional_requirements}}
Any specific requirements, constraints, or considerations for this policy.

When generating the Layer 3 policy:
1. Base the policy on the provided Layer 1 guidance, ensuring alignment with industry standards
2. Tailor the policy to the organization's context and risk appetite
3. Make the policy actionable and enforceable (not just high-level guidance)
4. Ensure the policy is risk-informed, considering the organization's specific risk profile
5. Structure the policy clearly with:
   - Policy statement and objectives
   - Scope and applicability
   - Roles and responsibilities
   - Requirements and controls
   - Compliance and enforcement mechanisms
   - Review and update procedures
6. Reference the Layer 1 guidance sources appropriately
7. Ensure the policy addresses the specified scope comprehensively

Generate the Layer 3 policy document:`,
			Variables: map[string]string{
				"scope":                   "The scope/domain for the policy (e.g., 'Cloud Infrastructure Security', 'API Security')",
				"layer1_guidance_source":  "The source of Layer 1 guidance (e.g., 'NIST Cybersecurity Framework', 'ISO 27001')",
				"layer1_guidance_content": "The specific Layer 1 guidance content to base the policy on",
				"organization_context":    "Context about the organization (industry, size, regulatory requirements)",
				"risk_appetite":           "The organization's risk appetite and tolerance",
				"additional_requirements": "Any additional requirements or constraints",
			},
			Category: "gemara",
			Tags:     []string{"gemara", "layer1", "layer3", "policy", "generation"},
		},
		{
			Name:        "gemara_layer_validator",
			Description: "Validates and provides context for Gemara layer references",
			Content: `You are a Gemara layer validation expert. Your task is to ensure proper usage of Gemara layer references.

Layer Reference: {{layer_reference}}

{{layer_validation_result}}

Task: {{task}}

Subject: {{subject}}

When validating layer usage:
1. Ensure the layer reference is valid (Layer 1-6, or layer name)
2. Verify that artifacts match the expected outputs for the specified layer
3. Check that layer relationships are correct (higher layers build on lower layers)
4. Confirm that the task is appropriate for the specified layer
5. Provide guidance if layer misuse is detected

{{additional_guidance}}

Validate and proceed:`,
			Variables: map[string]string{
				"layer_reference":         "The layer reference to validate (e.g., 'Layer 1', 'layer 2', 'Guidance')",
				"layer_validation_result": "The validation result and layer context (automatically populated)",
				"task":                    "The task being performed",
				"subject":                 "The subject matter",
				"additional_guidance":     "Additional guidance or warnings",
			},
			Category: "gemara",
			Tags:     []string{"gemara", "validation", "layer"},
		},
	},
	"user_facing": {
		{
			Name:        "create_layer3_policy_with_layer1_mappings",
			Description: "Create a Gemara Layer 3 policy conforming to the schema, scoped for a specific domain, and gather all applicable Layer 1 guidance mappings",
			Content: `You are a GRC expert working with the OpenSSF Gemara Project. Your task is to create a comprehensive Layer 3 (Policy) document that:

1. Conforms to the Gemara Layer 3 schema
2. Is scoped for: {{scope}}
3. Incorporates all relevant Layer 1 (Guidance) mappings that apply to the specified scope

Task: Create an OpenSSF Gemara Project Layer 3 policy that conforms to the schema and is scoped for {{scope}}, and gather all Gemara Layer 1 guidance mappings that apply to the specified scope.

Scope: {{scope}}
The scope defines the domain, technology, or area this policy will cover. Examples: "Cloud Infrastructure Security", "API Security", "Container Security", "Data Protection", "Network Security", "Identity and Access Management", "Incident Response", etc.

Organization Context: {{organization_context}}
Information about the organization to tailor the policy appropriately.

Risk Appetite: {{risk_appetite}}
The organization's risk tolerance and risk acceptance criteria.

Additional Requirements: {{additional_requirements}}
Any specific requirements, constraints, or considerations.

Instructions:
1. **Gather Layer 1 Guidance Mappings**: 
   - Identify all relevant Layer 1 guidance sources that apply to the scope "{{scope}}"
   - Common Layer 1 sources include: NIST Cybersecurity Framework, ISO 27001, PCI DSS, HIPAA, GDPR, CRA, CIS Benchmarks, etc.
   - For each relevant guidance source, identify the specific controls, requirements, or recommendations that apply to the scope
   - Document these mappings clearly, showing how Layer 1 guidance informs the Layer 3 policy

2. **Create Layer 3 Policy Document**:
   - Generate a complete Layer 3 policy document that conforms to the Gemara Layer 3 schema
   - The policy must be risk-informed and tailored to the organization
   - Structure the policy with:
     * Policy identifier and metadata
     * Policy statement and objectives
     * Scope and applicability (must match "{{scope}}")
     * Roles and responsibilities
     * Risk-informed requirements and controls
     * Compliance and enforcement mechanisms
     * Review and update procedures
     * References to Layer 1 guidance sources

3. **Schema Conformance**:
   - Ensure the policy document conforms to the Gemara Layer 3 schema structure
   - Include all required schema fields
   - Use proper schema formatting and structure

4. **Layer 1 Guidance Integration**:
   - Explicitly reference the Layer 1 guidance sources you identified
   - Show how each Layer 1 guidance element maps to Layer 3 policy requirements
   - Document the relationship between Layer 1 guidance and Layer 3 policy elements

5. **Output Format**:
   - Provide the Layer 3 policy document in a format that conforms to the Gemara Layer 3 schema
   - Include a separate section documenting all Layer 1 guidance mappings
   - Format mappings clearly showing: Layer 1 Source → Layer 1 Control/Guidance → Layer 3 Policy Requirement

Begin by gathering all applicable Layer 1 guidance mappings for scope "{{scope}}", then create the Layer 3 policy document.`,
			Variables: map[string]string{
				"scope":                   "The scope/domain for the policy (e.g., 'Cloud Infrastructure Security', 'API Security') - this is the key variable users can change",
				"organization_context":    "Context about the organization (industry, size, regulatory requirements)",
				"risk_appetite":           "The organization's risk appetite and tolerance",
				"additional_requirements": "Any additional requirements or constraints",
			},
			Category: "user_facing",
			Tags:     []string{"gemara", "layer1", "layer3", "policy", "user-facing", "chatbot", "schema"},
		},
		{
			Name:        "analyze_layer1_guidance_for_scope",
			Description: "Analyze and gather all Layer 1 guidance mappings that apply to a specific scope",
			Content: `You are a GRC expert analyzing Layer 1 (Guidance) sources for applicability to a specific scope.

Scope: {{scope}}

Task: Analyze all relevant Layer 1 guidance sources and identify which controls, requirements, or recommendations apply to the scope "{{scope}}".

Layer 1 Guidance Sources to Consider:
- NIST Cybersecurity Framework
- ISO 27001
- PCI DSS
- HIPAA
- GDPR
- CRA (Cyber Resilience Act)
- CIS Benchmarks
- Other relevant industry standards

For each Layer 1 guidance source:
1. Identify which controls/requirements apply to "{{scope}}"
2. Extract the specific guidance text or control descriptions
3. Note the relevance and applicability to the scope
4. Document any cross-references or mappings between different Layer 1 sources

Output Format:
- List each Layer 1 guidance source
- For each source, list applicable controls/requirements
- Include the specific guidance text
- Note how each applies to the scope
- Identify any gaps or areas where multiple sources overlap

Analyze Layer 1 guidance for scope "{{scope}}":`,
			Variables: map[string]string{
				"scope": "The scope/domain to analyze (e.g., 'Cloud Infrastructure Security', 'API Security')",
			},
			Category: "user_facing",
			Tags:     []string{"gemara", "layer1", "analysis", "user-facing", "chatbot"},
		},
		{
			Name:        "generate_layer3_policy_from_guidance",
			Description: "Generate a Layer 3 policy document from provided Layer 1 guidance, conforming to Gemara schema",
			Content: `You are a GRC expert creating a Gemara Layer 3 (Policy) document.

Task: Generate a complete Layer 3 policy document that:
1. Conforms to the Gemara Layer 3 schema
2. Is based on the provided Layer 1 guidance
3. Is scoped for: {{scope}}
4. Is tailored to the organization's context and risk appetite

Scope: {{scope}}

Layer 1 Guidance Mappings:
{{layer1_guidance_mappings}}
The Layer 1 guidance sources and their applicable controls/requirements that should inform this policy.

Organization Context: {{organization_context}}

Risk Appetite: {{risk_appetite}}

Policy Requirements:
1. **Schema Conformance**: The policy must conform to the Gemara Layer 3 schema structure
2. **Risk-Informed**: Policy must be tailored based on organizational risk appetite
3. **Actionable**: Policy must be specific and enforceable, not just high-level guidance
4. **Comprehensive**: Policy must address all aspects of the scope "{{scope}}"
5. **Traceable**: Policy must reference the Layer 1 guidance sources it's based on

Policy Structure (per Gemara Layer 3 schema):
- Policy metadata (identifier, version, date)
- Policy statement and objectives
- Scope and applicability
- Roles and responsibilities
- Risk-informed requirements
- Controls and safeguards
- Compliance mechanisms
- Enforcement procedures
- Review and update procedures
- References to Layer 1 guidance

Generate the Layer 3 policy document:`,
			Variables: map[string]string{
				"scope":                    "The scope/domain for the policy",
				"layer1_guidance_mappings": "The Layer 1 guidance mappings that apply to this scope",
				"organization_context":     "Context about the organization",
				"risk_appetite":            "The organization's risk appetite",
			},
			Category: "user_facing",
			Tags:     []string{"gemara", "layer3", "policy", "user-facing", "chatbot", "schema"},
		},
		{
			Name:        "customize_policy_scope",
			Description: "A user-friendly prompt for customizing policy scope via chatbot interface",
			Content: `You are helping a user create a Gemara Layer 3 policy with a customizable scope.

Current Scope: {{scope}}
The user can change this scope to generate policies for different domains.

User Request: {{user_request}}
The user's request or question about the policy.

Available Scopes (examples):
- Cloud Infrastructure Security
- API Security
- Container Security
- Data Protection and Privacy
- Network Security
- Identity and Access Management
- Incident Response
- Secure Software Development
- Supply Chain Security
- Or any other security domain

Instructions:
1. Understand the user's request
2. If the user wants to change scope, acknowledge the new scope
3. Explain what the policy will cover for the specified scope
4. Ask for any additional context needed (organization, risk appetite, etc.)
5. Proceed with policy generation when ready

Help the user with their request regarding scope "{{scope}}":`,
			Variables: map[string]string{
				"scope":        "The current or desired scope for the policy",
				"user_request": "The user's request or question",
			},
			Category: "user_facing",
			Tags:     []string{"gemara", "user-facing", "chatbot", "interactive"},
		},
	},
}

// GetDefaultPromptSets creates and returns default prompt sets organized by category
func GetDefaultPromptSets() *PromptSetGroup {
	group := NewPromptSetGroup()

	// Create prompt sets for each category
	for category, prompts := range SystemPromptDefinitions {
		ps := NewPromptSet(category, fmt.Sprintf("Prompts for %s tasks", category))
		for _, prompt := range prompts {
			_ = ps.AddPrompt(prompt)
		}
		_ = group.AddPromptSet(ps)
	}

	return group
}

// GetGemaraPromptSets creates and returns Gemara-specific prompt sets with layer validation
func GetGemaraPromptSets() *PromptSetGroup {
	group := NewPromptSetGroup()

	// Create Gemara prompt set with handlers for dynamic layer context
	gemaraSet := NewPromptSet("gemara", "Gemara GRC Engineering Model prompts")

	// Add all Gemara prompts
	for _, prompt := range SystemPromptDefinitions["gemara"] {
		_ = gemaraSet.AddPrompt(prompt)
	}

	// Add handler for gemara_architecture_explainer to auto-populate layers context
	_ = gemaraSet.AddHandler("gemara_architecture_explainer", func(ctx context.Context, req PromptRequest) (PromptResponse, error) {
		prompt, err := gemaraSet.GetPrompt("gemara_architecture_explainer")
		if err != nil {
			return PromptResponse{}, err
		}

		content := prompt.Content
		content = strings.ReplaceAll(content, "{{gemara_layers_context}}", GetAllLayersContext())

		// Handle other variables
		for key, value := range req.Variables {
			placeholder := fmt.Sprintf("{{%s}}", key)
			content = strings.ReplaceAll(content, placeholder, fmt.Sprintf("%v", value))
		}

		return PromptResponse{
			Content: content,
			Metadata: map[string]interface{}{
				"prompt_name": prompt.Name,
				"category":    prompt.Category,
			},
		}, nil
	})

	// Add handler for gemara_layer_context to validate and populate layer context
	_ = gemaraSet.AddHandler("gemara_layer_context", func(ctx context.Context, req PromptRequest) (PromptResponse, error) {
		prompt, err := gemaraSet.GetPrompt("gemara_layer_context")
		if err != nil {
			return PromptResponse{}, err
		}

		layerRef, ok := req.Variables["layer_reference"].(string)
		if !ok {
			return PromptResponse{}, fmt.Errorf("layer_reference is required and must be a string")
		}

		// Validate layer reference
		layerNum, err := ValidateLayerReference(layerRef)
		if err != nil {
			return PromptResponse{}, fmt.Errorf("invalid layer reference: %w", err)
		}

		// Get layer context
		layerContext, err := FormatLayerContext(layerNum)
		if err != nil {
			return PromptResponse{}, err
		}

		content := prompt.Content
		content = strings.ReplaceAll(content, "{{layer_reference}}", fmt.Sprintf("Layer %d: %s", layerNum, GemaraLayers[layerNum].Name))
		content = strings.ReplaceAll(content, "{{layer_context}}", layerContext)

		// Handle other variables
		for key, value := range req.Variables {
			if key != "layer_reference" {
				placeholder := fmt.Sprintf("{{%s}}", key)
				content = strings.ReplaceAll(content, placeholder, fmt.Sprintf("%v", value))
			}
		}

		return PromptResponse{
			Content: content,
			Metadata: map[string]interface{}{
				"prompt_name":  prompt.Name,
				"category":     prompt.Category,
				"layer_number": layerNum,
				"layer_name":   GemaraLayers[layerNum].Name,
			},
		}, nil
	})

	// Add handler for gemara_layer_validator
	_ = gemaraSet.AddHandler("gemara_layer_validator", func(ctx context.Context, req PromptRequest) (PromptResponse, error) {
		prompt, err := gemaraSet.GetPrompt("gemara_layer_validator")
		if err != nil {
			return PromptResponse{}, err
		}

		layerRef, ok := req.Variables["layer_reference"].(string)
		if !ok {
			return PromptResponse{}, fmt.Errorf("layer_reference is required and must be a string")
		}

		// Validate layer reference
		layerNum, err := ValidateLayerReference(layerRef)
		var validationResult string
		var additionalGuidance string

		if err != nil {
			validationResult = fmt.Sprintf("❌ INVALID LAYER REFERENCE: %v\n\nValid layer references:\n- 'Layer 1' through 'Layer 6'\n- Numbers 1-6\n- Layer names: Guidance, Controls, Policy, Evaluation, Enforcement, Audit", err)
			additionalGuidance = "Please correct the layer reference before proceeding."
		} else {
			layer, _ := GetLayer(layerNum)
			layerContext, _ := FormatLayerContext(layerNum)
			validationResult = fmt.Sprintf("✅ VALID LAYER REFERENCE\n\n%s", layerContext)
			additionalGuidance = fmt.Sprintf("Ensure your task and artifacts align with Layer %d (%s) expectations.", layerNum, layer.Name)
		}

		content := prompt.Content
		content = strings.ReplaceAll(content, "{{layer_reference}}", layerRef)
		content = strings.ReplaceAll(content, "{{layer_validation_result}}", validationResult)
		content = strings.ReplaceAll(content, "{{additional_guidance}}", additionalGuidance)

		// Handle other variables
		for key, value := range req.Variables {
			if key != "layer_reference" {
				placeholder := fmt.Sprintf("{{%s}}", key)
				content = strings.ReplaceAll(content, placeholder, fmt.Sprintf("%v", value))
			}
		}

		return PromptResponse{
			Content: content,
			Metadata: map[string]interface{}{
				"prompt_name":  prompt.Name,
				"category":     prompt.Category,
				"layer_number": layerNum,
				"valid":        err == nil,
			},
		}, nil
	})

	_ = group.AddPromptSet(gemaraSet)
	return group
}

// GetUserFacingPromptSets creates and returns user-facing prompts designed for chatbot/agent interfaces
func GetUserFacingPromptSets() *PromptSetGroup {
	group := NewPromptSetGroup()

	// Create user-facing prompt set
	userFacingSet := NewPromptSet("user_facing", "User-facing prompts for chatbot/agent interfaces")

	// Add all user-facing prompts
	for _, prompt := range SystemPromptDefinitions["user_facing"] {
		_ = userFacingSet.AddPrompt(prompt)
	}

	// Add handler for create_layer3_policy_with_layer1_mappings
	// This handler can be extended to actually fetch Layer 1 mappings from a database or API
	_ = userFacingSet.AddHandler("create_layer3_policy_with_layer1_mappings", func(ctx context.Context, req PromptRequest) (PromptResponse, error) {
		prompt, err := userFacingSet.GetPrompt("create_layer3_policy_with_layer1_mappings")
		if err != nil {
			return PromptResponse{}, err
		}

		scope, ok := req.Variables["scope"].(string)
		if !ok || scope == "" {
			return PromptResponse{}, fmt.Errorf("scope is required and must be a non-empty string")
		}

		content := prompt.Content

		// Replace all variables
		for key, value := range req.Variables {
			placeholder := fmt.Sprintf("{{%s}}", key)
			content = strings.ReplaceAll(content, placeholder, fmt.Sprintf("%v", value))
		}

		return PromptResponse{
			Content: content,
			Metadata: map[string]interface{}{
				"prompt_name": prompt.Name,
				"category":    prompt.Category,
				"scope":       scope,
				"user_facing": true,
			},
		}, nil
	})

	_ = group.AddPromptSet(userFacingSet)
	return group
}

// GetAllPromptSets returns both Gemara and user-facing prompt sets
func GetAllPromptSets() *PromptSetGroup {
	group := NewPromptSetGroup()

	// Add Gemara prompts
	gemaraGroup := GetGemaraPromptSets()
	for _, ps := range gemaraGroup.ListPromptSets() {
		_ = group.AddPromptSet(ps)
	}

	// Add user-facing prompts
	userFacingGroup := GetUserFacingPromptSets()
	for _, ps := range userFacingGroup.ListPromptSets() {
		_ = group.AddPromptSet(ps)
	}

	return group
}
