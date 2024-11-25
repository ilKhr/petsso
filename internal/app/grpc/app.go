package grpcapp

import (
	"fmt"
	"log/slog"
	"net"

	authgrpc "github.com/ilkhr/petsso/grpc/auth"
	healthgrpc "github.com/ilkhr/petsso/grpc/health"
	"github.com/ilkhr/petsso/internal/services/auth"
	"github.com/ilkhr/petsso/internal/services/health"
	"google.golang.org/grpc"
)

const (
	ErrExpectListener = "expect listener but have nil"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(
	log *slog.Logger,
	authService *auth.Auth,
	healthService *health.Health,
	port int,
) *App {
	gRPCServer := grpc.NewServer()

	authgrpc.Register(gRPCServer, authService)
	healthgrpc.Register(gRPCServer, healthService)

	return &App{log: log, gRPCServer: gRPCServer, port: port}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "grpcapp.Run"

	log := a.log.With(slog.String("op", op), slog.Int("port", a.port))

	log.Info("starting gRPC server")

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("grpc server is running", slog.String("addr", l.Addr().String()))

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *App) RunTest(l net.Listener) error {
	const op = "grpcapp.RunTest"

	log := a.log.With(slog.String("op", op), slog.Int("port", a.port))

	log.Info("starting gRPC server")

	if l == nil {
		return fmt.Errorf("%s: %s", op, ErrExpectListener)
	}

	log.Info("grpc server is running", slog.String("addr", l.Addr().String()))

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *App) MustRunTest(l net.Listener) {
	if err := a.RunTest(l); err != nil {
		panic(err)
	}
}

func (a *App) Stop() {
	const op = "grpcapp.Stop"

	a.log.With(slog.String("op", op)).Info("stoping gRPC server", slog.Int("port", a.port))

	a.gRPCServer.GracefulStop()
}
