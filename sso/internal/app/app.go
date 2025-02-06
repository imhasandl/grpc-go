package app

import (
	"log/slog"
	"time"

	"github.com/imhasandl/grpc-go/sso/internal/services/auth"
	"github.com/imhasandl/grpc-go/sso/internal/app/grpc"
	"github.com/imhasandl/grpc-go/sso/internal/storage/sqlite"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(log *slog.Logger, grpcPort int, storagePath string, tokenTTL time.Duration) *App {
	storage, err := sqlite.New(storagePath)
	if err != nil {
		panic(err)
	}

	authService := auth.New(log, storage, storage, storage, tokenTTL)
	
	grpcApp := grpcapp.New(log, authService,  grpcPort)

	return &App{
		GRPCSrv: grpcApp,
	}
}