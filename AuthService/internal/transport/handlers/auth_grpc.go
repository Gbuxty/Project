package handlers

import (
	"context"
	
	"fmt"
	"time"

	"AuthService/internal/domain/models"
	"AuthService/internal/service"
	"AuthService/proto/gen"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type AuthHandlers struct {
	service *service.AuthenticationService
	logger  *zap.Logger
	gen.UnimplementedAuthenticationServer
}

func NewAuthHandlers(service *service.AuthenticationService, logger *zap.Logger) *AuthHandlers {
	return &AuthHandlers{
		service: service,
		logger:  logger,
	}
}

func (h *AuthHandlers) Register(ctx context.Context, req *gen.RegisterRequest) (*gen.RegisterResponse, error) {
	h.logger.Info("Registering new user", zap.String("email", req.Email))

	if err := h.service.Register(ctx, req.Email, req.Password, req.RepeatPassword); err != nil {
		h.logger.Error("Failed to register user", zap.String("email", req.Email), zap.Error(err))
		return nil, fmt.Errorf("failed to register user: %w", err)
	}

	h.logger.Info("User registered successfully", zap.String("email", req.Email))


	return &gen.RegisterResponse{Success: true}, nil
}

func (h *AuthHandlers) Login(ctx context.Context, req *gen.LoginRequest) (*gen.LoginResponse, error) {
	h.logger.Info("Logging in user", zap.String("email", req.Email))

	
	user, accessToken, refreshToken, err := h.service.Login(ctx, req.Email, req.Password)
	if err != nil {
		h.logger.Error("Failed to login user", zap.String("email", req.Email), zap.Error(err))
		return nil, fmt.Errorf("failed to login user: %w", err)
	}

	response := mapUserToLoginResponse(user, accessToken, refreshToken, h.service.AccessTokenTTL, h.service.RefreshTokenTTL)

	h.logger.Info("User logged in successfully", zap.String("userID", user.ID.String()), zap.String("email", user.Email))
	return response, nil
	
}

func mapUserToLoginResponse(user *models.User, accessToken, refreshToken string, accessTokenTTL, refreshTokenTTL time.Duration) *gen.LoginResponse {
	return &gen.LoginResponse{
		User: &gen.User{
			Id:    user.ID.String(),
			Email: user.Email,
		},
		AccessToken: &gen.AccessToken{
			Token:     accessToken,
			ExpiresAt: time.Now().Add(accessTokenTTL).Unix(),
		},
		RefreshToken: &gen.RefreshToken{
			Token:     refreshToken,
			ExpiresAt: time.Now().Add(refreshTokenTTL).Unix(),
		},
	}
}
func (h *AuthHandlers) Logout(ctx context.Context, req *gen.LogoutRequest) (*gen.LogoutResponse, error) {
	userID, err := uuid.Parse(req.User.Id)
	if err != nil {
		h.logger.Error("Invalid user ID", zap.String("userID", req.User.Id), zap.Error(err))
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	h.logger.Info("Logging out user", zap.String("userID", userID.String()))

	
	if err := h.service.Logout(ctx, userID); err != nil {
		h.logger.Error("Failed to logout user", zap.String("userID", userID.String()), zap.Error(err))
		return nil, fmt.Errorf("failed to logout user: %w", err)
	}

	h.logger.Info("User logged out successfully", zap.String("userID", userID.String()))
	return &gen.LogoutResponse{Success: true}, nil
}



func (h *AuthHandlers) Refresh(ctx context.Context, req *gen.RefreshRequest) (*gen.RefreshResponse, error) {
	h.logger.Info("Refreshing tokens", zap.String("refreshToken", req.RefreshToken))

	user, accessToken, refreshToken, err := h.service.Refresh(ctx, req.RefreshToken)
	if err != nil {
		h.logger.Error("Failed to refresh tokens", zap.String("refreshToken", req.RefreshToken), zap.Error(err))
		return nil, fmt.Errorf("failed to refresh tokens: %w", err)
	}

	
	response := mapTokensToRefreshResponse(accessToken, refreshToken, h.service.AccessTokenTTL, h.service.RefreshTokenTTL)

	h.logger.Info("Tokens refreshed successfully", zap.String("userID", user.ID.String()))
	return response, nil
}

func mapTokensToRefreshResponse(accessToken, refreshToken string, accessTokenTTL, refreshTokenTTL time.Duration) *gen.RefreshResponse {
	return &gen.RefreshResponse{
		AccessToken: &gen.AccessToken{
			Token:     accessToken,
			ExpiresAt: time.Now().Add(accessTokenTTL).Unix(),
		},
		RefreshToken: &gen.RefreshToken{
			Token:     refreshToken,
			ExpiresAt: time.Now().Add(refreshTokenTTL).Unix(),
		},
	}
}

func (h *AuthHandlers) Me(ctx context.Context, req *gen.MeRequest) (*gen.MeResponse, error) {
	h.logger.Info("Fetching user info", zap.String("accessToken", req.AccessToken))

	user, err := h.service.Me(ctx, req.AccessToken)
	if err != nil {
		h.logger.Error("Failed to fetch user info", zap.String("accessToken", req.AccessToken), zap.Error(err))
		return nil, fmt.Errorf("failed to fetch user info: %w", err)
	}

	
	response := mapUserToMeResponse(user)

	h.logger.Info("User info fetched successfully", zap.String("userID", user.ID.String()), zap.String("email", user.Email))
	return response, nil
}

func mapUserToMeResponse(user *models.User)*gen.MeResponse{
	return &gen.MeResponse{
		User: &gen.User{
			Id:    user.ID.String(), 
			Email: user.Email,
		},
	}
}

func (h *AuthHandlers) ConfirmEmail(ctx context.Context, req *gen.ConfirmEmailRequest) (*gen.ConfirmEmailResponse, error) {
	h.logger.Info("Confirming email", zap.String("email", req.Email), zap.String("confirmation_code", req.ConfirmationCode))

	
	userID, err := h.service.ConfirmEmail(ctx, req.Email, req.ConfirmationCode)
	if err != nil {
		h.logger.Error("Failed to confirm email", zap.String("email", req.Email), zap.Error(err))
		return nil, fmt.Errorf("failed to confirm email: %w", err)
	}

	h.logger.Info("Email confirmed successfully", zap.String("userID", userID.String()))
	return &gen.ConfirmEmailResponse{Success: true}, nil
}
