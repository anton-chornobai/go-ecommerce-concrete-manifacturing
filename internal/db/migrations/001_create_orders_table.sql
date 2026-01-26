-- +goose Up
 
CREATE TABLE orders (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id TEXT NOT NULL,
    name TEXT NOT NULL,
    total INTEGER NOT NULL DEFAULT 0,
    status TEXT NOT NULL DEFAULT 'pending',
    order_name TEXT NOT NULL DEFAULT '',
    discount INTEGER NOT NULL DEFAULT 0,
    created_at TEXT NOT NULL DEFAULT (datetime('now','localtime')),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE users (
    id TEXT PRIMARY KEY,
    number TEXT NOT NULL UNIQUE,
    role TEXT NOT NULL,
    email TEXT DEFAULT '',
    password TEXT 
    address TEXT NOT NULL DEFAULT '',
    name TEXT NOT NULL DEFAULT '',
    surname TEXT NOT NULL DEFAULT '',
    created_at TEXT NOT NULL DEFAULT (datetime('now','localtime'))
);

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
