# **Conflict Resolution Guidance**

The primary conflicts stem from two major parallel developments:

1. **Robin Collins' Branch (Ours):** Refactored configuration by moving V1 logic out and simplifying dependency usage (`TBXark/confstore`), introduced Streamable HTTP support, and added the detailed `/healthCheck/` endpoint.
2. **Target Branch (Theirs):** Updated Go version, updated `mcp-go` to a newer version (`v0.39.1`), and maintained complex configuration loading logic (HTTP config, SSL skipping, environment expansion) using the `go-sphere/confstore` dependency structure.

The resolution strategy is to **keep the newer dependencies and advanced configuration loading from Theirs**, but **integrate Robin's feature logic (V1 adaptation, Health Check, Streamable HTTP)**.

---

## **1. `go.mod` Resolution**

The conflict involves different Go versions, different `mcp-go` versions, and different `confstore` dependencies.

| Area | Ours (Robin Collins) | Theirs (Target Branch) | Decision | Rationale |
| :--- | :--- | :--- | :--- | :--- |
| **Go Version** | `go 1.23.0` | `go 1.24.0` | **Keep Theirs** | Use the newer Go standard. |
| **`confstore`** | `github.com/TBXark/confstore v0.0.4` | `github.com/go-sphere/confstore v0.0.4` | **Keep `TBXark/confstore`** | Robin's refactoring relies on the specific structure of the `TBXark` fork (which is likely the project's current upstream). We will manually remove the `go-sphere` dependency. |
| **`mcp-go`** | `v0.31.0` | `v0.39.1` | **Keep Theirs (`v0.39.1`)** | Use the latest version, which should include the features required by `v0.31.0` (Streamable HTTP). |
| **`x/sync`** | `v0.14.0` | `v0.17.0` | **Keep Theirs** | Use the newer dependency version. |
| **Indirect Deps** | Minimal list | Extensive list | **Keep Theirs' list** | The extensive list corresponds to the newer `mcp-go` version and should be preserved. |

### **Recommended `go.mod` Final Content:**

```go
module github.com/TBXark/mcp-proxy

go 1.24.0

require (
    github.com/TBXark/confstore v0.0.4 // Retain TBXark fork
    github.com/TBXark/optional-go v0.0.1
    github.com/mark3labs/mcp-go v0.39.1 // Use newer version
    golang.org/x/sync v0.17.0
)

require (
    github.com/bahlo/generic-list-go v0.2.0 // indirect
    github.com/buger/jsonparser v1.1.1 // indirect
    github.com/google/uuid v1.6.0 // indirect
    github.com/invopop/jsonschema v0.13.0 // indirect
    github.com/mailru/easyjson v0.9.0 // indirect
    github.com/spf13/cast v1.9.2 // indirect
    github.com/wk8/go-ordered-map/v2 v2.1.8 // indirect
    github.com/yosida95/uritemplate/v3 v3.0.2 // indirect
    gopkg.in/yaml.v3 v3.0.1 // indirect
)
```

---

## **2. `go.sum` Resolution**

**Guidance:** **Discard both versions.** After resolving `go.mod`, run `go mod tidy` to generate the correct `go.sum` file based on the merged dependencies.

---

## **3. `config.go` Resolution**

The conflict here is structural: Robin's branch removed the complex configuration provider logic, while the target branch kept it. We must reintroduce Robin's V1 adaptation logic into the target branch's robust loading mechanism.

### **Step-by-Step Merge:** config.go

1. **Imports:** **Keep Theirs' imports** (including `crypto/tls`, `net/http`, `strings`, and the `confstore/provider` sub-packages). These are necessary for `newConfProvider`.
2. **`newConfProvider`:** **Keep Theirs' entire `newConfProvider` function.** This preserves the ability to load configs from remote URLs with custom options.
3. **`load` Function Signature:** **Keep Theirs' `load` function signature**, which accepts all provider options:

    ```go
    func load(path string, insecure, expandEnv bool, httpHeaders string, httpTimeout int) (*Config, error) {
    ```

