package app

import (
	"context"
	"log/slog"

	httpApp "github.com/passwordhash/jwt-test-task/internal/app/http"
	"github.com/passwordhash/jwt-test-task/internal/config"
	"github.com/passwordhash/jwt-test-task/internal/service/auth"
	authSvc "github.com/passwordhash/jwt-test-task/internal/service/auth"
	authStorage "github.com/passwordhash/jwt-test-task/internal/storage/postgres/auth"
	postgresPkg "github.com/passwordhash/jwt-test-task/pkg/postgres"
)

type App struct {
	HTTPSrv *httpApp.App
}

func New(
	ctx context.Context,
	log *slog.Logger,
	cfg *config.Config,
) *App {
	postgresPool, err := postgresPkg.NewPool(ctx, cfg.PG.DSN())
	if err != nil {
		panic("failed to create postgres pool: " + err.Error())
	}

	authStorage := authStorage.New(postgresPool)

	authService := authSvc.New(
		log.WithGroup("service"),
		authStorage,
		auth.RefreshTokenManager{},
		cfg.App.AccessTTL,
		cfg.App.JWTSecret,
	)

	httpApp := httpApp.New(
		ctx,
		log,
		cfg.HTTP,
		authService,
	)

	return &App{
		HTTPSrv: httpApp,
	}
}
