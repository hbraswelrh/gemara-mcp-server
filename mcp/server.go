package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/format"
	"cuelang.org/go/cue/load"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/ossf/gemara/layer1"
	"github.com/ossf/gemara/layer2"
	"github.com/ossf/gemara/layer3"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Version     string
	LogFilePath string
}

type Server struct {
	server *server.MCPServer
	config Config
}

func NewServer(cfg Config) (*Server, error) {
	mcpServer := server.NewMCPServer("gemara-mcp-server", cfg.Version)

	s := &Server{
		server: mcpServer,
		config: cfg,
	}

	s.registerTools()

	return s, nil
}

func (s *Server) registerTools() {
	// Tool: validate_cue - Validate CUE files
	s.server.AddTool(mcp.Tool{
		Name:        "validate_cue",
		Description: "Validates CUE files and returns any errors found",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"files": map[string]interface{}{
					"type":        "array",
					"description": "List of CUE file paths to validate",
					"items": map[string]interface{}{
						"type": "string",
					},
				},
				"content": map[string]interface{}{
					"type":        "string",
					"description": "CUE content to validate (alternative to files)",
				},
			},
		},
	}, s.handleValidateCUE)

	// Tool: evaluate_cue - Evaluate CUE expressions
	s.server.AddTool(mcp.Tool{
		Name:        "evaluate_cue",
		Description: "Evaluates CUE expressions and returns the result",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"content": map[string]interface{}{
					"type":        "string",
					"description": "CUE content to evaluate",
				},
				"expression": map[string]interface{}{
					"type":        "string",
					"description": "Optional CUE expression path to evaluate (e.g., 'foo.bar')",
				},
			},
			Required: []string{"content"},
		},
	}, s.handleEvaluateCUE)

	// Tool: format_cue - Format CUE files
	s.server.AddTool(mcp.Tool{
		Name:        "format_cue",
		Description: "Formats CUE code according to CUE style guidelines",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"content": map[string]interface{}{
					"type":        "string",
					"description": "CUE content to format",
				},
			},
			Required: []string{"content"},
		},
	}, s.handleFormatCUE)

	// Tool: unify_cue - Unify/merge CUE configurations
	s.server.AddTool(mcp.Tool{
		Name:        "unify_cue",
		Description: "Unifies (merges) multiple CUE configurations",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"configs": map[string]interface{}{
					"type":        "array",
					"description": "Array of CUE configuration strings to unify",
					"items": map[string]interface{}{
						"type": "string",
					},
				},
			},
			Required: []string{"configs"},
		},
	}, s.handleUnifyCUE)

	// Tool: export_cue - Export CUE to JSON/YAML
	s.server.AddTool(mcp.Tool{
		Name:        "export_cue",
		Description: "Exports CUE configuration to JSON or YAML format",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"content": map[string]interface{}{
					"type":        "string",
					"description": "CUE content to export",
				},
				"format": map[string]interface{}{
					"type":        "string",
					"description": "Output format: 'json' or 'yaml' (default: 'json')",
					"enum":        []string{"json", "yaml"},
				},
				"expression": map[string]interface{}{
					"type":        "string",
					"description": "Optional CUE expression path to export",
				},
			},
			Required: []string{"content"},
		},
	}, s.handleExportCUE)

	// Tool: import_to_cue - Convert JSON/YAML to CUE
	s.server.AddTool(mcp.Tool{
		Name:        "import_to_cue",
		Description: "Converts JSON or YAML data to CUE format",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"content": map[string]interface{}{
					"type":        "string",
					"description": "JSON or YAML content to convert",
				},
				"format": map[string]interface{}{
					"type":        "string",
					"description": "Input format: 'json' or 'yaml' (default: 'json')",
					"enum":        []string{"json", "yaml"},
				},
			},
			Required: []string{"content"},
		},
	}, s.handleImportToCUE)

	// Tool: import_guidelines_by_criteria - Import Layer 1 Guidelines by technology, sector, or jurisdiction
	s.server.AddTool(mcp.Tool{
		Name:        "import_guidelines_by_criteria",
		Description: "Imports Layer 1 Guidelines (Guidance Documents) filtered by technology domain, industry sector, or jurisdiction. Returns matching guidelines in JSON format.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"file_path": map[string]interface{}{
					"type":        "string",
					"description": "Path to the Layer 1 Guidance Document file (YAML or JSON). Supports file:/// or https:// URIs.",
				},
				"technology": map[string]interface{}{
					"type":        "string",
					"description": "Filter by technology domain (e.g., 'artificial-intelligence', 'cloud-computing')",
				},
				"sector": map[string]interface{}{
					"type":        "string",
					"description": "Filter by industry sector (e.g., 'financial-services', 'healthcare')",
				},
				"jurisdiction": map[string]interface{}{
					"type":        "string",
					"description": "Filter by jurisdiction (e.g., 'US', 'EU', 'HIPAA')",
				},
			},
			Required: []string{"file_path"},
		},
	}, s.handleImportGuidelinesByCriteria)

	// Tool: import_controls_by_label - Import Layer 2 Controls filtered by labels/applicability categories
	s.server.AddTool(mcp.Tool{
		Name:        "import_controls_by_label",
		Description: "Imports Layer 2 Controls filtered by labels (applicability categories). Returns matching controls in JSON format.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"file_path": map[string]interface{}{
					"type":        "string",
					"description": "Path to the Layer 2 Control Catalog file (YAML or JSON). Supports file:/// or https:// URIs.",
				},
				"labels": map[string]interface{}{
					"type":        "array",
					"description": "Array of label IDs to filter controls by (e.g., ['tlp_clear', 'PII-Data-Protection'])",
					"items": map[string]interface{}{
						"type": "string",
					},
				},
			},
			Required: []string{"file_path", "labels"},
		},
	}, s.handleImportControlsByLabel)

	// Tool: create_layer3_control_modifiers - Create Layer 3 Control Modifiers for harmonization
	s.server.AddTool(mcp.Tool{
		Name:        "create_layer3_control_modifiers",
		Description: "Creates Layer 3 Control Modifiers to harmonize Layer 2 Controls with Layer 1 Guidelines. Analyzes controls and guidelines to generate appropriate modifiers.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"guidelines_file": map[string]interface{}{
					"type":        "string",
					"description": "Path to Layer 1 Guidance Document file (YAML or JSON)",
				},
				"controls_file": map[string]interface{}{
					"type":        "string",
					"description": "Path to Layer 2 Control Catalog file (YAML or JSON)",
				},
				"modification_rationale": map[string]interface{}{
					"type":        "string",
					"description": "Rationale for the modifications (e.g., 'Harmonize controls with HIPAA requirements')",
				},
				"modification_type": map[string]interface{}{
					"type":        "string",
					"description": "Type of modification: 'alter', 'add', 'remove' (default: 'alter')",
					"enum":        []string{"alter", "add", "remove"},
				},
			},
			Required: []string{"guidelines_file", "controls_file", "modification_rationale"},
		},
	}, s.handleCreateLayer3ControlModifiers)
}

