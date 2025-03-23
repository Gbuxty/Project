package service

import (
	"context"
	"fmt"
	"time"

	"Project/AuthService/internal/domain/message"
	"Project/AuthService/internal/domain/models"
	"Project/AuthService/internal/kafka"
	"Project/AuthService/internal/logger"

	"Project/AuthService/internal/storage/postgres"
	"Project/AuthService/pkg/jwt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthenticationService struct {
	userStorage     UserStorage
	jwtSecretKey    string
	logger          *logger.Logger
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
	KafkapProducer  MessageBroker
	redisClient     RedisRepositories
}

type RedisRepositories interface {
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
}

type MessageBroker interface {
	SendMessage(ctx context.Context, key string, value interface{}) error
}

type UserStorage interface {
	CreateUser(ctx context.Context, email, password string) (uuid.UUID, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error)
	UserExists(ctx context.Context, email string) (bool, error)
	SaveTokens(ctx context.Context, userID uuid.UUID, accessToken string, accessExp time.Time, refreshToken string, refreshExp time.Time) error
	DeleteTokens(ctx context.Context, userID uuid.UUID) error
	ConfirmEmail(ctx context.Context, email, code string) (uuid.UUID, error)
	SaveConfirmationCode(ctx context.Context, userID uuid.UUID, confirmationCode string, confirmCodeExpiresAt time.Time) error
	GetAccessToken(ctx context.Context, userID uuid.UUID) (string, time.Time, error)
}

func NewAuthenticationService(
	userRepo *postgres.UserStorage,
	jwtSecretKey string,
	accessTokenTTL time.Duration,
	refreshTokenTTL time.Duration,
	log *logger.Logger,
	KafkaProducer *kafka.Producer,
	redisClient RedisRepositories,
) *AuthenticationService {
	return &AuthenticationService{
		userStorage:     userRepo,
		jwtSecretKey:    jwtSecretKey,
		AccessTokenTTL:  accessTokenTTL,
		RefreshTokenTTL: refreshTokenTTL,
		logger:          log,
		KafkapProducer:  KafkaProducer,
		redisClient:     redisClient,
	}
}

func (s *AuthenticationService) Register(ctx context.Context, email, password, repeatPassword string) error {
	if err := s.ValidateRegister(ctx, email, password, repeatPassword); err != nil {
		return fmt.Errorf("Failed validate register;%w", err)
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	userID, err := s.userStorage.CreateUser(ctx, email, string(passwordHash))
	if err != nil {
		return fmt.Errorf("failed to register user: %w", err)
	}

	confirmationCode := jwt.GenerateConfirmationCode()
	confirmCodeExpiresAt := time.Now().Add(24 * time.Hour)

	if err := s.userStorage.SaveConfirmationCode(ctx, userID, confirmationCode, confirmCodeExpiresAt); err != nil {
		return fmt.Errorf("failed to save confirmation code: %w", err)
	}

	message := message.ConfirmationMessage{
		ToEmail: email,
		Subject: "Welcome!",
		Body:    confirmationCode,
	}

	if err := s.KafkapProducer.SendMessage(ctx, email, message); err != nil {
		return fmt.Errorf("failed to send confrim code to kafka:%w", err)
	}

	return nil
}

func (s *AuthenticationService) ValidateRegister(ctx context.Context, email, password, repeatPassword string) error {
	if password != repeatPassword {
		return fmt.Errorf("passwords do not match")
	}

	exists, err := s.userStorage.UserExists(ctx, email)
	if err != nil {
		return fmt.Errorf("failed to check email uniqueness: %w", err)
	}

	if exists {
		return fmt.Errorf("email already exists")
	}
	return nil
}

func (s *AuthenticationService) Login(ctx context.Context, email, password string) (*models.User, string, string, error) {
	user, err := s.userStorage.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to get user: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, "", "", fmt.Errorf("invalid password")
	}

	accessToken, accessTokenExpiresAt, err := jwt.GenerateToken(user.ID, user.Email, s.jwtSecretKey, s.AccessTokenTTL)
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, refreshTokenExpiresAt, err := jwt.GenerateToken(user.ID, user.Email, s.jwtSecretKey, s.RefreshTokenTTL)
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	if err := s.userStorage.SaveTokens(ctx, user.ID, accessToken, accessTokenExpiresAt, refreshToken, refreshTokenExpiresAt); err != nil {
		return nil, "", "", fmt.Errorf("failed to save  token: %w", err)
	}

	userIDKey := user.ID.String()
	ttl := time.Until(accessTokenExpiresAt)
	if err := s.redisClient.Set(ctx, userIDKey, accessToken, ttl); err != nil {
		return nil, "", "", fmt.Errorf("failed to save access token to Redis: %w", err)
	}

	return user, accessToken, refreshToken, nil
}

