# MCP Proxy Setup

A centralized MCP (Model Context Protocol) proxy server configuration for managing multiple MCP servers through a single interface. This setup simplifies Claude Desktop configuration by aggregating all MCP servers behind one proxy.

## 🎯 Purpose

Instead of managing 15+ individual MCP server connections in Claude Desktop, this proxy allows you to:
- Connect to **one proxy server** that handles everything
- Centralize configuration and logging
- Implement security controls and tool filtering
- Easily add/remove MCP servers without changing Claude config

## 🏗️ Architecture

```
Claude Desktop → MCP Proxy → Your Individual MCP Servers
                     ↓
                [brave-search, github, filesystem, puppeteer,
                 memory, sequential-thinking, git, browser-tools, etc.]
```

## 🚀 Quick Start

### 1. Prerequisites

- Go 1.19+ (for building the proxy)
- Node.js and npm (for most MCP servers)
- uv/uvx (for Python-based MCP servers)

### 2. Clone and Build
## Features

- **Proxy Multiple MCP Clients**: Connects to multiple MCP resource servers and aggregates their tools and capabilities.
- **SSE / HTTP Streaming MCP Support**: Provides an SSE (Server-Sent Events) or HTTP streaming interface for real-time updates from MCP clients.
- **Flexible Configuration**: Supports multiple client types (`stdio`, `sse` or `streamable-http`) with customizable settings.

## Documentation

- Configuration: [docs/configuration.md](docs/CONFIGURATION.md)
- Usage: [docs/usage.md](docs/USAGE.md)
- Deployment: [docs/deployment.md](docs/DEPLOYMENT.md)
- Claude config converter: https://tbxark.github.io/mcp-proxy

```bash
git clone https://github.com/TBXark/mcp-proxy.git
cd mcp-proxy
make build
```

### Install via Go

```bash
go install github.com/TBXark/mcp-proxy@latest
```

### Docker

The image includes support for launching MCP servers via `npx` and `uvx`.

```bash
docker run -d -p 9090:9090 -v /path/to/config.json:/config/config.json ghcr.io/tbxark/mcp-proxy:latest
# or provide a remote config
docker run -d -p 9090:9090 ghcr.io/tbxark/mcp-proxy:latest --config https://example.com/config.json
```

More deployment options (including docker‑compose) are in [docs/deployment.md](docs/DEPLOYMENT.md).

## Configuration

See full configuration reference and examples in [docs/configuration.md](docs/CONFIGURATION.md).
An online Claude config converter is available at: https://tbxark.github.io/mcp-proxy

### 3. Configure

```bash
# Copy example configuration
cp proxy-config.example.json proxy-config.json
cp claude-config.example.json claude-config.json

# Edit proxy-config.json with your actual:
# - API keys
# - File paths
# - Environment variables
```

### 4. Update Claude Desktop

Replace your Claude Desktop configuration:

**On macOS:**
```bash
# Backup current config
cp "~/Library/Application Support/Claude/claude_desktop_config.json" \
   "~/Library/Application Support/Claude/claude_desktop_config.json.backup"

# Update with proxy config (edit paths first!)
cp claude-config.json "~/Library/Application Support/Claude/claude_desktop_config.json"
```

**On Windows:**
```bash
# Backup current config
cp "%APPDATA%/Claude/claude_desktop_config.json" \
   "%APPDATA%/Claude/claude_desktop_config.json.backup"

# Update with proxy config (edit paths first!)
cp claude-config.json "%APPDATA%/Claude/claude_desktop_config.json"
```

### 5. Start the Proxy

```bash
./build/mcp-proxy --config config.json
```

### 6. Restart Claude Desktop

## 📋 Included MCP Servers

This configuration includes popular MCP servers:

