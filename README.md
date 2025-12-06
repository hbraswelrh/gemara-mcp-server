# Gemara MCP Server

A Model Context Protocol (MCP) server for [Gemara](https://github.com/ossf/gemara) - the GRC Engineering Model for Automated Risk Assessment. This server provides tools and prompts for creating, validating, and managing Gemara artifacts (Layer 1 Guidance, Layer 2 Controls, and Layer 3 Policies).

## Overview

Gemara is a framework for representing cybersecurity guidance, controls, and policies in a structured, machine-readable format. This MCP server enables AI assistants to help users create and manage Gemara artifacts through a standardized interface.

## Features

### Tools

- **Storage & Validation**
  - `store_layer1_yaml` - Store and validate Layer 1 Guidance documents
  - `store_layer2_yaml` - Store and validate Layer 2 Control Catalogs
  - `store_layer3_yaml` - Store and validate Layer 3 Policy documents
  - `validate_gemara_yaml` - Validate YAML against Gemara schemas
  - `get_layer_schema_info` - Get schema information for any layer

- **Query & Discovery**
  - `list_layer1_guidance` - List all Layer 1 Guidance documents
  - `get_layer1_guidance` - Get detailed information about a guidance document
  - `search_layer1_guidance` - Search guidance by name, description, or author
  - `list_layer2_controls` - List Layer 2 Controls (with filtering)
  - `get_layer2_control` - Get detailed information about a control
  - `search_layer2_controls` - Search controls by name or description
  - `get_layer2_guideline_mappings` - Get Layer 1 guideline mappings for a control
  - `get_artifact_relationships` - Visualize cross-layer dependencies

- **Scoping & Applicability**
  - `find_applicable_artifacts` - Find artifacts applicable to a policy scope

- **File Loading**
  - `load_layer1_from_file` - Load Layer 1 Guidance from file
  - `load_layer2_from_file` - Load Layer 2 Control Catalog from file
  - `load_layer3_from_file` - Load Layer 3 Policy from file

### Prompts

- `gemara-system-prompt` - System-level context about Gemara
- `create-layer1-guidance` - Guide for creating Layer 1 Guidance documents
- `create-layer2-controls` - Guide for creating Layer 2 Control Catalogs
- `create-layer3-policies` - Guide for creating Layer 3 Policy documents
- `gemara-quick-start` - Quick start guide for creating your first artifacts

### Resources

- CUE schema resources for all layers (accessible via `gemara://schema/layer/{1-4}`)
- Common schema resources (base, metadata, mapping)

## Installation

```bash
go build ./cmd/gemara-mcp-server
```

## Usage

Run the MCP server:

```bash
./gemara-mcp-server
```

The server communicates via stdio and follows the [Model Context Protocol](https://modelcontextprotocol.io/) specification.

### Configuration

The server automatically looks for artifacts in an `artifacts/` directory relative to the working directory or executable location:

```
artifacts/
├── layer1/
│   └── *.yaml
├── layer2/
│   └── *.yaml
└── layer3/
    └── *.yaml
```

## Architecture

- **`mcp/`** - MCP server implementation and system prompts
- **`tools/`** - Gemara authoring tools, prompts, and utilities
- **`cmd/gemara-mcp-server/`** - Main application entry point

## License

See [LICENSE](LICENSE) file for details.

## Related Projects

- [Gemara](https://github.com/ossf/gemara) - The core Gemara framework
- [Model Context Protocol](https://modelcontextprotocol.io/) - The MCP specification
