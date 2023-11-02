package logger

import (
	"go.uber.org/zap"
)

type Logger struct {
	*zap.SugaredLogger
}

func New() (*Logger, error) {
	zapLogger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	sugar := zapLogger.Sugar()
	return &Logger{sugar}, nil
}

// Sync flushes any buffered log entries.
func (l *Logger) Sync() error {
	return l.SugaredLogger.Sync()
}
