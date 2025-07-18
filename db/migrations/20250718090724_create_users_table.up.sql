CREATE TABLE IF NOT EXISTS users (
    id             UUID PRIMARY KEY,
    name           VARCHAR(100),
    email          VARCHAR(100) UNIQUE,
    password_hash  VARCHAR(255),
    role           VARCHAR(20) CHECK (role IN ('admin', 'manager', 'user')) NOT NULL DEFAULT 'user',
    created_at     TIMESTAMP DEFAULT NOW(),
    updated_at     TIMESTAMP DEFAULT NOW()
);
