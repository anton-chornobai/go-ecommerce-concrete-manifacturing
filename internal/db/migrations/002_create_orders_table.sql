-- +goose Up
 
CREATE TABLE orders (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id TEXT NOT NULL,
    name TEXT NOT NULL,
    total INTEGER,
    status TEXT NOT NULL DEFAULT 'pending',
    order_name TEXT NOT NULL DEFAULT '',
    discount INTEGER NOT NULL DEFAULT 0,
    created_at TEXT NOT NULL DEFAULT (datetime('now','localtime')),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- +goose Down
