# 🚀 Quick Start for Collaborators

Welcome! This guide will help you get the Gemara MCP Server running quickly.

## ⚡ Fastest Path

```bash
# 1. Build the server
go build -o gemara-mcp-server ./cmd/gemara-mcp-server

# 2. Run setup (configures Cursor automatically)
./setup_cursor.sh

# 3. Restart Cursor and start using prompts!
```

## 📚 Documentation

- **`QUICK_START_CURSOR.md`** - Start here! Quick setup guide
- **`CURSOR_SETUP.md`** - Detailed Cursor integration guide
- **`COLLABORATOR_SETUP.md`** - Full setup instructions
- **`TEST_RESULTS.md`** - See all 15 available prompts

## 🧪 Testing

```bash
# Test all prompts
python3 test_prompts.py

# Test basic MCP server
python3 test_mcp_server.py --build
```

## 📄 Example Output

Check out the example Layer 3 policy:
- `container_security_policy_layer3.yaml` - Full policy
- `container_security_policy_summary.md` - Summary

This was generated using the `create_layer3_policy_with_layer1_mappings` prompt!

## 🎯 Using Prompts in Cursor

Once Cursor is restarted, just ask naturally:

```
"Create a Layer 3 policy for API Security"
"What is Gemara?"
"Explain the 6-layer model"
```

The prompts work automatically! 🎉

## ❓ Need Help?

- See `COMMIT_GUIDE.md` for what files are included
- Check `TEST_RESULTS.md` for available prompts
- Review `INTEGRATION_SUMMARY.md` for technical details

