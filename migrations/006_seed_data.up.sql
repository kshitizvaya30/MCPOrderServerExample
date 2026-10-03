TRUNCATE TABLE
    refunds,
    returns,
    order_items,
    orders,
    users
RESTART IDENTITY CASCADE;
-- Users
INSERT INTO users (name, email)
VALUES
    ('Kshitiz Vaya', 'kshitiz@example.com'),
    ('Alice Sharma', 'alice@example.com'),
    ('Bob Singh', 'bob@example.com');


-- Orders
INSERT INTO orders (order_number, user_id, status, total_amount)
VALUES
    ('ORD-1001', 1, 'DELIVERED', 1200.00),
    ('ORD-1002', 1, 'SHIPPED', 2500.00),
    ('ORD-1003', 2, 'DELIVERED', 899.00),
    ('ORD-1004', 1, 'DELIVERED', 1200.00),
    ('ORD-1005', 1, 'DELIVERED', 1230.00),
    ('ORD-1006', 1, 'DELIVERED', 1100.00),
    ('ORD-1007', 1, 'DELIVERED', 12000.00);


-- Order items
INSERT INTO order_items (order_id, product_name, quantity, unit_price)
VALUES
    (1, 'Wireless Headphones', 1, 1200.00),
    (2, 'Smart Watch', 1, 2500.00),
    (3, 'Keyboard', 1, 899.00),
    (4, 'Bottle', 1, 1200.00),
    (5, 'Car', 1, 1230.00),
    (6, 'bike', 1, 1100.00),
    (7, 'Blackboard', 1, 12000.00);