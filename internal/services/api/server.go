package api

import (
	"context"
	"errors"
	"net/http"
	"rate-limiter/internal/services/config"

	"github.com/gofiber/fiber/v2"
)

type App struct {
	app  *fiber.App
	addr string
}

func NewApp(cfg *config.App) *App {
	app := fiber.New(fiber.Config{
		ReadTimeout:  cfg.ReadTimeout.Time,
		WriteTimeout: cfg.WriteTimeout.Time,
		IdleTimeout:  cfg.IdleTimeout.Time,
	})

	return &App{
		app:  app,
		addr: cfg.Addr,
	}
}

func (a *App) Run() error {
	err := a.app.Listen(a.addr)
	if err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}
	return nil
}

func (a *App) Stop(ctx context.Context) error {
	err := a.app.ShutdownWithContext(ctx)
	if err != nil {
		return err
	}
	return nil
}
