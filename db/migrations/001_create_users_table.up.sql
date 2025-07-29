-- Create users table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE
);

-- Create index on name for faster lookups
CREATE INDEX idx_users_name ON users(name); 