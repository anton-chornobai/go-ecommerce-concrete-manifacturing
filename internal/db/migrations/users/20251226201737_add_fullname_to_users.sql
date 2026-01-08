-- +goose Up
ALTER TABLE users ADD COLUMN fullname TEXT NOT NULL DEFAULT '';

-- +goose Down

CREATE TABLE users_temp (
    id TEXT PRIMARY KEY,
    number TEXT NOT NULL UNIQUE,
    role TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL DEFAULT '',
    address TEXT NOT NULL DEFAULT '',
    name TEXT NOT NULL DEFAULT '',
    surname TEXT NOT NULL DEFAULT '',
    created_at TEXT NOT NULL DEFAULT (datetime('now','localtime'))
);

INSERT INTO users_temp (id, number, role, email, address, name, surname, created_at)
SELECT id, number, role, email, address, name, surname, created_at
FROM users;

DROP TABLE users;

ALTER TABLE users_temp RENAME TO users;