func (s *Server) handleValidateCUE(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var params struct {
		Files   []string `json:"files,omitempty"`
		Content string   `json:"content,omitempty"`
	}

	if err := request.BindArguments(&params); err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("Error parsing arguments: %v", err),
				},
			},
			IsError: true,
		}, nil
	}

	ctx2 := cuecontext.New()
	var errors []string

	if len(params.Files) > 0 {
		// Validate files
		instances := load.Instances(params.Files, nil)
		for _, inst := range instances {
			if inst.Err != nil {
				errors = append(errors, fmt.Sprintf("load error: %v", inst.Err))
				continue
			}
			value := ctx2.BuildInstance(inst)
			if err := value.Validate(); err != nil {
				errors = append(errors, fmt.Sprintf("%s: %v", inst.Dir, err))
			}
		}
	} else if params.Content != "" {
		// Validate content string
		value := ctx2.CompileString(params.Content)
		if err := value.Err(); err != nil {
			errors = append(errors, err.Error())
		} else if err := value.Validate(); err != nil {
			errors = append(errors, err.Error())
		}
	} else {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: "Either 'files' or 'content' must be provided",
				},
			},
			IsError: true,
		}, nil
	}

	result := map[string]interface{}{
		"valid":  len(errors) == 0,
		"errors": errors,
	}

	resultJSON, _ := json.MarshalIndent(result, "", "  ")
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: string(resultJSON),
			},
		},
		StructuredContent: result,
	}, nil
}

