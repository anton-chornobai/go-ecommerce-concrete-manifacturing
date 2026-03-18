-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    number TEXT UNIQUE,
    role TEXT NOT NULL DEFAULT 'user',
    email TEXT UNIQUE,
    password TEXT,
    address TEXT NOT NULL DEFAULT '',
    name TEXT NOT NULL DEFAULT '',
    surname TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    verified BOOLEAN NOT NULL DEFAULT FALSE,
    verification_hash TEXT NOT NULL DEFAULT '', 
    verification_expires_at TIMESTAMP
);

CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    price INTEGER NOT NULL DEFAULT 0 CHECK (price >= 0),
    title VARCHAR(50) NOT NULL UNIQUE,
    product_type TEXT NOT NULL,
    image_url TEXT,
    color TEXT,
    description TEXT,
    status TEXT NOT NULL DEFAULT 'archived' CHECK (status IN ('archived', 'displayed')), 
    stock_quantity INTEGER NOT NULL DEFAULT 0 CHECK (stock_quantity >= 0),
    weight_grams INTEGER CHECK (weight_grams >= 0),
    rating SMALLINT CHECK (rating >= 0 AND rating <= 5),
    size_width INTEGER CHECK (size_width >= 0),
    size_height INTEGER CHECK (size_height >= 0)
);

CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    order_name TEXT NOT NULL DEFAULT '',
    total INTEGER NOT NULL DEFAULT 0 CHECK (total >= 0),
    deposit INTEGER NOT NULL DEFAULT 0 CHECK (deposit >= 0),
    status TEXT NOT NULL DEFAULT 'pending' 
        CHECK (status IN ('pending', 'confirmed', 'shipped', 'delivered', 'cancelled')),
    payment_status TEXT NOT NULL DEFAULT 'unpaid'
        CHECK (payment_status IN ('deposited','unpaid', 'paid', 'failed', 'refunded')),
    discount INTEGER NOT NULL DEFAULT 0 CHECK (discount >= 0),
    description VARCHAR(500) NOT NULL DEFAULT '',
    shipping_address TEXT NOT NULL,
    shipping_city TEXT NOT NULL,
    shipping_postal_code TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE order_item (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    order_id INT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    product_id INT NOT NULL REFERENCES products(id),
    title TEXT NOT NULL,
    unit_price INTEGER NOT NULL CHECK (unit_price >= 0),
    type TEXT,
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    color TEXT,
    height INTEGER CHECK (height >= 0),
    width INTEGER CHECK (width >= 0),
    material TEXT,
    thickness INTEGER CHECK (thickness >= 0)
);

-- +goose Down

DROP TABLE IF EXISTS order_item;
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS products;
DROP TABLE IF EXISTS users;