| Server | Description | Tools |
|--------|-------------|-------|
| **brave-search** | Web search capabilities | Web search, local search |
| **github** | GitHub repository management | Repo operations, code search |
| **filesystem** | Safe file system access | File operations, directory listing |
| **puppeteer** | Browser automation | Screenshots, navigation, interaction |
| **memory** | Knowledge graph storage | Entity management, relations |
| **sequential-thinking** | Step-by-step reasoning | Structured thinking |
| **git** | Git repository operations | Status, commit, branch management |
| **iterm** | Terminal integration | Command execution |
| **obsidian** | Note management | Vault operations |
| **browser-tools** | Browser dev tools | Debugging, auditing |

## 🔧 Configuration

### Proxy Settings (`proxy-config.json`)

```json
{
  "mcpProxy": {
    "baseURL": "http://localhost:9090",
    "addr": ":9090",
    "name": "MCP Proxy",
    "version": "1.0.0",
    "type": "streamable-http",// The transport type of the MCP proxy server, can be `streamable-http`, `sse`. By default, it is `sse`.
    "options": {
      "logEnabled": true,
      "authTokens": ["your-auth-token"]
    }
  },
  "mcpServers": {
    "server-name": {
      "command": "npx",
      "args": ["-y", "package-name"],
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
  }
}
```

### Claude Desktop Config (`claude-config.json`)

```json
{
  "mcpServers": {
    "mcp-proxy": {
      "command": "/path/to/mcp-proxy/build/mcp-proxy",
      "args": ["--config", "/path/to/proxy-config.json"]
    }
  }
}
```

## 🔒 Security Features

- **Tool Filtering**: Each server has explicit allow-lists
- **Authentication**: Auth tokens for proxy access
- **Path Restrictions**: Filesystem access limited to safe directories
- **Environment Isolation**: Each MCP server runs separately

## 🛠️ Management

### Adding New Servers

1. Add to `proxy-config.json`:
```json
"new-server": {
  "command": "npx",
  "args": ["-y", "new-mcp-package"],
  "options": {
    "toolFilter": {
      "mode": "allow",
      "list": ["allowed-tool"]
    }
  }
}
```

2. Restart the proxy

### Removing Servers

1. Delete the server entry from `proxy-config.json`
2. Restart the proxy

### Viewing Logs

The proxy provides detailed logs showing:
- Server connection status
- Tool registrations
- Request routing
- Error details

## 🚨 Troubleshooting

### Common Issues

**Proxy won't start**
- Check if port 9090 is available: `lsof -i :9090`
- Verify all MCP packages are installed
- Check file paths in configuration

**Claude can't connect**
- Ensure proxy is running
- Verify claude-config.json paths are correct
- Restart Claude Desktop after config changes

**Missing tools**
- Check proxy logs for "Adding tool" messages
- Verify tool names in filter lists match exactly
- Ensure API keys are valid

## 📚 Resources

