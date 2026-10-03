package tool

import (
	"context"
	"encoding/json"
	"fmt"

	"Supportagent/backend/internal/model"
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
		context.Context, string) (*model.Refund, error)
}

type OrderTool struct {
	orderReader OrderReader
}

type GetOrderInput struct {
	OrderNumber string `json:"order_number"`
}

type CreateReturnInput struct {
	OrderNumber string `json:"order_number"`
	Reason      string `json:"reason"`
}

type CreateRefundInput struct {
	OrderNumber string `json:"order_number"`
}

////////////////////////////
////////////////////////////

func NewOrderTool(orderReader OrderReader) *OrderTool {
	return &OrderTool{
		orderReader: orderReader,
	}
}

func (t *OrderTool) GetOrder(
	ctx context.Context,
	input GetOrderInput,
) (string, error) {

	order, err := t.orderReader.GetOrder(
		ctx,
		input.OrderNumber,
	)
	if err != nil {
		return "", err
	}

	result := model.Order{
		ID:          order.ID,
		OrderNumber: order.OrderNumber,
		UserID:      order.UserID,
		Status:      order.Status,
		TotalAmount: order.TotalAmount,
		CreatedAt:   order.CreatedAt,
	}

	data, err := json.Marshal(result)
	if err != nil {
		return "", fmt.Errorf("failed to serialize order: %w", err)
	}

	return string(data), nil
}

func (t *OrderTool) CreateReturn(
	ctx context.Context,
	input CreateReturnInput,
) (string, error) {

	order, err := t.orderReader.GetOrder(
		ctx,
		input.OrderNumber,
	)
	if err != nil {
		return "", fmt.Errorf("failed to get order: %w", err)
	}

	if order.Status != "DELIVERED" {
		return "", fmt.Errorf(
			"order %s is not eligible for return because its status is %s",
			order.OrderNumber,
			order.Status,
		)
	}

	result, err := t.orderReader.CreateReturn(
		ctx,
		input.OrderNumber,
		input.Reason,
	)
	if err != nil {
		return "", fmt.Errorf("failed to create return: %w", err)
	}

	data, err := json.Marshal(result)
	if err != nil {
		return "", fmt.Errorf(
			"failed to serialize return: %w",
			err,
		)
	}

	return string(data), nil
}

func (t *OrderTool) CreateRefund(
	ctx context.Context,
	input CreateRefundInput,
) (string, error) {

	result, err := t.orderReader.CreateRefund(
		ctx,
		input.OrderNumber,
	)
	if err != nil {
		return "", err
	}

	data, err := json.Marshal(result)
	if err != nil {
		return "", fmt.Errorf(
			"failed to serialize refund: %w",
			err,
		)
	}

	return string(data), nil
}
