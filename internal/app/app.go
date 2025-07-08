package app

import (
	"context"
	"log/slog"

	httpApp "github.com/passwordhash/jwt-test-task/internal/app/http"
	"github.com/passwordhash/jwt-test-task/internal/config"
	"github.com/passwordhash/jwt-test-task/internal/service/auth"
)

type App struct {
	HTTPSrv *httpApp.App
}

func New(
	ctx context.Context,
	log *slog.Logger,
	cfg *config.Config,
) *App {
	// TODO: connect to database

	authService := auth.NewService(
		log.WithGroup("service"),
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
