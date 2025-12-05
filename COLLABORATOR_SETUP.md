# Setup Guide for Collaborators

## Quick Start

1. **Clone and build:**
   ```bash
   git clone <repo-url>
   cd gemara-mcp-server
   go build -o gemara-mcp-server ./cmd/gemara-mcp-server
   ```

2. **Run setup script:**
   ```bash
   ./setup_cursor.sh
   ```

3. **Restart Cursor** and start using prompts!

## Example Output

See the example Layer 3 policy that was generated:
- `container_security_policy_layer3.yaml` - Full Gemara Layer 3 policy
- `container_security_policy_summary.md` - Executive summary

These were generated using the `create_layer3_policy_with_layer1_mappings` prompt with scope="Container Security".

## Testing

```bash
# Test all prompts
python3 test_prompts.py

# Test basic MCP functionality
python3 test_mcp_server.py --build
```

## Documentation

- `QUICK_START_CURSOR.md` - Quick start guide
- `CURSOR_SETUP.md` - Detailed Cursor setup
- `TEST_RESULTS.md` - Test results and available prompts
- `INTEGRATION_SUMMARY.md` - Technical integration details

## Available Prompts

15 prompts are registered, including:
- `create_layer3_policy_with_layer1_mappings` ⭐ - Main prompt with dynamic scope
- `gemara-system-prompt` - System context
- All Gemara layer prompts (1-6)

See `TEST_RESULTS.md` for the complete list.

