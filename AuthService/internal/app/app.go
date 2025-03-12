package app

import (
	"Project/AuthService/internal/config"
	"Project/AuthService/internal/kafka"
	"Project/AuthService/internal/logger"
	"Project/AuthService/internal/service"
	"Project/AuthService/internal/storage/postgres"
	"Project/AuthService/internal/transport/handlers"
	"Project/AuthService/internal/transport/server"
	"Project/AuthService/pkg/database"
	
	"fmt"
)

type App struct {
	GRPCServer *server.GRPCServer
}

func New(log *logger.Logger, cfg *config.Config,kafkaProducer *kafka.Producer) (*App,error) {
	db, err := database.ConnectToDB(cfg.Postgres.StoragePath)
	if err != nil {
		return nil,fmt.Errorf("Failed connect to db:%w",err)
	}

	userStorage, err := postgres.NewUserStorage(db)
	if err != nil {
		return nil,fmt.Errorf("Failed to init user repositories:%w",err)
	}

	authService := service.NewAuthenticationService(
		userStorage,
		cfg.Auth.SecretKey,
		cfg.Auth.AccessTokenTTL,
		cfg.Auth.RefreshTokenTTL,
		log,
		kafkaProducer,
		
	)

	authHandlers := handlers.NewAuthHandlers(authService, log.Logger)

	grpcServer := server.NewGRPCServer(
		log.Logger,
		cfg.Grpc.Port,
		authHandlers,
	)

	return &App{
		GRPCServer: grpcServer,
	},nil
}

func (a *App) RunGRPC() error {
	return a.GRPCServer.Run()
}

func (a *App) StopGRPC() {
	a.GRPCServer.Stop()
}
