package mcp

import (
	"context"

	"Supportagent/backend/internal/model"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type OrderReader interface {
	GetOrder(
		ctx context.Context,
		orderNumber string,
	) (*model.Order, error)

	CreateReturn(
		ctx context.Context,
		orderNumber string,
		reason string,
	) (*model.Return, error)

	CreateRefund(
		ctx context.Context,
		orderNumber string,
	) (*model.Refund, error)
}

type GetOrderInput struct {
	OrderNumber string `json:"order_number"`
}

type GetOrderOutput struct {
	OrderNumber string  `json:"order_number"`
	Status      string  `json:"status"`
	TotalAmount float64 `json:"total_amount"`
}

type CreateReturnInput struct {
	OrderNumber string `json:"order_number"`
	Reason      string `json:"reason"`
}

type CreateRefundInput struct {
	OrderNumber string `json:"order_number"`
}

type Server struct {
	orderReader OrderReader
}

func NewServer(orderReader OrderReader) *mcp.Server {

	server := mcp.NewServer(
		&mcp.Implementation{
			Name:    "support-ai-mcp",
			Version: "1.0.0",
		},
		nil,
	)

	// -----------------------------------------
	// get_order
	// -----------------------------------------

	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "get_order",
			Description: "Get order details by order number",
		},
		func(
			ctx context.Context,
			req *mcp.CallToolRequest,
			input GetOrderInput,
		) (*mcp.CallToolResult, GetOrderOutput, error) {

			order, err := orderReader.GetOrder(
				ctx,
				input.OrderNumber,
			)

			if err != nil {
				return nil, GetOrderOutput{}, err
			}

			return nil, GetOrderOutput{
				OrderNumber: order.OrderNumber,
				Status:      order.Status,
				TotalAmount: order.TotalAmount,
			}, nil
		},
	)

	// -----------------------------------------
	// create_return
	// -----------------------------------------

	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "create_return",
			Description: "Create a return request for a delivered order",
		},
		func(
			ctx context.Context,
			req *mcp.CallToolRequest,
			input CreateReturnInput,
		) (*mcp.CallToolResult, model.Return, error) {

			result, err := orderReader.CreateReturn(
				ctx,
				input.OrderNumber,
				input.Reason,
			)

			if err != nil {
				return nil, model.Return{}, err
			}

			return nil, *result, nil
		},
	)

	// -----------------------------------------
	// create_refund
	// -----------------------------------------

	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "create_refund",
			Description: "Create a full refund for a delivered order",
		},
		func(
			ctx context.Context,
			req *mcp.CallToolRequest,
			input CreateRefundInput,
		) (*mcp.CallToolResult, model.Refund, error) {

			result, err := orderReader.CreateRefund(
				ctx,
				input.OrderNumber,
			)

			if err != nil {
				return nil, model.Refund{}, err
			}

			return nil, *result, nil
		},
	)

	return server
}
