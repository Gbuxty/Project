package logger

import (
	"go.uber.org/zap"
)

type Logger struct {
	*zap.Logger
}

func New() (*Logger, error) {

	zapLogger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	return &Logger{Logger: zapLogger}, nil
}

func (l *Logger) Sync() {
	_ = l.Logger.Sync()
}
