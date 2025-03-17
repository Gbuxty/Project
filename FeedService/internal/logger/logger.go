package logger

import (
	"fmt"
	"os"
	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)


type Logger struct {
	*zap.Logger
}


func colorCode(level zapcore.Level) string {
	switch level {
	case zapcore.DebugLevel:
		return "\033[36m" 
	case zapcore.InfoLevel:
		return "\033[32m" 
	case zapcore.WarnLevel:
		return "\033[33m" 
	case zapcore.ErrorLevel, zapcore.DPanicLevel, zapcore.PanicLevel, zapcore.FatalLevel:
		return "\033[31m" 
	default:
		return "\033[0m" 
	}
}


type coloredEncoder struct {
	zapcore.Encoder
}


func (e *coloredEncoder) EncodeEntry(entry zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {

	entry.Message = fmt.Sprintf("%s%s\033[0m", colorCode(entry.Level), entry.Message)
	return e.Encoder.EncodeEntry(entry, fields)
}


func New() (*Logger, error) {
	
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder      
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder 

	
	encoder := zapcore.NewConsoleEncoder(config.EncoderConfig)
	colored := &coloredEncoder{Encoder: encoder}

	
	core := zapcore.NewCore(colored, zapcore.AddSync(os.Stdout), zapcore.DebugLevel)
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	return &Logger{Logger: logger}, nil
}


func (l *Logger) Sync() {
	_ = l.Logger.Sync()
}
