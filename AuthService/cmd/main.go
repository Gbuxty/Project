package main

import (
	"Project/AuthService/internal/app"
	"Project/AuthService/internal/config"
	"Project/AuthService/internal/kafka"
	"Project/AuthService/internal/logger"
	"context"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

func main() {
	configPath:=config.InitFlags()

	log, err := logger.New()
	if err != nil {
		log.Fatal("Failed to logger", zap.Error(err)) //вот тут логфатал не сработает если ошибка .он не инициализировался
	}

	defer log.Sync()

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatal("Failed to load config", zap.Error(err))
	}

	kafkaProducer := kafka.NewProducer([]string{cfg.Kafka.Broker},
		cfg.Kafka.Topic,
		log,
	)

	application,err := app.New(log, cfg, kafkaProducer)
	if err!=nil{
		log.Fatal("Failed to load application", zap.Error(err))
	}

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
