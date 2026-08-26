package middleware

import (
	"context"
	"log/slog"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type contextKey string

const LoggerKey contextKey = "logger"

func UnaryLoggerInterceptor(logger *slog.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		ctx = context.WithValue(ctx, LoggerKey, logger)
		start := time.Now()

		resp, err := handler(ctx, req)

		duration := time.Since(start)
		statusCode := status.Code(err)

		// Подготавливаем атрибуты для логирования
		attrs := []any{
			slog.String("method", info.FullMethod),
			slog.Duration("duration", duration),
			slog.String("status", statusCode.String()),
		}

		// Логируем результат
		if err != nil {
			attrs = append(attrs, slog.String("error", err.Error()))
			logger.ErrorContext(ctx, "gRPC request failed", attrs...)
		} else {
			logger.InfoContext(ctx, "gRPC request completed", attrs...)
		}

		return resp, err
	}
}

func LoggerFromContext(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(LoggerKey).(*slog.Logger); ok && logger != nil {
		return logger
	}
	// Возвращаем дефолтный логгер, если не найден
	return slog.Default()
}
