package grpcapp

import (
	"fmt"
	"log/slog"
	"net"

	"github.com/imhasandl/grpc-go/sso/internal/grpc/auth"
	"google.golang.org/grpc"
)

type App struct {
	log  *slog.Logger
	gRPCServer *grpc.Server
	port int
}

func New(log *slog.Logger, authService auth.Auth, port int) *App {
	gRPCServer := grpc.NewServer()

	auth.Register(gRPCServer, authService)

	return &App{
		log:  log,
		gRPCServer: gRPCServer,
		port: port,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "grpcapp.Run"

	log := a.log.With(
		slog.String("op", op),
		slog.Int("port", a.port),
	)

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port)) //nolint:gosec // this is not a production code
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("grpc server is running", slog.String("addr", l.Addr().String()))

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *App) Stop() {
	const op ="grpcapp.Stop"

	a.log.With(slog.String("op", op)).
		Info("Stopping grpc server", slog.Int("port", a.port))

	a.gRPCServer.GracefulStop()
}
