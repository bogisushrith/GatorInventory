-- Initial database schema for Inventory Management System

-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL
);

-- Create products table
CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    quantity INTEGER NOT NULL,
    category VARCHAR(100) NOT NULL
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_products_category ON products(category);

-- Insert sample user (password is hashed version of 'Test1234')
-- Note: You should hash passwords properly in your application
INSERT INTO users (username, password, role) 
VALUES ('admin', '$2a$10$YourHashedPasswordHere', 'ADMIN')
ON CONFLICT (username) DO NOTHING;
