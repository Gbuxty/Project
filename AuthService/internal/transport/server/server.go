package server

import (
	"fmt"
	"net"

	"AuthService/internal/transport/handlers"
	"AuthService/proto/gen"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type GRPCServer struct {
	log    *zap.Logger
	server *grpc.Server
	port   int
}

func NewGRPCServer(
	log *zap.Logger,
	port int,
	authHandlers *handlers.AuthHandlers,
) *GRPCServer {
	grpcServer := grpc.NewServer()
	gen.RegisterAuthenticationServer(grpcServer, authHandlers)
	
	return &GRPCServer{
		log:    log,
		server: grpcServer,
		port:   port,
	}
}

func (s *GRPCServer) Run() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	s.log.Info("Starting gRPC server", zap.Int("port", s.port))

	if err := s.server.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %w", err)
	}
	return nil
}

func (s *GRPCServer) Stop() {
	s.log.Info("Stopping gRPC server")
	s.server.GracefulStop()
}