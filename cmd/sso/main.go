package main

import (
	"log/slog"
	"os"
	"os/signal"
	"sso/internal/app"
	"sso/internal/config"
	repository_auth "sso/internal/repository"
	service_auth "sso/internal/services/auth"
	"syscall"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)
	log.Info("starting application", slog.String("env", cfg.Env))
	storagePath := os.Getenv("STORAGE_PATH")

	authRepository := repository_auth.NewAuthRepository(storagePath)
	authService := service_auth.NewAuthService(log, cfg.TokenTTL, authRepository)
	application := app.NewApp(log, cfg.GRPC.Port, cfg.StoragePath, cfg.TokenTTL, authService)

	go func() {
		if err := application.GRPCSrv.Run(); err != nil {
			log.Error("gRPC server error", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	application.GRPCSrv.Stop()

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
