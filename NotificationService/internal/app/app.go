package app

import (
	"NotificationService/internal/config"
	"NotificationService/internal/logger"
	"NotificationService/internal/mailer"
)

type App struct {
	logger *logger.Logger
	config *config.Config
	mailer *mailer.Mailer
}

func New(logger *logger.Logger, cfg *config.Config, mailer *mailer.Mailer) *App {
	return &App{
		logger: logger,
		config: cfg,
		mailer: mailer,
	}

}


func (a *App) Run() error {
	a.logger.Info("Mailer service started")
	return nil
}