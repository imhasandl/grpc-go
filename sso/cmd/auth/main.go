package main

import (
	"log/slog"
	"os"

	"github.com/imhasandl/grpc-go/sso/internal/app"
	"github.com/imhasandl/grpc-go/sso/internal/config"
	"github.com/imhasandl/grpc-go/sso/internal/lib/logger/handlers/slogpretty"
)

const (
	envLocal = "local"
	envDev = "dev"
	envProd = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("starting application", slog.Any("config", cfg))

	app := app.New(log, cfg.GRPC.Port, cfg.StoragePath, cfg.TokenTTL)
	
	app.GRPCSrv.MustRun()
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal: 
		log = setupPrettySlog()
	case envDev: 
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd: 
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
