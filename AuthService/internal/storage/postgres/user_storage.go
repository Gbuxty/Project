package postgres

import (
	"Project/AuthService/internal/domain/models"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserStorage struct {
	db *pgxpool.Pool
}

func NewUserStorage(db *pgxpool.Pool) (*UserStorage, error) {
	return &UserStorage{
		db: db,
	}, nil
}

func (r *UserStorage) CloseDb() {
	r.db.Close()
}

func (r *UserStorage) CreateUser(ctx context.Context, email, password string) (uuid.UUID, error) {
	var id uuid.UUID

	query := `INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id`
	err := r.db.QueryRow(ctx, query, email, password).Scan(&id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create user: %w", err)
	}

	return id, nil
}

func (r *UserStorage) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User

	query := `SELECT id, email, password_hash FROM users WHERE email = $1`
	err := r.db.QueryRow(ctx, query, email).Scan(&user.ID, &user.Email, &user.PasswordHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return &user, nil
}

func (r *UserStorage) GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	var user models.User

	query := `SELECT id, email, password_hash FROM users WHERE id = $1`
	err := r.db.QueryRow(ctx, query, userID).Scan(&user.ID, &user.Email, &user.PasswordHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	return &user, nil
}

func (r *UserStorage) UserExists(ctx context.Context, email string) (bool, error) {
	var exists bool

	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`
	err := r.db.QueryRow(ctx, query, email).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check user existence: %w", err)
	}

	return exists, nil
}

func (r *UserStorage) SaveTokens(ctx context.Context, userID uuid.UUID, accessToken string, accessTokenExpiresAt time.Time, refreshToken string, refreshExp time.Time) error {
	query := `
        INSERT INTO users_tokens 
            (user_id, access_token, access_token_expires_at, refresh_token, refresh_token_expires_at) 
        VALUES 
            ($1, $2, $3, $4, $5)
        ON CONFLICT (user_id) 
        DO UPDATE SET 
            access_token = $2,
            access_token_expires_at = $3,
            refresh_token = $4,
            refresh_token_expires_at = $5
    `
	_, err := r.db.Exec(ctx, query, userID, accessToken, accessTokenExpiresAt, refreshToken, refreshExp)
	if err != nil {
		return fmt.Errorf("failed to save tokens: %w", err)
	}

	return nil
}

func (r *UserStorage) DeleteTokens(ctx context.Context, userID uuid.UUID) error {
	query := `
        UPDATE users_tokens 
        SET 
            access_token = NULL,
            access_token_expires_at = NULL,
            refresh_token = NULL,
            refresh_token_expires_at = NULL,
            deleted_at = NOW()
        WHERE user_id = $1
    `
	_, err := r.db.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete tokens: %w", err)
	}

	return nil
}

func (r *UserStorage) ConfirmEmail(ctx context.Context, email, code string) (uuid.UUID, error) {
	var userID uuid.UUID

	query := `
		UPDATE users
		SET email_confirmed = true
		WHERE id = (
			SELECT user_id
			FROM users_code
			WHERE confirmation_code = $1 AND confirmation_code_expires_at > NOW()
		)
		AND email = $2
		RETURNING id
	`
	err := r.db.QueryRow(ctx, query, code, email).Scan(&userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return uuid.Nil, fmt.Errorf("invalid confirmation code")
		}
		return uuid.Nil, fmt.Errorf("failed to confirm email: %w", err)
	}

	return userID, nil
}

func (r *UserStorage) SaveConfirmationCode(ctx context.Context, userID uuid.UUID, confirmationCode string, confirmCodeExpiresAt time.Time) error {
	query := `
		INSERT INTO users_code (user_id, confirmation_code, confirmation_code_expires_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id)
		DO UPDATE SET confirmation_code = $2, confirmation_code_expires_at = $3
	`
	_, err := r.db.Exec(ctx, query, userID, confirmationCode, confirmCodeExpiresAt)
	if err != nil {
		return fmt.Errorf("failed to save confirmation code: %w", err)
	}

	return nil
}

func (r *UserStorage) GetAccessToken(ctx context.Context, userID uuid.UUID) (string, time.Time, error) {
	var (
		accessToken    string
		accessTokenExp time.Time
	)

	query := `
		SELECT access_token, access_token_expires_at
		FROM users_tokens
		WHERE user_id = $1
	`
	err := r.db.QueryRow(ctx, query, userID).Scan(&accessToken, &accessTokenExp)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", time.Time{}, fmt.Errorf("access token not found")
		}
		return "", time.Time{}, fmt.Errorf("failed to get access token: %w", err)
	}

	return accessToken, accessTokenExp, nil
}
