package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

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
	
	go app.GRPCSrv.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	ch := <-stop

	log.Info("Stopping application", slog.String("signal", ch.String()))

	app.GRPCSrv.Stop()

	log.Info("application stopped")
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
