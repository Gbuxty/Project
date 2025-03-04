package main

import (
	"AuthService/internal/app"
	"AuthService/internal/config"
	"AuthService/internal/kafka"
	"AuthService/internal/logger"
	"context"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

func main() {

	log, err := logger.New()
	if err != nil {
		log.Fatal("Failed to logger", zap.Error(err))
	}
	defer log.Sync()

	cfg, err := config.LoadConfig("config/local.yaml")
	if err != nil {
		log.Fatal("Failed to load config", zap.Error(err))
	}

	kafkaProducer := kafka.NewProducer([]string{cfg.Kafka.Broker},
		cfg.Kafka.Topic,
		log,
	)

	application := app.New(log, cfg,kafkaProducer)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := application.GRPCServer.Run(); err != nil {
			log.Error("gRPC server failed", zap.Error(err))
			cancel()
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
