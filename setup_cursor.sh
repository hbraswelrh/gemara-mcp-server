#!/bin/bash
# Setup script for Cursor MCP integration

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BINARY_PATH="$SCRIPT_DIR/gemara-mcp-server"
CURSOR_CONFIG="$SCRIPT_DIR/.cursor/mcp.json"

echo "🔧 Setting up Gemara MCP Server for Cursor"
echo ""

# Check if binary exists
if [ ! -f "$BINARY_PATH" ]; then
    echo "❌ Binary not found. Building..."
    cd "$SCRIPT_DIR"
    go build -o gemara-mcp-server ./cmd/gemara-mcp-server
    if [ $? -ne 0 ]; then
        echo "❌ Build failed!"
        exit 1
    fi
    echo "✅ Binary built successfully"
fi

# Make sure binary is executable
chmod +x "$BINARY_PATH"
echo "✅ Binary is executable: $BINARY_PATH"

# Create .cursor directory if it doesn't exist
mkdir -p "$SCRIPT_DIR/.cursor"

# Create/update mcp.json
cat > "$CURSOR_CONFIG" << EOF
{
  "mcpServers": {
    "gemara-mcp-server": {
      "command": "$BINARY_PATH",
      "args": [],
      "env": {}
    }
  }
}
EOF

echo "✅ MCP configuration updated: $CURSOR_CONFIG"
echo ""
echo "📋 Configuration:"
cat "$CURSOR_CONFIG"
echo ""
echo "🎯 Next steps:"
echo "1. Restart Cursor completely (close and reopen)"
echo "2. The MCP server will automatically connect"
echo "3. Start using prompts in Cursor chat!"
echo ""
echo "📚 See CURSOR_SETUP.md for detailed usage instructions"

