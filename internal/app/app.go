package app

import (
	"log/slog"
	app_grpc "sso/internal/app/grpc"
	auth_grpc "sso/internal/grpc/auth"
	"time"
)

type App struct {
	GRPCSrv *app_grpc.App
}

func NewApp(
	log *slog.Logger,
	grpcPort int,
	storagePath string,
	tokenTTL time.Duration,
	authService auth_grpc.AuthService,
) *App {
	grpcApp := app_grpc.NewApp(log, grpcPort, authService)

	return &App{
		GRPCSrv: grpcApp,
	}
}