func (s *Server) handleEvaluateCUE(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var params struct {
		Content    string `json:"content"`
		Expression string `json:"expression,omitempty"`
	}

	if err := request.BindArguments(&params); err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("Error parsing arguments: %v", err),
				},
			},
			IsError: true,
		}, nil
	}

	ctx2 := cuecontext.New()
	value := ctx2.CompileString(params.Content)

	if err := value.Err(); err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("Compilation error: %v", err),
				},
			},
			IsError: true,
		}, nil
	}

	if err := value.Validate(); err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("Validation error: %v", err),
				},
			},
			IsError: true,
		}, nil
	}

	// Evaluate specific expression if provided
	if params.Expression != "" {
		value = value.LookupPath(cue.ParsePath(params.Expression))
		if err := value.Err(); err != nil {
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					mcp.TextContent{
						Type: "text",
						Text: fmt.Sprintf("Expression error: %v", err),
					},
				},
				IsError: true,
			}, nil
		}
	}

	// Export to JSON
	var result interface{}
	if err := value.Decode(&result); err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("Export error: %v", err),
				},
			},
			IsError: true,
		}, nil
	}

	valueStr, _ := value.String()
	resultJSON, _ := json.MarshalIndent(map[string]interface{}{
		"result": result,
		"value":  valueStr,
	}, "", "  ")

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: string(resultJSON),
			},
		},
		StructuredContent: map[string]interface{}{
			"result": result,
			"value":  valueStr,
		},
	}, nil
}

func (s *Server) handleFormatCUE(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var params struct {
		Content string `json:"content"`
	}

	if err := request.BindArguments(&params); err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("Error parsing arguments: %v", err),
				},
			},
			IsError: true,
		}, nil
	}

	ctx2 := cuecontext.New()
	value := ctx2.CompileString(params.Content)

	if err := value.Err(); err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("Compilation error: %v", err),
				},
			},
			IsError: true,
		}, nil
	}

	// Format the CUE
	formatted, err := format.Node(value.Syntax())
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("Format error: %v", err),
				},
			},
			IsError: true,
		}, nil
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: string(formatted),
			},
		},
		StructuredContent: map[string]interface{}{
			"formatted": string(formatted),
		},
	}, nil
}

