package model

import "time"

type Order struct {
	ID          int64     `json:"id"`
	OrderNumber string    `json:"order_number"`
	UserID      int64     `json:"user_id"`
	Status      string    `json:"status"`
	TotalAmount float64   `json:"total_amount"`
	CreatedAt   time.Time `json:"created_at"`
}
