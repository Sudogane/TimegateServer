package logger

import (
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Logger struct {
	*zap.SugaredLogger
}

func NewLogger() (*Logger, error) {
	l, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}
	newLogger := &Logger{
		l.Sugar(),
	}

	return newLogger, nil
}

func (l *Logger) WithSession(sessionId string, playerId uuid.UUID) *zap.SugaredLogger {
	return l.With("session_id", sessionId, "player_id", playerId)
}

func (l *Logger) Log(level string, message string) {
	switch level {
	case "ERROR":
		l.Errorw(message)
	case "WARN":
		l.Warnw(message)
	default:
		l.Infow(message)
	}
}
