#!/bin/bash
# Helper script to prepare files for commit

echo "📦 Preparing files for collaborator testing..."
echo ""

# Add documentation
echo "📚 Adding documentation files..."
git add CURSOR_SETUP.md QUICK_START_CURSOR.md QUICK_TEST_GUIDE.md
git add TEST_RESULTS.md INTEGRATION_SUMMARY.md
git add COLLABORATOR_SETUP.md COMMIT_GUIDE.md

# Add test scripts
echo "🧪 Adding test scripts..."
git add test_prompts.py test_mcp_server.py setup_cursor.sh

# Add example output
echo "📄 Adding example output..."
git add container_security_policy_layer3.yaml container_security_policy_summary.md

# Add updated configs
echo "⚙️  Adding configuration updates..."
git add .cursor/mcp.json .gitignore

echo ""
echo "✅ Files staged for commit!"
echo ""
echo "📋 Review what will be committed:"
git status --short

echo ""
echo "💡 To commit, run:"
echo "   git commit -m 'docs: Add setup guides, test scripts, and example output'"
echo ""
echo "💡 To see what's excluded (binary, cache, etc.):"
echo "   git status --ignored"

