package main

import (
	"Project/FeedService/internal/app"
	"Project/FeedService/internal/config"
	"Project/FeedService/internal/logger"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

func main() {
	configPath := config.InitFlags()

	log, err := logger.New()
	if err != nil {
		fmt.Errorf("Failed load to logger %w", err)
		return
	}

	defer log.Sync()

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatal("Failed to load config", zap.Error(err))
	}

	application, err := app.NewApp(log, cfg)
	if err != nil {
		log.Fatal("Failed to initialize application", zap.Error(err))
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := application.GRPCServer.Run(); err != nil {
			log.Error("GRPC server failed", zap.Error(err))
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	select {
	case <-stop:
		log.Info("Shutting down gracefully...")
		application.GRPCServer.Stop()
	case <-ctx.Done():
		log.Info("Forced shutdown...")
	}

}
