package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path"
	"strings"
	"syscall"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"golang.org/x/sync/errgroup"
)

type MiddlewareFunc func(http.Handler) http.Handler

func chainMiddleware(h http.Handler, middlewares ...MiddlewareFunc) http.Handler {
	for _, mw := range middlewares {
		h = mw(h)
	}
	return h
}

func newAuthMiddleware(tokens []string) MiddlewareFunc {
	tokenSet := make(map[string]struct{}, len(tokens))
	for _, token := range tokens {
		tokenSet[token] = struct{}{}
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if len(tokens) != 0 {
				token := r.Header.Get("Authorization")
				token = strings.TrimSpace(strings.TrimPrefix(token, "Bearer "))
				if token == "" {
					http.Error(w, "Unauthorized", http.StatusUnauthorized)
					return
				}
				if _, ok := tokenSet[token]; !ok {
					http.Error(w, "Unauthorized", http.StatusUnauthorized)
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}

func loggerMiddleware(prefix string) MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("<%s> Request [%s] %s", prefix, r.Method, r.URL.Path)
			next.ServeHTTP(w, r)
		})
	}
}

func recoverMiddleware(prefix string) MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					log.Printf("<%s> Recovered from panic: %v", prefix, err)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

var (
	activeClients = make(map[string]*Client)
)

// healthCheckHandler returns a JSON health status, mcp_servers, and process uptime
func healthCheckHandler(config *Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uptime := time.Since(startTime)
		mcpServers := make([]interface{}, 0, len(config.McpServers))

		for name, serverConfig := range config.McpServers {
			serverInfo := map[string]interface{}{
				"Status": "unknown",
				"tools":  []string{},
				"error":  "",
			}

			// Check if we have an active client
			if client, exists := activeClients[name]; exists {
				// Try to get tools list from existing client
				toolsRequest := mcp.ListToolsRequest{}
				if tools, err := client.client.ListTools(r.Context(), toolsRequest); err == nil {
					toolNames := make([]string, 0, len(tools.Tools))
					for _, tool := range tools.Tools {
						toolNames = append(toolNames, tool.Name)
					}
					serverInfo["tools"] = toolNames
					serverInfo["Status"] = "running"
				} else {
					errorMsg := fmt.Sprintf("Failed to list tools: %v", err)
					log.Printf("<%s> Health check error: %s", name, errorMsg)
					serverInfo["Status"] = "error"
					serverInfo["error"] = errorMsg
				}
			} else {
				// Create new client if it doesn't exist
				if client, err := newMCPClient(name, serverConfig); err == nil {
					activeClients[name] = client
					serverInfo["Status"] = "initializing"
					log.Printf("<%s> Health check: Client created, initializing", name)
				} else {
					errorMsg := fmt.Sprintf("Failed to create client: %v", err)
					log.Printf("<%s> Health check error: %s", name, errorMsg)
					serverInfo["Status"] = "error"
					serverInfo["error"] = errorMsg
				}
			}

			mcpServers = append(mcpServers, name)
			mcpServers = append(mcpServers, serverInfo)
		}

		response := map[string]interface{}{
			"status":      "healthy",
			"mcp_servers": mcpServers,
			"uptime":      uptime.String(),
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf("Failed to encode response: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}

func startHTTPServer(config *Config) error {
	baseURL, uErr := url.Parse(config.McpProxy.BaseURL)
	if uErr != nil {
		return uErr
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var errorGroup errgroup.Group
	httpMux := http.NewServeMux()
	httpServer := &http.Server{
		Addr:    config.McpProxy.Addr,
		Handler: httpMux,
	}
	info := mcp.Implementation{
		Name: config.McpProxy.Name,
	}

	for name, clientConfig := range config.McpServers {
		mcpClient, err := newMCPClient(name, clientConfig)
		if err != nil {
			return err
		}
		// Store the client in activeClients
		activeClients[name] = mcpClient

		server, err := newMCPServer(name, config.McpProxy, clientConfig)
		if err != nil {
			return err
		}
		errorGroup.Go(func() error {
			log.Printf("<%s> Connecting", name)
			addErr := mcpClient.addToMCPServer(ctx, info, server.mcpServer)
			if addErr != nil {
				log.Printf("<%s> Failed to add client to server: %v", name, addErr)
				if clientConfig.Options.PanicIfInvalid.OrElse(false) {
					return addErr
				}
				return nil
			}
			log.Printf("<%s> Connected", name)

			middlewares := make([]MiddlewareFunc, 0)
			middlewares = append(middlewares, recoverMiddleware(name))
			if clientConfig.Options.LogEnabled.OrElse(false) {
				middlewares = append(middlewares, loggerMiddleware(name))
			}
			if len(clientConfig.Options.AuthTokens) > 0 {
				middlewares = append(middlewares, newAuthMiddleware(clientConfig.Options.AuthTokens))
			}
			mcpRoute := path.Join(baseURL.Path, name)
			if !strings.HasPrefix(mcpRoute, "/") {
				mcpRoute = "/" + mcpRoute
			}
			if !strings.HasSuffix(mcpRoute, "/") {
				mcpRoute += "/"
			}
			log.Printf("<%s> Handling requests at %s", name, mcpRoute)
			httpMux.Handle(mcpRoute, chainMiddleware(server.handler, middlewares...))
			httpServer.RegisterOnShutdown(func() {
				log.Printf("<%s> Shutting down", name)
				_ = mcpClient.Close()
				delete(activeClients, name)
			})
			return nil
		})
	}

	// Register /healthCheck/ endpoint (no auth, no middleware)
	httpMux.HandleFunc("/healthCheck/", healthCheckHandler(config))

	go func() {
		err := errorGroup.Wait()
		if err != nil {
			log.Fatalf("Failed to add clients: %v", err)
		}
		log.Printf("All clients initialized")
	}()

	go func() {
		log.Printf("Starting SSE server")
		log.Printf("SSE server listening on %s", config.McpProxy.Addr)
		hErr := httpServer.ListenAndServe()
		if hErr != nil && !errors.Is(hErr, http.ErrServerClosed) {
			log.Fatalf("Failed to start server: %v", hErr)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	log.Println("Shutdown signal received")

	shutdownCtx, shutdownCancel := context.WithTimeout(ctx, 5*time.Second)
	defer shutdownCancel()

	err := httpServer.Shutdown(shutdownCtx)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}
