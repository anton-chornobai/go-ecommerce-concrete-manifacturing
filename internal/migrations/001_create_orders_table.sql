-- +goose Up

CREATE TABLE users (
    id UUID PRIMARY KEY,
    number TEXT UNIQUE,
    role TEXT NOT NULL,
    email TEXT UNIQUE,
    password TEXT,
    address TEXT,
    name TEXT NOT NULL DEFAULT '',
    surname TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    price INTEGER NOT NULL DEFAULT 0,
    title VARCHAR(50) NOT NULL UNIQUE,
    product_type TEXT NOT NULL,
    image_url TEXT,
    color TEXT,
    description TEXT,
    stock_quantity INTEGER,
    weight_grams INTEGER,
    rating SMALLINT,
    size_width INTEGER,
    size_height INTEGER
);

CREATE TABLE orders (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    total INTEGER NOT NULL DEFAULT 0,
    status TEXT NOT NULL DEFAULT 'pending',
    order_name TEXT NOT NULL DEFAULT '',
    discount INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE order_items (
    id UUID PRIMARY KEY,
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    product_id UUID NOT NULL,
    title TEXT NOT NULL,
    unit_price INTEGER NOT NULL,
    type TEXT,
    quantity INTEGER NOT NULL,
    color TEXT,
    height INTEGER,
    width INTEGER,
    material TEXT,
    thickness INTEGER
);

-- +goose Down

