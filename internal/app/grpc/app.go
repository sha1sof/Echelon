package grpcapp

import (
	"fmt"
	previewgrpc "github.com/sha1sof/Echelon-/internal/grpc/preview"
	"google.golang.org/grpc"
	"log/slog"
	"net"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

// New конструктор сервера
func New(log *slog.Logger, previewService previewgrpc.Thumbnail, port int) *App {
	gRPCServer := grpc.NewServer()

	previewgrpc.Register(gRPCServer, previewService)

	return &App{log, gRPCServer, port}
}

// MustRun метод для запуска сервера. Который паникует, если сервер не запустился.
func (a *App) MustRun() {
	if err := a.run(); err != nil {
		panic(err)
	}
}

// run запуск сервера.
func (a *App) run() error {
	const op = "grpcapp.Run"

	log := a.log.With(slog.String("op", op))

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("start gRPC server", slog.String("addr", lis.Addr().String()))

	if err := a.gRPCServer.Serve(lis); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// Stop остановка сервера.
func (a *App) Stop() {
	const op = "grpcapp.Stop"

	a.log.With(slog.String("op", op)).Info("stop gRPC server")

	a.gRPCServer.GracefulStop()
}
