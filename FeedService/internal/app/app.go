package app

import (
	"Project/FeedService/internal/config"
	"Project/FeedService/internal/repositories/postgres"
	"Project/FeedService/internal/service"
	"Project/FeedService/internal/transport/handlers"
	"Project/FeedService/internal/transport/server"
	"Project/FeedService/pkg/database"
	"fmt"

	"Project/FeedService/internal/logger"
)

type App struct {
	GRPCServer *server.GRPCServer
}

func NewApp(log *logger.Logger, cfg *config.Config) (*App, error) {
	db, err := database.ConnectToDB(cfg.Postgres.StoragePath)
	if err != nil {
		return nil, fmt.Errorf("Failed connect to db:%w", err)
	}

	feedRepo, err := postgres.NewPostRepositories(db)
	if err != nil {
		return nil, fmt.Errorf("Failed to init user repositories:%w", err)
	}

	feedService:=service.NewFeedService(feedRepo,cfg.Auth.SecretKey)

	feedHandlers:=handlers.NewFeedHandlers(feedService,log.Logger,cfg.Auth.SecretKey)

	grpsServer:=server.NewGRPCServer(log.Logger,cfg.Grpc.Port,feedHandlers)

	return &App{GRPCServer: grpsServer}, nil
}
