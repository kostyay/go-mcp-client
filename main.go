package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/ThinkInAIXYZ/go-mcp/client"
	"github.com/ThinkInAIXYZ/go-mcp/protocol"
	"github.com/ThinkInAIXYZ/go-mcp/transport"
)

// BearerAuthTransport is an http.RoundTripper that adds a Bearer token
// to the Authorization header of requests.
type BearerAuthTransport struct {
	Token     string           // Bearer token
	Transport http.RoundTripper // Base transport (e.g., http.DefaultTransport)
}

// RoundTrip executes a single HTTP transaction, adding the Bearer token.
func (t *BearerAuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Get the base transport if not set
	transport := t.Transport
	if transport == nil {
		transport = http.DefaultTransport
	}

	// Clone the request to avoid modifying the original
	clonedReq := req.Clone(req.Context())

	// Add the Authorization header
	clonedReq.Header.Set("Authorization", "Bearer "+t.Token)

	// Delegate to the base transport
	return transport.RoundTrip(clonedReq)
}

type stdoutLogger struct {
}

func (s stdoutLogger) Debugf(format string, a ...any) {
	//TODO implement me
	fmt.Printf(format, a...)
}

func (s stdoutLogger) Infof(format string, a ...any) {
	//TODO implement me
	fmt.Printf(format, a...)
}

func (s stdoutLogger) Warnf(format string, a ...any) {
	//TODO implement me
	fmt.Printf(format, a...)
}

func (s stdoutLogger) Errorf(format string, a ...any) {
	//TODO implement me
	fmt.Printf(format, a...)
}

var (
	serverURL   = flag.String("serverURL", "", "URL of the directory service")
	bearerToken = flag.String("token", "", "Bearer token for authentication")
)

func main() {
	flag.Parse() // Parse flags first

	if *serverURL == "" {
		log.Fatalf("serverURL is required")
	}

	if *bearerToken == "" {
		log.Fatalf("bearerToken is required")
	}

	l := &stdoutLogger{}
	// Create transport client (using SSE in this example)
	fmt.Printf("Creating transport client for %s\n", *serverURL)

	httpClient := &http.Client{
		// Use the custom BearerAuthTransport
		Transport: &BearerAuthTransport{
			Token: *bearerToken,
		},
	}

	// Use the configured httpClient for SSE transport options if needed
	// (check transport.NewSSEClientTransport documentation for options
	// like transport.WithSSEClientOptionHTTPClient(httpClient))
	transportClient, err := transport.NewSSEClientTransport(
		fmt.Sprintf("%s/sse", *serverURL),
		transport.WithSSEClientOptionLogger(l),
		transport.WithSSEClientOptionHTTPClient(httpClient),
	)
	if err != nil {
		log.Fatalf("Failed to create transport client: %v", err)
	}

	// Create MCP client using transport
	mcpClient, err := client.NewClient(transportClient, client.WithClientInfo(protocol.Implementation{
		Name:    "example MCP client",
		Version: "1.0.0",
	}), client.WithLogger(l))
	if err != nil {
		log.Fatalf("Failed to create MCP client: %v", err)
	}
	defer mcpClient.Close()

	// List available tools
	toolsResult, err := mcpClient.ListTools(context.Background())
	if err != nil {
		log.Fatalf("Failed to list tools: %v", err)
	}
	b, _ := json.Marshal(toolsResult.Tools)
	fmt.Printf("Available tools: %+v\n", string(b))
	//
	//// Call tool
	//callResult, err := mcpClient.CallTool(
	//	context.Background(),
	//	protocol.NewCallToolRequest("current time", map[string]interface{}{
	//		"timezone": "UTC",
	//	}))
	//if err != nil {
	//	log.Fatalf("Failed to call tool: %v", err)
	//}
	//b, _ = json.Marshal(callResult)
	//fmt.Printf("Tool call result: %+v\n", string(b))
}
