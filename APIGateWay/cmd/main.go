package main

import (
	"Project/APIGateWay/internal/config"
	"Project/APIGateWay/internal/logger"
	"Project/APIGateWay/internal/server"
	"Project/APIGateWay/internal/service"
	"fmt"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	pathToConfig := config.InitFlags()

	log, err := logger.New()
	if err != nil {
		fmt.Errorf("Failed load logger")
		return
	}

	cfg, err := config.LoadConfig(pathToConfig)
	if err != nil {
		log.Fatal("Failed load config")
	}

	authConn, err := grpc.NewClient(cfg.AuthServiceAdress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("Failed to connect to AuthService: ", zap.Error(err))
	}
	defer authConn.Close()

	feedConn, err := grpc.NewClient(cfg.FeedServiceAdress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("Failed to connect to FeedService: ", zap.Error(err))
	}
	defer feedConn.Close()

	authService := service.NewAuthService(authConn)
	feedService := service.NewFeedService(feedConn)

	server := server.NewServer(authService, feedService)

	log.Info("Starting API Gateway on :8080")
	if err := server.Start(cfg.HttpServerAdress); err != nil {
		log.Fatal("Failed to start server: ", zap.Error(err))
	}

}