func (s *Server) handleUnifyCUE(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var params struct {
		Configs []string `json:"configs"`
	}

	if err := request.BindArguments(&params); err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("Error parsing arguments: %v", err),
				},
			},
			IsError: true,
		}, nil
	}

	if len(params.Configs) == 0 {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: "At least one config must be provided",
				},
			},
			IsError: true,
		}, nil
	}

	ctx2 := cuecontext.New()
	var unified cue.Value

	for i, config := range params.Configs {
		value := ctx2.CompileString(config)
		if err := value.Err(); err != nil {
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					mcp.TextContent{
						Type: "text",
						Text: fmt.Sprintf("Compilation error in config %d: %v", i+1, err),
					},
				},
				IsError: true,
			}, nil
		}

		if i == 0 {
			unified = value
		} else {
			unified = unified.Unify(value)
			if err := unified.Err(); err != nil {
				return &mcp.CallToolResult{
					Content: []mcp.Content{
						mcp.TextContent{
							Type: "text",
							Text: fmt.Sprintf("Unification error with config %d: %v", i+1, err),
						},
					},
					IsError: true,
				}, nil
			}
		}
	}

	if err := unified.Validate(); err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("Validation error: %v", err),
				},
			},
			IsError: true,
		}, nil
	}

	// Format the unified result
	formatted, err := format.Node(unified.Syntax())
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("Format error: %v", err),
				},
			},
			IsError: true,
		}, nil
	}

	// Also export as JSON
	var jsonResult interface{}
	if err := unified.Decode(&jsonResult); err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("Export error: %v", err),
				},
			},
			IsError: true,
		}, nil
	}

	unifiedStr, _ := unified.String()
	resultJSON, _ := json.MarshalIndent(map[string]interface{}{
		"unified": string(formatted),
		"json":    jsonResult,
		"value":   unifiedStr,
	}, "", "  ")

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: string(resultJSON),
			},
		},
		StructuredContent: map[string]interface{}{
			"unified": string(formatted),
			"json":    jsonResult,
			"value":   unifiedStr,
		},
	}, nil
}

func (s *Server) handleExportCUE(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var params struct {
		Content    string `json:"content"`
		Format     string `json:"format,omitempty"`
		Expression string `json:"expression,omitempty"`
	}

	if err := request.BindArguments(&params); err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("Error parsing arguments: %v", err),
				},
			},
			IsError: true,
		}, nil
	}

	if params.Format == "" {
		params.Format = "json"
	}

	ctx2 := cuecontext.New()
	value := ctx2.CompileString(params.Content)

	if err := value.Err(); err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("Compilation error: %v", err),
				},
			},
			IsError: true,
		}, nil
	}

	if err := value.Validate(); err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("Validation error: %v", err),
				},
			},
			IsError: true,
		}, nil
	}

	// Evaluate specific expression if provided
	if params.Expression != "" {
		value = value.LookupPath(cue.ParsePath(params.Expression))
		if err := value.Err(); err != nil {
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					mcp.TextContent{
						Type: "text",
						Text: fmt.Sprintf("Expression error: %v", err),
					},
				},
				IsError: true,
			}, nil
		}
	}

	var result interface{}
	if err := value.Decode(&result); err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("Export error: %v", err),
				},
			},
			IsError: true,
		}, nil
	}

	resultJSON, _ := json.MarshalIndent(result, "", "  ")

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: string(resultJSON),
			},
		},
		StructuredContent: map[string]interface{}{
			"format": params.Format,
			"data":   result,
		},
	}, nil
}

func (s *Server) handleImportToCUE(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var params struct {
		Content string `json:"content"`
		Format  string `json:"format,omitempty"`
	}

	if err := request.BindArguments(&params); err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("Error parsing arguments: %v", err),
				},
			},
			IsError: true,
		}, nil
	}

	if params.Format == "" {
		params.Format = "json"
	}

	ctx2 := cuecontext.New()
	var value cue.Value

	if params.Format == "json" {
		var data interface{}
		if err := json.Unmarshal([]byte(params.Content), &data); err != nil {
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					mcp.TextContent{
						Type: "text",
						Text: fmt.Sprintf("JSON parse error: %v", err),
					},
				},
				IsError: true,
			}, nil
		}
		value = ctx2.Encode(data)
	} else if params.Format == "yaml" {
		// For YAML, we'd need gopkg.in/yaml.v3
		// For now, try to parse as JSON first
		var data interface{}
		if err := json.Unmarshal([]byte(params.Content), &data); err != nil {
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					mcp.TextContent{
						Type: "text",
						Text: fmt.Sprintf("YAML parsing not fully supported, try JSON format: %v", err),
					},
				},
				IsError: true,
			}, nil
		}
		value = ctx2.Encode(data)
	} else {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: "Format must be 'json' or 'yaml'",
				},
			},
			IsError: true,
		}, nil
	}

	if err := value.Err(); err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("CUE encoding error: %v", err),
				},
			},
			IsError: true,
		}, nil
	}

	formatted, err := format.Node(value.Syntax())
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("Format error: %v", err),
				},
			},
			IsError: true,
		}, nil
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: string(formatted),
			},
		},
		StructuredContent: map[string]interface{}{
			"cue": string(formatted),
		},
	}, nil
}

