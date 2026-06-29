CREATE TABLE IF NOT EXISTS account (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    account_id CHAR(27) NOT NULL,
    total_price MONEY NOT NULL DEFAULT 0.0,
);

CREATE TABLE IF NOT EXISTS order_products (
    order_id CHAR(27) REFERENCES orders(order_id) ON DELETE CASCADE,
    product_id CHAR(27),
    quantity INT NOT NULL,
    PRIMARY KEY (order_id, product_id)
)