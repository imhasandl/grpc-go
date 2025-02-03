package app

import (
	"log/slog"

	"github.com/imhasandl/grpc-go/sso/internal/grpc/auth"
	"google.golang.org/grpc"
)

type App struct {
	log  *slog.Logger
	grpc *grpc.Server
	port int
}

func New(log *slog.Logger, port int) *App {
	gRPCServer := grpc.NewServer()

	auth.RegisterAuthServer(gRPCServer, &auth.Server{})

	return &App{
		log: log,
		grpc: gRPCServer,
		port: port,
	}
}
