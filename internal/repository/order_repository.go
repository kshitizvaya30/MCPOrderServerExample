package repository

import (
	"context"
	"errors"
	"fmt"

	"Supportagent/backend/internal/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderRepository struct {
	db *pgxpool.Pool
}

func NewOrderRepository(db *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{
		db: db,
	}
}

func (r *OrderRepository) GetByOrderNumber(
	ctx context.Context,
	orderNumber string,
) (*model.Order, error) {

	query := `
		SELECT id, order_number, user_id, status, total_amount, created_at
		FROM orders
		WHERE order_number = $1
	`

	var order model.Order

	err := r.db.QueryRow(
		ctx,
		query,
		orderNumber,
	).Scan(
		&order.ID,
		&order.OrderNumber,
		&order.UserID,
		&order.Status,
		&order.TotalAmount,
		&order.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrOrderNotFound
		}

		return nil, fmt.Errorf("failed to get order %s: %w", orderNumber, err)
	}

	return &order, nil
}

func (r *OrderRepository) GetByUserID(
	ctx context.Context,
	userID int64,
) ([]model.Order, error) {

	query := `
		SELECT
			id, 
			order_number, 
			user_id, status,
			total_amount, 
			created_at
		FROM orders
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders for user %d: %w", userID, err)
	}
	defer rows.Close()

	var orders []model.Order

	for rows.Next() {
		var order model.Order

		err := rows.Scan(
			&order.ID,
			&order.OrderNumber,
			&order.UserID,
			&order.Status,
			&order.TotalAmount,
			&order.CreatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}

		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed while reading orders: %w", err)
	}

	return orders, nil
}

func (r *OrderRepository) CreateReturn(
	ctx context.Context,
	orderID int64,
	reason string,
) (*model.Return, error) {

	query := `
		INSERT INTO returns 
			(order_id, reason, status)
		VALUES 
			($1, $2, 'REQUESTED')
		RETURNING 
			id, 
			order_id, 
			reason, 
			status, 
			created_at
	`

	var returnRequest model.Return

	err := r.db.QueryRow(
		ctx,
		query,
		orderID,
		reason,
	).Scan(
		&returnRequest.ID,
		&returnRequest.OrderID,
		&returnRequest.Reason,
		&returnRequest.Status,
		&returnRequest.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create return: %w", err)
	}

	return &returnRequest, nil
}

func (r *OrderRepository) CreateRefundAndMarkRefunded(
	ctx context.Context,
	orderID int64,
	amount float64,
) (*model.Refund, error) {

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		_ = tx.Rollback(ctx)
	}()

	refundQuery := `
		INSERT INTO refunds 
			(order_id, amount, status)
		VALUES 
			($1, $2, 'COMPLETED')
		RETURNING 
			id, 
			order_id, 
			amount, 
			status, 
			created_at
	`

	var refund model.Refund

	err = tx.QueryRow(
		ctx,
		refundQuery,
		orderID,
		amount,
	).Scan(
		&refund.ID,
		&refund.OrderID,
		&refund.Amount,
		&refund.Status,
		&refund.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create refund: %w", err)
	}

	updateOrderQuery := `
		UPDATE orders
		SET status = 'REFUNDED'
		WHERE id = $1
	`

	result, err := tx.Exec(
		ctx,
		updateOrderQuery,
		orderID,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to cancel order: %w", err)
	}

	if result.RowsAffected() != 1 {
		return nil, errors.New("failed to cancel order")
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &refund, nil
}

var ErrOrderNotFound = errors.New("order not found")
