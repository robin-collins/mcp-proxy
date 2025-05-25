#!/bin/bash

# MCP Proxy Startup Script
# This script starts the MCP Proxy server with your configuration

cd /Users/jacques/DevFolder/mcp-proxy

echo "🚀 Starting MCP Proxy Server..."
echo "Configuration: proxy-config.json"
echo "Listening on: http://localhost:9090"
echo ""

# Start the proxy server
./build/mcp-proxy --config proxy-config.json

echo "✅ MCP Proxy Server stopped"
