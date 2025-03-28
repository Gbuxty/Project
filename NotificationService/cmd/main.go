package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"Project/NotificationService/internal/app"
	"Project/NotificationService/internal/config"
	"Project/NotificationService/internal/logger"

	"go.uber.org/zap"
)

func main() {
	configPath:=config.InitFlags()

	log, err := logger.New()
	if err != nil {
		fmt.Errorf("Failed to load logger %w",err)
		return
	}

	defer log.Sync()

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatal("Failed to load config", zap.Error(err))
	}

	app, err := app.New(log, cfg)
	if err != nil {
		log.Fatal("Failed to load app", zap.Error(err))
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := app.Run(ctx); err != nil {
		log.Fatal("Failed to start app", zap.Error(err))
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	select {
	case <-stop:
		log.Info("Shutting down gracefully...")
		app.Stop()
	case <-ctx.Done():
		log.Info("Forced shutdown...")
	}

	log.Info("Notification stopped")
}
 