func (s *AuthenticationService) Logout(ctx context.Context, userID uuid.UUID) error {
	userIDKey := userID.String()
	if err := s.redisClient.Delete(ctx, userIDKey); err != nil {
		return fmt.Errorf("failed to delete access token from Redis: %w", err)
	}

	if err := s.userStorage.DeleteTokens(ctx, userID); err != nil {
		return fmt.Errorf("failed to delete refresh token: %w", err)
	}

	return nil
}

func (s *AuthenticationService) Refresh(ctx context.Context, refreshToken string) (*models.User, string, string, error) {
	userID, err := jwt.ExtractUserIDFromToken(refreshToken, s.jwtSecretKey)
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to extract user_id from token")
	}
	user, err := s.userStorage.GetUserByID(ctx, userID)
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to fetch user: %w", err)
	}

	accessToken, accessTokenExpiresAt, err := jwt.GenerateToken(user.ID, user.Email, s.jwtSecretKey, s.AccessTokenTTL)
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, refreshTokenExpiresAt, err := jwt.GenerateToken(user.ID, user.Email, s.jwtSecretKey, s.RefreshTokenTTL)
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	if err := s.userStorage.SaveTokens(ctx, user.ID, accessToken, accessTokenExpiresAt, refreshToken, refreshTokenExpiresAt); err != nil {
		return nil, "", "", fmt.Errorf("failed to save access token: %w", err)
	}

	userIDKey := user.ID.String()
	ttl := time.Until(accessTokenExpiresAt)
	if err := s.redisClient.Set(ctx, userIDKey, accessToken, ttl); err != nil {
		return nil, "", "", fmt.Errorf("failed to save access token to Redis: %w", err)
	}

	return user, accessToken, refreshToken, nil
}

func (s *AuthenticationService) Me(ctx context.Context, accessToken string) (*models.User, error) {
	userID, err := jwt.ExtractUserIDFromToken(accessToken, s.jwtSecretKey)
	if err != nil {
		return nil, fmt.Errorf("failed to extract user ID: %w", err)
	}

	userIDKey := userID.String()
	storedAccessToken, err := s.redisClient.Get(ctx, userIDKey)
	if err != nil {
		storedAccessToken, expiresAt, err := s.userStorage.GetAccessToken(ctx, userID)
		if err != nil {
			return nil, fmt.Errorf("invalid or expired access token: %w", err)
		}

		ttl := time.Until(expiresAt)
		if err := s.redisClient.Set(ctx, userIDKey, storedAccessToken, ttl); err != nil {
			return nil, fmt.Errorf("failed to save access token to Redis: %w", err)
		}
	}

	if storedAccessToken != accessToken {
		return nil, fmt.Errorf("invalid access token")
	}

	user, err := s.userStorage.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	return user, nil
}

func (s *AuthenticationService) ConfirmEmail(ctx context.Context, email, confirmationCode string) (uuid.UUID, error) {
	userID, err := s.userStorage.ConfirmEmail(ctx, email, confirmationCode)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to confirm email: %w", err)
	}

	return userID, nil
}
