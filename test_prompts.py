#!/usr/bin/env python3
"""
Enhanced test script to verify all registered prompts are available
"""

import json
import subprocess
import sys
import time
from typing import Dict, Any, Optional, List


class PromptTester:
    def __init__(self, server_path: str = "./gemara-mcp-server"):
        self.server_path = server_path
        self.process: Optional[subprocess.Popen] = None
        self.request_id = 1

    def start_server(self):
        """Start the MCP server process."""
        print(f"🚀 Starting MCP server: {self.server_path}")
        self.process = subprocess.Popen(
            [self.server_path],
            stdin=subprocess.PIPE,
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
            text=True,
            bufsize=0
        )
        time.sleep(0.5)  # Give server time to start
        print("✅ Server started\n")

    def stop_server(self):
        """Stop the MCP server process."""
        if self.process:
            self.process.terminate()
            try:
                self.process.wait(timeout=5)
            except subprocess.TimeoutExpired:
                self.process.kill()
            print("\n🛑 Server stopped")

    def send_request(
        self, method: str, params: Optional[Dict[str, Any]] = None
    ) -> Dict[str, Any]:
        """Send a JSON-RPC request to the server."""
        request = {
            "jsonrpc": "2.0",
            "id": self.request_id,
            "method": method
        }
        if params:
            request["params"] = params

        self.request_id += 1
        request_json = json.dumps(request) + "\n"

        if not self.process or not self.process.stdin:
            raise RuntimeError("Server not started")
        
        self.process.stdin.write(request_json)
        self.process.stdin.flush()

        # Read response
        response_line = self.process.stdout.readline()
        if not response_line:
            raise RuntimeError("No response from server")
        
        try:
            response = json.loads(response_line.strip())
            return response
        except json.JSONDecodeError as e:
            print(f"❌ Error parsing JSON: {e}")
            print(f"Raw response: {response_line}")
            raise

    def test_initialize(self) -> bool:
        """Test the initialize method."""
        print("=" * 70)
        print("TEST 1: Initialize Server")
        print("=" * 70)
        
        params = {
            "protocolVersion": "2024-11-05",
            "capabilities": {},
            "clientInfo": {
                "name": "test-client",
                "version": "1.0.0"
            }
        }
        
        response = self.send_request("initialize", params)
        
        if response.get("result"):
            server_info = response["result"].get("serverInfo", {})
            print("✅ Initialize successful")
            server_name = server_info.get('name', 'unknown')
            server_version = server_info.get('version', 'unknown')
            print(f"   Server: {server_name} v{server_version}")
            return True
        else:
            error = response.get('error', 'Unknown error')
            print(f"❌ Initialize failed: {error}")
            return False

    def test_list_prompts(self) -> tuple[bool, List[str]]:
        """Test listing available prompts."""
        print("\n" + "=" * 70)
        print("TEST 2: List All Prompts")
        print("=" * 70)
        
        response = self.send_request("prompts/list")
        
        if response.get("result"):
            prompts = response["result"].get("prompts", [])
            print(f"✅ Found {len(prompts)} prompt(s):\n")
            
            prompt_names = []
            for prompt in prompts:
                name = prompt.get('name', 'unknown')
                desc = prompt.get('description', 'No description')
                args = prompt.get('arguments', [])
                prompt_names.append(name)
                print(f"   📝 {name}")
                print(f"      {desc[:80]}...")
                if args:
                    arg_names = [a.get('name', '?') for a in args]
                    print(f"      Arguments: {', '.join(arg_names)}")
                print()
            
            return True, prompt_names
        else:
            print(f"❌ List prompts failed: {response.get('error', 'Unknown error')}")
            return False, []

    def test_get_prompt(
        self, prompt_name: str, arguments: Optional[Dict[str, str]] = None
    ) -> bool:
        """Test getting a specific prompt."""
        print("=" * 70)
        print(f"TEST: Get Prompt '{prompt_name}'")
        print("=" * 70)
        
        params = {"name": prompt_name}
        if arguments:
            params["arguments"] = arguments
        
        response = self.send_request("prompts/get", params)
        
        if response.get("result"):
            prompt = response["result"]
            messages = prompt.get("messages", [])
            print(f"✅ Prompt retrieved: {prompt_name}")
            print(f"   Messages: {len(messages)}")
            
            for i, msg in enumerate(messages):
                role = msg.get("role", "unknown")
                content = msg.get("content", [])
                if content:
                    # Handle both list and dict formats
                    if isinstance(content, list) and len(content) > 0:
                        if isinstance(content[0], dict):
                            text = content[0].get("text", "")
                        else:
                            text = str(content[0])
                    elif isinstance(content, dict):
                        text = content.get("text", str(content))
                    else:
                        text = str(content)
                    if text:
                        preview = text[:150].replace('\n', ' ')
                    else:
                        preview = "(empty)"
                    print(f"   [{i+1}] {role}: {preview}...")
                else:
                    print(f"   [{i+1}] {role}: (no content)")
            return True
        else:
            error = response.get('error', {})
            print(f"❌ Get prompt failed: {error.get('message', 'Unknown error')}")
            return False

    def test_key_prompts(self, prompt_names: List[str]) -> bool:
        """Test key prompts that should be registered."""
        print("\n" + "=" * 70)
        print("TEST 3: Verify Key Prompts")
        print("=" * 70)
        
        expected_prompts = [
            "gemara-system-prompt",
            "create_layer3_policy_with_layer1_mappings",
            "gemara_layer1_to_layer3_policy",
        ]
        
        all_found = True
        for expected in expected_prompts:
            if expected in prompt_names:
                print(f"✅ Found: {expected}")
            else:
                print(f"❌ Missing: {expected}")
                all_found = False
        
        return all_found

    def test_dynamic_scope(self) -> bool:
        """Test the dynamic scope prompt with different scopes."""
        print("\n" + "=" * 70)
        print("TEST 4: Test Dynamic Scope Prompt")
        print("=" * 70)
        
        prompt_name = "create_layer3_policy_with_layer1_mappings"
        
        # Test with first scope
        print("\n📋 Testing with scope: 'API Security'")
        success1 = self.test_get_prompt(prompt_name, {
            "scope": "API Security",
            "organization_context": "Test organization",
            "risk_appetite": "Moderate"
        })

        # Test with different scope
        print("\n📋 Testing with scope: 'Container Security'")
        success2 = self.test_get_prompt(prompt_name, {
            "scope": "Container Security",
            "organization_context": "Test organization",
            "risk_appetite": "Moderate"
        })
        
        return success1 and success2

    def run_all_tests(self):
        """Run all tests."""
        try:
            self.start_server()
            
            results = []
            
            # Test 1: Initialize
            results.append(("Initialize", self.test_initialize()))
            
            # Test 2: List prompts
            success, prompt_names = self.test_list_prompts()
            results.append(("List Prompts", success))
            
            # Test 3: Verify key prompts
            if success:
                results.append(("Key Prompts", self.test_key_prompts(prompt_names)))
            
            # Test 4: Test dynamic scope
            key_prompt = "create_layer3_policy_with_layer1_mappings"
            if key_prompt in prompt_names:
                results.append(("Dynamic Scope", self.test_dynamic_scope()))
            
            # Print summary
            print("\n" + "=" * 70)
            print("TEST SUMMARY")
            print("=" * 70)
            for test_name, passed in results:
                status = "✅ PASS" if passed else "❌ FAIL"
                print(f"{status}: {test_name}")
            
            total = len(results)
            passed = sum(1 for _, p in results if p)
            print(f"\n📊 Total: {passed}/{total} tests passed")
            
            return passed == total
            
        except Exception as e:
            print(f"\n❌ Test failed with error: {e}")
            import traceback
            traceback.print_exc()
            return False
        finally:
            self.stop_server()


def main():
    """Main entry point."""
    import argparse
    
    parser = argparse.ArgumentParser(
        description="Test Gemara MCP Server Prompts"
    )
    parser.add_argument(
        "--server",
        default="./gemara-mcp-server",
        help="Path to MCP server binary (default: ./gemara-mcp-server)"
    )
    
    args = parser.parse_args()
    
    tester = PromptTester(args.server)
    success = tester.run_all_tests()
    sys.exit(0 if success else 1)


if __name__ == "__main__":
    main()
