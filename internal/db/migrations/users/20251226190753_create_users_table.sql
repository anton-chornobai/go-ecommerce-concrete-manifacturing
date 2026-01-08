-- +goose Up
CREATE TABLE users (
    id TEXT PRIMARY KEY,
    number TEXT NOT NULL UNIQUE,
    role TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL DEFAULT '',
    address TEXT NOT NULL DEFAULT '',
    name TEXT NOT NULL DEFAULT '',
    surname TEXT NOT NULL DEFAULT '',
    created_at TEXT NOT NULL DEFAULT (datetime('now','localtime'))
);

-- +goose Down
-- No operation, table will not be dropped
