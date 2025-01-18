CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    token VARCHAR(255) NULL, -- Added token field
    name VARCHAR(255), -- Nullable
    preference VARCHAR(255), -- Nullable
    weight_unit VARCHAR(255), -- Nullable
    height_unit VARCHAR(255), -- Nullable
    weight FLOAT DEFAULT 1, -- Default to 1
    height FLOAT DEFAULT 2, -- Default to 2
    image_uri VARCHAR(255), -- Nullable
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email ON users (email);
CREATE INDEX IF NOT EXISTS idx_users_email_password ON users (email, password);