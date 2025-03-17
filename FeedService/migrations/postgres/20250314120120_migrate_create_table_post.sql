-- +goose Up
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- Создаем таблицу posts с UUID
CREATE TABLE posts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), 
    user_id UUID NOT NULL, 
    content TEXT NOT NULL,
    image_url TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE INDEX idx_posts_user_id ON posts(user_id);
CREATE INDEX idx_posts_created_at ON posts(created_at);

-- +goose Down
DROP INDEX IF EXISTS idx_posts_user_id;
DROP INDEX IF EXISTS idx_posts_created_at;

DROP TABLE IF EXISTS posts;

DROP EXTENSION IF EXISTS pgcrypto;