4. **`load` Function Body:** Merge the logic:
    * Keep Theirs' logic for setting up the provider (`pro, err := newConfProvider(...)`) and loading the config (`confstore.Load[FullConfig](pro, codec.JsonCodec())`).
    * Immediately after loading the config, insert Robin's V1 adaptation call.

### **Recommended `config.go` Final Content (Conflict Area Only):**

```go
// ... (Keep all imports from Theirs)

// ... (Keep FullConfig struct definition)

// Keep Theirs' entire newConfProvider function
func newConfProvider(path string, insecure, expandEnv bool, httpHeaders string, httpTimeout int) (provider.Provider, error) {
// ... (Theirs' implementation)
}

// Use Theirs' signature, but implement Robin's logic
func load(path string, insecure, expandEnv bool, httpHeaders string, httpTimeout int) (*Config, error) {
    pro, err := newConfProvider(path, insecure, expandEnv, httpHeaders, httpTimeout)
    if err != nil {
        return nil, err
    }
    conf, err := confstore.Load[FullConfig](pro, codec.JsonCodec())
    if err != nil {
        return nil, err
    }
    // --- MERGE: Insert Robin's V1 Adaptation Call ---
    adaptMCPClientConfigV1ToV2(conf)
    // -------------------------------------------------

    if conf.McpProxy == nil {
        return nil, errors.New("mcpProxy is required")
    }
// ... (Rest of the function, which handles validation and setting default type)
```

---

## **4. `http.go` Resolution**

The conflict in `http.go` is straightforward: Robin's branch introduced the entire health check feature and client management logic, which is missing in the target branch.

### **Step-by-Step Merge:** http.go

1. **Health Check Logic:** **Keep all of Ours (Robin's code)** for the `activeClients` map declaration and the entire `healthCheckHandler` function. Insert this code block immediately after `recoverMiddleware`.
2. **`startHTTPServer` Integration:**
    * In the client loop, **keep Ours' logic** to store the client in `activeClients`.
    * In the shutdown hook, **keep Ours' logic** to `delete(activeClients, name)`.
    * After the client loop, **keep Ours' logic** to register the health check endpoint: `httpMux.HandleFunc("/healthCheck/", healthCheckHandler(config))`.
    * In the final `go func()` block, **keep Ours' dynamic logging** which uses `config.McpProxy.Type`.

### **Recommended `http.go` Final Content (Conflict Areas Only):**

```go
// After recoverMiddleware function:
var (
    activeClients = make(map[string]*Client)
)

// healthCheckHandler returns a JSON health status, mcp_servers, and process uptime
func healthCheckHandler(config *Config) http.HandlerFunc {
// ... (Keep Robin's full implementation of healthCheckHandler)
}

func startHTTPServer(config *Config) error {
// ...
    for name, clientConfig := range config.McpServers {
        mcpClient, err := newMCPClient(name, clientConfig)
        if err != nil {
            return err
        }
        // Store the client in activeClients <-- KEEP OURS
        activeClients[name] = mcpClient

        server, err := newMCPServer(name, config.McpProxy, clientConfig)
// ...
            httpServer.RegisterOnShutdown(func() {
                log.Printf("<%s> Shutting down", name)
                _ = mcpClient.Close()
                delete(activeClients, name) // <-- KEEP OURS
            })
            return nil
        })
    }

    // Register /healthCheck/ endpoint (no auth, no middleware) <-- KEEP OURS
    httpMux.HandleFunc("/healthCheck/", healthCheckHandler(config))

    go func() {
        // Use dynamic logging from Robin's branch
        log.Printf("Starting %s server", config.McpProxy.Type)
        log.Printf("%s server listening on %s", config.McpProxy.Type, config.McpProxy.Addr)
        hErr := httpServer.ListenAndServe()
// ...
```
