# Quick Start: Using Gemara MCP Server in Cursor

## 🚀 Setup (One-Time)

### Option 1: Automatic Setup (Recommended)
```bash
./setup_cursor.sh
```

This will:
- Build the server if needed
- Update `.cursor/mcp.json` with the correct path
- Verify everything is ready

### Option 2: Manual Setup

1. **Update `.cursor/mcp.json`** with your absolute path:
   ```json
   {
     "mcpServers": {
       "gemara-mcp-server": {
         "command": "/home/hbraswel/GIT/PSCE/Baklava/developer/gem-mcp/gemara-mcp-server/gemara-mcp-server",
         "args": [],
         "env": {}
       }
     }
   }
   ```

2. **Restart Cursor completely**

## ✅ Verify It's Working

After restarting Cursor, test it:

```
Ask Cursor: "What is Gemara?"
```

If the MCP server is connected, Cursor will use the `gemara-system-prompt` to provide context.

## 🎯 Using Your Prompts

### Main Prompt: Create Layer 3 Policy

**Try this in Cursor chat:**

```
Create a Layer 3 policy for API Security. 
My organization is a fintech company with moderate risk appetite.
```

Cursor will automatically:
1. Use `create_layer3_policy_with_layer1_mappings` prompt
2. Set scope="API Security"
3. Generate a Gemara-compliant Layer 3 policy

### Change Scope Dynamically

```
Now create the same policy but for Container Security instead
```

The prompt will be reused with the new scope!

### Other Useful Prompts

```
"What Layer 1 guidance applies to cloud security?"
→ Uses: analyze_layer1_guidance_for_scope

"Explain the Gemara 6-layer model"
→ Uses: gemara_architecture_explainer

"Validate this Layer 3 policy reference"
→ Uses: gemara_layer_validator
```

## 📋 All Available Prompts

You have **15 prompts** registered:

1. `gemara-system-prompt` - System context
2. `create_layer3_policy_with_layer1_mappings` ⭐ - Main user prompt
3. `analyze_layer1_guidance_for_scope`
4. `generate_layer3_policy_from_guidance`
5. `customize_policy_scope`
6. `gemara_layer1_to_layer3_policy`
7. `gemara_layer_1_guidance` through `gemara_layer_6_audit` (6 prompts)
8. `gemara_layer_context`
9. `gemara_layer_validator`
10. `gemara_architecture_explainer`

## 🐛 Troubleshooting

**Server not connecting?**
- Run `./setup_cursor.sh` again
- Check path in `.cursor/mcp.json` is absolute
- Restart Cursor completely

**Prompts not working?**
- Verify server is running: `python3 test_prompts.py`
- Check Cursor's MCP logs in output panel

**Need help?**
- See `CURSOR_SETUP.md` for detailed guide
- See `TEST_RESULTS.md` for test information

## 🎉 You're Ready!

Just restart Cursor and start asking for policies. The prompts will work automatically!

