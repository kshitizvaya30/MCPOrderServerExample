package service

import (
	"context"
	"errors"

	"Supportagent/backend/internal/model"
	"Supportagent/backend/internal/repository"
)

type OrderService struct {
	orderRepository *repository.OrderRepository
}

func NewOrderService(
	orderRepository *repository.OrderRepository,
) *OrderService {
	return &OrderService{
		orderRepository: orderRepository,
	}
}

func (s *OrderService) GetOrder(
	ctx context.Context,
	orderNumber string,
) (*model.Order, error) {

	if orderNumber == "" {
		return nil, errors.New("order number is required")
	}

	return s.orderRepository.GetByOrderNumber(ctx, orderNumber)
}

func (s *OrderService) GetOrdersByUserID(
	ctx context.Context,
	userID int64,
) ([]model.Order, error) {

	if userID <= 0 {
		return nil, errors.New("invalid user id")
	}

	return s.orderRepository.GetByUserID(ctx, userID)
}

func (s *OrderService) CreateReturn(
	ctx context.Context,
	orderNumber string,
	reason string,
) (*model.Return, error) {

	if orderNumber == "" {
		return nil, errors.New("order number is required")
	}

	if reason == "" {
		return nil, errors.New("return reason is required")
	}

	order, err := s.orderRepository.GetByOrderNumber(
		ctx,
		orderNumber,
	)

	if err != nil {
		return nil, err
	}

	if order.Status != "DELIVERED" {
		return nil, errors.New(
			"return can only be requested for a delivered order",
		)
	}

	return s.orderRepository.CreateReturn(
		ctx,
		order.ID,
		reason,
	)
}

func (s *OrderService) CreateRefund(
	ctx context.Context,
	orderNumber string,
) (*model.Refund, error) {

	if orderNumber == "" {
		return nil, errors.New("order number is required")
	}

	order, err := s.orderRepository.GetByOrderNumber(
		ctx,
		orderNumber,
	)

	if err != nil {
		return nil, err
	}

	if order.Status != "DELIVERED" {
		return nil, errors.New(
			"refund can only be requested for a delivered order",
		)
	}

	return s.orderRepository.CreateRefundAndMarkRefunded(
		ctx,
		order.ID,
		order.TotalAmount,
	)
}
