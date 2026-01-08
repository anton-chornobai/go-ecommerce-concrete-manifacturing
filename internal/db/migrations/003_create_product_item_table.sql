-- +goose Up

CREATE TABLE order_items (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    order_id INTEGER NOT NULL,
    product_id INTEGER NOT NULL,
    title TEXT NOT NULL,
    unit_price INTEGER NOT NULL,
    type TEXT,
    quantity INTEGER NOT NULL,
    color TEXT,
    height INTEGER,
    width INTEGER,
    material TEXT,
    thickness INTEGER,
    FOREIGN KEY (order_id) REFERENCES orders(id)
);

-- +goose Down