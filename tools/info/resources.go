package info

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/mark3labs/mcp-go/mcp"
)

// handleBaseSchemaResource returns the base.cue schema content
func (g *GemaraInfoTools) handleBaseSchemaResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	schemaContent, err := g.getCommonCUESchema("base.cue")
	if err != nil {
		return nil, fmt.Errorf("failed to load base schema: %w", err)
	}
	return []mcp.ResourceContents{
		&mcp.TextResourceContents{
			URI:      request.Params.URI,
			MIMEType: "text/x-cue",
			Text:     schemaContent,
		},
	}, nil
}

// handleMetadataSchemaResource returns the metadata.cue schema content
func (g *GemaraInfoTools) handleMetadataSchemaResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	schemaContent, err := g.getCommonCUESchema("metadata.cue")
	if err != nil {
		return nil, fmt.Errorf("failed to load metadata schema: %w", err)
	}
	return []mcp.ResourceContents{
		&mcp.TextResourceContents{
			URI:      request.Params.URI,
			MIMEType: "text/x-cue",
			Text:     schemaContent,
		},
	}, nil
}

// handleMappingSchemaResource returns the mapping.cue schema content
func (g *GemaraInfoTools) handleMappingSchemaResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	schemaContent, err := g.getCommonCUESchema("mapping.cue")
	if err != nil {
		return nil, fmt.Errorf("failed to load mapping schema: %w", err)
	}
	return []mcp.ResourceContents{
		&mcp.TextResourceContents{
			URI:      request.Params.URI,
			MIMEType: "text/x-cue",
			Text:     schemaContent,
		},
	}, nil
}

// handleLayer1SchemaResource returns the layer-1.cue schema content
func (g *GemaraInfoTools) handleLayer1SchemaResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	schemaContent, err := g.getCUESchema(1)
	if err != nil {
		return nil, fmt.Errorf("failed to load layer 1 schema: %w", err)
	}
	return []mcp.ResourceContents{
		&mcp.TextResourceContents{
			URI:      request.Params.URI,
			MIMEType: "text/x-cue",
			Text:     schemaContent,
		},
	}, nil
}

// handleLayer2SchemaResource returns the layer-2.cue schema content
func (g *GemaraInfoTools) handleLayer2SchemaResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	schemaContent, err := g.getCUESchema(2)
	if err != nil {
		return nil, fmt.Errorf("failed to load layer 2 schema: %w", err)
	}
	return []mcp.ResourceContents{
		&mcp.TextResourceContents{
			URI:      request.Params.URI,
			MIMEType: "text/x-cue",
			Text:     schemaContent,
		},
	}, nil
}

// handleLayer3SchemaResource returns the layer-3.cue schema content
func (g *GemaraInfoTools) handleLayer3SchemaResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	schemaContent, err := g.getCUESchema(3)
	if err != nil {
		return nil, fmt.Errorf("failed to load layer 3 schema: %w", err)
	}
	return []mcp.ResourceContents{
		&mcp.TextResourceContents{
			URI:      request.Params.URI,
			MIMEType: "text/x-cue",
			Text:     schemaContent,
		},
	}, nil
}

// handleLayer4SchemaResource returns the layer-4.cue schema content
func (g *GemaraInfoTools) handleLayer4SchemaResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	schemaContent, err := g.getCUESchema(4)
	if err != nil {
		return nil, fmt.Errorf("failed to load layer 4 schema: %w", err)
	}
	return []mcp.ResourceContents{
		&mcp.TextResourceContents{
			URI:      request.Params.URI,
			MIMEType: "text/x-cue",
			Text:     schemaContent,
		},
	}, nil
}

// handleLayer5SchemaResource returns the layer-5.cue schema content
func (g *GemaraInfoTools) handleLayer5SchemaResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	schemaContent, err := g.getCUESchema(5)
	if err != nil {
		return nil, fmt.Errorf("failed to load layer 5 schema: %w", err)
	}
	return []mcp.ResourceContents{
		&mcp.TextResourceContents{
			URI:      request.Params.URI,
			MIMEType: "text/x-cue",
			Text:     schemaContent,
		},
	}, nil
}

// handleLayer6SchemaResource returns the layer-6.cue schema content
func (g *GemaraInfoTools) handleLayer6SchemaResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	schemaContent, err := g.getCUESchema(6)
	if err != nil {
		return nil, fmt.Errorf("failed to load layer 6 schema: %w", err)
	}
	return []mcp.ResourceContents{
		&mcp.TextResourceContents{
			URI:      request.Params.URI,
			MIMEType: "text/x-cue",
			Text:     schemaContent,
		},
	}, nil
}

