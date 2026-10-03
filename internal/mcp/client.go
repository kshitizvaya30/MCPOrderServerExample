package mcp

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type Client struct {
	session *mcp.ClientSession
}

func NewClient(ctx context.Context) (*Client, error) {
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
		return nil, err
	}

	return &Client{
		session: session,
	}, nil
}

func (c *Client) ListTools(
	ctx context.Context,
) (*mcp.ListToolsResult, error) {
	return c.session.ListTools(ctx, nil)
}

func (c *Client) CallTool(
	ctx context.Context,
	name string,
	arguments map[string]any,
) (*mcp.CallToolResult, error) {

	return c.session.CallTool(
		ctx,
		&mcp.CallToolParams{
			Name:      name,
			Arguments: arguments,
		},
	)
}

func (c *Client) Close() error {
	return c.session.Close()
}

func (c *Client) GetOrder(
	ctx context.Context,
	orderNumber string,
) (map[string]any, error) {

	result, err := c.CallTool(
		ctx,
		"get_order",
		map[string]any{
			"order_number": orderNumber,
		},
	)

	if err != nil {
		return nil, err
	}

	if result.IsError {
		return nil, fmt.Errorf("get_order tool failed")
	}

	data, ok := result.StructuredContent.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("unexpected get_order result format")
	}

	return data, nil
}
