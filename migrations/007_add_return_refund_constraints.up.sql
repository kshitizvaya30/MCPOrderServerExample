CREATE UNIQUE INDEX idx_returns_order_id
ON returns(order_id);

CREATE UNIQUE INDEX idx_refunds_order_id
ON refunds(order_id);