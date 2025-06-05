# MCP Proxy

## File Tree

```filetree.txt
.
├── CHANGELOG.md
├── client.go
├── config_deprecated.go
├── config.go
├── config.json
├── docker-compose.yaml
├── Dockerfile
├── docs
│   └── index.html
├── FILETREE.md
├── go.mod
├── go.sum
├── http.go
├── LICENSE
├── main.go
├── Makefile
└── README.md
```

## File short descriptions

### `CHANGELOG.md`

Documents changes made to the project

### `client.go`

Implements MCP client functionality and server aggregation

### `config.go`

Defines configuration structures and parsing logic

### `config.json`

Contains the runtime configuration for the proxy

### `.cursorignore`

Specifies files to be ignored by the Cursor IDE

### `docker-compose.yaml`

Docker Compose configuration for orchestrating containers

### `Dockerfile`

Docker image build instructions

### `.dockerignore`

Specifies files to be excluded from Docker builds

### `docs/index.html`

Web interface for converting configurations for Claude

### `FILETREE.md`

Documents the project file structure and purpose

### `.gitattributes`

Defines Git attributes for specific file paths

### `.github/workflows/docker.yml`

GitHub Actions workflow for Docker image builds

### `.gitignore`

Specifies files to be ignored by Git

### `.golangci.yml`

Configuration for the golangci-lint tool

### `go.mod`

Go module definition and dependencies

### `go.sum`

Checksums for Go module dependencies

### `http.go`

HTTP server implementation with middlewares and endpoints

### `LICENSE`

MIT license for the project

### `main.go`

Application entry point with command-line parsing

### `Makefile`

Build automation for the project

### `README.md`

Project documentation and usage instructions
