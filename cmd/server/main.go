package main

import (
	"github.com/sha1sof/Echelon-/internal/app"
	"github.com/sha1sof/Echelon-/internal/config"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.MustLoadPath("./config/prod.yaml")

	log := setupLogger(cfg.Env)

	app := app.New(log, cfg.GRPCServer.Port, cfg.Storage.Path)

	go app.GRPCServ.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	sign := <-stop

	log.Info("stopping app", slog.String("signal", sign.String()))

	app.GRPCServ.Stop()

	log.Info("app stop")
}

// setupLogger инициализация логера.
func setupLogger(env string) *slog.Logger {
	var logger *slog.Logger

	switch env {
	case "prod":
		logger = slog.New(
			slog.NewJSONHandler(
				os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	case "local":
		logger = slog.New(
			slog.NewTextHandler(
				os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	}

	return logger
}