// Helper function to load YAML/JSON files (similar to gemara's internal loaders)
func loadLayer1Document(filePath string, doc *layer1.GuidanceDocument) error {
	parsedURL, err := url.Parse(filePath)
	if err != nil {
		// If not a URL, treat as file path
		return loadLayer1FromFile(filePath, doc)
	}

	if parsedURL.Scheme == "https" || parsedURL.Scheme == "http" {
		return loadLayer1FromURL(parsedURL.String(), doc)
	}
	if parsedURL.Scheme == "file" {
		filePath = strings.TrimPrefix(parsedURL.String(), "file://")
		return loadLayer1FromFile(filePath, doc)
	}
	return fmt.Errorf("unsupported scheme: %s", parsedURL.Scheme)
}

func loadLayer1FromFile(filePath string, doc *layer1.GuidanceDocument) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	ext := filepath.Ext(filePath)
	if ext == ".json" {
		decoder := json.NewDecoder(file)
		if err := decoder.Decode(doc); err != nil {
			return fmt.Errorf("error decoding JSON: %w", err)
		}
	} else {
		decoder := yaml.NewDecoder(file)
		if err := decoder.Decode(doc); err != nil {
			return fmt.Errorf("error decoding YAML: %w", err)
		}
	}
	return nil
}

func loadLayer1FromURL(urlStr string, doc *layer1.GuidanceDocument) error {
	resp, err := http.Get(urlStr)
	if err != nil {
		return fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch URL; response status: %v", resp.Status)
	}

	// Try YAML first, then JSON
	decoder := yaml.NewDecoder(resp.Body)
	if err := decoder.Decode(doc); err != nil {
		// Reset body and try JSON
		resp.Body.Close()
		resp, err = http.Get(urlStr)
		if err != nil {
			return fmt.Errorf("failed to fetch URL for JSON: %w", err)
		}
		defer resp.Body.Close()
		decoder := json.NewDecoder(resp.Body)
		if err := decoder.Decode(doc); err != nil {
			return fmt.Errorf("error decoding JSON: %w", err)
		}
	}
	return nil
}

func (s *Server) handleImportGuidelinesByCriteria(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var params struct {
		FilePath    string `json:"file_path"`
		Technology  string `json:"technology,omitempty"`
		Sector      string `json:"sector,omitempty"`
		Jurisdiction string `json:"jurisdiction,omitempty"`
	}

	if err := request.BindArguments(&params); err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("Error parsing arguments: %v", err),
				},
			},
			IsError: true,
		}, nil
	}

	// Load the guidance document
	doc := &layer1.GuidanceDocument{}
	if err := loadLayer1Document(params.FilePath, doc); err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("Error loading guidance document: %v", err),
				},
			},
			IsError: true,
		}, nil
	}

	// Filter by criteria if provided
	var filteredCategories []layer1.Category
	for _, category := range doc.Categories {
		var matchingGuidelines []layer1.Guideline
		for _, guideline := range category.Guidelines {
			// Check if document applicability matches criteria
			matches := true
			if doc.Metadata.Applicability != nil {
				if params.Technology != "" {
					found := false
					for _, tech := range doc.Metadata.Applicability.TechnologyDomains {
						if strings.EqualFold(tech, params.Technology) {
							found = true
							break
						}
					}
					if !found {
						matches = false
					}
				}
				if params.Sector != "" {
					found := false
					for _, sec := range doc.Metadata.Applicability.IndustrySectors {
						if strings.EqualFold(sec, params.Sector) {
							found = true
							break
						}
					}
					if !found {
						matches = false
					}
				}
				if params.Jurisdiction != "" {
					found := false
					for _, jur := range doc.Metadata.Applicability.Jurisdictions {
						if strings.EqualFold(jur, params.Jurisdiction) {
							found = true
							break
						}
					}
					if !found {
						matches = false
					}
				}
			} else if params.Technology != "" || params.Sector != "" || params.Jurisdiction != "" {
				// If filters are specified but document has no applicability, skip
				matches = false
			}

			if matches {
				matchingGuidelines = append(matchingGuidelines, guideline)
			}
		}

		if len(matchingGuidelines) > 0 {
			filteredCategory := category
			filteredCategory.Guidelines = matchingGuidelines
			filteredCategories = append(filteredCategories, filteredCategory)
		}
	}

	// Create filtered document
	filteredDoc := *doc
	filteredDoc.Categories = filteredCategories

	// Convert to JSON
	resultJSON, err := json.MarshalIndent(filteredDoc, "", "  ")
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("Error marshaling to JSON: %v", err),
				},
			},
			IsError: true,
		}, nil
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: string(resultJSON),
			},
		},
		StructuredContent: filteredDoc,
	}, nil
}

