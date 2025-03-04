package app

import (
	"AuthService/internal/config"
	"AuthService/internal/kafka"
	"AuthService/internal/logger"
	"AuthService/internal/service"
	"AuthService/internal/storage/postgres"
	"AuthService/internal/transport/handlers"
	"AuthService/internal/transport/server"

	"go.uber.org/zap"
)

type App struct {
	GRPCServer *server.GRPCServer
}

func New(log *logger.Logger, cfg *config.Config,kafkaProducer *kafka.Producer) *App {


	db, err := postgres.ConnectToDB(cfg.Postgres.StoragePath)
	if err != nil {
		log.Error("Failed to connect to database", zap.Error(err))
		return nil
	}

	userStorage, err := postgres.NewUserStorage(db)
	if err != nil {
		log.Error("Failed to init user storage", zap.Error(err))
		return nil
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
	if authHandlers == nil {
		return nil
	}

	grpcServer := server.NewGRPCServer(
		log.Logger,
		cfg.Grpc.Port,
		authHandlers,
	)

	return &App{
		GRPCServer: grpcServer,
	}
}

func (a *App) RunGRPC() error {

	return a.GRPCServer.Run()
}

func (a *App) StopGRPC() {

	a.GRPCServer.Stop()
}
