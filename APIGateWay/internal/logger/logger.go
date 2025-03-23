package logger

import (
	"fmt"

	"go.uber.org/zap"
)

type Logger struct {
	*zap.Logger
}

func New() (*Logger, error) {

	zapLogger, err := zap.NewProduction()
	if err != nil {
		return nil, fmt.Errorf("Failed init new Logger")
	}
	return &Logger{Logger: zapLogger}, nil
}
