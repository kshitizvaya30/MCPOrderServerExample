CREATE TABLE refunds (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL,
    amount NUMERIC(12, 2) NOT NULL,
    status VARCHAR(30) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_refunds_order
        FOREIGN KEY (order_id)
        REFERENCES orders(id),

    CONSTRAINT chk_refunds_amount
        CHECK (amount > 0)
);

