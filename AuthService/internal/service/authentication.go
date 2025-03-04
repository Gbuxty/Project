package service

import (
	"context"
	"fmt"
	"time"

	"AuthService/internal/domain/models"
	"AuthService/internal/kafka"
	"AuthService/internal/logger"

	"AuthService/internal/storage/postgres"
	"AuthService/pkg/jwt"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type AuthenticationService struct {
	userStorage     *postgres.UserStorage
	jwtSecretKey    string
	logger          *logger.Logger
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
	KafkapProducer  *kafka.Producer
}

func NewAuthenticationService(
	userRepo *postgres.UserStorage,
	jwtSecretKey string,
	accessTokenTTL time.Duration,
	refreshTokenTTL time.Duration,
	log *logger.Logger,
	KafkaProducer *kafka.Producer,

) *AuthenticationService {
	return &AuthenticationService{
		userStorage:     userRepo,
		jwtSecretKey:    jwtSecretKey,
		AccessTokenTTL:  accessTokenTTL,
		RefreshTokenTTL: refreshTokenTTL,
		logger:          log,
		KafkapProducer:  KafkaProducer,
	}
}

func (s *AuthenticationService) Register(ctx context.Context, email, password, repeatPassword string) error {
	s.logger.Info("Registering new user", zap.String("email", email))

	if err := s.ValidateRegister(ctx, email, password, repeatPassword); err != nil {
		s.logger.Error("Validation failed", zap.String("email", email), zap.Error(err))
		return err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("Failed to hash password", zap.String("email", email), zap.Error(err))
		return fmt.Errorf("failed to hash password: %w", err)
	}

	userID, err := s.userStorage.CreateUser(ctx, email, string(passwordHash))
	if err != nil {
		s.logger.Error("Failed to create user", zap.String("email", email), zap.Error(err))
		return fmt.Errorf("failed to register user: %w", err)
	}

	confirmationCode := jwt.GenerateConfirmationCode()
	confirmCodeExpiresAt := time.Now().Add(24 * time.Hour)

	if err := s.userStorage.SaveConfirmationCode(ctx, userID, confirmationCode, confirmCodeExpiresAt); err != nil {
		s.logger.Error("Failed to save confirmation code", zap.Error(err))
		return fmt.Errorf("failed to save confirmation code: %w", err)
	}

	message := map[string]string{
		"to_email": email,
		"subject":  "Welcome!",
		"body":     confirmationCode,
	}

	if err := s.KafkapProducer.SendMessage(ctx, email, message); err != nil {
		s.logger.Error("failed to send Confirm code to Kafka")
		return fmt.Errorf("failed to send confrim code to kafka%w", err)
	}

	s.logger.Info("Confirmation code sent to Kafka", zap.String("email", email))
	s.logger.Info("User registered successfully", zap.String("email", email))
	return nil
}

func (s *AuthenticationService) ValidateRegister(ctx context.Context, email, password, repeatPassword string) error {
	s.logger.Info("Validating registration request", zap.String("email", email))

	if password != repeatPassword {
		s.logger.Error("Passwords do not match", zap.String("email", email))
		return fmt.Errorf("passwords do not match")
	}

	exists, err := s.userStorage.UserExists(ctx, email)
	if err != nil {
		s.logger.Error("Failed to check email uniqueness", zap.String("email", email), zap.Error(err))
		return fmt.Errorf("failed to check email uniqueness: %w", err)
	}

	if exists {
		s.logger.Error("Email already exists", zap.String("email", email))
		return fmt.Errorf("email already exists")
	}

	s.logger.Info("Registration validation successful", zap.String("email", email))
	return nil
}

func (s *AuthenticationService) Login(ctx context.Context, email, password string) (*models.User, string, string, error) {
	s.logger.Info("Logging in user", zap.String("email", email))

	user, err := s.userStorage.GetUserByEmail(ctx, email)
	if err != nil {
		s.logger.Error("Failed to fetch user by email", zap.String("email", email), zap.Error(err))
		return nil, "", "", fmt.Errorf("failed to get user: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		s.logger.Error("Invalid password", zap.String("email", email))
		return nil, "", "", fmt.Errorf("invalid password")
	}

	accessToken, _, err := jwt.GenerateToken(user.ID, user.Email, s.jwtSecretKey, s.AccessTokenTTL)
	if err != nil {
		s.logger.Error("Failed to generate access token", zap.String("userID", user.ID.String()), zap.Error(err))
		return nil, "", "", fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, refreshTokenExpiresAt, err := jwt.GenerateToken(user.ID, user.Email, s.jwtSecretKey, s.RefreshTokenTTL)
	if err != nil {
		s.logger.Error("Failed to generate refresh token", zap.String("userID", user.ID.String()), zap.Error(err))
		return nil, "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	if err := s.userStorage.SaveTokens(ctx, user.ID, accessToken, refreshToken, refreshTokenExpiresAt); err != nil {
		s.logger.Error("Failed to save  token", zap.String("userID", user.ID.String()), zap.Error(err))
		return nil, "", "", fmt.Errorf("failed to save  token: %w", err)
	}

	s.logger.Info("User logged in successfully", zap.String("userID", user.ID.String()), zap.String("email", user.Email))
	return user, accessToken, refreshToken, nil
}

func (s *AuthenticationService) Logout(ctx context.Context, userID uuid.UUID) error {
	s.logger.Info("Logging out user", zap.String("userID", userID.String()))

	if err := s.userStorage.DeleteTokens(ctx, userID); err != nil {
		s.logger.Error("Failed to delete refresh token", zap.String("userID", userID.String()), zap.Error(err))
		return fmt.Errorf("failed to delete refresh token: %w", err)
	}

	s.logger.Info("User logged out successfully", zap.String("userID", userID.String()))
	return nil
}
func (s *AuthenticationService) Refresh(ctx context.Context, refreshToken string) (*models.User, string, string, error) {
	s.logger.Info("Refreshing tokens", zap.String("refreshToken", refreshToken))

	userID, expiresAt, err := s.userStorage.GetUserIDByRefreshToken(ctx, refreshToken)
	if err != nil {
		s.logger.Error("Invalid refresh token", zap.String("refreshToken", refreshToken), zap.Error(err))
		return nil, "", "", fmt.Errorf("invalid refresh token: %w", err)
	}

	if time.Now().After(expiresAt) {
		s.logger.Error("Refresh token expired", zap.String("userID", userID.String()))
		if err := s.userStorage.DeleteTokens(ctx, userID); err != nil {
			s.logger.Error("Failed to delete expired refresh token", zap.String("userID", userID.String()), zap.Error(err))
		}
		return nil, "", "", fmt.Errorf("refresh token expired")
	}

	user, err := s.userStorage.GetUserByID(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to fetch user by ID", zap.String("userID", userID.String()), zap.Error(err))
		return nil, "", "", fmt.Errorf("failed to fetch user: %w", err)
	}

	accessToken, _, err := jwt.GenerateToken(user.ID, user.Email, s.jwtSecretKey, s.AccessTokenTTL)
	if err != nil {
		s.logger.Error("Failed to generate access token", zap.String("userID", user.ID.String()), zap.Error(err))
		return nil, "", "", fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, refreshTokenExpiresAt, err := jwt.GenerateToken(user.ID, user.Email, s.jwtSecretKey, s.RefreshTokenTTL)
	if err != nil {
		s.logger.Error("Failed to generate refresh token", zap.String("userID", user.ID.String()), zap.Error(err))
		return nil, "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	if err := s.userStorage.SaveTokens(ctx, user.ID, accessToken, refreshToken, refreshTokenExpiresAt); err != nil {
		s.logger.Error("Failed to save access token", zap.String("userID", user.ID.String()), zap.Error(err))
		return nil, "", "", fmt.Errorf("failed to save access token: %w", err)
	}

	s.logger.Info("Tokens refreshed successfully", zap.String("userID", user.ID.String()))
	return user, accessToken, refreshToken, nil
}

func (s *AuthenticationService) Me(ctx context.Context, accessToken string) (*models.User, error) {
	s.logger.Info("Fetching user info", zap.String("accessToken", accessToken))

	userID, err := jwt.ExtractUserIDFromToken(accessToken, s.jwtSecretKey)
	if err != nil {
		s.logger.Error("Failed to extract user ID from token", zap.String("accessToken", accessToken), zap.Error(err))
		return nil, fmt.Errorf("failed to extract user ID: %w", err)
	}

	user, err := s.userStorage.GetUserByID(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to fetch user by ID", zap.String("userID", userID.String()), zap.Error(err))
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	s.logger.Info("User info fetched successfully", zap.String("userID", user.ID.String()), zap.String("email", user.Email))
	return user, nil
}

func (s *AuthenticationService) ConfirmEmail(ctx context.Context, email, confirmationCode string) (uuid.UUID, error) {
	s.logger.Info("Confirming email", zap.String("email", email), zap.String("confirmation code", confirmationCode))

	userID, err := s.userStorage.ConfirmEmail(ctx, email, confirmationCode)
	if err != nil {
		s.logger.Error("Failed to confirm email", zap.Error(err))
		return uuid.Nil, fmt.Errorf("failed to confirm email: %w", err)
	}

	s.logger.Info("Email confirmed successfully", zap.String("userID", userID.String()))
	return userID, nil
}