- [Model Context Protocol Specification](https://modelcontextprotocol.io/)
- [Official MCP Servers](https://github.com/modelcontextprotocol/servers)
- [TBXark MCP Proxy](https://github.com/TBXark/mcp-proxy) (base proxy)

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test with your setup
5. Submit a pull request

## 📝 License

This project follows the MIT License. See individual MCP server licenses for their terms.

## 🙏 Acknowledgments

- [TBXark](https://github.com/TBXark) for the excellent mcp-proxy implementation
- [Anthropic](https://github.com/anthropic) for Claude and MCP specification
- The MCP community for creating amazing servers

---

**Need help?** Open an issue or check the troubleshooting section above.

### **`options`**
Common options for `mcpProxy` and `mcpServers`.

- `panicIfInvalid`: If true, the server will panic if the client is invalid.
- `logEnabled`: If true, the server will log the client's requests.
- `authTokens`: A list of authentication tokens for the client. The `Authorization` header will be checked against this list.
- `toolFilter`: Optional tool filtering configuration. **This configuration is only effective in `mcpServers`.**
  - `mode`: Specifies the filtering mode. Must be explicitly set to `allow` or `block` if `list` is provided. If `list` is present but `mode` is missing or invalid, the filter will be ignored for this server.
  - `list`: A list of tool names to filter (either allow or block based on the `mode`).
  > **Tip:** If you don't know the exact tool names, run the proxy once without any `toolFilter` configured. The console will log messages like `<server_name> Adding tool <tool_name>` for each successfully registered tool. You can use these logged names in your `toolFilter` list.

> In the new configuration, the `authTokens` of `mcpProxy` is not a global authentication token, but rather the default authentication token for `mcpProxy`. When `authTokens` is set in `mcpServers`, the value of `authTokens` in `mcpServers` will be used instead of the value in `mcpProxy`. In other words, the `authTokens` of `mcpProxy` serves as a default value and is only applied when `authTokens` is not set in `mcpServers`.

> Other fields are the same.

### **`mcpProxy`**
Proxy HTTP server configuration
- `baseURL`: The public accessible URL of the server. This is used to generate the URLs for the clients.
- `addr`: The address the server listens on.
- `name`: The name of the server.
- `version`: The version of the server.
- `type`: The transport type of the MCP proxy server. Can be `streamable-http` or `sse`. By default, it is `sse`.
  - `streamable-http`: The MCP proxy server supports HTTP streaming.
  - `sse`: The MCP proxy server supports Server-Sent Events (SSE).
- `options`: Default options for the `mcpServers`.

### **`mcpServers`**
MCP server configuration, Adopt the same configuration format as other MCP Clients.
- `transportType`: The transport type of the MCP client. Except for `streamable-http`, which requires manual configuration, the rest will be automatically configured according to the content of the configuration file.
  - `stdio`: The MCP client is a command line tool that is run in a subprocess.
  - `sse`: The MCP client is a server that supports SSE (Server-Sent Events).
  - `streamable-http`: The MCP client is a server that supports HTTP streaming.

For stdio mcp servers, the `command` field is required.
- `command`: The command to run the MCP client.
- `args`: The arguments to pass to the command.
- `env`: The environment variables to set for the command.
- `options`: Options specific to the client.

For sse mcp servers, the `url` field is required. When the current `url` exists, `sse` will be automatically configured.
- `url`: The URL of the MCP client.
- `headers`: The headers to send with the request to the MCP client.

For http streaming mcp servers, the `url` field is required. and `transportType` need to manually set to `streamable-http`.
- `url`: The URL of the MCP client.
- `headers`: The headers to send with the request to the MCP client.
- `timeout`: The timeout for the request to the MCP client.


## Usage

Command‑line flags, endpoints, and auth examples are documented in [docs/usage.md](docs/USAGE.md).

```
Usage of mcp-proxy:
  -config string
        path to config file or a http(s) url (default "config.json")
  -help
        print help and exit
  -version
        print version and exit
```
1. The server will start and aggregate the tools and capabilities of the configured MCP clients.
2. When MCP Server type is `sse`, You can access the server at `http(s)://{baseURL}/{clientName}/sse`. (e.g., `https://mcp.example.com/fetch/sse`, based on the example configuration)
3. When MCP Server type is `streamable-http`, You can access the server at `http(s)://{baseURL}/{clientName}/mcp`. (e.g., `https://mcp.example.com/fetch/mcp`, based on the example configuration)
4. If your MCP client does not support custom request headers., you can change the key in `clients` such as `fetch` to `fetch/{authToken}`, and then access it via `fetch/{authToken}`.

## Thanks

- This project was inspired by the [adamwattis/mcp-proxy-server](https://github.com/adamwattis/mcp-proxy-server) project
- If you have any questions about deployment, you can refer to  [《在 Docker 沙箱中运行 MCP Server》](https://miantiao.me/posts/guide-to-running-mcp-server-in-a-sandbox/)([@ccbikai](https://github.com/ccbikai))

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
