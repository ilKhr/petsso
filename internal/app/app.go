package app

import (
	"log/slog"
	"time"

	grpcapp "github.com/ilkhr/petsso/internal/app/grpc"
	"github.com/ilkhr/petsso/internal/services/auth"
	"github.com/ilkhr/petsso/internal/services/health"
	sqlite "github.com/ilkhr/petsso/internal/storage/sqllite"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(
	log *slog.Logger,
	grpcPort int,
	storagePath string,
	tokenTTL time.Duration,
) *App {
	storage, err := sqlite.New(storagePath)

	if err != nil {
		panic(err)
	}

	authService := auth.New(log, storage, storage, storage, tokenTTL)
	healthService := health.New(log)

	grpcApp := grpcapp.New(log, authService, healthService, grpcPort)

	return &App{
		GRPCSrv: grpcApp,
	}
}