// getSchemaResourceContent retrieves schema content via resource URI
// This provides a consistent way for tools to access schemas through the resource system
func (g *GemaraInfoTools) getSchemaResourceContent(uri string) (string, error) {
	// Find the resource handler for this URI
	var handler func(context.Context, mcp.ReadResourceRequest) ([]mcp.ResourceContents, error)

	switch uri {
	case "gemara://schema/common/base":
		handler = g.handleBaseSchemaResource
	case "gemara://schema/common/metadata":
		handler = g.handleMetadataSchemaResource
	case "gemara://schema/common/mapping":
		handler = g.handleMappingSchemaResource
	case "gemara://schema/layer/1":
		handler = g.handleLayer1SchemaResource
	case "gemara://schema/layer/2":
		handler = g.handleLayer2SchemaResource
	case "gemara://schema/layer/3":
		handler = g.handleLayer3SchemaResource
	case "gemara://schema/layer/4":
		handler = g.handleLayer4SchemaResource
	case "gemara://schema/layer/5":
		handler = g.handleLayer5SchemaResource
	case "gemara://schema/layer/6":
		handler = g.handleLayer6SchemaResource
	default:
		return "", fmt.Errorf("unknown schema resource URI: %s", uri)
	}

	// Call the resource handler
	request := mcp.ReadResourceRequest{
		Params: mcp.ReadResourceParams{
			URI: uri,
		},
	}

	contents, err := handler(context.Background(), request)
	if err != nil {
		return "", fmt.Errorf("failed to load schema resource %s: %w", uri, err)
	}

	if len(contents) == 0 {
		return "", fmt.Errorf("no content returned for schema resource %s", uri)
	}

	// Extract text content
	if textContent, ok := contents[0].(*mcp.TextResourceContents); ok {
		return textContent.Text, nil
	}

	return "", fmt.Errorf("unexpected content type for schema resource %s", uri)
}

// getLayerSchemaResourceURI returns the resource URI for a given layer
func (g *GemaraInfoTools) getLayerSchemaResourceURI(layer int) string {
	return fmt.Sprintf("gemara://schema/layer/%d", layer)
}

// getCommonSchemaResourceURI returns the resource URI for a common schema
func (g *GemaraInfoTools) getCommonSchemaResourceURI(schemaName string) string {
	return fmt.Sprintf("gemara://schema/common/%s", schemaName)
}

// getCUESchema fetches or returns cached CUE schema for a layer
func (g *GemaraInfoTools) getCUESchema(layer int) (string, error) {
	// Create version-aware cache key
	cacheKey := fmt.Sprintf("%s:layer:%d", g.schemaVersion, layer)
	
	// Check cache first
	if schema, ok := g.schemaCache[cacheKey]; ok {
		return schema, nil
	}

	// Fetch schema from GitHub using the configured version
	schemaURL := fmt.Sprintf("https://raw.githubusercontent.com/ossf/gemara/%s/schemas/layer-%d.cue", g.schemaVersion, layer)

	resp, err := http.Get(schemaURL)
	if err != nil {
		return "", fmt.Errorf("failed to fetch schema from %s: %w", schemaURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch schema: HTTP %d", resp.StatusCode)
	}

	schemaBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read schema: %w", err)
	}

	schemaContent := string(schemaBytes)

	// Cache the schema with version-aware key
	g.schemaCache[cacheKey] = schemaContent

	return schemaContent, nil
}

// getCommonCUESchema fetches or returns cached common CUE schema files
func (g *GemaraInfoTools) getCommonCUESchema(schemaName string) (string, error) {
	// Create version-aware cache key
	cacheKey := fmt.Sprintf("%s:common:%s", g.schemaVersion, schemaName)

	// Check cache first
	if schema, ok := g.schemaCache[cacheKey]; ok {
		return schema, nil
	}

	// Fetch schema from GitHub using the configured version
	schemaURL := fmt.Sprintf("https://raw.githubusercontent.com/ossf/gemara/%s/schemas/%s", g.schemaVersion, schemaName)

	resp, err := http.Get(schemaURL)
	if err != nil {
		return "", fmt.Errorf("failed to fetch schema from %s: %w", schemaURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch schema: HTTP %d", resp.StatusCode)
	}

	schemaBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read schema: %w", err)
	}

	schemaContent := string(schemaBytes)

	// Cache the schema
	g.schemaCache[cacheKey] = schemaContent

	return schemaContent, nil
}

// handleLexiconResource returns the Gemara lexicon content
func (g *GemaraInfoTools) handleLexiconResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	// Check cache first
	cacheKey := "lexicon:gemara.openssf.org"
	if lexicon, ok := g.schemaCache[cacheKey]; ok {
		return []mcp.ResourceContents{
			&mcp.TextResourceContents{
				URI:      request.Params.URI,
				MIMEType: "text/html",
				Text:     lexicon,
			},
		}, nil
	}

	// Fetch lexicon from Gemara website
	lexiconURL := "https://gemara.openssf.org/lexicon.html"
	resp, err := http.Get(lexiconURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch lexicon from %s: %w", lexiconURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch lexicon: HTTP %d", resp.StatusCode)
	}

	lexiconBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read lexicon: %w", err)
	}

	lexiconContent := string(lexiconBytes)

	// Cache the lexicon
	g.schemaCache[cacheKey] = lexiconContent

	return []mcp.ResourceContents{
		&mcp.TextResourceContents{
			URI:      request.Params.URI,
			MIMEType: "text/html",
			Text:     lexiconContent,
		},
	}, nil
}