func (s *Server) handleImportControlsByLabel(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var params struct {
		FilePath string   `json:"file_path"`
		Labels   []string `json:"labels"`
	}

	if err := request.BindArguments(&params); err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("Error parsing arguments: %v", err),
				},
			},
			IsError: true,
		}, nil
	}

	if len(params.Labels) == 0 {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: "At least one label must be provided",
				},
			},
			IsError: true,
		}, nil
	}

	// Load the control catalog
	catalog := &layer2.Catalog{}
	if err := catalog.LoadFile(params.FilePath); err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("Error loading control catalog: %v", err),
				},
			},
			IsError: true,
		}, nil
	}

	// Create a map of label IDs for quick lookup
	labelMap := make(map[string]bool)
	for _, label := range params.Labels {
		labelMap[strings.ToLower(label)] = true
	}

	// Filter controls by labels
	var filteredFamilies []layer2.ControlFamily
	for _, family := range catalog.ControlFamilies {
		var matchingControls []layer2.Control
		for _, control := range family.Controls {
			// Check assessment requirements for matching labels
			matches := false
			for _, req := range control.AssessmentRequirements {
				for _, app := range req.Applicability {
					if labelMap[strings.ToLower(app)] {
						matches = true
						break
					}
				}
				if matches {
					break
				}
			}

			// Also check metadata applicability categories if available
			if !matches && len(catalog.Metadata.ApplicabilityCategories) > 0 {
				for _, cat := range catalog.Metadata.ApplicabilityCategories {
					if labelMap[strings.ToLower(cat.Id)] {
						matches = true
						break
					}
				}
			}

			if matches {
				matchingControls = append(matchingControls, control)
			}
		}

		if len(matchingControls) > 0 {
			filteredFamily := family
			filteredFamily.Controls = matchingControls
			filteredFamilies = append(filteredFamilies, filteredFamily)
		}
	}

	// Create filtered catalog
	filteredCatalog := *catalog
	filteredCatalog.ControlFamilies = filteredFamilies

	// Convert to JSON
	resultJSON, err := json.MarshalIndent(filteredCatalog, "", "  ")
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("Error marshaling to JSON: %v", err),
				},
			},
			IsError: true,
		}, nil
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: string(resultJSON),
			},
		},
		StructuredContent: filteredCatalog,
	}, nil
}

