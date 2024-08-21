package app

import (
	grpcapp "github.com/sha1sof/Echelon-/internal/app/grpc"
	"github.com/sha1sof/Echelon-/internal/services/preview"
	"github.com/sha1sof/Echelon-/internal/storage/sqlite"
	"log/slog"
	"time"
)

type App struct {
	GRPCServ *grpcapp.App
}

// New конструктор приложения.
func New(log *slog.Logger, grpcPort int, storagePath string) *App {
	storage, err := sqlite.New(storagePath)
	if err != nil {
		panic(err)
	}

	previewService := preview.New(log, storage, storage, storage, time.Duration(10))

	grpcApp := grpcapp.New(log, previewService, grpcPort)

	return &App{
		GRPCServ: grpcApp,
	}
}
