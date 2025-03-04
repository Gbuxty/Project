-- +goose Up
ALTER TABLE users_tokens
DROP COLUMN access_token_expires_at;
-- +goose Down

ALTER TABLE users_tokens
ADD COLUMN access_token_expires_at  TIMESTAMP;
 