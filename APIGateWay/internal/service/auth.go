package service

import (
	"Project/APIGateWay/internal/domain"
	"Project/proto/gen"
	"context"

	"google.golang.org/grpc"
)

type AuthService struct {
	client gen.AuthenticationClient
}

func NewAuthService(conn *grpc.ClientConn) *AuthService {
	return &AuthService{
		client: gen.NewAuthenticationClient(conn),
	}
}

func (s *AuthService) Register(ctx context.Context, req *domain.RegisterRequest) (bool, error) {
	grpcReq := &gen.RegisterRequest{
		Email:          req.Email,
		Password:       req.Password,
		RepeatPassword: req.Password,
	}

	res, err := s.client.Register(ctx, grpcReq)
	if err != nil {
		return false, err
	}

	return res.Success, nil
}

func (s *AuthService) Login(ctx context.Context, req *domain.LoginRequest) (*domain.TokenPair, *domain.User, error) {
	grpcReq := &gen.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	res, err := s.client.Login(ctx, grpcReq)
	if err != nil {
		return nil, nil, err
	}

	tokenPair := &domain.TokenPair{
		AccessToken:  res.AccessToken.Token,
		RefreshToken: res.RefreshToken.Token,
	}

	user := &domain.User{
		ID:    res.User.Id,
		Email: res.User.Email,
	}

	return tokenPair, user, nil
}

func (s *AuthService) Logout(ctx context.Context, userID string) (bool, error) {
	grpcReq := &gen.LogoutRequest{
		Id: userID,
	}

	res, err := s.client.Logout(ctx, grpcReq)
	if err != nil {
		return false, err
	}

	return res.Success, nil
}

func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (*domain.TokenPair, error) {
	grpcReq := &gen.RefreshRequest{
		RefreshToken: refreshToken,
	}

	res, err := s.client.Refresh(ctx, grpcReq)
	if err != nil {
		return nil, err
	}

	return &domain.TokenPair{
		AccessToken:  res.AccessToken.Token,
		RefreshToken: res.RefreshToken.Token,
	}, nil
}

func (s *AuthService) Me(ctx context.Context, accessToken string) (*domain.User, error) {
	grpcReq := &gen.MeRequest{
		AccessToken: accessToken,
	}

	res, err := s.client.Me(ctx, grpcReq)
	if err != nil {
		return nil, err
	}

	return &domain.User{
		ID:    res.User.Id,
		Email: res.User.Email,
	}, nil
}

func (s *AuthService) ConfirmEmail(ctx context.Context, req *domain.ConfirmationRequest) (bool, error) {
	grpcReq := &gen.ConfirmEmailRequest{
		Email:            req.Email,
		ConfirmationCode: req.ConfirmationCode,
	}

	res, err := s.client.ConfirmEmail(ctx, grpcReq)
	if err != nil {
		return false, err
	}

	return res.Success, nil
}
