# Jacques' MCP Proxy Setup

## 🎯 What This Is

This MCP (Model Context Protocol) proxy centralizes **all your MCP servers** into a single interface that Claude can connect to. Instead of managing 15+ individual MCP server connections, Claude now only needs to connect to **one proxy server** that handles everything.

## 🏗️ Architecture

```
Claude Desktop → MCP Proxy → Your Individual MCP Servers
                     ↓
                [brave-search, github, filesystem, puppeteer, 
                 memory, sequential-thinking, git, iterm, 
                 n8n, obsidian, whatsapp, browser-tools]
```

## 📁 Files Created

- `proxy-config.json` - Configuration for all your MCP servers
- `claude-config.json` - Simplified Claude Desktop configuration  
- `start-proxy.sh` - Easy startup script
- `build/mcp-proxy` - The compiled proxy server binary

## 🚀 Quick Start

### 1. Start the Proxy Server
```bash
cd /Users/jacques/DevFolder/mcp-proxy
./start-proxy.sh
```

### 2. Update Claude Desktop Configuration

Replace your current `claude_desktop_config.json` with the contents of `claude-config.json`:

**Location**: `~/Library/Application Support/Claude/claude_desktop_config.json`

```bash
# Backup your current config
cp "~/Library/Application Support/Claude/claude_desktop_config.json" "~/Library/Application Support/Claude/claude_desktop_config.json.backup"

# Copy the new simplified config
cp "/Users/jacques/DevFolder/mcp-proxy/claude-config.json" "~/Library/Application Support/Claude/claude_desktop_config.json"
```

### 3. Restart Claude Desktop

## 🔧 Configuration Details

### Proxy Settings
- **URL**: http://localhost:9090  
- **Auth Token**: `mcp-proxy-token-2025`
- **Tool Filtering**: Enabled with allow-lists for security

### Included MCP Servers
- **brave-search**: Web search capabilities
- **github**: GitHub repository management
- **filesystem**: File system access with safe paths
- **puppeteer**: Browser automation
- **memory**: Knowledge graph and memory
- **sequential-thinking**: Step-by-step reasoning
- **git**: Git repository operations  
- **iterm**: Terminal integration
- **n8n**: Workflow automation
- **obsidian**: Note management
- **whatsapp**: WhatsApp integration
- **browser-tools**: Browser development tools

## 🔒 Security Features

- **Tool Filtering**: Each server has explicit allow-lists of permitted tools
- **Authentication**: Uses auth tokens for proxy access
- **Safe Paths**: Filesystem access restricted to approved directories only
- **Environment Isolation**: Each MCP server runs in its own process

## 🛠️ Management Commands

### Start Proxy
```bash
cd /Users/jacques/DevFolder/mcp-proxy
./start-proxy.sh
```

### Stop Proxy
Press `Ctrl+C` in the terminal running the proxy

### Check Status
Visit: http://localhost:9090/health (if health endpoint exists)

### View Logs
The proxy outputs detailed logs showing:
- Server connections
- Tool registrations  
- Request routing
- Error debugging

## 🔧 Customization

### Adding New MCP Servers
Edit `proxy-config.json` and add to the `mcpServers` section:

```json
"your-new-server": {
  "command": "npx",
  "args": ["-y", "your-mcp-package"],
  "env": {
    "API_KEY": "your-api-key"
  },
  "options": {
    "toolFilter": {
      "mode": "allow",
      "list": ["tool1", "tool2"]
    }
  }
}
```

### Removing Servers
Simply delete the server entry from `proxy-config.json`

### Updating Tool Filters
Modify the `toolFilter.list` arrays to add/remove allowed tools

## 🚨 Troubleshooting

### Proxy Won't Start
- Check if port 9090 is available: `lsof -i :9090`
- Verify all MCP commands are installed (npx, uvx, etc.)
- Check file paths in configuration

### Claude Can't Connect
- Ensure proxy is running: `./start-proxy.sh`
- Verify claude-config.json is in the right location
- Restart Claude Desktop after config changes

### Individual Servers Failing
- Check the proxy logs for specific error messages
- Verify environment variables and API keys
- Test individual servers directly if needed

### Missing Tools
- Check proxy logs for "Adding tool" messages
- Verify tool names in toolFilter lists match exactly
- Some servers may need additional setup or API keys

## 🎉 Benefits

✅ **Simplified Management**: One config file instead of 15+  
✅ **Better Security**: Centralized tool filtering and access control  
✅ **Easier Debugging**: Single log output for all servers  
✅ **Performance**: Reduced Claude startup time  
✅ **Maintainability**: Easy to add/remove servers  
✅ **Reliability**: Failed servers don't break others  

## 📞 Support

If you need to modify this setup:
1. Edit `proxy-config.json` for server changes
2. Run `./start-proxy.sh` to restart
3. Check logs for any issues
4. Update tool filters as needed

---

**Created**: May 25, 2025  
**Version**: 1.0  
**Proxy**: TBXark/mcp-proxy v1.0.0
