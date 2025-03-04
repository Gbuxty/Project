package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"NotificationService/internal/app"
	"NotificationService/internal/config"
	"NotificationService/internal/kafka"
	"NotificationService/internal/logger"
	"NotificationService/internal/mailer"

	"go.uber.org/zap"
)

func main() {
	log, err := logger.New()
	if err != nil {
		fmt.Println("Failed to initialize logger:", err)
		return
	}
	defer log.Sync()

	cfg, err := config.LoadConfig("config/local.yaml")
	if err != nil {
		log.Fatal("Failed to load config", zap.Error(err))
	}

	

	notificationService := mailer.NewMailer(cfg.Mailer.ApiURL, cfg.Mailer.ApiToken, cfg.Mailer.FromEmail, log)

	kafkaConsumer := kafka.NewConsumer(
		[]string{cfg.Kafka.Broker},
		cfg.Kafka.Topic,
		cfg.Kafka.GroupID,
		notificationService,
		log,
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)


	var wg sync.WaitGroup
	wg.Add(1) 

	
	go func() {
		defer wg.Done() 
		kafkaConsumer.Start(ctx)
	}()
		

	app := app.New(log, cfg, notificationService)
	if err := app.Run(); err != nil {
		log.Fatal("failed to run app", zap.Error(err))
	}

	
	<-stop
	log.Info("Shutting down gracefully...")

	
	
	wg.Wait()
	log.Info("Service stopped")
}