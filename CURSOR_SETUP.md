# Setting Up Gemara MCP Server with Cursor

## ✅ Quick Setup

The MCP server is already configured in `.cursor/mcp.json`. Just follow these steps:

### Step 1: Verify the Binary Path

The configuration points to:
```
/home/hbraswel/GIT/PSCE/Baklava/developer/gem-mcp/gemara-mcp-server/gemara-mcp-server
```

Make sure the binary exists and is executable:
```bash
ls -lh gemara-mcp-server
chmod +x gemara-mcp-server  # If needed
```

### Step 2: Restart Cursor

1. **Close Cursor completely** (not just the window)
2. **Reopen Cursor** in this workspace
3. Cursor will automatically load the MCP server configuration

### Step 3: Verify Connection

1. Open Cursor's command palette (`Cmd/Ctrl + Shift + P`)
2. Look for MCP-related commands or check the status bar
3. The server should appear as connected

You can also check Cursor's MCP logs:
- Look for "gemara-mcp-server" in the output panel
- Check for any connection errors

## 🎯 Using Prompts in Cursor

Once connected, you have **15 prompts** available! Here's how to use them:

### Available Prompts

#### System Prompt
- **`gemara-system-prompt`** - Provides system-level context about Gemara

#### User-Facing Prompts (Main ones to use!)
- **`create_layer3_policy_with_layer1_mappings`** ⭐ - **Main prompt for creating Layer 3 policies with dynamic scope**
  - Arguments: `scope`, `organization_context`, `risk_appetite`, `additional_requirements`
  
- **`analyze_layer1_guidance_for_scope`** - Analyze Layer 1 guidance for a scope
- **`generate_layer3_policy_from_guidance`** - Generate policy from guidance
- **`customize_policy_scope`** - Customize policy scope

#### Gemara Layer Prompts
- `gemara_layer1_to_layer3_policy` - Convert Layer 1 to Layer 3
- `gemara_layer_1_guidance` through `gemara_layer_6_audit` - Layer-specific prompts
- `gemara_layer_context` - Get layer context
- `gemara_layer_validator` - Validate layer references
- `gemara_architecture_explainer` - Explain the architecture

### How to Use Prompts in Cursor Chat

#### Method 1: Direct Prompt Usage

In Cursor chat, you can reference prompts directly. Cursor will automatically use the MCP prompts when relevant.

**Example conversation:**
```
You: "Create a Layer 3 policy for API Security using the create_layer3_policy_with_layer1_mappings prompt"

Cursor will:
1. Call the prompt with scope="API Security"
2. Generate the policy based on Layer 1 guidance
3. Return a Gemara-compliant Layer 3 policy
```

#### Method 2: Natural Language (Recommended)

Just ask naturally - Cursor will use the appropriate prompts:

```
You: "I need a Layer 3 policy for container security. 
     My organization has moderate risk appetite and 
     focuses on cloud-native deployments."

Cursor will:
- Use create_layer3_policy_with_layer1_mappings
- Set scope="Container Security"
- Include your organization context
- Generate the policy
```

#### Method 3: Explicit Scope Changes

You can change scope dynamically:

```
You: "Now create the same policy but for API Security instead"

Cursor will:
- Reuse the prompt with new scope
- Generate a new policy for API Security
```

### Example Use Cases

#### 1. Create a Policy for a Specific Scope

```
User: "Create a Gemara Layer 3 policy for API Security. 
      My organization is a fintech startup with moderate risk appetite."

Cursor uses: create_layer3_policy_with_layer1_mappings
Arguments:
  - scope: "API Security"
  - organization_context: "fintech startup"
  - risk_appetite: "moderate"
```

#### 2. Analyze Guidance for a Scope

```
User: "What Layer 1 guidance applies to cloud infrastructure security?"

Cursor uses: analyze_layer1_guidance_for_scope
Arguments:
  - scope: "cloud infrastructure security"
```

#### 3. Understand Gemara Layers

```
User: "Explain the difference between Layer 2 and Layer 3 in Gemara"

Cursor uses: gemara_architecture_explainer or gemara_layer_context
```

#### 4. Validate Layer References

```
User: "Is this a valid Layer 3 policy reference?"

Cursor uses: gemara_layer_validator
```

## 🔍 Verifying Prompts Are Available

To verify all prompts are registered, you can:

1. **Check Cursor's MCP panel** (if available)
2. **Use the test script:**
   ```bash
   python3 test_prompts.py
   ```
   This will list all 15 registered prompts.

3. **Ask Cursor directly:**
   ```
   "What MCP prompts are available from the gemara-mcp-server?"
   ```

## 🎨 Best Practices

### 1. Be Specific About Scope

**Good:**
```
"Create a Layer 3 policy for API Security in a healthcare organization"
```

**Better:**
```
"Create a Layer 3 policy with:
- scope: API Security
- organization_context: Healthcare organization handling PHI
- risk_appetite: Low (strict compliance requirements)"
```

### 2. Use Natural Language

Cursor will automatically map your requests to the right prompts. You don't need to mention prompt names explicitly.

### 3. Iterate on Scope

You can easily change scope:
```
"Create a policy for API Security"
"Now do the same for Container Security"
"Also create one for Database Security"
```

Each will use the same prompt with different scope values.

### 4. Combine with Context

Provide organization context for better results:
```
"Create a Layer 3 policy for cloud infrastructure. 
We're a SaaS company using AWS, have SOC 2 compliance, 
and prefer automated enforcement."
```

## 🐛 Troubleshooting

### Server Not Connecting

1. **Check binary path:**
   ```bash
   ls -lh /home/hbraswel/GIT/PSCE/Baklava/developer/gem-mcp/gemara-mcp-server/gemara-mcp-server
   ```

2. **Verify binary is executable:**
   ```bash
   chmod +x gemara-mcp-server
   ```

3. **Test server manually:**
   ```bash
   ./gemara-mcp-server
   # Should start without errors
   ```

4. **Check Cursor logs:**
   - Open Cursor's output panel
   - Look for MCP-related errors

### Prompts Not Appearing

1. **Restart Cursor completely**
2. **Verify server is running:**
   ```bash
   python3 test_prompts.py
   ```

3. **Check MCP configuration:**
   - Verify `.cursor/mcp.json` has correct path
   - Path must be absolute (not relative)

### Prompt Not Working

1. **Check prompt name:** Use exact names from the list above
2. **Verify arguments:** Some prompts require specific arguments
3. **Test with Python script:**
   ```bash
   python3 test_prompts.py
   ```

## 📚 Next Steps

1. ✅ **Server is configured** - `.cursor/mcp.json` is set up
2. ✅ **15 prompts are registered** - All available via MCP
3. 🔄 **Restart Cursor** - To load the configuration
4. 🎯 **Start using prompts** - Ask Cursor to create policies!

## Example Conversation Flow

```
You: "I need to create a security policy for our API infrastructure"

Cursor: [Uses create_layer3_policy_with_layer1_mappings]
        [Generates Layer 3 policy for API Security]

You: "Can you also create one for our container deployments?"

Cursor: [Reuses prompt with scope="Container Security"]
        [Generates new policy]

You: "What Layer 1 guidance did you use for the API policy?"

Cursor: [Uses analyze_layer1_guidance_for_scope]
        [Lists relevant Layer 1 guidance]
```

Enjoy using Gemara prompts in Cursor! 🚀

