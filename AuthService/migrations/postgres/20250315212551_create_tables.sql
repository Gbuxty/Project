-- +goose Up
DROP TABLE IF EXISTS users_tokens;
DROP TABLE IF EXISTS users;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";


CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    email_confirmed BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE users_tokens (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    access_token VARCHAR(255),
    refresh_token VARCHAR(255),
    refresh_token_expires_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    access_token_expires_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);


CREATE TABLE users_code (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    confirmation_code VARCHAR(255),
    confirmation_code_expires_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE INDEX idx_users_tokens_user_id ON users_tokens(user_id);
CREATE INDEX idx_users_code_user_id ON users_code (user_id);
CREATE INDEX idx_users_email ON users(email);

-- +goose Down
DROP TABLE IF EXISTS users_code;
DROP TABLE IF EXISTS users_tokens;
DROP TABLE IF EXISTS users;