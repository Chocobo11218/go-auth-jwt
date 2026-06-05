package logger

import (
	"context"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type contextKey string

const (
	RequestIDKey contextKey = "request_id"
	UserIDKey    contextKey = "user_id"
)

var (
	log  *zap.Logger
	once sync.Once
)

func Initialize() {
	once.Do(func() {
		config := zap.NewProductionConfig()

		config.EncoderConfig.TimeKey = "timestamp"
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

		var err error

		log, err = config.Build(
			zap.AddCaller(),       // แสดง caller
			zap.AddCallerSkip(1),  // ข้าม logger.Info() wrapper
		)

		if err != nil {
			panic(err)
		}
	})
}

func getLogger() *zap.Logger {
	if log == nil {
		Initialize()
	}

	return log
}

func Info(ctx context.Context, message string, fields ...zap.Field) {
	getLogger().Info(message, appendContextFields(ctx, fields...)...)
}

func Debug(ctx context.Context, message string, fields ...zap.Field) {
	getLogger().Debug(message, appendContextFields(ctx, fields...)...)
}

func Fatal(ctx context.Context, message string, fields ...zapcore.Field) {
	getLogger().Fatal(message, appendContextFields(ctx, fields...)...)
}

func Warn(ctx context.Context, message string, fields ...zap.Field) {
	getLogger().Warn(message, appendContextFields(ctx, fields...)...)
}

func Error(ctx context.Context, message string, fields ...zap.Field) {
	getLogger().Error(message, appendContextFields(ctx, fields...)...)
}

func appendContextFields(ctx context.Context, fields ...zap.Field) []zap.Field {
	if ctx == nil {
		return fields
	}

	if requestID, ok := ctx.Value(RequestIDKey).(string); ok && requestID != "" {
		fields = append(fields, zap.String("request_id", requestID))
	}

	if userID, ok := ctx.Value(UserIDKey).(string); ok && userID != "" {
		fields = append(fields, zap.String("user_id", userID))
	}

	return fields
}

func Sync() {
	if log != nil {
		_ = log.Sync()
	}
}