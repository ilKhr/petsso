package suite

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"

	"math/rand"

	ssov1 "github.com/ilKhr/petprotos/gen/go/sso"
	"github.com/ilkhr/petsso/internal/app"
	"github.com/ilkhr/petsso/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/test/bufconn"

	"github.com/golang-migrate/migrate/v4"

	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Suite struct {
	*testing.T
	Cfg        *config.Config
	AuthClient ssov1.AuthClient
}

func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()

	cfg := config.MustLoadByPath("../config/test.yaml")

	dsn := getUniqueDsn(cfg.StoragePath)

	ctx, cancelCtx := context.WithTimeout(context.Background(), cfg.GPRC.Timeout)

	log := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)

	l := bufconn.Listen(1024 * 1024)

	migrationsUp(dsn)

	application := app.New(log, cfg.GPRC.Port, dsn, cfg.TokenTTL)

	go application.GRPCSrv.MustRunTest(l)

	t.Cleanup(func() {
		t.Helper()
		cancelCtx()
		application.GRPCSrv.Stop()
	})

	// https://stackoverflow.com/a/78485897/14274230
	resolver.SetDefaultScheme("passthrough")

	cc, err := grpc.NewClient("bufnet", grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
		return l.Dial()
	}), grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		t.Fatalf("grpc create client failed: %v", err)
	}

	return ctx, &Suite{
		T:          t,
		Cfg:        cfg,
		AuthClient: ssov1.NewAuthClient(cc),
	}
}

func getUniqueDsn(storagePath string) string {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	randInt := rand.Intn(10000)

	dsn := fmt.Sprintf(storagePath, randInt)

	return dsn
}

func migrationsUp(storagePath string) {
	path, err := filepath.Abs("../migrations/")

	if err != nil {
		panic(err)
	}

	m, err := migrate.New("file:"+path, fmt.Sprintf("sqlite3://%s", storagePath))

	if err != nil {
		panic(err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no migrations to apply")

			return
		}

		panic(err)
	}

	fmt.Println("migrations applied successfully")
}
