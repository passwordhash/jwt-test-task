package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/passwordhash/jwt-test-task/internal/app"
	"github.com/passwordhash/jwt-test-task/internal/config"
)

const shutdownTimeout = 5 * time.Second

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(),
		os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	cfg := config.MustLoad()

	log := config.SetupLogger(cfg.Env)

	application := app.New(ctx, log, cfg)

	go application.HTTPSrv.MustRun()

	<-ctx.Done()

	log.Info("received signal stop signal")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	application.HTTPSrv.Stop(shutdownCtx)

	log.Info("application stopped gracefully")
}
