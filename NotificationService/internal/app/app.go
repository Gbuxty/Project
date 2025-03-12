package app

import (
	"Project/NotificationService/internal/config"
	"Project/NotificationService/internal/kafka"
	"Project/NotificationService/internal/logger"
	"Project/NotificationService/internal/mailer"
	"context"
	

	
)

type App struct {
	logger *logger.Logger
	config *config.Config
	mailer *mailer.Mailer
	consumer *kafka.Consumer
}

func New(logger *logger.Logger, cfg *config.Config) (*App, error) {
	notificationService := mailer.NewMailer(
		cfg.Mailer.ApiURL,
		cfg.Mailer.ApiToken,
		cfg.Mailer.FromEmail,
	)

	kafkaConsumer := kafka.NewConsumer(
		[]string{cfg.Kafka.Broker},
		cfg.Kafka.Topic,
		cfg.Kafka.GroupID,
		notificationService,
		logger,
	)

	return &App{
		logger: logger,
		config: cfg,
		mailer: notificationService,
		consumer: kafkaConsumer,
	}, nil

}



func (a *App) Run(ctx context.Context) error {
    go a.consumer.Start(ctx)
    a.logger.Info("Application started")
    return nil
}
func (a *App) Stop(){
	a.consumer.Close()
	a.logger.Info("Application stopped")
}