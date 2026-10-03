package main

import (
	"context"
	"log"
	"os/exec"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	ctx := context.Background()

	client := mcp.NewClient(
		&mcp.Implementation{
			Name:    "support-ai-client",
			Version: "1.0.0",
		},
		nil,
	)

	transport := &mcp.CommandTransport{
		Command: exec.Command(
			"go",
			"run",
			"./cmd/mcp-server",
		),
	}

	session, err := client.Connect(ctx, transport, nil)
	if err != nil {
		log.Fatal(err)
	}
	result, err := session.CallTool(
		ctx,
		&mcp.CallToolParams{
			Name: "get_order",
			Arguments: map[string]any{
				"order_number": "ORD-1007",
			},
		},
	)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("MCP RESULT: %+v", result)
	defer session.Close()

	// Ask MCP server what tools it provides.
	tools, err := session.ListTools(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	for _, tool := range tools.Tools {
		log.Printf("MCP TOOL: %s - %s", tool.Name, tool.Description)
	}
}
