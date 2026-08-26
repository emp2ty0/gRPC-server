package app_grpc

import (
	"fmt"
	"log/slog"
	"net"
	auth_grpc "sso/internal/grpc/auth"

	"google.golang.org/grpc"
)

type App struct {
	log       *slog.Logger
	gRPCSever *grpc.Server
	port      int
}

func NewApp(
	log *slog.Logger,
	port int,
	authService auth_grpc.AuthService,
) *App {
	gRPCServer := grpc.NewServer()
	auth_grpc.Register(gRPCServer, authService)

	return &App{
		log:       log,
		gRPCSever: gRPCServer,
		port:      port,
	}
}

func (a *App) Run() error {
	const op = "app_grpc.Run"

	log := a.log.With(slog.String("op", op), slog.Int("port", a.port))

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	log.Info("gRPC server is running", slog.String("addr", l.Addr().String()))

	if err := a.gRPCSever.Serve(l); err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	return nil
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic("failed to run gRPC server:" + err.Error())
	}
}
func (a *App) Stop() {
	const op = "app_grpc.Stop"

	a.log.Info("shutting down gracefully...", slog.String("op", op))

	a.gRPCSever.GracefulStop()
}
