package service

import (
	"context"
	"Project/proto/gen"
	
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


func (s *AuthService) Register(ctx context.Context, req *gen.RegisterRequest) (*gen.RegisterResponse, error) {
	return s.client.Register(ctx, req)
}


func (s *AuthService) Login(ctx context.Context, req *gen.LoginRequest) (*gen.LoginResponse, error) {
	return s.client.Login(ctx, req)
}


func (s *AuthService) Logout(ctx context.Context, req *gen.LogoutRequest) (*gen.LogoutResponse, error) {
	return s.client.Logout(ctx, req)
}


func (s *AuthService) Refresh(ctx context.Context, req *gen.RefreshRequest) (*gen.RefreshResponse, error) {
	return s.client.Refresh(ctx, req)
}


func (s *AuthService) Me(ctx context.Context, req *gen.MeRequest) (*gen.MeResponse, error) {
	return s.client.Me(ctx, req)
}


func (s *AuthService) ConfirmEmail(ctx context.Context, req *gen.ConfirmEmailRequest) (*gen.ConfirmEmailResponse, error) {
	return s.client.ConfirmEmail(ctx, req)
}