func (s *Server) handleCreateLayer3ControlModifiers(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var params struct {
		GuidelinesFile        string `json:"guidelines_file"`
		ControlsFile          string `json:"controls_file"`
		ModificationRationale string `json:"modification_rationale"`
		ModificationType      string `json:"modification_type,omitempty"`
	}

	if err := request.BindArguments(&params); err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("Error parsing arguments: %v", err),
				},
			},
			IsError: true,
		}, nil
	}

	if params.ModificationType == "" {
		params.ModificationType = "alter"
	}

	// Load guidelines
	guidelinesDoc := &layer1.GuidanceDocument{}
	if err := loadLayer1Document(params.GuidelinesFile, guidelinesDoc); err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("Error loading guidelines document: %v", err),
				},
			},
			IsError: true,
		}, nil
	}

	// Load controls
	controlsCatalog := &layer2.Catalog{}
	if err := controlsCatalog.LoadFile(params.ControlsFile); err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("Error loading controls catalog: %v", err),
				},
			},
			IsError: true,
		}, nil
	}

	// Create a map of guideline IDs for quick lookup
	guidelineMap := make(map[string]layer1.Guideline)
	for _, category := range guidelinesDoc.Categories {
		for _, guideline := range category.Guidelines {
			guidelineMap[guideline.Id] = guideline
		}
	}

	// Generate control modifiers by analyzing controls and their guideline mappings
	var controlModifiers []layer3.ControlModifier
	for _, family := range controlsCatalog.ControlFamilies {
		for _, control := range family.Controls {
			// Check if control has guideline mappings
			for _, mapping := range control.GuidelineMappings {
				// For each guideline mapping, check if we need to create a modifier
				for _, entry := range mapping.Entries {
					// Check if the guideline exists in our loaded guidelines
					if _, exists := guidelineMap[entry.ReferenceId]; exists {
						// Create a modifier for this control
						modifier := layer3.ControlModifier{
							TargetId:                control.Id,
							ModType:                 layer3.ModType(params.ModificationType),
							ModificationRationale:   params.ModificationRationale,
							Title:                   control.Title,
							Objective:               control.Objective,
						}
						controlModifiers = append(controlModifiers, modifier)
						break // Only create one modifier per control
					}
				}
			}
		}
	}

	// Create a Layer 3 mapping structure with the modifiers
	mapping := layer3.Mapping{
		ReferenceId:            controlsCatalog.Metadata.Id,
		ControlModifications:    controlModifiers,
		InScope:                layer3.Scope{},
		OutOfScope:             layer3.Scope{},
		AssessmentRequirementModifications: []layer3.AssessmentRequirementModifier{},
		GuidelineModifications:             []layer3.GuidelineModifier{},
	}

	// Convert to JSON
	resultJSON, err := json.MarshalIndent(mapping, "", "  ")
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("Error marshaling to JSON: %v", err),
				},
			},
			IsError: true,
		}, nil
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: string(resultJSON),
			},
		},
		StructuredContent: mapping,
	}, nil
}

func (s *Server) Start() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	stdioServer := server.NewStdioServer(s.server)

	var slogHandler slog.Handler
	var logOutput io.Writer
	if s.config.LogFilePath != "" {
		file, err := os.OpenFile(s.config.LogFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
		if err != nil {
			return fmt.Errorf("failed to open log file: %w", err)
		}
		logOutput = file
		slogHandler = slog.NewTextHandler(logOutput, &slog.HandlerOptions{Level: slog.LevelDebug})
	} else {
		logOutput = os.Stderr
		slogHandler = slog.NewTextHandler(logOutput, &slog.HandlerOptions{Level: slog.LevelInfo})
	}
	logger := slog.New(slogHandler)
	logger.Info("starting Gemara CUE MCP server", "version", s.config.Version)
	stdLogger := log.New(logOutput, "gemara-mcp-server: ", 0)
	stdioServer.SetErrorLogger(stdLogger)

	errC := make(chan error, 1)
	go func() {
		errC <- stdioServer.Listen(ctx, os.Stdin, os.Stdout)
	}()

	_, _ = fmt.Fprintf(os.Stderr, "Gemara CUE MCP Server running on stdio\n")

	select {
	case <-ctx.Done():
		logger.Info("shutting down server", "signal", "context done")
	case err := <-errC:
		if err != nil {
			logger.Error("error running server", "error", err)
			return fmt.Errorf("error running server: %w", err)
		}
	}

	return nil
}
