package main

import (
	"context"
	"log"
	"os/signal"
	"rate-limiter/internal/services/api"
	"rate-limiter/internal/services/config"
	"syscall"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	app := api.NewApp(&cfg.App)

	errChan := make(chan error, 1)

	go func() {
		errChan <- app.Run()
	}()

	notifyCtx, notifyCancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer notifyCancel()

	select {
	case <-notifyCtx.Done():
		graceCtx, graceCancel := context.WithTimeout(context.Background(), cfg.Core.GracefulTimeout.Time)
		defer graceCancel()

		err := app.Stop(graceCtx)
		if err != nil {
			log.Println("error while stopping fiber application", err.Error())
		}
		log.Println("shutdown app by signal")
	case err := <-errChan:
		if err != nil {
			log.Println("shutdown - err from err chan", err.Error())
		} else {
			log.Println("server stopped")
		}
	}

}
