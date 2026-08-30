package middleware

import (
	"context"
	"log/slog"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func LoggingInterceptor(logger *slog.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()

		// Логируем начало запроса
		logger.InfoContext(ctx, "gRPC запрос начат",
			slog.String("method", info.FullMethod),
			slog.Any("request", req), // Будьте осторожны с большими запросами!
		)

		// Выполняем обработчик
		resp, err := handler(ctx, req)

		// Логируем результат
		duration := time.Since(start)
		attrs := []any{
			slog.String("method", info.FullMethod),
			slog.Duration("duration", duration),
		}

		if err != nil {
			st, _ := status.FromError(err)
			attrs = append(attrs,
				slog.String("error", err.Error()),
				slog.Int("code", int(st.Code())),
				slog.String("code_name", st.Code().String()),
			)
			logger.ErrorContext(ctx, "gRPC запрос завершен с ошибкой", attrs...)
		} else {
			logger.InfoContext(ctx, "gRPC запрос успешно завершен", attrs...)
		}

		return resp, err
	}
}
