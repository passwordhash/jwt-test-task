package app

import (
	"context"
	"log/slog"

	httpApp "github.com/passwordhash/jwt-test-task/internal/app/http"
	"github.com/passwordhash/jwt-test-task/internal/config"
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

	httpApp := httpApp.New(
		ctx,
		log,
		cfg.HTTP,
	)

	return &App{
		HTTPSrv: httpApp,
	}
}
