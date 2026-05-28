package logger

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type contextKey string

const (
	RequestIDKey contextKey = "request_id"
	UserIDKey    contextKey = "user_id"
)

var log *zap.Logger

func Initialize() {
	config := zap.NewProductionConfig()

	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	var err error

	log, err = config.Build()
	if err != nil {
		panic(err)
	}
}

func Info(ctx context.Context, message string, fields ...zap.Field) {
	log.Info(message, appendContextFields(ctx, fields...)...)
}

func Debug(ctx context.Context, message string, fields ...zap.Field) {
	log.Debug(message, appendContextFields(ctx, fields...)...)
}

func Warn(ctx context.Context, message string, fields ...zap.Field) {
	log.Warn(message, appendContextFields(ctx, fields...)...)
}

func Error(ctx context.Context, message string, fields ...zap.Field) {
	log.Error(message, appendContextFields(ctx, fields...)...)
}

func appendContextFields(ctx context.Context, fields ...zap.Field) []zap.Field {

	if ctx == nil {
		return fields
	}

	if requestID, ok := ctx.Value(RequestIDKey).(string); ok {
		fields = append(fields, zap.String("request_id", requestID))
	}

	if userID, ok := ctx.Value(UserIDKey).(string); ok {
		fields = append(fields, zap.String("user_id", userID))
	}

	return fields
}
