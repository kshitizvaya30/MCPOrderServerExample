package model

import "time"

type Return struct {
	ID        int64     `json:"id"`
	OrderID   int64     `json:"order_id"`
	Reason    string    `json:"reason"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// REQUESTED
// APPROVED
// REJECTED
// COMPLETED
